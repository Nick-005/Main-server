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
	GetVacancy(ID int) (ResponseVac, error)
	GetAllVacs() ([]ResponseVac, error)

	GetAllEmps() ([]RequestEmployee, error)
	GetEmployee(ID int) (RequestEmployee, error)
}

type Request struct {
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
}
type RequestEmployee struct {
	ID               int    `json:"ID"`
	Limit            int    `json:"limit"`
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	Geography        string `json:"geography"`
	About            string `json:"about"`
}

type ResponseVac struct {
	// resp.Response
	ID         int    `json:"ID"`
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
}

type Response struct {
	resp.Response
	ID int    `json:"id"`
	Ko string `json:"text"`
}
type ResponseError struct {
	resp.Response
	Info string `json:"info"`
}

func GetEmployeeByID(log *slog.Logger, getReq GetRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(ID)
		if err != nil {
			fmt.Println(err)
		}

		const op = "storage.sqlite.GetEmployeeByID"
		res, err := getReq.GetEmployee(id)
		log.Info(op)
		if err != nil {
			w.WriteHeader(452)
			render.JSON(w, r, ResponseError{
				Response: resp.Error("Employee doesn't exist!"),
				Info:     "Такого работадателя не существует! Перепроверьте на наличие ошибок в запросе!",
			})
			return
		}
		// fmt.Println(res, err, op)
		render.JSON(w, r, res)
	}
}

func GetVacancyByID(log *slog.Logger, getReq GetRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ID := chi.URLParam(r, "id")
		id, err := strconv.Atoi(ID)
		if err != nil {
			fmt.Println(err)
		}

		const op = "storage.sqlite.GetVacancyByID"
		res, err := getReq.GetVacancy(id)
		log.Info(op)
		if err != nil {
			w.WriteHeader(452)
			render.JSON(w, r, ResponseError{
				Response: resp.Error("Vacancy doesn't exist!"),
				Info:     "Такой вакансии не существует! Перепроверьте на наличие ошибок запрос!",
			})
			return
		}

		// fmt.Println(res, err, op)
		w.WriteHeader(200)
		render.JSON(w, r, ResponseVac{
			ID:         res.ID,
			Emp_ID:     res.Emp_ID,
			Vac_Name:   res.Vac_Name,
			Price:      res.Price,
			Location:   res.Location,
			Experience: res.Experience,
		})
	}
}

func GetAllVacancy(log *slog.Logger, getReq GetRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "storage.sqlite.GetAllVacancy"
		log.Info(op)
		res, err := getReq.GetAllVacs()
		if err != nil {
			w.WriteHeader(452)
			render.JSON(w, r, ResponseError{
				Response: resp.Error("Vacancy doesn't exist!"),
				Info:     "Вакансий не существует! Перепроверьте на наличие ошибок запрос!",
			})
			return
		}
		// fmt.Println(res, err, op)
		render.JSON(w, r, res)
	}
}
func GetAllEmployees(log *slog.Logger, getReq GetRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "storage.sqlite.GetAllEmployees"
		log.Info(op)
		res, err := getReq.GetAllEmps()
		if err != nil {
			w.WriteHeader(452)
			render.JSON(w, r, ResponseError{
				Response: resp.Error("Vacancy doesn't exist!"),
				Info:     "Вакансий не существует! Перепроверьте на наличие ошибок запрос!",
			})
			return
		}
		// fmt.Println(res, err, op)
		render.JSON(w, r, res)
	}
}
