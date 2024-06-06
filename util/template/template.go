package template

import (
	"io"
)

type Template interface {
	ExecuteTemplate(io.Writer, string, any) error
}
