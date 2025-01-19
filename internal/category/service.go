package category

import (
	"context"
	"database/sql"
	"post-backend/internal/custom"
	"post-backend/internal/helper"
)

type CategoryService interface {
	Create(ctx context.Context, input CreateInputCategory) (Category, error)
	Get(ctx context.Context, input GetInputCategory) (Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, inputData CreateInputCategory, inputId GetInputCategory) (Category, error)
	Delete(ctx context.Context, input GetInputCategory) error
}

type CategoryServiceImpl struct {
	CategoryRepository CategoryRepository
	DB                 *sql.DB
}

func NewCategoryService(categoryRepository CategoryRepository, DB *sql.DB) CategoryService {
	return &CategoryServiceImpl{CategoryRepository: categoryRepository, DB: DB}
}

func (c *CategoryServiceImpl) Create(ctx context.Context, input CreateInputCategory) (Category, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return Category{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	category := Category{
		Name: input.Name,
	}

	category, err = c.CategoryRepository.Insert(ctx, tx, category)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (c *CategoryServiceImpl) Delete(ctx context.Context, input GetInputCategory) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	category, err := c.CategoryRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if category.Id <= 0 {
		return custom.ErrNotFound
	}

	err = c.CategoryRepository.Delete(ctx, tx, category.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *CategoryServiceImpl) Get(ctx context.Context, input GetInputCategory) (Category, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return Category{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	category, err := c.CategoryRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return Category{}, err
	}
	if category.Id <= 0 {
		return Category{}, custom.ErrNotFound
	}

	return category, nil
}

func (c *CategoryServiceImpl) GetAll(ctx context.Context) ([]Category, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return []Category{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	categories, err := c.CategoryRepository.FindAll(ctx, tx)
	if err != nil {
		return []Category{}, err
	}

	return categories, nil
}

func (c *CategoryServiceImpl) Update(ctx context.Context, inputData CreateInputCategory, inputId GetInputCategory) (Category, error) {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return Category{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	category, err := c.CategoryRepository.FindById(ctx, tx, inputId.Id)
	if err != nil {
		return Category{}, err
	}
	if category.Id <= 0 {
		return Category{}, custom.ErrNotFound
	}
	if category.Id != inputId.Id {
		return Category{}, custom.ErrNotFound
	}

	category.Name = inputData.Name
	category, err = c.CategoryRepository.Update(ctx, tx, category)
	if err != nil {
		return Category{}, err
	}

	return category, nil
}
