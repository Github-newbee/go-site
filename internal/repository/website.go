package repository

import (
	"context"
	v1 "go-my-demo/api/v1"
	"go-my-demo/internal/model"
	"go-my-demo/pkg/db"

	"gorm.io/gorm"
)

type WebsiteRepository interface {
	GetWebsiteById(ctx context.Context, id string) (*model.Website, error)
	Create(ctx context.Context, website *model.Website) (*model.Website, error)
	Get(ctx context.Context, req v1.GetWebsiteRequest) ([]v1.WebsiteResponse, error)
	Update(ctx context.Context, website *model.Website) error
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

func (r *websiteRepository) GetWebsiteById(ctx context.Context, id string) (website *model.Website, err error) {
	if err := r.DB(ctx).Where("id = ?", id).First(&website).Error; err != nil {
		return nil, err
	}
	return website, nil
}

func (r *websiteRepository) Get(ctx context.Context, req v1.GetWebsiteRequest) (website []v1.WebsiteResponse, err error) {

	// 只返回关联表的指定内容且不嵌套对象 category_name, category_desc，结构体为v1.WebsiteResponse
	qs := r.DB(ctx).Model(&model.Website{}).Session(&gorm.Session{}).Scopes(db.FilterByQuery(req))
	if err := qs.Select("websites.*, categories.category_name AS category_name, categories.description AS category_desc").
		Joins("left join categories on categories.id = websites.category_id").
		Limit(req.Limit).Offset(req.Skip).Scan(&website).Error; err != nil {
		return nil, err
	}

	// 嵌套对象返回全部关联表内容，注意结构体嵌套的对象的json不能设置为 '-'
	// qs := r.DB(ctx).Model(&model.Website{}).Session(&gorm.Session{}).Scopes(db.FilterByQuery(req))
	// if err := qs.Select("id", "name", "url", "category_id").Preload("Category").Limit(req.Limit).Offset(req.Skip).Find(&website).Error; err != nil {
	// 	return nil, err
	// }
	return website, nil
}

func (r *websiteRepository) Update(ctx context.Context, website *model.Website) error {
	if err := r.DB(ctx).Save(website).Error; err != nil {
		return err
	}
	return nil
}
