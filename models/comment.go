package models

//fetch comments
func GetComments() ([]string, error) {
	return client.LRange("comments", 0, 10).Result()
}

//post comment
func PostComment(comment string) error {
	return client.LPush("comments", comment).Err()
}
