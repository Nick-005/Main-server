package auth

import (
	"log/slog"
	"net/http"
	resp "server/internal/api"
	"server/internal/lib/logger/slogf"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

/*
Имя
Почта
Телефон
Пароль

// Инн для работадателя
*/
type AddRequest interface {
	AddUser(email string, password string, name string, phoneNumber string) error
}

type RequestAdd struct {
	Email       string `json:"email" `
	Password    string `json:"password"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone"`
}

type Response struct {
	resp.Response
}

func NewUser(log *slog.Logger, addReq AddRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.New.User"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestAdd

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", slogf.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body success decoded", slog.Any("request", req))

		err = addReq.AddUser(req.Email, req.Password, req.Name, req.PhoneNumber)
		if err != nil {
			log.Error("failed to add new user", slogf.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}
		render.JSON(w, r, Response{
			Response: resp.OK(),
		})

	}
}
