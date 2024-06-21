package error

// CustomError defines a struct for custom error, containing message and http status code
type CustomError struct {
	Message        string
	HttpStatusCode int
}

type Msg string

const (
	UnreachableURLError Msg = "The URL provided is not reachable, " +
		"please check your internet connection and ensure that the URL is correct"
	InvalidURLError Msg = "Invalid URL provided, please ensure the URL format is correct, " +
		"for example: https://www.google.com"
)
