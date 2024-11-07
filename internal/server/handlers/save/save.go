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
	SaveURL(employee_id int, name string, price int, org string, location string, experience string) (int64, error)
}

type Request struct {
	Emp_ID       int    `json:"emp_id"`
	Vac_Name     string `json:"vac_name"`
	Price        int    `json:"price"`
	Organization string `json:"org"`
	Location     string `json:"location"`
	Experience   string `json:"exp"`
}

type Response struct {
	resp.Response
	ID int `json:"id"`
}

// const aliasLength int = 10

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
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

		log.Info("request body decoded", slog.Any("request", req))

		// if err := validator.New().Struct(req); err != nil {
		// 	validateErr := err.(validator.ValidationErrors)

		// 	log.Error("invalid request", slogf.Err(err))

		// 	render.JSON(w, r, resp.ValidationError(validateErr))

		// 	return
		// }
		id, err := urlSaver.SaveURL(req.Emp_ID, req.Vac_Name, req.Price, req.Organization, req.Location, req.Experience)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("vac", req.Vac_Name))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
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
