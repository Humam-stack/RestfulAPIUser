package repository

import (
	"errors"
	"restfulapi/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id uint) error
	FindById(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) Update(user *models.User) error {
	return u.db.Save(user).Error
}

func (u *userRepository) Delete(id uint) error {
	return u.db.Delete(&models.User{}, id).Error
}

func (u *userRepository) FindById(id uint) (*models.User, error) {
	var user models.User
	err := u.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record not found")
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("record nott found")
		}
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := u.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
