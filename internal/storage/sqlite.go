package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

type Todo struct {
	ID          int64
	Description string
	Completed   bool
	Tags        []string
}

func New(storagePath string) (*Storage, error) {
	const operationPath string = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operationPath, err)
	}

	statement, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS todo(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT,
		completed BOOLEAN DEFAULT FALSE
	);
	`)
	if err != nil {
		return nil, err
	}

	_, err = statement.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operationPath, err)
	}

	secondStatement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS tags(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		);

		CREATE TABLE IF NOT EXISTS todo_tags(
			todo_id INTEGER,
			tag_id INTEGER,
			PRIMARY KEY (todo_id, tag_id),
			FOREIGN KEY (todo_id) REFERENCES todo(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	_, err = secondStatement.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operationPath, err)
	}

	return &Storage{db: db}, nil
}
