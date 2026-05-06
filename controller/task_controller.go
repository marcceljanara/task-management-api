package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type TaskController interface {
	CreateTask(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetAllTasks(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	GetTaskById(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateTask(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	DeleteTask(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}