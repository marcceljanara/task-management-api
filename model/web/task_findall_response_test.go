package web

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTasksResponseMarshalMatchesFindAllContract(t *testing.T) {
	response := TasksResponse{
		Data: []TaskFindAllResponse{
			{
				Id:       "task-id",
				Title:    "Memasak Mie Goreng",
				Status:   "pending",
				Priority: "high",
				DueDate:  "07-05-2026 12:30:00",
			},
		},
		Pagination: Pagination{
			Page:       1,
			Limit:      10,
			TotalRows:  1,
			TotalPages: 1,
			HasNext:    false,
			HasPrev:    false,
		},
	}

	result, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	jsonString := string(result)
	for _, disallowedField := range []string{"description", "created_at", "updated_at"} {
		if strings.Contains(jsonString, disallowedField) {
			t.Fatalf("expected response to not contain %q, got %s", disallowedField, jsonString)
		}
	}

	for _, expectedField := range []string{"id", "title", "status", "priority", "due_date", "pagination"} {
		if !strings.Contains(jsonString, expectedField) {
			t.Fatalf("expected response to contain %q, got %s", expectedField, jsonString)
		}
	}
}
