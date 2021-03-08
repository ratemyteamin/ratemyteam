package common

import (
	"context"
	"errors"
)

func (s *Session) CreateNewUser(ctx context.Context, tx Transaction, user *User) error{
	log := GetLogFromCtx(ctx)
	log.Infof("Creating new User with the request")
	status, err := s.Store.StoreUser(user)
	if err != nil {
		log.Errorf("unable to create user: %v", err)
		return err
	}
	if status{
		log.Debugf("stored user with user id %s", user.UserId)
	}else {
		log.Debugf("Failed to Store User with ", user.UserId)
		return errors.New("Failed to Store in DB")
	}
	return nil
}
