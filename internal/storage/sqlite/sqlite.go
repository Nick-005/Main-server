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

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL UNIQUE);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
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
	stmt, err := s.db.Prepare("SELECT * FROM url WHERE url.id = ?")
	if err != nil {
		fmt.Printf("%s: preparing statement  %w", op, storage.ErrURLNotFound)

	}
	var result take.Response
	err = stmt.QueryRow(ID).Scan(&result.ID, &result.Alias, &result.URL)
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
