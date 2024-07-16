package services

import "fmt"

type ErrServ struct {
	Service string
	Message string
	Err     error
}

func (e *ErrServ) Error() string {
	return fmt.Sprintf("[%s] - %s: %e", e.Service, e.Message, e.Err)
}

func (e *ErrServ) Unwrap() error {
	return e.Err
}
