package user

import (
	"fmt"

	"azure.com/ecovo/user-service/pkg/entity"
)

// UseCase is an interface representing the ability to handle the business
// logic that involves users.
type UseCase interface {
	Register(u *entity.User) (*entity.User, error)
	Update(modifiedUser *entity.User) error
	FindByID(ID entity.ID) (*entity.User, error)
	FindBySubID(subID string) (*entity.User, error)
	Delete(ID entity.ID) error
}

// A Service handles the business logic related to users.
type Service struct {
	repo Repository
}

// NewService creates a user service to handle business logic and manipulate
// users through a repository.
func NewService(repo Repository) *Service {
	return &Service{repo}
}

// Register validates the user's personal informartion, makes it move on to the
// next sign up phase, and persists it in the repository.
func (s *Service) Register(u *entity.User) (*entity.User, error) {
	if u == nil {
		return nil, fmt.Errorf("user.Service: user is nil")
	}

	_, err := s.FindBySubID(u.SubID)
	if err == nil {
		return nil, AlreadyExistsError{fmt.Sprintf("user.Service: user already exists with ID \"%s\"", u.SubID)}
	}

	u.SignUpPhase = entity.SignUpPhasePreferences

	err = u.Validate()
	if err != nil {
		return nil, err
	}

	u.ID, err = s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindByID retrieves the user with the given ID in the repository, if it
// exists.
func (s *Service) FindByID(ID entity.ID) (*entity.User, error) {
	u, err := s.repo.FindByID(ID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return u, nil
}

// FindBySubID retrieves the user with the given subscription ID in the
// repository, if it exists.
func (s *Service) FindBySubID(subID string) (*entity.User, error) {
	u, err := s.repo.FindBySubID(subID)
	if err != nil {
		return nil, NotFoundError{err.Error()}
	}

	return u, nil
}

// Update validates that the user contains all the required personal
// information, that all values are correct and well formatted, and persists
// the modified user in the repository.
func (s *Service) Update(modifiedUser *entity.User) error {
	if modifiedUser == nil {
		return fmt.Errorf("user.Service: modified user is nil")
	}

	u, err := s.repo.FindByID(entity.ID(modifiedUser.ID))
	if err != nil {
		return NotFoundError{err.Error()}
	}

	if modifiedUser.FirstName != "" {
		u.FirstName = modifiedUser.FirstName
	}

	if modifiedUser.LastName != "" {
		u.LastName = modifiedUser.LastName
	}

	if !modifiedUser.DateOfBirth.IsZero() {
		u.DateOfBirth = modifiedUser.DateOfBirth
	}

	if modifiedUser.PhoneNumber != "" {
		u.PhoneNumber = modifiedUser.PhoneNumber
	}

	if modifiedUser.Gender != "" {
		u.Gender = modifiedUser.Gender
	}

	if modifiedUser.Photo != "" {
		u.Photo = modifiedUser.Photo
	}

	if modifiedUser.Description != "" {
		u.Description = modifiedUser.Description
	}

	if modifiedUser.Preferences != nil {
		if u.Preferences == nil {
			u.Preferences = &entity.Preferences{}
		}

		u.Preferences.Smoking = modifiedUser.Preferences.Smoking
		u.Preferences.Conversation = modifiedUser.Preferences.Conversation
		u.Preferences.Music = modifiedUser.Preferences.Music
	}

	if modifiedUser.SignUpPhase != "" {
		u.SignUpPhase = modifiedUser.SignUpPhase
	}

	if modifiedUser.UserRating <= 0 && modifiedUser.UserRating >= 5 {
		u.UserRating = modifiedUser.UserRating
	}

	if modifiedUser.DriverRating <= 0 && modifiedUser.DriverRating >= 5 {
		u.DriverRating = modifiedUser.DriverRating
	}

	err = u.Validate()
	if err != nil {
		return err
	}

	err = s.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

// Delete erases the user from the repository.
func (s *Service) Delete(ID entity.ID) error {
	err := s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
