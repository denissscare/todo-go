package savetodo

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Request struct {
	Description string   `json:"description"`
	Tags        []string `json:"tags,omitempty"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type TodoSaver interface {
	AddTodo(description string, tags []string) (int64, error)
}

func New(todoSaver TodoSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.saveTodo.New"
		_ = op

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, Response{Status: "Error", Error: "Ошибка чтения запроса"})
			return
		}
		fmt.Printf("\n\nRequest: %s", req.Tags)
	}
}
