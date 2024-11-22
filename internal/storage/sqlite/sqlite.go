package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"server/internal/lib/logger/slogf"
	"server/internal/server/handlers/take"
	"server/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func CreateEmployeeTable(storagPath string) (*Storage, error) {
	const op = "storage.sqlite.Emp"
	db, err := sql.Open("sqlite3", storagPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	stmtEmp, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS employee(
		id INTEGER PRIMARY KEY,
		limitVac INTEGER,
		nameOrganization TEXT NOT NULL UNIQUE,
		phoneNumber TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE ,
		geography TEXT NOT NULL,
		about TEXT);
		CREATE INDEX IF NOT EXISTS about ON employee(about);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtEmp.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func CreateVacancyTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	stmtVacancy, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS vacancy(
		id INTEGER PRIMARY KEY,
		employee_id INTEGER,
		name TEXT NOT NULL,
		price INTEGER,
		location TEXT NOT NULL,
		experience TEXT);
		CREATE INDEX IF NOT EXISTS price ON vacancy(price);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtVacancy.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddVacancy(employee_id int, name string, price int, location string, experience string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmtVacancy, err := s.db.Prepare("INSERT INTO vacancy(employee_id,name ,price,location,experience) VALUES (?,?,?,?,?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	limit := s.GetLimit(employee_id)
	if limit != 0 {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrVACLimitIsOver)
	}

	_, err = stmtVacancy.Exec(employee_id, name, price, location, experience)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrVACExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return 1, nil
}

func (s *Storage) AddEmployee(limitIsOver int, nameOrganization string, phoneNumber string, email string, geography string, about string) (int64, error) {
	const op = "storage.sqlite.AddEmp"
	stmt, err := s.db.Prepare("INSERT INTO employee(limitVac ,nameOrganization,phoneNumber,email,geography,about) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(limitIsOver, nameOrganization, phoneNumber, email, geography, about)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrVACSomething)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// fmt.Println(id)
	return id, nil
}

func (s *Storage) GetLimit(ID int) int {

	stmtCount, err := s.db.Prepare("SELECT limitVac FROM employee WHERE id = ?")
	if err != nil {
		return -1
	}
	var count int
	err = stmtCount.QueryRow(ID).Scan(&count)
	if err != nil {
		return -1
	}
	if count >= 10 {
		return -1
	}
	update := count + 1
	stmtUpdate, err := s.db.Prepare("UPDATE employee SET limitVac = ? WHERE id = ?")
	if err != nil {
		return -1
	}
	_, err = stmtUpdate.Exec(update, ID)
	if err != nil {
		return -1
	}
	return 0
}

func (s *Storage) GetVacancy(ID int) (take.ResponseVac, error) {
	const op = "storage.sqlite.GetVacancyByIDs"
	var result take.ResponseVac
	stmtVacancy, err := s.db.Prepare("SELECT * FROM vacancy WHERE id = ?")
	if err != nil {
		return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)
	}
	_ = stmtVacancy
	// row, err := stmtVacancy.Query.Query("SELECT * FROM vacancy")
	err = s.db.QueryRow("SELECT * FROM vacancy WHERE id = ?", ID).Scan(&result.ID, &result.Emp_ID, &result.Vac_Name, &result.Price, &result.Location, &result.Experience)
	// fmt.Println(result.Emp_ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Errorf("failed to decode request body", slogf.Err(err))
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		} else {
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		}
	}

	return result, nil
}

func (s *Storage) GetEmployee(ID int) (take.RequestEmployee, error) {
	const op = "storage.sqlite.GetEmployeeByIDs"
	var result take.RequestEmployee
	stmtVacancy, err := s.db.Prepare("SELECT * FROM employee WHERE id = ?")
	if err != nil {
		return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)
	}
	_ = stmtVacancy
	// row, err := stmtVacancy.Query.Query("SELECT * FROM vacancy")
	err = s.db.QueryRow("SELECT * FROM employee WHERE id = ?", ID).Scan(&result.ID, &result.Limit, &result.NameOrganization, &result.PhoneNumber, &result.Email, &result.Geography, &result.About)
	// fmt.Println(result.Emp_ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Errorf("failed to decode request body", slogf.Err(err))
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		} else {
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		}
	}

	return result, nil
}

func (s *Storage) GetAllVacs() ([]take.ResponseVac, error) {
	const op = "storage.sqlite.GetAllVacancy"
	_, err := s.db.Prepare("SELECT * FROM vacancy")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		return nil, fmt.Errorf("ERROR IN CREATING REQUEST OT DB")
	}
	result := []take.ResponseVac{}
	row, err := s.db.Query("SELECT * FROM vacancy")
	if err != nil {
		fmt.Println(err, "Error")
		return nil, nil
	}
	for row.Next() {
		r := take.ResponseVac{}
		err := row.Scan(&r.ID, &r.Emp_ID, &r.Vac_Name, &r.Price, &r.Location, &r.Experience)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// r.Status = resp.OK().Status
		result = append(result, r)
	}
	fmt.Println()
	return result, nil
}

func (s *Storage) GetAllEmps() ([]take.RequestEmployee, error) {
	const op = "storage.sqlite.GetAllEmployees"
	_, err := s.db.Prepare("SELECT * FROM employee")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		return nil, fmt.Errorf("ERROR IN CREATING REQUEST OT DB")
	}
	result := []take.RequestEmployee{}
	row, err := s.db.Query("SELECT * FROM employee")
	if err != nil {
		fmt.Println(err, "Error")
		return nil, nil
	}
	for row.Next() {
		r := take.RequestEmployee{}
		err := row.Scan(&r.ID, &r.Limit, &r.NameOrganization, &r.PhoneNumber, &r.Email, &r.Geography, &r.About)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, r)
	}
	fmt.Println()
	return result, nil
}
