package repository

import (
	"context"
	"go-my-demo/internal/model"
)

type WebsiteRepository interface {
	GetWebsiteById(ctx context.Context, id string) (*model.Website, error)
	Create(ctx context.Context, website *model.Website) (*model.Website, error)
}

func NewWebsiteRepository(
	repository *Repository,
) WebsiteRepository {
	return &websiteRepository{
		Repository: repository,
	}
}

type websiteRepository struct {
	*Repository
}

func (r *websiteRepository) Create(ctx context.Context, website *model.Website) (*model.Website, error) {
	if err := r.DB(ctx).Create(website).Error; err != nil {
		return nil, err
	}
	return website, nil
}

func (r *websiteRepository) GetWebsiteById(ctx context.Context, id string) (*model.Website, error) {
	var website model.Website
	if err := r.DB(ctx).Where("id = ?", id).First(&website).Error; err != nil {
		return nil, err
	}
	return &website, nil
}
