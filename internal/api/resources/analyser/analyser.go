package analyser

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
)

type Analyser struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Analyser {
	return &Analyser{
		logger: logger,
	}
}

// Index serves the index  page.
func (a *Analyser) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("index")
	//h.logger.Error().Str(logger.KeyReqID, uCtx.RequestID(r.Context())).Err(errors.New("test")).Msg("")
	//e.ServerError(w, e.Error)
}
