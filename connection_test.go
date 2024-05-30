package pg_test

import (
	"testing"

	"github.com/x-ethr/pg"
)

func Test(t *testing.T) {
	t.Run("Test", func(t *testing.T) {
		t.Logf("URI: %s", pg.DSN())
	})
}
