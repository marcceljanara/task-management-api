package web

type TaskFindAllResponse struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
	DueDate  string `json:"due_date"`
}
