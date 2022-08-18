package store

import (
	"context"
	"database/sql"

	"github.com/Rindrics/go_todo_app/entity"
)

func (r *Repository) ListTasks(
	ctx context.Context, db *sql.DB,
) (entity.Tasks, error) {
	sql := `SELECT
              id, title,
              status, created, modified
            FROM task;`
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks entity.Tasks
	for rows.Next() {
		t := &entity.Task{}
		if err := rows.Scan(
			&(t.ID), &(t.Title),
			&(t.Status), &(t.Created), &(t.Modified),
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
