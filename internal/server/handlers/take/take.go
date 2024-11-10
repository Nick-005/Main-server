package take

import (
	"fmt"
	"log/slog"
	"net/http"
	resp "server/internal/api"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type SaveRequest interface {
	AddVacancy(employee_id int, name string, price int, location string, experience string) (int64, error)
	AddEmployee(limitIsOver int, nameOrganization string, phoneNumber string, email string, geography string, about string) (int64, error)
}

type GetRequest interface {
	GetURL(ID int) (Response, error)
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
	ID int    `json:"id"`
	Ko string `json:"text"`
}

func GetAll(log *slog.Logger, getReq GetRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(ID)
		if err != nil {
			fmt.Println(err)
		}

		const op = "storage.sqlite.GetURL"
		res, err := getReq.GetURL(id)
		log.Info("sdf")
		fmt.Println(res, err, op)
		render.JSON(w, r, Response{
			Response: resp.OK(),
			ID:       res.ID,
			Ko:       res.Ko,
		})
	}
}
