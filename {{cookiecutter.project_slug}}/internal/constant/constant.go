package constant

const (
	ServiceName = "{{ cookiecutter.project_slug }}"

	MaxPageSize = 1000
)

// Version, Commit, BuiltAt подставляются при сборке (Makefile, ldflags): манифест сервиса для
// pulse, build.commit — полный SHA коммита в проде.
var (
	Version = "dev"
	Commit  = ""
	BuiltAt = ""
)
