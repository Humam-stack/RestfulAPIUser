package services

import (
	"errors"
	"restfulapi/models"
	"restfulapi/repository"
	"restfulapi/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req models.UserLoginRequest) (string, error)
	Register(req models.UserRegisterRequest) (*models.UserResponse, error)
	GetAllUsers() ([]models.UserResponse, error)
	GetUserById(id uint) (*models.UserResponse, error)
	UpdateUser(id uint, req models.UserRegisterRequest) (*models.UserResponse, error)
	DeleteUser(id uint) error
}

type authService struct {
	userRepo repository.UserRepository
}

// constructor

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func toUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

func (s *authService) Register(req models.UserRegisterRequest) (*models.UserResponse, error) {

	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	//hashing password

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// buat user baru

	user := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hash),
	}

	// simpan ke db via repo
	if err := s.userRepo.Create(&user); err != nil {
		return nil, err
	}

	//retur respose
	respose := toUserResponse(&user)
	return respose, nil
}

func (s *authService) Login(req models.UserLoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	//perbandingang password yang sudah di hash dengan yang ada di database!
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	//gemerate jwt TokEn

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *authService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	//konversi ke response
	var response []models.UserResponse
	for _, user := range users {
		response = append(response, *toUserResponse(&user))
	}
	return response, nil
}

func (s *authService) GetUserById(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, errors.New("User tidak ditemukan")
	}

	response := toUserResponse(user)
	return response, nil
}

func (s *authService) UpdateUser(id uint, req models.UserRegisterRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return nil, errors.New("user tidak ditemuka")
	}

	//updated filed filed yang ada
	if req.Email != user.Email {
		exst, _ := s.userRepo.FindByEmail(req.Email)
		if exst != nil && exst.ID != id {
			return nil, errors.New("Email sudah terpakai")
		}
	}

	user.Name = req.Name
	user.Email = req.Email

	// kalau password diisi maka
	if req.Password != "" {
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedpassword)
	}

	// simpan ke db
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	response := toUserResponse(user)
	return response, nil
}

func (s *authService) DeleteUser(id uint) error {
	_, err := s.userRepo.FindById(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	return s.userRepo.Delete(id)
}
