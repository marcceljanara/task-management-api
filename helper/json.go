package helper

import (
	"encoding/json"
	"errors"
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/model/web"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) error {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		return err
	}
	return nil
}

func WriteErrorResponse(writer http.ResponseWriter, err error) error {
	statusCode := http.StatusInternalServerError
	message := "internal server error"

	var appErr *exception.AppError
	if errors.As(err, &appErr) && appErr.Message != "" && !errors.Is(err, exception.ErrInternal) {
		message = appErr.Message
	}

	switch {
	case errors.Is(err, exception.ErrValidation), errors.Is(err, exception.ErrBadRequest):
		statusCode = http.StatusBadRequest
		if message == "internal server error" {
			message = "bad request"
		}
	case errors.Is(err, exception.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		if message == "internal server error" {
			message = "unauthorized"
		}
	case errors.Is(err, exception.ErrConflict):
		statusCode = http.StatusConflict
		if message == "internal server error" {
			message = "conflict"
		}
	case errors.Is(err, exception.ErrNotFound):
		statusCode = http.StatusNotFound
		if message == "internal server error" {
			message = "not found"
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)

	return json.NewEncoder(writer).Encode(web.WebResponse{
		Code:   statusCode,
		Status: message,
		Data:   nil,
	})
}
