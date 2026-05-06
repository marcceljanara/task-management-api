package web

type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalRows  int  `json:"total_rows"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type TasksResponse struct {
	Data       []TaskResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
