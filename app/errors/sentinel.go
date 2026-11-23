package errors

import "errors"

const (
	ErrNilInputMsg               = "nil input"
	ErrContextIsNilMsg           = "context is nil"
	ErrInvalidInputMsg           = "invalid input"
	ErrInvalidClaimsMsg          = "invalid claims"
	ErrMissingClaimsOnContextMsg = "missing claims on context"
	ErrMissingAuthorizationMsg   = "missing authorization header"
	ErrUnauthorizedMsg           = "unauthorized"
)

var (
	ErrNilInput               = errors.New(ErrNilInputMsg)
	ErrInvalidInput           = errors.New(ErrInvalidInputMsg)
	ErrContextIsNil           = errors.New(ErrContextIsNilMsg)
	ErrInvalidClaims          = errors.New(ErrInvalidClaimsMsg)
	ErrMissingClaimsOnContext = errors.New(ErrMissingClaimsOnContextMsg)
	ErrMissingAuthorization   = errors.New(ErrMissingAuthorizationMsg)
	ErrUnauthorized           = errors.New(ErrUnauthorizedMsg)
)
