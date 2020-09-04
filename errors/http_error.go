package errors

import (
	"net/http"
)

type HttpError interface {
	error
	Code() int
}

type statusBadRequest struct {
	error
}

func (e statusBadRequest) Code() int {
	return http.StatusBadRequest
}

type statusUnauthorized struct {
	error
}

func (e statusUnauthorized) Code() int {
	return http.StatusUnauthorized
}

type statusForbidden struct {
	error
}

func (e statusForbidden) Code() int {
	return http.StatusForbidden
}

type statusNotFound struct {
	error
}

func (e statusNotFound) Code() int {
	return http.StatusNotFound
}

type statusInternalServerError struct {
	error
}

func (e statusInternalServerError) Code() int {
	return http.StatusInternalServerError
}

func ErrParameter(err error) error {
	return statusBadRequest{err}
}

func ErrUnknown(err error) error {
	return statusInternalServerError{err}
}
