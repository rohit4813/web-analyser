package error

import "net/http"

var Error = []byte("Error")

func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(error)
}

func BadRequest(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write(error)
}

func ValidationErrors(w http.ResponseWriter, reps []byte) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(reps)
}
