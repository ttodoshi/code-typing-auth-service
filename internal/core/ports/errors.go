package ports

import "errors"

var (
	BadRequestError     = errors.New("bad request")
	UnauthorizedError   = errors.New("unauthorized")
	ForbiddenError      = errors.New("forbidden")
	NotFoundError       = errors.New("not found")
	InternalServerError = errors.New("internal server error")
)
