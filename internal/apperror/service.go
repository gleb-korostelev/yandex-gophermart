package apperror

import "errors"

var (
	ErrNoServerAddress        = errors.New("server address is empty")
	ErrNoDatabaseDestination  = errors.New("database destination is empty")
	ErrNoAccuralSystemAddress = errors.New("no address for accural system")
	ErrTokenInvalid           = errors.New("token is not valid")
	ErrLoginExists            = errors.New("user is already exists")
	ErrWrongPassword          = errors.New("wrong password")
	ErrGone                   = errors.New("user was deleted")
	ErrNoFunds                = errors.New("insufficient funds")
	ErrNotFound               = errors.New("order not found")
)
