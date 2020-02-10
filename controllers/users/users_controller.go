package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/narinjtp/bookstore_users-api/domain/users"
	"github.com/narinjtp/bookstore_users-api/logger"
	"github.com/narinjtp/bookstore_users-api/services"
	"github.com/narinjtp/bookstore_users-api/utils/errors"
	"net/http"
	"strconv"
)

var(
	counter int
)
func getUserId(userIdParam string)(int64, *errors.RestErr){
	userId, userErr := strconv.ParseInt(userIdParam,10,64)
	if userErr != nil{
		return 0, errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}
func Get(c *gin.Context){
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil{
		c.JSON(idErr.Status,idErr)
		return
	}

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}

func Create(c *gin.Context){
	//buf := make([]byte, 1024)
	//num, _ := c.Request.Body.Read(buf)
	//reqBody := string(buf[0:num])
	//logger.Info("request :" + reqBody)

	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil{
		logger.Error("invalid json body", err)
			restErr := errors.NewBadRequestError("invalid json body")
		//TODO: Handle json error
		c.JSON(http.StatusOK, restErr)
		return
	}

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		c.JSON(http.StatusOK, saveErr)
		return
	}
	fmt.Println(user)
	//fmt.Println(string(bytes))
	// fmt.Println(err)
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Update(c *gin.Context){
	userId, userErr := strconv.ParseInt(c.Param("user_id"),10,64)
	if userErr != nil{
		err := errors.NewBadRequestError("user id should be a number")
		c.JSON(err.Status,err)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil{
		fmt.Println(err.Error())
		restErr := errors.NewBadRequestError("invalid json body")
		//TODO: Handle json error
		c.JSON(restErr.Status, restErr)
		return
	}
	user.Id = userId
	isPartial := c.Request.Method == http.MethodPatch
	result, err := services.UsersService.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(c.GetHeader("X-Public") == "true"))
}

func Delete(c *gin.Context){
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil{
		c.JSON(idErr.Status,idErr)
		return
	}
	if err := services.UsersService.DeleteUser(userId); err != nil{
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status":"deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.SearchUser(status)
	if err != nil {
		c.JSON(err.Status,err)
		return
	}
	c.JSON(http.StatusOK,users.Marshall(c.GetHeader("X-Public") == "true"))
}

func Login(c *gin.Context){
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status,restErr)
		return
	}
	user, err := services.UsersService.LoginUser(request)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(c.GetHeader("X-Public") == "true"))
}