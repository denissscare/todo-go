package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
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

type Tag struct {
	ID   int64
	Name string
}

func New(storagePath string) (*Storage, error) {
	const op string = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
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
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddTodo(description string, tags []string) (int64, error) {
	const op = "storage.sqlite.AddTodo"

	statement, err := s.db.Prepare("INSERT INTO todo (description) VALUES (?)")
	if err != nil {
		return 0, fmt.Errorf("%s: Не удалось подготовить выражение: %w", op, err)
	}

	result, err := statement.Exec(description)
	if err != nil {
		return 0, fmt.Errorf("%s: Не удалось добавить todo %w", op, err)
	}

	lastTodoID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: Не удалось получить TODO ID: %w", op, err)
	}

	for _, tag := range tags {
		secondStatement, err := s.db.Prepare("INSERT OR IGNORE INTO tags (name) VALUES (?)")
		if err != nil {
			return 0, fmt.Errorf("%s: Не подготовить выражение: %w", op, err)
		}

		_, err = secondStatement.Exec(tag)
		if err != nil {
			return 0, fmt.Errorf("%s: Не удалось добавить тэг: %w", op, err)
		}

		_, err = s.db.Exec(`
		INSERT INTO todo_tags (todo_id, tag_id) 
		VALUES (?, (SELECT id FROM tags WHERE name = ?))`,
			lastTodoID, tag)

		if err != nil {
			return 0, fmt.Errorf("%s: Не удалось связать todo с тэгом:%w", op, err)
		}
	}
	return lastTodoID, nil
}
