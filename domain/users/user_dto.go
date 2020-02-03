package users

import (
	"github.com/narinjtp/bookstore_users-api/utils/errors"
	"strings"
)

type User struct{
	Id int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string  `json:"last_name"`
	Email string  `json:"email"`
	DateCreated string  `json:"date_created"`
}
/*
func Validate() *errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == ""{
		return errors.NewBadRequestError("invalid email address")
	}
	return nil
}
*/
func (user *User) Validate() *errors.RestErr{
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == "" {
		return errors.NewBadRequestError("invalid email address")
	}
	return nil
}