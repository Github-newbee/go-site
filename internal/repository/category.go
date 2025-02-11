package repository

import (
	"context"
	"errors"
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/model"
	"go-my-demo/pkg/db"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	GetCategoryById(ctx context.Context, id string) (*model.Category, error)
	GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) ([]model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
}

func NewCategoryRepository(
	repository *Repository,
) CategoryRepository {
	return &categoryRepository{
		Repository: repository,
	}
}

type categoryRepository struct {
	*Repository
}

func (r *categoryRepository) GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) (results []model.Category, err error) {
	qs := r.DB(ctx).Model(&model.Category{}).Session(&gorm.Session{}).Scopes(db.FilterByQuery(req))
	if err := qs.Limit(req.Limit).Offset(req.Skip).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *categoryRepository) GetCategoryById(ctx context.Context, id string) (*model.Category, error) {
	var category model.Category
	if err := r.DB(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrDataNotFound
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	// gorm在创建数据时会自动填充传入结构体指针中的字段
	if err := r.DB(ctx).Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *model.Category) error {
	if err := r.DB(ctx).Save(category).Error; err != nil {
		return err
	}
	return nil
}
