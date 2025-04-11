package savetodo

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	Description string   `json:"description" validate:"required,min=1"`
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

func New(log *slog.Logger, todoSaver TodoSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.saveTodo.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, Message{Status: "Error", Error: "Ошибка чтения запроса"})
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.String("error", err.Error()))
			render.JSON(w, r, Message{Status: "Error", Error: err.Error()})
			return
		}

		id, err := todoSaver.AddTodo(req.Description, req.Tags)
		if err != nil {
			log.Error("Error adding to the database", slog.String("error", err.Error()))
			render.JSON(w, r, Message{Status: "Error", Error: err.Error()})
			return
		}
		log.Info("Success adding to the database", slog.Int64("id", id))
		render.JSON(w, r, Response{Message: Message{Status: "OK"}, Id: id})
	}
}
