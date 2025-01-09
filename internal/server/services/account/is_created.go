package service_account

import "context"

func (s *Service) isAccountCreated(ctx context.Context, tgID uint64) (bool, error) {
	acc, err := s.dao.Find(ctx, tgID)
	if err != nil {
		return false, err
	}
	return len(acc) > 0, nil
}
