package routing

// ResponseStatus - struct representing error response from server
type ResponseStatus struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

// ErrorRenderer - returns ResponseStatus for given error err, with needed statusText & statusCode
func ErrorRenderer(err error, statusText string, statusCode int) *ResponseStatus {
	return &ResponseStatus{
		Err:        err,
		StatusCode: statusCode,
		StatusText: statusText,
		Message:    err.Error(),
	}
}

// ErrorRendererDefault - returns ResponseStatus with status code 400 - Bad request error
func ErrorRendererDefault(err error) *ResponseStatus {
	return &ResponseStatus{
		Err:        err,
		StatusCode: 400,
		StatusText: "Bad request",
		Message:    err.Error(),
	}
}
