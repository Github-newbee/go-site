package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// ConfigValidator 配置验证器
type ConfigValidator struct {
	conf *viper.Viper
}

// NewConfigValidator 创建配置验证器
func NewConfigValidator(conf *viper.Viper) *ConfigValidator {
	return &ConfigValidator{conf: conf}
}

// Validate 验证配置的合理性
func (cv *ConfigValidator) Validate() error {
	if err := cv.validateApp(); err != nil {
		return fmt.Errorf("app config validation failed: %w", err)
	}

	if err := cv.validateHTTP(); err != nil {
		return fmt.Errorf("http config validation failed: %w", err)
	}

	if err := cv.validateSecurity(); err != nil {
		return fmt.Errorf("security config validation failed: %w", err)
	}

	if err := cv.validateDatabase(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	if err := cv.validateRedis(); err != nil {
		return fmt.Errorf("redis config validation failed: %w", err)
	}

	if err := cv.validateLog(); err != nil {
		return fmt.Errorf("log config validation failed: %w", err)
	}

	return nil
}

// validateApp 验证应用配置
func (cv *ConfigValidator) validateApp() error {
	if cv.conf.GetString("app.name") == "" {
		return fmt.Errorf("app.name is required")
	}

	if cv.conf.GetString("app.version") == "" {
		return fmt.Errorf("app.version is required")
	}

	env := cv.conf.GetString("env")
	if env != "local" && env != "test" && env != "production" {
		return fmt.Errorf("env must be one of: local, test, production")
	}

	return nil
}

// validateHTTP 验证HTTP配置
func (cv *ConfigValidator) validateHTTP() error {
	port := cv.conf.GetInt("http.port")
	if port <= 0 || port > 65535 {
		return fmt.Errorf("http.port must be between 1 and 65535, got: %d", port)
	}

	host := cv.conf.GetString("http.host")
	if host == "" {
		return fmt.Errorf("http.host is required")
	}

	// 验证超时配置
	readTimeout := cv.conf.GetDuration("http.read_timeout")
	if readTimeout < time.Second {
		return fmt.Errorf("http.read_timeout should be at least 1s, got: %v", readTimeout)
	}

	writeTimeout := cv.conf.GetDuration("http.write_timeout")
	if writeTimeout < time.Second {
		return fmt.Errorf("http.write_timeout should be at least 1s, got: %v", writeTimeout)
	}

	return nil
}

// validateSecurity 验证安全配置
func (cv *ConfigValidator) validateSecurity() error {
	// JWT 配置验证
	jwtKey := cv.conf.GetString("security.jwt.key")
	if len(jwtKey) < 32 {
		return fmt.Errorf("security.jwt.key should be at least 32 characters long")
	}

	// API 签名配置验证
	apiKey := cv.conf.GetString("security.api_sign.app_key")
	if len(apiKey) < 6 {
		return fmt.Errorf("security.api_sign.app_key should be at least 6 characters long")
	}

	apiSecret := cv.conf.GetString("security.api_sign.app_security")
	if len(apiSecret) < 6 {
		return fmt.Errorf("security.api_sign.app_security should be at least 6 characters long")
	}

	// 生产环境安全检查
	if cv.conf.GetString("env") == "production" {
		if apiKey == "123456" || apiSecret == "123456" {
			return fmt.Errorf("production environment should not use default API keys")
		}
		if jwtKey == "QQYnRFerJTSEcrfB89fw8prOaObmrch8" {
			return fmt.Errorf("production environment should not use default JWT key")
		}
	}

	return nil
}

// validateDatabase 验证数据库配置
func (cv *ConfigValidator) validateDatabase() error {
	driver := cv.conf.GetString("data.db.user.driver")
	if driver == "" {
		return fmt.Errorf("data.db.user.driver is required")
	}

	dsn := cv.conf.GetString("data.db.user.dsn")
	if dsn == "" {
		return fmt.Errorf("data.db.user.dsn is required")
	}

	// 连接池配置验证
	maxIdleConns := cv.conf.GetInt("data.db.user.max_idle_conns")
	maxOpenConns := cv.conf.GetInt("data.db.user.max_open_conns")

	if maxIdleConns < 0 {
		return fmt.Errorf("data.db.user.max_idle_conns should be >= 0")
	}

	if maxOpenConns <= 0 {
		return fmt.Errorf("data.db.user.max_open_conns should be > 0")
	}

	if maxIdleConns > maxOpenConns {
		return fmt.Errorf("data.db.user.max_idle_conns should be <= max_open_conns")
	}

	return nil
}

// validateRedis 验证Redis配置
func (cv *ConfigValidator) validateRedis() error {
	addr := cv.conf.GetString("data.redis.addr")
	if addr == "" {
		return fmt.Errorf("data.redis.addr is required")
	}

	db := cv.conf.GetInt("data.redis.db")
	if db < 0 || db > 15 {
		return fmt.Errorf("data.redis.db must be between 0 and 15")
	}

	return nil
}

// validateLog 验证日志配置
func (cv *ConfigValidator) validateLog() error {
	logLevel := cv.conf.GetString("log.log_level")
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}

	if !validLevels[logLevel] {
		return fmt.Errorf("log.log_level must be one of: debug, info, warn, error, fatal")
	}

	encoding := cv.conf.GetString("log.encoding")
	if encoding != "json" && encoding != "console" {
		return fmt.Errorf("log.encoding must be either 'json' or 'console'")
	}

	maxSize := cv.conf.GetInt("log.max_size")
	if maxSize <= 0 {
		return fmt.Errorf("log.max_size should be > 0")
	}

	return nil
}

// PrintConfigSummary 打印配置摘要
func (cv *ConfigValidator) PrintConfigSummary() {
	fmt.Printf("=== Configuration Summary ===\n")
	fmt.Printf("Environment: %s\n", cv.conf.GetString("env"))
	fmt.Printf("App: %s v%s\n", cv.conf.GetString("app.name"), cv.conf.GetString("app.version"))
	fmt.Printf("HTTP: %s:%d\n", cv.conf.GetString("http.host"), cv.conf.GetInt("http.port"))
	fmt.Printf("Database: %s\n", cv.conf.GetString("data.db.user.driver"))
	fmt.Printf("Redis: %s (DB:%d)\n", cv.conf.GetString("data.redis.addr"), cv.conf.GetInt("data.redis.db"))
	fmt.Printf("Log Level: %s\n", cv.conf.GetString("log.log_level"))
	fmt.Printf("============================\n")
}
