package users

import (
	"fmt"
	"github.com/narinjtp/bookstore_users-api/datasources/mysql/users_db"
	"github.com/narinjtp/bookstore_users-api/utils/errors"
	"github.com/narinjtp/bookstore_users-api/utils/mysql_utils"
	"log"
)

const(
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES (?, ?, ?, ?, ?, ?);"
	queryGetUser = "SELECT id, first_name, last_name, email, date_created, status FROM users where id = ?;"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users where status = ?;"
	)
var(
	usersDB = make(map[int64]*User)
)
func (user *User)Get() *errors.RestErr{
	stmt, err := users_db.Client.Prepare(queryGetUser)
	log.Println(stmt)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
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
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)

	if saveErr != nil{
		log.Println(saveErr.Error())
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
		return errors.NewInternalServerError(err.Error())
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil{
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName,user.LastName,user.Email,user.Id)//dont care result, just error
	if err != nil {
		return mysql_utils.ParseError(err)
	}
	return  nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil{
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	if _, err = stmt.Exec(user.Id); err != nil {
		return mysql_utils.ParseError(err)
	}
	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr){
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next(){
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil{
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}