package errs

type Err string

func (e Err) Error() string {
	return string(e)
}

// common errors
const (
	NoRows            = Err("err_no_rows")
	ServiceNA         = Err("service_not_available")
	NotAuthorized     = Err("not_authorized")
	ObjectNotFound    = Err("object_not_found")
	IncorrectPageSize = Err("incorrect_page_size")
)

// ErrFull
type ErrFull struct {
	Err    error
	Desc   string
	Fields map[string]string
}

func (e ErrFull) Error() string {
	return e.Err.Error() + ", desc: " + e.Desc
}
