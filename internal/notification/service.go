package notification

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
	"post-backend/internal/helper"
)

type NotificationService interface {
	GetAll(ctx context.Context) ([]Notification, error)
	GetById(ctx context.Context, input GetByIdNotificationInput) (Notification, error)
	UpdateStatus(ctx context.Context, inputId GetByIdNotificationInput, inputData UpdateNotificationInput) error
}

type NotificationServiceImpl struct {
	NotificationRepository NotificationRepository
	DB                     *sql.DB
}

func NewNotificationService(notificationRepository NotificationRepository, DB *sql.DB) NotificationService {
	return &NotificationServiceImpl{
		NotificationRepository: notificationRepository,
		DB:                     DB,
	}
}

func (n *NotificationServiceImpl) GetAll(ctx context.Context) ([]Notification, error) {
	tx, err := n.DB.BeginTx(ctx, nil)
	if err != nil {
		return []Notification{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	notifications, err := n.NotificationRepository.FindAll(ctx, tx)
	if err != nil {
		return []Notification{}, err
	}

	return notifications, nil
}

func (n *NotificationServiceImpl) GetById(ctx context.Context, input GetByIdNotificationInput) (Notification, error) {
	tx, err := n.DB.BeginTx(ctx, nil)
	if err != nil {
		return Notification{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	notification, err := n.NotificationRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return Notification{}, err
	}
	if notification.Id < 0 {
		return Notification{}, custom.ErrNotFound
	}

	return notification, nil
}

func (n *NotificationServiceImpl) UpdateStatus(ctx context.Context, inputId GetByIdNotificationInput, inputData UpdateNotificationInput) error {
	tx, err := n.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	notification, err := n.NotificationRepository.FindById(ctx, tx, inputId.Id)
	if err != nil {
		return err
	}
	if notification.Id < 0 {
		return custom.ErrNotFound
	}

	err = n.NotificationRepository.UpdateStatus(ctx, tx, inputId.Id, inputData.Status)
	if err != nil {
		return err
	}

	return nil
}
