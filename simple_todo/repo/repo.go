package repo

import (
	"context"
	"database/sql"
	"simple_todo/model"
)

type Repo struct {
	db *sql.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db}
}

func (r *Repo) GetAll(ctx context.Context) ([]model.Todo, error) {
	q := "SELECT id, task, done, created_at, updated_at FROM todos"
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []model.Todo
	for rows.Next() {
		t, err := scan(rows)
		if err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, rows.Err()
}

func (r *Repo) GetByID(ctx context.Context, id int) (model.Todo, error) {
	q := "SELECT id, task, done, created_at, updated_at FROM todos WHERE id = ?"
	return scan(r.db.QueryRowContext(ctx, q, id))
}

func (r *Repo) Add(ctx context.Context, task string) error {
	q := "INSERT INTO todos (task) VALUES ( ? )"
	_, err := r.db.ExecContext(ctx, q, task)
	return err
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	q := "DELETE FROM todos WHERE id = ?"
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *Repo) EditTask(ctx context.Context, id int, task string) error {
	q := "UPDATE todos SET task = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, q, task, id)
	return err
}

func (r *Repo) EditStatus(ctx context.Context, id int, done bool) error {
	q := "UPDATE todos SET done = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?"
	_, err := r.db.ExecContext(ctx, q, done, id)
	return err
}

func scan(s interface{ Scan(dest ...any) error }) (model.Todo, error) {
	var r model.Todo
	err := s.Scan(&r.ID, &r.Task, &r.Done, &r.CreatedAt, &r.UpdatedAt)
	return r, err
}
