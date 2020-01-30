package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/narinjtp/bookstore_users-api/domain/users"
	"github.com/narinjtp/bookstore_users-api/services"
	"github.com/narinjtp/bookstore_users-api/utils/errors"
	"net/http"
	"strconv"
)

var(
	counter int
)
func GetUser(c *gin.Context){
	userId, userErr := strconv.ParseInt(c.Param("user_id"),10,64)
	if userErr != nil{
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status,err)
		return
	}

	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context){
	var user users.User
	//fmt.Println(user)
	//bytes, err := ioutil.ReadAll(c.Request.Body)
	//if err != nil{
	//	//TODO: Handle error
	//	return
	//}
	//if err := json.Unmarshal(bytes, &user); err != nil{
	//	fmt.Println(err.Error())
	//	//TODO: Handle json error
	//	return
	//}
	if err := c.ShouldBindJSON(&user); err != nil{
		fmt.Println(err.Error())
			restErr := errors.NewBadRequestError("invalid json body")
		//TODO: Handle json error
		c.JSON(http.StatusOK, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(http.StatusOK, saveErr)
		return
	}
	fmt.Println(user)
	//fmt.Println(string(bytes))
	// fmt.Println(err)
	c.JSON(http.StatusOK, result)
}

