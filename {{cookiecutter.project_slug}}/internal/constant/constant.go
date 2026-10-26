package constant

const (
	ServiceName = "{{ cookiecutter.project_slug }}"

	MaxPageSize = 1000
)

// Version, Commit, BuiltAt подставляются при сборке (Makefile, ldflags): что именно запущено —
// версия, полный SHA коммита, время сборки (пишутся в лог при старте).
var (
	Version = "dev"
	Commit  = ""
	BuiltAt = ""
)
