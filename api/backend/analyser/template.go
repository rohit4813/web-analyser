package analyser

import (
	"io"
)

type Template interface {
	ExecuteTemplate(io.Writer, string, any) error
}
