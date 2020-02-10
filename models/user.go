package models

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("Invalid Login")
)

type User struct {
	key string
}

func NewUser(username string, hash []byte) (*User, error) {
	id, err := client.Incr("user:next-id").Result()

	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", username)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", username, id)

	_, err = pipe.Exec()

	if err != nil {
		return nil, err
	}

	return &user{key}, nil
}

func (user *User) GetUsername() (string, error) {
	return client.HGet(user.key, "username").Result()
}

func (user *User) GetHash() ([]byte, error) {
	return client.HGet(user.key, "hash").Bytes()
}

func (user *User) Authenticate(password string) error {

}

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
