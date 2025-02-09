package collector

import "github.com/daemondxx/lks_back/entity"

type ErrLimitAttempt struct {
	Credentials []*entity.Credential
}

func (e *ErrLimitAttempt) Error() string {
	return "attempt limit has been reached"
}

func newErrLimitAttempt(u []*entity.Credential) *ErrLimitAttempt {
	e := &ErrLimitAttempt{
		Credentials: u,
	}
	return e
}
