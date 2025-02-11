package service

import (
	"context"
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/model"
	"go-my-demo/internal/repository"

	"github.com/jinzhu/copier"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, category *v1.CategoryRequest) (*model.Category, error)
	GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) ([]model.Category, error)
	UpdateCategory(ctx context.Context, id string, req *v1.CategoryRequest) error
}

func NewCategoryService(
	service *Service,
	categoryRepository repository.CategoryRepository,
) CategoryService {
	return &categoryService{
		Service:            service,
		categoryRepository: categoryRepository,
	}
}

type categoryService struct {
	*Service
	categoryRepository repository.CategoryRepository
}

func (s *categoryService) GetAllCategory(req v1.GetCategoryRequest, ctx context.Context) ([]model.Category, error) {
	res, err := s.categoryRepository.GetAllCategory(req, ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *categoryService) GetCategoryById(ctx context.Context, id string) (*model.Category, error) {
	category, err := s.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *v1.CategoryRequest) error {
	category, err := s.categoryRepository.GetCategoryById(ctx, id)
	if err != nil {
		return err
	}

	// 使用 copier.CopyWithOption 合并请求参数和查询结果
	if err := copier.CopyWithOption(category, req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return err
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.categoryRepository.UpdateCategory(ctx, category); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *categoryService) CreateCategory(ctx context.Context, req *v1.CategoryRequest) (*model.Category, error) {

	var category model.Category
	if err := copier.Copy(&category, req); err != nil {
		return nil, err
	}

	var createdCategory *model.Category
	// Transaction demo
	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		var err error
		createdCategory, err = s.categoryRepository.CreateCategory(ctx, &category)
		if err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return createdCategory, err
}
