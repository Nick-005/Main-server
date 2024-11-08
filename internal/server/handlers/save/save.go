package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "server/internal/api"
	"server/internal/lib/logger/slogf"
	"server/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLSaver interface {
	AddVacancy(employee_id int, name string, price int, location string, experience string) (int64, error)
	AddEmployee(limitIsOver int, nameOrganization string, phoneNumber string, email string, geography string, about string) (int64, error)
}

type Request struct {
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
}
type RequestEmployee struct {
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	Geography        string `json:"geography"`
	About            string `json:"about"`
}

type Response struct {
	resp.Response
	ID int `json:"id"`
}

// const aliasLength int = 10

func NewVac(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", slogf.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body success decoded", slog.Any("request", req))

		id, err := urlSaver.AddVacancy(req.Emp_ID, req.Vac_Name, req.Price, req.Location, req.Experience)

		if errors.Is(err, storage.ErrVACExists) {
			log.Info("vacancy already exists", slog.String("vac", req.Vac_Name))
			render.JSON(w, r, resp.Error("vacancy already exists"))
			return
		}
		if err != nil {
			w.WriteHeader(452)
			log.Error("failed to add vacancy", slogf.Err(err))
			render.JSON(w, r, resp.Error("failed to add vacancy"))
			return
		}

		log.Info("Success! Vacancy has been added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       int(id),
			// Alias:    alias,
		})

	}
}

func NewEmp(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.save.NewEmp"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestEmployee

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", slogf.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body success decoded", slog.Any("request", req))

		id, err := urlSaver.AddEmployee(0, req.NameOrganization, req.PhoneNumber, req.Email, req.Geography, req.About)
		if err != nil {
			w.WriteHeader(461)
			log.Info("employee already exists", slog.String("email", req.Email))
			render.JSON(w, r, resp.Error("employee already exists"))
			return
		}
		log.Info("Success! Enployee has been added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       int(id),
		})
	}
}
