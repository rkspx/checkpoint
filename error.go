package checkpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrObjectLocked     = errors.New("Object is locked by another session")
	ErrObjectNameExists = errors.New("Object with the same name exists")
	ErrIPAddressExists  = errors.New("Object with the same ip address exists")
	ErrObjectNotFound   = errors.New("Object not found")

	ErrTimeoutExceeded = errors.New("Timeout Exceeded")
)

type FailResponse struct {
	Message        string        `json:"message"`
	Warnings       []ErrorObject `json:"warnings"`
	Errors         []ErrorObject `json:"errors"`
	BlockingErrors []ErrorObject `json:"blocking-errors"`
	Code           string        `json:"code"`
	StatusCode     int
}

type ErrorObject struct {
	CurrentSession bool   `json:"current-session"`
	Message        string `json:"message"`
}

func (r FailResponse) Error() string {
	return fmt.Sprintf("%#v", r)
}

func (r *FailResponse) fromReader(i io.Reader) error {
	return json.NewDecoder(i).Decode(r)
}

func parseError(err FailResponse) error {
	if isObjectLockedError(err) {
		return ErrObjectLocked
	}

	if isObjectNameExistError(err) {
		return ErrObjectNameExists
	}

	if isIPAddressEixtsError(err) {
		return ErrIPAddressExists
	}

	if isObjectNotFoundError(err) {
		return ErrObjectNotFound
	}

	return err
}

func isObjectLockedError(err FailResponse) bool {
	return err.Code == "generic_err_object_locked" || (err.Code == "generic_error" && strings.Contains(err.Message, "is locked by another session."))
}

func isObjectNameExistError(err FailResponse) bool {
	return err.Code == "err_validation_failed" && err.StatusCode == 0 && strings.Contains(err.Message, "More than one object named")
}

func isIPAddressEixtsError(err FailResponse) bool {
	return err.Code == "err_validation_failed" && err.StatusCode == 0 && strings.Contains(err.Message, "Multiple objects have the same IP address ")
}

func isObjectNotFoundError(err FailResponse) bool {
	return err.Code == "generic_err_object_not_found" && err.StatusCode == 0
}

func isSessionLimitReached(err FailResponse) bool {
	return strings.Contains(err.Message, "sk113955")
}
