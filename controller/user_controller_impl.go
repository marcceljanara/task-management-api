package controller

import (
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
	helper.ReadFromRequestBody(request, &userCreateRequest)

	userResponse := controller.UserService.Register(request.Context(), userCreateRequest)
	webResponse := web.WebResponse{
		Code: 200,
		Status: "Berhasil mendaftar akun",
		Data: userResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

func (controller *UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userLoginRequest := web.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &userLoginRequest)

	userResponse := controller.UserService.Login(request.Context(), userLoginRequest)

	cookie := &http.Cookie{
		Name: "access_token",
		Value: userResponse,
		Path: "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure: false,
	}

	http.SetCookie(writer, cookie)

	webResponse := web.WebResponse{
		Code: 200,
		Status: "Login berhasil!",
	}

	helper.WriteToResponseBody(writer, webResponse)
}
