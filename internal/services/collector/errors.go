package collector

import "github.com/daemondxx/lks_back/entity"

type ErrLimitAttempt struct {
	Users []entity.User
}

func (e *ErrLimitAttempt) Error() string {
	return "attempt limit has been reached"
}

func newErrLimitAttempt(u []*entity.User) *ErrLimitAttempt {
	e := &ErrLimitAttempt{
		Users: make([]entity.User, 0, len(u)),
	}

	for _, i := range u {
		e.Users = append(e.Users, *i)
	}

	return e
}
