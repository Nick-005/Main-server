package sqlite

import (
	"database/sql"
	"fmt"
	"server/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
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
		org TEXT NOT NULL,
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

func (s *Storage) SaveURL(employee_id int, name string, price int, org string, location string, experience string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmtVacancy, err := s.db.Prepare("INSERT INTO vacancy(employee_id,name ,price,org,location,experience) VALUES (?, ?,?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmtVacancy.Exec(employee_id, name, price, org, location, experience)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
	// return res.LastInsertId(), nil
}

/*
func (s *Storage) GetURL(ID int) (take.Response, error) {
	const op = "storage.sqlite.GetURL"
	stmtVacancy, err := s.db.Prepare("SELECT * FROM url WHERE url.id = ?")
	if err != nil {
		fmt.Printf("%s: preparing statement  %w", op, storage.ErrURLNotFound)

	}
	var result take.Response
	err = stmtVacancy.QueryRow(ID).Scan(&result.ID, &result.Alias, &result.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("%s: preparing statement  %w", op, storage.ErrURLNotFound)
			return result, fmt.Errorf("This is bad", slogf.Err(err))
		} else {
			fmt.Printf("%s: scan statement %w", op, storage.ErrURLNotFound)
			return result, fmt.Errorf("This is bad", slogf.Err(err))
		}
	}
	return result, nil
}

func (s *Storage) GetAll() ([]take.Response, error) {
	const op = "storage.sqlite.GetAll"
	_, err := s.db.Prepare("SELECT * FROM url")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		return nil, fmt.Errorf("ERROR IN CREATING REQUEST OT DB")
	}
	result := []take.Response{}
	row, err := s.db.Query("SELECT * FROM url")
	if err != nil {
		fmt.Println(err, "Error")
		return nil, nil
	}
	for row.Next() {
		r := take.Response{}
		err := row.Scan(&r.ID, &r.Alias, &r.URL)
		if err != nil {
			fmt.Println(err)
			continue
		}
		r.Status = resp.OK().Status
		result = append(result, r)
	}
	fmt.Println()
	return result, nil
}
*/
