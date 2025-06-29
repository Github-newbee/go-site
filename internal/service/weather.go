package service

import (
	"encoding/json"
	"fmt"
	v1 "go-site/api/v1"
	"go-site/internal/model"
	"go-site/internal/repository"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type WeatherService interface {
	GetWeather(weatherReq v1.WeatherRequest) (*model.Weather, error)
}

func NewWeatherService(
	service *Service,
	weatherRepo repository.WeatherRepository,
) WeatherService {
	return &weatherService{
		weatherRepo: weatherRepo,
		Service:     service,
	}
}

type weatherService struct {
	weatherRepo repository.WeatherRepository
	*Service
}

func (s *weatherService) GetWeather(weatherReq v1.WeatherRequest) (*model.Weather, error) {
	city := weatherReq.City
	cacheKey := fmt.Sprintf("weather:%s", city)

	// 1. 首先尝试从Redis缓存获取
	cachedData, err := s.rdb.Get(s.ctx, cacheKey).Result()
	if err == nil {
		var weather model.Weather
		if err := json.Unmarshal([]byte(cachedData), &weather); err == nil {
			// 检查缓存数据是否为当天
			if s.isToday(weather.CreatedAt) {
				s.logger.Info("Weather data retrieved from cache",
					zap.String("city", city),
					zap.Time("cache_time", weather.CreatedAt))
				return &weather, nil
			} else {
				s.logger.Info("Cached weather data is outdated, will fetch new data",
					zap.String("city", city),
					zap.Time("cache_time", weather.CreatedAt))
			}
		}
	}

	// 2. 缓存未命中，尝试从数据库获取
	dbWeather, err := s.weatherRepo.GetTodayWeatherByCity(city)
	if err == nil && dbWeather != nil {
		// 检查数据库数据是否为当天
		if s.isToday(dbWeather.CreatedAt) {
			// 将当天数据重新缓存到Redis（缓存30分钟）
			weatherJSON, _ := json.Marshal(dbWeather)
			if err := s.rdb.Set(s.ctx, cacheKey, weatherJSON, 30*time.Minute).Err(); err != nil {
				s.logger.Error("Failed to cache database data to Redis",
					zap.Error(err))
			}

			return dbWeather, nil
		} else {
			s.logger.Info("Database weather data is outdated, will fetch new data",
				zap.String("city", city),
				zap.Time("db_time", dbWeather.CreatedAt))
		}
	}

	// 3. 数据库也没有数据，从第三方API获取
	weatherData, err := s.fetchWeatherFromAPI(city)
	if err != nil {
		s.logger.Error("Failed to fetch weather from API",
			zap.String("city", city),
			zap.Error(err))

		// API 调用失败，返回数据库中的最新数据（即使不是当天的）
		if dbWeather != nil {
			s.logger.Info("API failed, returning latest database data as fallback",
				zap.String("city", city),
				zap.Time("fallback_time", dbWeather.CreatedAt))
			return dbWeather, nil
		}
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}

	// 4. 保存到数据库
	if err := s.weatherRepo.CreateOrUpdateTodayWeather(weatherData); err != nil {
		s.logger.Error("Failed to save weather to database",
			zap.Error(err))
		// 不返回错误，因为我们已经有了天气数据
	} else {
		s.logger.Info("Weather data saved to database",
			zap.Reflect("weatherData", weatherData))
	}

	// 5. 保存到Redis缓存（缓存30分钟）
	weatherJSON, _ := json.Marshal(weatherData)
	if err := s.rdb.Set(s.ctx, cacheKey, weatherJSON, 30*time.Minute).Err(); err != nil {
		s.logger.Error("Failed to cache weather data",
			zap.String("city", city),
			zap.Error(err))
	} else {
		s.logger.Info("Weather data cached to Redis",
			zap.String("city", city))
	}

	return weatherData, nil
}

// 检查给定时间是否为今天
func (s *weatherService) isToday(t time.Time) bool {
	now := time.Now()
	year1, month1, day1 := t.Date()
	year2, month2, day2 := now.Date()
	return year1 == year2 && month1 == month2 && day1 == day2
}

// 从高德地图API获取天气数据
func (s *weatherService) fetchWeatherFromAPI(city string) (*model.Weather, error) {
	// 这里需要配置高德地图的API Key
	apiKey := "8609964f96739dbe24d66e1ab8abff11" // 需要在配置中添加

	url := fmt.Sprintf("https://restapi.amap.com/v3/weather/weatherInfo?key=%s&city=%s&extensions=base", apiKey, city)

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	fmt.Println("获取高德接口----------")
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var amapResp v1.AmapWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&amapResp); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	// 检查API响应状态
	if amapResp.Status != "1" {
		return nil, fmt.Errorf("API returned error: %s (code: %s)", amapResp.Info, amapResp.InfoCode)
	}

	if len(amapResp.Lives) == 0 {
		return nil, fmt.Errorf("no weather data found for city: %s", city)
	}

	// 转换为内部数据结构
	live := amapResp.Lives[0]
	weather := &model.Weather{
		AdCode:        live.AdCode,
		City:          live.City,
		Province:      live.Province,
		Weather:       live.Weather,
		Temperature:   live.Temperature,
		WindDirection: live.WindDirection,
		WindPower:     live.WindPower,
		Humidity:      live.Humidity,
		ReportTime:    live.ReportTime,
	}

	return weather, nil
}
