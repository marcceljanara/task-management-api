package web

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const TaskDateTimeLayout = "02-01-2006 15:04:05"

type taskRequestJSON struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

func (request *TaskCreateRequest) UnmarshalJSON(data []byte) error {
	var raw taskRequestJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	dueDate, err := parseTaskDueDate(raw.DueDate)
	if err != nil {
		return err
	}

	*request = TaskCreateRequest{
		Title:       raw.Title,
		Description: raw.Description,
		Priority:    raw.Priority,
		Status:      raw.Status,
		DueDate:     dueDate,
	}
	return nil
}

func (request *TaskUpdateRequest) UnmarshalJSON(data []byte) error {
	var raw taskRequestJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	dueDate, err := parseTaskDueDate(raw.DueDate)
	if err != nil {
		return err
	}

	*request = TaskUpdateRequest{
		Id:          raw.Id,
		Title:       raw.Title,
		Description: raw.Description,
		Priority:    raw.Priority,
		Status:      raw.Status,
		DueDate:     dueDate,
	}
	return nil
}

func parseTaskDueDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}

	if dueDate, err := time.ParseInLocation(TaskDateTimeLayout, value, time.Local); err == nil {
		return dueDate, nil
	}

	if dueDate, err := time.Parse(time.RFC3339, value); err == nil {
		return dueDate, nil
	}

	return time.Time{}, fmt.Errorf("due_date must use format %s", TaskDateTimeLayout)
}
