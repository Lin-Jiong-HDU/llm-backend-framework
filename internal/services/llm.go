package services

import (
	"context"

	"github.com/Lin-Jiong-HDU/llm-backend-framework/internal/config"
	"github.com/Lin-Jiong-HDU/llm-backend-framework/internal/models"
	"github.com/google/uuid"
	"github.com/yankeguo/zhipu"
)

// ZhipuClient 客户端封装
type ZhipuClient struct {
	client *zhipu.Client
	config *config.Config
}

// ZhipuService 服务封装
type ZhipuService struct {
	service    *zhipu.ChatCompletionService
	request_id string
}

// 服务管理
type ZhipuServiceManager struct {
	services map[string]*ZhipuService
}

// NewZhipuClient 创建新的Zhipu客户端实例
func NewZhipuClient(cfg *config.Config) (*ZhipuClient, error) {
	client, err := zhipu.NewClient(zhipu.WithAPIKey(cfg.LLM.APIKey))
	if err != nil {
		return nil, err
	}
	return &ZhipuClient{
		client: client,
		config: cfg,
	}, nil
}

// NewZhipuService 创建新的Zhipu服务实例
func (z *ZhipuClient) NewZhipuService(model string) *ZhipuService {
	service := z.client.ChatCompletion(z.config.LLM.Model)
	service = service.SetTemperature(z.config.LLM.Temperature)
	service = service.SetMaxTokens(z.config.LLM.MaxTokens)
	request_id := uuid.New().String()
	service = service.SetRequestID(request_id)
	//
	// service = service.SetMessages()
	//
	return &ZhipuService{
		service:    service,
		request_id: request_id,
	}
}

// GetRequestID 获取请求ID
func (z *ZhipuService) GetRequestID() string {
	return z.request_id
}

// GetService 获取服务实例
func (z *ZhipuService) GetService() *zhipu.ChatCompletionService {
	return z.service
}

// ChatCompletion 执行聊天完成请求
func (z *ZhipuService) ChatCompletion(message_model *models.MessageModel) (*zhipu.ChatCompletionResponse, error) {
	z.service.SetMessages(message_model.GetMessages())
	serviceResponse, err := z.service.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &serviceResponse, nil
}
