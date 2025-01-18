package setting

import (
	"context"
	"database/sql"
	"post-backend/internal/helper"
)

type SettingService interface {
	GetAll(ctx context.Context) ([]Setting, error)
}

type SettingServiceImpl struct {
	SettingRepository SettingRepository
	DB                *sql.DB
}

func NewSettingService(settingRepository SettingRepository, DB *sql.DB) SettingService {
	return &SettingServiceImpl{SettingRepository: settingRepository, DB: DB}
}

func (s *SettingServiceImpl) GetAll(ctx context.Context) ([]Setting, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return []Setting{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	settings, err := s.SettingRepository.FindAll(ctx, tx)
	if err != nil {
		return []Setting{}, err
	}

	return settings, nil
}
