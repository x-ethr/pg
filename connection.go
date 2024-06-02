package pg

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DSN represents the postgresql connection string.
// pool_max_conns: integer greater than 0
// pool_min_conns: integer 0 or greater
// pool_max_conn_lifetime: duration string
// pool_max_conn_idle_time: duration string
// pool_health_check_period: duration string
// pool_max_conn_lifetime_jitter: duration string
//
//   - https://www.postgresql.org/docs/current/libpq-envars.html
//   - https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-PARAMKEYWORDS
func DSN() (v string) {
	uri := url.URL{
		Scheme: "postgresql",
		Host:   os.Getenv("PGHOST"),
	}

	if uri.Host == "" {
		uri.Host = "localhost"
	}

	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	port := os.Getenv("PGPORT")
	if port == "" {
		port = "5432"
	}

	timeout := os.Getenv("PGCONNECT_TIMEOUT")
	if timeout == "" {
		timeout = "10"
	}

	application := os.Getenv("PGAPPNAME")

	sslmode := os.Getenv("PGSSLMODE")
	root := os.Getenv("PGSSLROOTCERT")

	maxconnections := os.Getenv("PGPOOLMAXCONNECTIONS")
	if maxconnections == "" {
		value, cpu := 4, runtime.NumCPU()
		if value < cpu {
			value = cpu
		}

		maxconnections = strconv.Itoa(value)
	}

	minconnections := os.Getenv("PGPOOLMINCONNECTIONS")
	if minconnections == "" {
		minconnections = strconv.Itoa(1)
	}

	tz := os.Getenv("PGTZ")
	if tz == "" {
		tz = "UTC"
	}

	query := uri.Query()
	query.Add("user", user)
	query.Add("password", password)
	query.Add("port", port)
	query.Add("connect_timeout", timeout)

	query.Add("application_name", application)

	query.Add("pool_max_conns", maxconnections)
	query.Add("pool_min_conns", minconnections)

	query.Add("sslmode", sslmode)
	query.Add("sslrootcert", root)

	for key, values := range query {
		if len(values) >= 1 && strings.TrimSpace(values[0]) == "" {
			query.Del(key)
		}
	}

	uri.RawQuery = query.Encode()

	return uri.String()
}

var Pool atomic.Pointer[pgxpool.Pool]

// Connection establishes a connection to the database using pgxpool.
// If a connection pool does not exist, a new one is created and stored in the pool variable.
// Returns a connection from the connection pool.
// If an error occurs during connection creation, nil and the error are returned.
func Connection(ctx context.Context, uri string) (*pgxpool.Conn, error) {
	if Pool.Load() == nil {
		configuration, e := pgxpool.ParseConfig(uri)
		if e != nil {
			slog.ErrorContext(ctx, "Unable to Generate Configuration from DSN String", slog.String("error", e.Error()))
			return nil, e
		}

		instance, e := pgxpool.NewWithConfig(ctx, configuration)
		if e != nil {
			slog.ErrorContext(ctx, "Unable to Establish Pool Connection to Database", slog.String("error", e.Error()))
			return nil, e
		}

		Pool.Store(instance)
	}

	return Pool.Load().Acquire(ctx)
}

// Disconnect closes the transaction and releases the connection back to the pool.
// If `tx` is not nil, it rolls back the transaction and logs any error.
// If `connection` is not nil, it releases the connection back to the pool.
func Disconnect(ctx context.Context, connection *pgxpool.Conn, tx pgx.Tx) {
	if tx != nil {
		e := tx.Rollback(ctx)
		if e != nil && !(errors.Is(e, pgx.ErrTxClosed)) {
			slog.ErrorContext(ctx, "Error Rolling Back Transaction", slog.String("error", e.Error()))
		} else if e != nil && (errors.Is(e, pgx.ErrTxClosed)) {
			slog.InfoContext(ctx, "Successfully Committed Database Transaction")
		} else if e == nil {
			slog.InfoContext(ctx, "Successfully Rolled Back Database Transaction")
		}
	}

	if connection != nil {
		connection.Release()
	}
}
