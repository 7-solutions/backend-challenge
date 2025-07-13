package usecases

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/testGolang/backend-challenge/application/ports"
	"github.com/testGolang/backend-challenge/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo ports.UserRepository
}

func NewUserUseCase(r ports.UserRepository) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) Register(name, email, password string) (*domain.User, error) {
	user := &domain.User{ID: primitive.NewObjectID().Hex(), Name: name, Email: email, Password: password, CreatedAt: time.Now()}
	// รับ plainPwd เป็นรหัสผ่านที่ผู้ใช้กรอกมา
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}
	// แปลงเป็น string แล้วเก็บลง struct or ฐานข้อมูล
	user.Password = string(hashedBytes)
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *UserUseCase) Authenticate(email, plainPwd string) (string, error) {
	// fetch the user by email
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	// verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPwd)); err != nil {
		return "", errors.New("invalid credentials")
	}
	// create JWT token
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (uc *UserUseCase) List() ([]*domain.User, error) {
	return uc.repo.GetAll()
}

func (uc *UserUseCase) Get(id string) (*domain.User, error) {
	return uc.repo.GetByID(id)
}

func (uc *UserUseCase) Update(id, name, email, password string) (*domain.User, error) {
	u, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		u.Password = string(hashed)
	}
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	if err := u.Validate(); err != nil {
		return nil, err
	}
	err = uc.repo.Update(u)
	return u, err
}

func (uc *UserUseCase) Delete(id string) error {
	return uc.repo.Delete(id)
}
