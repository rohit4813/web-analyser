package health

import "net/http"

// Read implements health api end point
func Read(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("OK!"))
}
