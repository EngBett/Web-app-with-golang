package models

import "fmt"

type Update struct {
	key string
}

func NewUpdate(userId int64, body string) (*User, error) {
	id, err := client.Incr("update:next-id").Result()

	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("update:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "user_id", userId)
	pipe.HSet(key, "body", body)
	pipe.LPush("updates", id)

	_, err = pipe.Exec()

	if err != nil {
		return nil, err
	}

	return &User{key}, nil
}

//fetch comments
func GetUpdates() ([]string, error) {
	return client.LRange("updates", 0, 10).Result()
}

//post comment
func PostUpdate(body string) error {
	return client.LPush("updates", body).Err()
}
