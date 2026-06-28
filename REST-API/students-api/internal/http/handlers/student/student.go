package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/storage"
	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/types"
	"github.com/Suyash-Kamath/Golang-Coders-Gyan/REST-API/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// crud

// Create/ New is convention and this func returns handler func

// See, how we will use storage , as discussed , Student New function hai , uske andhar database use karna hai , so as a dependency receive karna padega , this is dependency injection
// we will use Storage interface , no concrete implementation , with this out system will be plug in type


func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a Student")

		var student types.Student

		// json ka decoder io.Reader naam ka interface chata hai and ye r *http.Request me interface implemented hota hai 
		err:=json.NewDecoder(r.Body).Decode(&student)
		// could have used nil but wanted to know explicitly ki kya wo wo error hai jo ham soch rahe hai , so error package hai , isme err pass karte hai , and type mention , io.EOF , means no more input is available (means body empty ,decode karne keliye kuch hai hi nahi )
		if errors.Is(err,io.EOF){
			// response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err)) // returns error , but hum handle nahi karoge 
			// Custom message 
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("Empty Body")))
			return
		}

		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}


		// Request ko validate karna hi karna hai , 0 trust policy , har ek ko test karna hai kya send kiya wo kiya jo chahiye tha woh mila ?

		// Request Validation

		// could have done manually , but not recommened , we have package , validator

		// struct tags daalo `validate:"required"` naam se 

		if err:=validator.New().Struct(student);err!=nil{
			// yahape alag type hai ValidationError , so typecast
			validateErrs:= err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validateErrs))
			return
		}

		lastId,err:=storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("User Created successfully",slog.String("UserID",fmt.Sprint(lastId)))

		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,err)
			return
		}	

		// response.WriteJson(w, http.StatusCreated,map[string]string{"success":"OK"})
		response.WriteJson(w, http.StatusCreated,map[string]int64{"id":lastId})
	}
}
