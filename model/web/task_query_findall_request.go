package web

type TaskQueryFindAllRequest struct {
	Title  string `json:"title"`
	Status string `json:"status"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}
