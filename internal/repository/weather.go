package repository

import (
	"go-site/internal/model"
	"go-site/pkg/log"
	"time"

	"gorm.io/gorm"
)

type WeatherRepository interface {
	CreateWeather(weather *model.Weather) error
	GetWeatherByCity(city string) (*model.Weather, error)
	UpdateWeather(weather *model.Weather) error
	CreateOrUpdateTodayWeather(weather *model.Weather) error
	GetTodayWeatherByCity(city string) (*model.Weather, error)
}

type weatherRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewWeatherRepository(db *gorm.DB, logger *log.Logger) WeatherRepository {
	return &weatherRepository{
		db:     db,
		logger: logger,
	}
}

func (r *weatherRepository) CreateWeather(weather *model.Weather) error {
	return r.db.Create(weather).Error
}

// 创建或更新当天的天气数据
func (r *weatherRepository) CreateOrUpdateTodayWeather(weather *model.Weather) error {
	// 查找今天该城市的数据
	var existing model.Weather
	err := r.db.Where("ad_code = ?",
		weather.AdCode).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 今天没有数据，创建新记录
		weather.CreatedAt = time.Now()
		weather.UpdatedAt = time.Now()
		return r.db.Create(weather).Error
	} else if err != nil {
		// 其他错误
		return err
	} else {
		// 今天已有数据，更新记录
		weather.Id = existing.Id
		weather.CreatedAt = time.Now() // 保持原创建时间
		weather.UpdatedAt = time.Now()
		return r.db.Save(weather).Error
	}
}

// 获取指定城市当天的天气数据
func (r *weatherRepository) GetTodayWeatherByCity(city string) (*model.Weather, error) {
	// 获取今天的开始和结束时间
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var weather model.Weather
	err := r.db.Where("ad_code = ? AND created_at >= ? AND created_at < ?",
		city, startOfDay, endOfDay).Order("created_at DESC").First(&weather).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 没有找到当天数据
		}
		return nil, err
	}
	return &weather, nil
}

func (r *weatherRepository) GetWeatherByCity(city string) (*model.Weather, error) {
	var weather model.Weather
	err := r.db.Where("ad_code = ?", city).Order("created_at DESC").First(&weather).Error
	if err != nil {
		return nil, err
	}
	return &weather, nil
}

func (r *weatherRepository) UpdateWeather(weather *model.Weather) error {
	return r.db.Save(weather).Error
}
