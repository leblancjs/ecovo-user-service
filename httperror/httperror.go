package httperror

import (
	"encoding/json"
	"log"
	"net/http"

	"azure.com/ecovo/user-service/requestid"
)

type HTTPError struct {
	StatusCode int    `json:"code"`
	Error      string `json:"error"`
	Cause      error  `json:"-"`
	RequestID  string `json:"requestId"`
}

const (
	ErrInternalServerError = "Something went wrong processing your request. Please contact a system administrator."
	ErrUserNotFound        = "No user exists with the given ID"
	ErrBadRequest          = "Malformed request body"
)

func New(statusCode int, err string, cause error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Error:      err,
		Cause:      cause,
	}
}

func NewInternalServerError(err string, cause error) *HTTPError {
	return New(http.StatusInternalServerError, err, cause)
}

func NewNotFoundError(err string, cause error) *HTTPError {
	return New(http.StatusNotFound, err, cause)
}

func NewBadRequestError(err string, cause error) *HTTPError {
	return New(http.StatusBadRequest, err, cause)
}

func Handler(w http.ResponseWriter, r *http.Request, httpErr *HTTPError) {
	requestID, err := requestid.FromContext(r.Context())
	if err != nil {
		requestID = "hidethepain"
	}
	httpErr.RequestID = requestID

	log.Printf("[Request ID %s] %s", requestID, httpErr.Cause.Error())

	w.WriteHeader(httpErr.StatusCode)
	json.NewEncoder(w).Encode(httpErr)
}
