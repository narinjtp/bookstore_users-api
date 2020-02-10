package users

import (
	"fmt"
	"github.com/narinjtp/bookstore_users-api/datasources/mysql/users_db"
	"github.com/narinjtp/bookstore_users-api/logger"
	"github.com/narinjtp/bookstore_users-api/utils/errors"
	"github.com/narinjtp/bookstore_users-api/utils/mysql_utils"
	"log"
	"strings"
)

const(
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created, status FROM users where id = ?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id=?;"
	queryFindByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users where status = ?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users where email = ? and password = ? and status = ?;"
	)
var(
	usersDB = make(map[int64]*User)
)
func (user *User)Get() *errors.RestErr{
	stmt, err := users_db.Client.Prepare(queryGetUser)
	log.Println(stmt)
	if err != nil {
		logger.Error("error when try to prepare get user statement",err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	//select all rows
	//results, _:= stmt.Query(user.Id)
	//if err != nil {
	//	return errors.NewInternalServerError(err.Error())
	//}
	//defer results.Close()

	if getError := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getError != nil {
		logger.Error("error when try to get user statement",getError)
		return mysql_utils.ParseError(getError)
		//sqlErr, ok := getError.(*mysql.MySQLError)
		//if !ok {
		//	return errors.NewInternalServerError(getError.Error())
		//}
		//fmt.Println(sqlErr.Number)
		//fmt.Println(sqlErr.Message)
		//if strings.Contains(err.Error(),errorNoRow){
		//	return errors.NewInternalServerError(
		//		fmt.Sprintf("user %d not found",user.Id))
		//}
		//fmt.Println(err)
		//return errors.NewInternalServerError(
		//	fmt.Sprintf("error when trying to get user %d: %s",user.Id,err.Error()))
	}
	//result := usersDB[user.Id]
	//if result == nil {
	//	return errors.NewNotFoundError(fmt.Sprintf(("user %d not found"),user.Id))
	//}
	//user.Id = result.Id
	//user.FirstName = result.FirstName
	//user.LastName = result.LastName
	//user.Email = result.Email
	//user.DateCreated = result.DateCreated
	return nil
}

func (user *User)Save() *errors.RestErr{
	//if nil == users_db.Client {
	//	log.Println("userdb nil")
	//	return errors.NewInternalServerError("userdb nil")
	//}
	//this stmt try to insert
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	//log.Println(stmt)
	if err != nil {
		logger.Error("error when try to prepare save user statement",err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)

	if saveErr != nil{
		//log.Println(saveErr.Error())
		logger.Error("error when try to save user statement",saveErr)
		return mysql_utils.ParseError(saveErr)
		//another catch error method
		//if strings.Contains(err.Error(),indexUniqueEmail){
		//	return errors.NewBadRequestError(fmt.Sprintf("email %s already exists",user.Email))
		//}
		//return errors.NewInternalServerError(err.Error())
	}
	//end

	//but this stmt insert only not recommended
	//insertResult, err := users_db.Client.Exec(queryInsertUser,user.FirstName, user.LastName, user.Email, user.DateCreated)


	userId, err := insertResult.LastInsertId()
	if err != nil{
		logger.Error("error when try to get last insert id after creating a new user",err)
		return errors.NewInternalServerError("database error")
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil{
		logger.Error("error when try to prepare update user statement",err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName,user.LastName,user.Email,user.Id)//dont care result, just error
	if err != nil {
		logger.Error("error when try to update user statement",err)
		return mysql_utils.ParseError(err)
	}
	return  nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil{
		logger.Error("error when try to prepare delete user statement",err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		logger.Error("error when try to prepare update user statement",err)
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr){
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when try to prepare find user statement",err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		logger.Error("error when try to find user statement",err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next(){
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil{
			logger.Error("error when scan user row into user struct",err)
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User)FindByEmailAndPassword() *errors.RestErr{
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	log.Println(stmt)
	if err != nil {
		logger.Error("error when try to prepare get user by email amd password",err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if getError := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getError != nil {
		if strings.Contains(getError.Error(), mysql_utils.ErrorNoRows){
			return errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when try to get user by email amd password",getError)
		return errors.NewInternalServerError("database error")
	}
	return nil
}