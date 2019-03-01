package vehicule

import (
	"azure.com/ecovo/user-service/pkg/entity"
)

// Repository is an interface representing the ability to perform CRUD
// operations on vehicules in a database.
type Repository interface {
	FindByID(ID entity.ID) (*entity.Vehicule, error)
	FindByUserID(userID entity.ID) ([]*entity.Vehicule, error)
	Create(user *entity.Vehicule) (entity.ID, error)
	Delete(ID entity.ID) error
}
