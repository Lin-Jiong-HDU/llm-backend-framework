package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Config 应用程序配置结构
// 包含服务器配置、LLM配置、数据库配置等
type Config struct {
	Server   ServerConfig   `json:"server"`
	LLM      LLMConfig      `json:"llm"`
	Database DatabaseConfig `json:"database"`
	Auth     AuthConfig     `json:"auth"`
	Prompt   InitialPrompt  `json:"prompt"` // 初始系统提示列表
}

// ServerConfig HTTP服务器配置
type ServerConfig struct {
	Port         int `json:"port"`          // 服务器端口
	ReadTimeout  int `json:"read_timeout"`  // 读取超时时间(秒)
	WriteTimeout int `json:"write_timeout"` // 写入超时时间(秒)
	IdleTimeout  int `json:"idle_timeout"`  // 空闲超时时间(秒)
}

// LLMConfig LLM服务配置
type LLMConfig struct {
	Provider    string  `json:"provider"`    // LLM提供商 (openai, zhipu, etc.)
	APIKey      string  `json:"api_key"`     // API密钥
	BaseURL     string  `json:"base_url"`    // API基础URL
	Model       string  `json:"model"`       // 默认模型名称
	Temperature float64 `json:"temperature"` // 采样温度
	MaxTokens   int     `json:"max_tokens"`  // 最大生成长度
	Timeout     int     `json:"timeout"`     // 请求超时时间(秒)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `json:"driver"`   // 数据库驱动 (postgres, mysql, sqlite)
	Host     string `json:"host"`     // 数据库主机
	Port     int    `json:"port"`     // 数据库端口
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
	Database string `json:"database"` // 数据库名
	SSLMode  string `json:"ssl_mode"` // SSL模式
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret     string `json:"jwt_secret"`     // JWT密钥
	TokenExpiry   int    `json:"token_expiry"`   // Token过期时间(小时)
	RefreshExpiry int    `json:"refresh_expiry"` // 刷新Token过期时间(小时)
	APIKeyHeader  string `json:"api_key_header"` // API密钥头部名称
}

type InitialPrompt struct {
	Prompt string `json:"prompt"` // 初始系统提示内容
}

// Load 从环境变量和配置文件加载配置
// 优先级: 环境变量 > 配置文件 > 默认值
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  60,
		},
		LLM: LLMConfig{
			Provider:    "zhipu",
			BaseURL:     "https://open.bigmodel.cn/api/paas/v4",
			APIKey:      "",
			Model:       "glm-4.5",
			Timeout:     60,
			Temperature: 0.7,
			MaxTokens:   8192,
		},
		Database: DatabaseConfig{
			Driver:  "postgres",
			Host:    "localhost",
			Port:    5432,
			SSLMode: "disable",
		},
		Auth: AuthConfig{
			TokenExpiry:   24,
			RefreshExpiry: 168, // 7 days
			APIKeyHeader:  "X-API-Key",
			JWTSecret:     "your-secret-key",
		},
		Prompt: InitialPrompt{
			Prompt: "你是一个乐于助人的AI助手。你的名字叫做小智。",
		},
	}

	// 尝试从配置文件加载
	if err := loadFromFile(config); err != nil {
		// 配置文件不存在或读取失败，使用环境变量和默认值
		fmt.Printf("Config file not found or invalid, using environment variables: %v\n", err)
	}

	// 从环境变量覆盖配置
	loadFromEnv(config)

	// 验证必需的配置
	if err := validate(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// loadFromFile 从config.json文件加载配置
func loadFromFile(config *Config) error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(config)
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) {
	// 服务器配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Server.Port = p
		}
	}

	// LLM配置
	if provider := os.Getenv("LLM_PROVIDER"); provider != "" {
		config.LLM.Provider = provider
	}
	if apiKey := os.Getenv("LLM_API_KEY"); apiKey != "" {
		config.LLM.APIKey = apiKey
	}
	if baseURL := os.Getenv("LLM_BASE_URL"); baseURL != "" {
		config.LLM.BaseURL = baseURL
	}
	if model := os.Getenv("LLM_MODEL"); model != "" {
		config.LLM.Model = model
	}

	// 数据库配置
	if host := os.Getenv("DB_HOST"); host != "" {
		config.Database.Host = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Database.Port = p
		}
	}
	if username := os.Getenv("DB_USERNAME"); username != "" {
		config.Database.Username = username
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		config.Database.Password = password
	}
	if database := os.Getenv("DB_DATABASE"); database != "" {
		config.Database.Database = database
	}

	// 认证配置
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.Auth.JWTSecret = jwtSecret
	}
}

// validate 验证配置的有效性
func validate(config *Config) error {
	if config.LLM.APIKey == "" {
		return fmt.Errorf("LLM_API_KEY is required")
	}
	if config.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}
