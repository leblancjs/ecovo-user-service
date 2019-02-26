package user

import (
	"azure.com/ecovo/user-service/pkg/entity"
)

// Repository is an interface representing the ability to perform CRUD
// operations on users in a database.
type Repository interface {
	FindByID(ID entity.ID) (*entity.User, error)
	FindBySubID(subID string) (*entity.User, error)
	Create(user *entity.User) (entity.ID, error)
	Update(user *entity.User) error
	Delete(ID entity.ID) error
}
