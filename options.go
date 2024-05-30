package pg

type Warnings struct {
	Empty   bool // Empty represents a logging option to warn if a given environment variable is set to any empty string. Requires [Options.Variables]. Defaults to false.
	Missing bool // Missing represents a logging option to warn if a given environment variable isn't found. Requires [Options.Variables]. Defaults to false.
}

// Options is the configuration structure optionally mutated via the [Variadic] constructor used throughout the package.
type Options struct {
	Variables []string  // Variables represents an array of environment variables (as returned by [os.Environ]), to selectively log.
	Warnings  *Warnings // Warnings represents logging options relating to [slog.Warn] logs. Defaults to a non-nil [Warnings] reference with all attributes set to false.
}

// Variadic represents a functional constructor for the [Options] type. Typical callers of Variadic won't need to perform
// nil checks as all implementations first construct an [Options] reference using packaged default(s).
//
//   - However, see [Settings] for construction of an [Options] type.
type Variadic func(o *Options)

// Settings represents a default constructor for [Options].
func Settings() *Options {
	return &Options{ // default Options constructor
		Variables: make([]string, 0),
		Warnings: &Warnings{
			Empty:   false,
			Missing: false,
		},
	}
}
