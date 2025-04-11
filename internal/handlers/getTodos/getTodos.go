package gettodos

import (
	"fmt"
	"net/http"

	sqlite "github.com/denissscare/todo-go/internal/storage"
	"github.com/go-chi/render"
)

type Message struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type Response struct {
	Message
	Todos []sqlite.Todo `json:"todos,omitempty"`
}

type TodosGetter interface {
	GetAllTodo() ([]sqlite.Todo, error)
}

func listRespons(todos []sqlite.Todo) []render.Renderer {
	list := []render.Renderer{}
	for _, todo := range todos {
		list = append(list, todo)
	}
	return list
}

func New(todosGetter TodosGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getTodos.New"
		_ = op

		todos, err := todosGetter.GetAllTodo()
		fmt.Println(todos)
		if err != nil {
			er := fmt.Sprintf("%v", err)
			render.JSON(w, r, Message{Status: "Error", Error: er})
			return
		}

		if err := render.RenderList(w, r, listRespons(todos)); err != nil {
			return
		}
	}
}
