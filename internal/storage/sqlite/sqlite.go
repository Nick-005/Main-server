package sqlite

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"server/internal/server/handlers/auth"
	"server/internal/server/handlers/take"
	"server/internal/storage"
	"time"

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

func CreateTokenTable(storagPath string) (*Storage, error) {
	const op = "storage.sqlite.Token"
	db, err := sql.Open("sqlite3", storagPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	stmtEmp, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS token(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		active_token TEXT NOT NULL, 
		is_active INTEGER NOT NULL CHECK (is_active IN (0,1))
		);
		CREATE INDEX IF NOT EXISTS about ON token(user_id);
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

func CreateTableUser(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.User"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		phoneNumber TEXT NOT NULL UNIQUE
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.Exec()
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

type Header struct {
	Alg string `json:"alg"` // Алгоритм подписи
	Typ string `json:"typ"` // Тип токена
}

type Payload struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"` // Subject (обычно идентификатор пользователя)
	Iat int64  `json:"iat"` // Issued at - время в которое был выдан токен
	Exp int64  `json:"exp"` // Время истечения токена (в Unix timestamp)
}

func (s *Storage) CreateRefreshToken(email string) (string, error) {
	var secretKEY string = "super-nice-SECRETKEY-for-backendPART of 12341213 years from colleges"

	var header Header
	header.Alg = "HS256"
	header.Typ = "JWT"

	var payload Payload
	payload.Iss = "Nick005-aka-monkeyZV-nikita"
	payload.Sub = email
	payload.Iat = time.Now().Unix()
	payload.Exp = time.Now().Add(time.Hour * 72).Unix()

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "error", fmt.Errorf("error in converting HEADER to JSON")
	}
	headerBASE64 := base64.RawURLEncoding.Strict().EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "error", fmt.Errorf("error in converting PAYLOAD to JSON")
	}
	payloadBASE64 := base64.RawURLEncoding.Strict().EncodeToString(payloadJSON)

	// создаем подпись для JWTшки
	signaturePayAndHeader := fmt.Sprintf("%s.%s", headerBASE64, payloadBASE64)

	h := hmac.New(sha256.New, []byte(secretKEY))
	h.Write([]byte(signaturePayAndHeader))
	var signature string = base64.RawStdEncoding.EncodeToString(h.Sum(nil))

	var tokenJWT string = fmt.Sprintf("%s.%s.%s", headerBASE64, payloadBASE64, signature)

	return tokenJWT, nil
}

func (s *Storage) CreateAccessToken(email string) (string, error) {
	var secretKEY string = "ISP-7-21-borodinna"

	var header Header
	header.Alg = "HS256"
	header.Typ = "JWT"

	var payload Payload
	payload.Iss = "Nick005-aka-monkeyZV"
	payload.Sub = email
	payload.Iat = time.Now().Unix()
	payload.Exp = time.Now().Add(time.Second * 60).Unix()

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "error", fmt.Errorf("error in converting HEADER to JSON")
	}
	headerBASE64 := base64.RawURLEncoding.Strict().EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "error", fmt.Errorf("error in converting PAYLOAD to JSON")
	}
	payloadBASE64 := base64.RawURLEncoding.Strict().EncodeToString(payloadJSON)

	// создаем подпись для JWTшки
	signaturePayAndHeader := fmt.Sprintf("%s.%s", headerBASE64, payloadBASE64)

	h := hmac.New(sha256.New, []byte(secretKEY))
	h.Write([]byte(signaturePayAndHeader))
	var signature string = base64.RawStdEncoding.EncodeToString(h.Sum(nil))

	var tokenJWT string = fmt.Sprintf("%s.%s.%s", headerBASE64, payloadBASE64, signature)

	return tokenJWT, nil
}

func (s *Storage) AddUser(email string, password string, name string, phoneNumber string) error {
	const op = "storage.sqlite.Add.User"
	stmtUser, err := s.db.Prepare("INSERT INTO user(email, password, name , phoneNumber) VALUES (?,?,?,?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmtUser.Exec(email, password, name, phoneNumber)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return fmt.Errorf("%s: %w", op, storage.ErrUSERExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetLoginWithPassword(uEmail string, uPassword string) (auth.RequestAuth, error) {
	const op = "storage.sqlite.Get.VacancyByIDs"
	var result auth.RequestAuth
	stmtVacancy, err := s.db.Prepare("SELECT * FROM vacancy WHERE id = ?")
	if err != nil {
		return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrUSERNotFound)
	}
	_ = stmtVacancy
	// row, err := stmtVacancy.Query.Query("SELECT * FROM vacancy")
	err = s.db.QueryRow("SELECT email, password  FROM user WHERE email = ? and password = ?", uEmail, uPassword).Scan(&result.Email, &result.Password)
	// fmt.Println(result.Emp_ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrUSERSomething)

		} else {
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrUSERNotFound)

		}
	}

	return result, nil
}

func (s *Storage) AddVacancy(employee_id int, name string, price int, location string, experience string) (int64, error) {
	const op = "storage.sqlite.Add.Vacancy"

	stmtVacancy, err := s.db.Prepare("INSERT INTO vacancy(employee_id,name ,price,location,experience) VALUES (?,?,?,?,?)")

	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	limit := s.GetLimit(employee_id)
	if limit != 0 {
		return -1, fmt.Errorf("%s: %w", op, storage.ErrVACLimitIsOver)
	}

	sqlResult, err := stmtVacancy.Exec(employee_id, name, price, location, experience)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrVACExists)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	vac_id, err := sqlResult.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error in take index inside method")
	}
	return vac_id, nil
}

func (s *Storage) AddEmployee(limitIsOver int, nameOrganization string, phoneNumber string, email string, geography string, about string) (int64, error) {
	const op = "storage.sqlite.Add.Emp"
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
	const op = "storage.sqlite.Get.VacancyByIDs"
	var result take.ResponseVac
	_, err := s.db.Prepare("SELECT * FROM vacancy WHERE id = ?")
	if err != nil {
		return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)
	}

	// row, err := stmtVacancy.Query.Query("SELECT * FROM vacancy")
	err = s.db.QueryRow("SELECT * FROM vacancy WHERE id = ?", ID).Scan(&result.ID, &result.Emp_ID, &result.Vac_Name, &result.Price, &result.Location, &result.Experience)
	// fmt.Println(result.Emp_ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// fmt.Errorf("failed to decode request body", slogf.Err(err))
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		} else {
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		}
	}

	return result, nil
}

func (s *Storage) GetEmployee(ID int) (take.RequestEmployee, error) {
	const op = "storage.sqlite.Get.EmployeeByIDs"
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
			// fmt.Errorf("failed to decode request body", slogf.Err(err))
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		} else {
			return result, fmt.Errorf("%s: preparing statement  %w", op, storage.ErrVACNotFound)

		}
	}

	return result, nil
}

func (s *Storage) GetAllVacsForEmployee(emp_id int) ([]take.ResponseVac, error) {
	const op = "storage.sqlite.Get.AllVacancy"
	_, err := s.db.Prepare("SELECT * FROM vacancy WHERE employee_id = ?")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		return nil, fmt.Errorf("ERROR IN CREATING REQUEST OT DB")
	}
	result := []take.ResponseVac{}
	row, err := s.db.Query("SELECT * FROM vacancy WHERE employee_id = ?", emp_id)
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

func (s *Storage) GetAllVacs() ([]take.ResponseVac, error) {
	const op = "storage.sqlite.Get.AllVacancy"
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
	const op = "storage.sqlite.Get.AllEmployees"
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
