package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rindrics/go_todo_app/entity"
	"github.com/Rindrics/go_todo_app/store"
	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Message string   `json"message"`
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {

	w.Header().Set("contet-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error response error: %v", err)
		}
		return
	}

	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}
}

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}
	t := &entity.Task{
		Title:   b.Title,
		Status:  entity.TaskStatusTodo,
		Created: time.Now(),
	}
	id, err := at.Store.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
