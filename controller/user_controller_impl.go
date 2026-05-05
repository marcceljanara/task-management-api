package controller

import (
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/helper"
	"marcceljanara/task-management-api/model/web"
	"marcceljanara/task-management-api/service"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userCreateRequest := web.UserCreateRequest{}
	err := helper.ReadFromRequestBody(request, &userCreateRequest)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, exception.Wrap(exception.ErrBadRequest, "invalid JSON request body", err))
		return
	}

	userResponse, err := controller.UserService.Register(request.Context(), userCreateRequest)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Berhasil mendaftar akun",
		Data:   userResponse,
	}

	_ = helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userLoginRequest := web.UserLoginRequest{}
	err := helper.ReadFromRequestBody(request, &userLoginRequest)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, exception.Wrap(exception.ErrBadRequest, "invalid JSON request body", err))
		return
	}

	tokenString, err := controller.UserService.Login(request.Context(), userLoginRequest)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}

	http.SetCookie(writer, cookie)

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Login berhasil!",
		Data:   nil,
	}

	_ = helper.WriteToResponseBody(writer, webResponse)
}
