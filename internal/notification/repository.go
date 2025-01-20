package notification

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
)

type NotificationRepository interface {
	Insert(ctx context.Context, tx *sql.Tx, notification Notification) (Notification, error)
	FindAll(ctx context.Context, tx *sql.Tx) ([]Notification, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) (Notification, error)
	UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error
}

type NotificationRepositoryImpl struct {
}

func NewNotificationRepository() NotificationRepository {
	return &NotificationRepositoryImpl{}
}

func (s *NotificationRepositoryImpl) UpdateStatus(ctx context.Context, tx *sql.Tx, id int, status string) error {
	sql := "update notifications set status = ? where id = ?"
	_, err := tx.ExecContext(ctx, sql, status, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *NotificationRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]Notification, error) {
	sql := "SELECT id, title, type, message, status, created_at, updated_at FROM notifications ORDER BY created_at DESC"

	rows, err := tx.QueryContext(ctx, sql)
	if err != nil {
		return []Notification{}, err
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.Id, &notification.Title, &notification.Type, &notification.Message, &notification.Status, &notification.CreatedAt, &notification.UpdatedAt); err != nil {
			return []Notification{}, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (s *NotificationRepositoryImpl) Insert(ctx context.Context, tx *sql.Tx, notification Notification) (Notification, error) {
	sql := "INSERT INTO notifications(title, type, message, status) VALUES (?,?,?,?)"
	result, err := tx.ExecContext(ctx, sql, notification.Title, notification.Type, notification.Message, notification.Status)
	if err != nil {
		return Notification{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Notification{}, err
	}

	notification.Id = int(id)
	return notification, nil
}

func (s *NotificationRepositoryImpl) FindById(ctx context.Context, tx *sql.Tx, id int) (Notification, error) {
	sql := "SELECT id, title, type, message, status, created_at, updated_at FROM notifications WHERE id = ?"

	row, err := tx.QueryContext(ctx, sql, id)
	if err != nil {
		return Notification{}, err
	}
	defer row.Close()

	var notification Notification
	if row.Next() {
		if err := row.Scan(&notification.Id, &notification.Title, &notification.Type, &notification.Message, &notification.Status, &notification.CreatedAt, &notification.UpdatedAt); err != nil {
			return Notification{}, err
		}

		return notification, nil
	}

	return notification, custom.ErrNotFound
}
