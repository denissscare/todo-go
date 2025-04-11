package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

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
	const operationPath string = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operationPath, err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todo(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				description TEXT,
				completed BOOLEAN DEFAULT FALSE
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create todo table: %w", operationPath, err)
	}

	_, err = db.Exec(`
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
		return nil, fmt.Errorf("%s: failed to create tags tables: %w", operationPath, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) AddTodo(description string, tags []string) (int64, error) {
	const op = "storage.AddTodo"

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	res, err := tx.Exec(
		"INSERT INTO todo (description, completed) VALUES (?, ?)",
		description,
		false,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to insert todo: %w", op, err)
	}

	todoID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	for _, tagName := range tags {
		_, err := tx.Exec(
			"INSERT OR IGNORE INTO tags (name) VALUES (?)",
			tagName,
		)
		if err != nil {
			return 0, fmt.Errorf("%s: failed to insert tag: %w", op, err)
		}

		_, err = tx.Exec(
			`INSERT INTO todo_tags (todo_id, tag_id)
             VALUES (?, (SELECT id FROM tags WHERE name = ?))`,
			todoID,
			tagName,
		)
		if err != nil {
			return 0, fmt.Errorf("%s: failed to link tag to todo: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return int64(todoID), nil
}

func (s *Storage) GetAllTodo() ([]Todo, error) {
	const op = "storage.sqlite.GetAllTodo"

	rows, err := s.db.Query(`
        SELECT todo.id, todo.description, todo.completed, group_concat(tags.name, ', ') 
        FROM todo
        LEFT JOIN todo_tags ON todo_tags.id = todo_tags.todo_id
        LEFT JOIN tags ON todo_tags.tag_id = tags.id
        GROUP BY todo.id
        ORDER BY todo.id
    `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var tags sql.NullString

		err := rows.Scan(
			&todo.ID,
			&todo.Description,
			&todo.Completed,
			&tags,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if tags.Valid {
			todo.Tags = strings.Split(tags.String, ", ")
		} else {
			todo.Tags = []string{}
		}

		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return todos, nil
}
