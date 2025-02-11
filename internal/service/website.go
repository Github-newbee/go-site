package service

import (
	"context"
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/model"
	"go-my-demo/internal/repository"

	"github.com/jinzhu/copier"
)

type WebsiteService interface {
	GetWebsite(ctx context.Context, id string) (*model.Website, error)
	CreateWebsite(ctx context.Context, website *v1.WebsiteRequest) (*model.Website, error)
}

func NewWebsiteService(
	service *Service,
	websiteRepository repository.WebsiteRepository,
	categoryRepository repository.CategoryRepository,
) WebsiteService {
	return &websiteService{
		Service:            service,
		websiteRepository:  websiteRepository,
		categoryRepository: categoryRepository,
	}
}

type websiteService struct {
	*Service
	websiteRepository  repository.WebsiteRepository
	categoryRepository repository.CategoryRepository
}

func (s *websiteService) CreateWebsite(ctx context.Context, req *v1.WebsiteRequest) (*model.Website, error) {
	var website model.Website
	if err := copier.Copy(&website, req); err != nil {
		return nil, err
	}
	if _, err := s.categoryRepository.GetCategoryById(ctx, req.CategoryID); err != nil {
		return nil, err
	}

	var created *model.Website
	// Transaction demo
	err := s.tm.Transaction(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.websiteRepository.Create(ctx, &website)
		if err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return created, err
}

func (s *websiteService) GetWebsite(ctx context.Context, id string) (*model.Website, error) {
	website, err := s.websiteRepository.GetWebsiteById(ctx, id)
	return website, err
}
