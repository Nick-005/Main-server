package auth

import (
	"log/slog"
	"net/http"
	resp "server/internal/api"
	"server/internal/lib/logger/slogf"
	"strings"

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
	AddUser(email string, password string, name string, phoneNumber string) (int, error)
	GetLoginWithPassword(uEmail string, uPassword string) (RequestAuth, error)
	CreateAccessToken(email string, uid int) (int, string, error)
	CreateRefreshToken(email string) (string, error)
}

type RequestAdd struct {
	Email       string `json:"email" `
	Password    string `json:"password"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone"`
}

type RequestAuth struct {
	Email    string `json:"email" `
	Password string `json:"password"`
}

type RequestToken struct {
	Email   string `json:"email" `
	JWToken string `json:"token"`
}

type ResponseRegistration struct {
	resp.Response
	UID     int    `json:"user_id"`
	Email   string `json:"email" `
	JWToken string `json:"token"`
}

type ResponseErr struct {
	resp.Response
	Message string `json:"error"`
}

type Response struct {
	resp.Response
}

func TakeToken(log *slog.Logger, addReq AddRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")
		render.JSON(w, r, RequestToken{
			Email:   "TakeToken@test.ru",
			JWToken: token,
		})

	}
}

func CreateOrUpdateAccessToken(log *slog.Logger, addReq AddRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		uid, answer, err := addReq.CreateAccessToken("CreateOrUpdateToken@test.ru", 0)
		if err != nil {
			render.JSON(w, r, "error")
			return
		}
		render.JSON(w, r, ResponseRegistration{
			UID:     uid,
			Email:   "nice_email@bk.ru",
			JWToken: answer,
		})
	}
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

		uid, err := addReq.AddUser(req.Email, req.Password, req.Name, req.PhoneNumber)
		if err != nil {
			log.Error("failed to add new user", slogf.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}
		uid, token, err := addReq.CreateAccessToken(req.Email, uid)
		if err != nil {
			log.Error("failed to create token for new user", slogf.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}
		render.JSON(w, r, ResponseRegistration{
			Response: resp.OK(),
			UID:      uid,
			Email:    req.Email,
			JWToken:  token,
		})

	}
}

func AuthUser(log *slog.Logger, addReq AddRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.Auth.User"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestAuth

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", slogf.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body success decoded", slog.Any("request", req))

		uData, err := addReq.GetLoginWithPassword(req.Email, req.Password)
		if err != nil {
			log.Error("failed to add new user", slogf.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		if uData.Password == req.Password && uData.Email == req.Email {
			render.JSON(w, r, Response{
				Response: resp.OK(),
			})
			return
		} else {
			render.JSON(w, r, ResponseErr{
				Response: resp.OK(),
				Message:  "the data does not match",
			})
			return

		}

	}
}
