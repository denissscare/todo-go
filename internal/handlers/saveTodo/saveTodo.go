package savetodo

import "net/http"

type Request struct {
	Description string   `json:"description"`
	Tags        []string `json:"alias,omitempty"`
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
	}
}
