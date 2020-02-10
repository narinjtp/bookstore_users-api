package mysql_utils

import (
	"github.com/go-sql-driver/mysql"
	"github.com/narinjtp/bookstore_users-api/utils/errors"
	"strings"
)
const(
	ErrorNoRows = "no rows in result set"
)

func ParseError(err error) *errors.RestErr{
	sqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		//fmt.Printf(sqlError.Error())
		if strings.Contains(err.Error(),ErrorNoRows){
			return errors.NewNotFoundError("no record matching given id")
		}
		return  errors.NewInternalServerError("error parsing database response")
	}

	switch sqlError.Number{
	case 1062:
		return errors.NewBadRequestError("invalid data")
	}
	return  errors.NewInternalServerError("error processing request")
}
