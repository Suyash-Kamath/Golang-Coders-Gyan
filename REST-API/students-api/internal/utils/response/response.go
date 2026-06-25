package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// `json:status` is struct tags , it means jab json ke andhar convert hoga na tab aisa dikhna chahiye
// jabhi sturct serialize hoga json ke andhar toh Status ke bajaaye status aayega
type Response struct{
	Status string `json:"status"`
	Error string  `json:"error"`
}

const (
	StatusOk="OK"
	StatusError="Error"
)

func WriteJson(w http.ResponseWriter,status int , data interface{}) error{

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)


	return json.NewEncoder(w).Encode(data) // error is returned agar error aaya hai toh

}


func GeneralError(err error) Response{
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}

}

// package sends the error list , and we receieve it here
func ValidationError(errs validator.ValidationErrors )Response{

	var errMsgs []string

	for _,err:=range errs{
		switch err.ActualTag(){
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field",err.Field()))
	    default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid",err.Field()))
		}
	}

	return Response{
		Status:StatusError,
		Error:strings.Join(errMsgs,", "),
	}
 
}
