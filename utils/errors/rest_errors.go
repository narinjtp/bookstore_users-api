package errors

import "net/http"

type RestErr struct{
	Messsage string `json:"message"`
	Status int `json:"status"`
	Error string `json:"error"`
}

func NewBadRequestError(message string)*RestErr  {
	return &RestErr{
		Messsage: message,
		Status: http.StatusBadRequest  ,
		Error:    "bad_request",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Messsage: message,
		Status: http.StatusNotFound  ,
		Error:    "not_found",
	}
}

func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Messsage: message,
		Status: http.StatusInternalServerError  ,
		Error:    "internal_server_error",
	}
}