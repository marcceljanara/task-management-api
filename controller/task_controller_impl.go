package controller

import (
	"marcceljanara/task-management-api/exception"
	"marcceljanara/task-management-api/helper"
	"marcceljanara/task-management-api/middleware"
	"marcceljanara/task-management-api/model/web"
	"marcceljanara/task-management-api/service"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type TaskControllerImpl struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) TaskController {
	return &TaskControllerImpl{
		TaskService: taskService,
	}
}

func (controller *TaskControllerImpl) CreateTask(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		_ = helper.WriteErrorResponse(writer, exception.New(exception.ErrUnauthorized, "invalid user session"))
		return
	}

	taskCreateRequest := web.TaskCreateRequest{}
	err := helper.ReadFromRequestBody(request, &taskCreateRequest)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, exception.Wrap(exception.ErrBadRequest, "invalid JSON request body", err))
		return
	}

	taskResponse, err := controller.TaskService.Save(request.Context(), taskCreateRequest, userID)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Berhasil membuat task",
		Data:   taskResponse,
	}

	_ = helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaskControllerImpl) GetAllTasks(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		_ = helper.WriteErrorResponse(writer, exception.New(exception.ErrUnauthorized, "invalid user session"))
		return
	}

	query := request.URL.Query()
	page, err := parseOptionalIntQuery(query.Get("page"), "page")
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}
	limit, err := parseOptionalIntQuery(query.Get("limit"), "limit")
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	taskQueryRequest := web.TaskQueryFindAllRequest{
		Title:  query.Get("title"),
		Status: query.Get("status"),
		Page:   page,
		Limit:  limit,
	}

	tasksResponse, err := controller.TaskService.FindAll(request.Context(), taskQueryRequest, userID)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	response := struct {
		Code       int                `json:"code"`
		Status     string             `json:"status"`
		Data       []web.TaskResponse `json:"data"`
		Pagination web.Pagination     `json:"pagination"`
	}{
		Code:       http.StatusOK,
		Status:     "Berhasil mengambil daftar task",
		Data:       tasksResponse.Data,
		Pagination: tasksResponse.Pagination,
	}

	_ = helper.WriteToResponseBody(writer, response)
}

func (controller *TaskControllerImpl) GetTaskById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		_ = helper.WriteErrorResponse(writer, exception.New(exception.ErrUnauthorized, "invalid user session"))
		return
	}

	taskID := params.ByName("taskId")
	taskResponse, err := controller.TaskService.FindById(request.Context(), taskID, userID)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Berhasil mengambil detail task",
		Data:   taskResponse,
	}

	_ = helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaskControllerImpl) UpdateTask(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		_ = helper.WriteErrorResponse(writer, exception.New(exception.ErrUnauthorized, "invalid user session"))
		return
	}

	taskUpdateRequest := web.TaskUpdateRequest{}
	err := helper.ReadFromRequestBody(request, &taskUpdateRequest)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, exception.Wrap(exception.ErrBadRequest, "invalid JSON request body", err))
		return
	}
	taskUpdateRequest.Id = params.ByName("taskId")

	err = controller.TaskService.Update(request.Context(), taskUpdateRequest, userID)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	taskResponse, err := controller.TaskService.FindById(request.Context(), taskUpdateRequest.Id, userID)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Berhasil mengubah task",
		Data:   taskResponse,
	}

	_ = helper.WriteToResponseBody(writer, webResponse)
}

func (controller *TaskControllerImpl) DeleteTask(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userID, ok := middleware.UserIDFromContext(request.Context())
	if !ok {
		_ = helper.WriteErrorResponse(writer, exception.New(exception.ErrUnauthorized, "invalid user session"))
		return
	}

	err := controller.TaskService.Delete(request.Context(), params.ByName("taskId"), userID)
	if err != nil {
		_ = helper.WriteErrorResponse(writer, err)
		return
	}

	webResponse := web.WebResponse{
		Code:   http.StatusOK,
		Status: "Berhasil menghapus task",
		Data:   nil,
	}

	_ = helper.WriteToResponseBody(writer, webResponse)
}

func parseOptionalIntQuery(value string, fieldName string) (int, error) {
	if value == "" {
		return 0, nil
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, exception.Wrap(exception.ErrBadRequest, fieldName+" must be a valid integer", err)
	}

	return result, nil
}
