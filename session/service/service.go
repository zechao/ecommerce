package user

import (
	"context"

	"github.com/zechao158/ecomm/session"
)

type Repository interface {
	Create(ctx context.Context, user *session.User) error
	CreateHistory(ctx context.Context, history *session.History) error
}

type Service struct {
	session    session.Session
	repository Repository
}

func NewService(session session.Session, repository Repository) *Service {
	return &Service{
		session:    session,
		repository: repository,
	}
}

// Register create a new user with an associated "register" history.
func (s *Service) Register(ctx context.Context, user *session.User) (*session.User, error) {

	err := s.session.Transaction(ctx, func(ctx context.Context) error {
		// You can also call another service from here, not necessarily a repository.
		var err error
		err = s.repository.Create(ctx, user)
		if err != nil {
			return err
		}

		history := &session.History{
			UserID: user.ID,
			Action: "register",
		}
		err = s.repository.CreateHistory(ctx, history)
		return err
	})
	if err != nil {
		return user, nil
	}

	return user, nil
}
