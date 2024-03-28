package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	ID     int    `json:"id"`
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginClaims struct {
	UserID        int   `json:"userID"`
	ExpiresAt     int   `json:"expiresAt"`
	AccountNumber int64 `json:"accountNumber"`
	jwt.StandardClaims
}

type TransferRequest struct {
	FromAccount int `json:"fromAccount"`
	ToAccount   int `json:"toAccount"`
	Amount      int `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Number            int64     `json:"number"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (a *Account) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(password)) == nil
}

func ValidateStrLength(value, field string, length int) error {
	if len(value) < length {
		return fmt.Errorf("%s must be greater than %v characters", field, length)
	}

	return nil
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	err := ValidateStrLength(firstName, "First Name", 3)
	if err != nil {
		return nil, err
	}
	err = ValidateStrLength(lastName, "Last Name", 3)
	if err != nil {
		return nil, err
	}
	err = ValidateStrLength(password, "Password", 8)
	if err != nil {
		return nil, err
	}

	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1000000)),
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
