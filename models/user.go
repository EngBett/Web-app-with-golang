package models

import (
	"errors"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("Invalid Login")
)

func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	if err != nil {
		return err
	}

	return client.Set("User:"+username, hash, 0).Err()
}

func AuthenticatesUser(username, password string) error {
	hash, err := client.Get("User:" + username).Bytes()

	//user not found || error 500
	if err == redis.Nil {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))

	// Wrong password
	if err != nil {
		return ErrInvalidLogin
	}

	return nil
}
