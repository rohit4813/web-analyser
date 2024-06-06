package error

// Error defines a struct for custom error, containing msg and http status code
type Error struct {
	Msg            string
	HTTPStatusCode int
}

type Message string

const (
	UnreachableURLError Message = "The URL provided is not reachable, " +
		"please check your internet connection and ensure that the URL is correct"
	InvalidURLError Message = "Invalid URL provided, please ensure the URL format is correct, " +
		"for example: https://www.google.com"
)
