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

type Message struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Response struct {
	Message
	Id int64 `json:"id,omitempty"`
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
			render.JSON(w, r, Message{Status: "Error", Error: "Ошибка чтения запроса"})
			return
		}

		id, err := todoSaver.AddTodo(req.Description, req.Tags)
		if err != nil {
			er := fmt.Sprintf("%v", err)
			render.JSON(w, r, Message{Status: "Error", Error: er})
			return
		}
		render.JSON(w, r, Response{Message: Message{Status: "OK"}, Id: id})
	}
}
