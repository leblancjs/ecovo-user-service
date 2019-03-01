package vehicule

import (
	"fmt"

	"azure.com/ecovo/user-service/pkg/entity"
	"azure.com/ecovo/user-service/pkg/user"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves vehicules.
type UseCase interface {
	Register(v *entity.Vehicule, subID string) (*entity.Vehicule, error)
	FindByID(ID entity.ID) (*entity.Vehicule, error)
	FindByUserID(userID entity.ID) ([]*entity.Vehicule, error)
	Delete(ID entity.ID, userID entity.ID, subID string) error
}

// A Service handles the business logic related to vehicules.
type Service struct {
	repo     Repository
	uService user.UseCase
}

// NewService creates a vehicule service to handle business logic and manipulate
// vehicules through a repository.
func NewService(repo Repository, uService user.UseCase) *Service {
	return &Service{repo, uService}
}

// Register validates the vehicule's informartion and persists it in the repository.
func (s *Service) Register(v *entity.Vehicule, subID string) (*entity.Vehicule, error) {
	if v == nil {
		return nil, fmt.Errorf("vehicule.Service: vehicule is nil")
	}

	u, err := s.uService.FindBySubID(subID)
	if err != nil {
		return nil, err
	}

	if v.UserID != u.ID {
		return nil, WrongUserError{fmt.Sprintf("vehicule.Service: cannot add a vehicule to another user \"%s\"", v.ID)}
	}

	err = v.Validate()
	if err != nil {
		return nil, err
	}

	v.UserID = entity.ID(v.UserID)
	v.ID, err = s.repo.Create(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// FindByID retrieves the vehicule with the given ID in the repository, if it
// exists.
func (s *Service) FindByID(ID entity.ID) (*entity.Vehicule, error) {
	v, err := s.repo.FindByID(ID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return v, nil
}

// FindByUserID retrieves the multiple vehicules with the given user ID in the
// repository, if some exists.
func (s *Service) FindByUserID(userID entity.ID) ([]*entity.Vehicule, error) {
	v, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return v, nil
}

// Delete erases the vehicule from the repository.
func (s *Service) Delete(ID entity.ID, userID entity.ID, subID string) error {
	u, err := s.uService.FindBySubID(subID)
	if err != nil {
		return err
	}

	if userID != u.ID {
		return WrongUserError{fmt.Sprintf("vehicule.Service: cannot delete a vehicule of another user \"%s\"", ID)}
	}

	err = s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
