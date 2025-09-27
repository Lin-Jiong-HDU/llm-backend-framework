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
	service       *zhipu.ChatCompletionService
	message_model *models.MessageModel
	request_id    string
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
func (z *ZhipuClient) NewZhipuService() *ZhipuService {

	service := z.client.ChatCompletion(z.config.LLM.Model)
	service = service.SetTemperature(z.config.LLM.Temperature)
	service = service.SetMaxTokens(z.config.LLM.MaxTokens)
	request_id := uuid.New().String()
	service = service.SetRequestID(request_id)
	messages := models.NewMessageModelWithInitialPrompt(z.config, request_id)

	return &ZhipuService{
		service:       service,
		message_model: messages,
		request_id:    request_id,
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
func (z *ZhipuService) ChatCompletion(content string) (*zhipu.ChatCompletionResponse, error) {
	z.service.SetMessages(*z.message_model.GetMessages())
	z.service.AddMessage(zhipu.ChatCompletionMessage{
		Role:    "user",
		Content: content,
	})
	serviceResponse, err := z.service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	// Convert the generic []any slice returned by GetMessages into []zhipu.ChatCompletionMessage
	if msgHistory := z.service.GetMessages(); msgHistory != nil {
		raw := *msgHistory // raw is []any
		converted := make([]zhipu.ChatCompletionMessage, 0, len(raw))
		for _, v := range raw {
			if m, ok := v.(zhipu.ChatCompletionMessage); ok {
				converted = append(converted, m)
			}
		}
		z.message_model.SetMessages(converted)
		z.message_model.AddTokenUsage(serviceResponse.Usage.TotalTokens)
	}

	return &serviceResponse, nil
}

// GetTokenUsage 获取token使用量
func (z *ZhipuService) GetTokenUsage() int64 {
	return z.message_model.GetTokenUsage()
}

// NewZhipuServiceManager 创建新的Zhipu服务管理实例
func NewZhipuServiceManager() *ZhipuServiceManager {
	return &ZhipuServiceManager{
		services: make(map[string]*ZhipuService),
	}
}

// CreateService 创建并存储新的Zhipu服务实例
func (m *ZhipuServiceManager) CreateService(client *ZhipuClient) *ZhipuService {
	service := client.NewZhipuService()
	m.services[service.GetRequestID()] = service
	return service
}

// GetServiceByID 根据请求ID获取Zhipu服务实例
func (m *ZhipuServiceManager) GetServiceByID(request_id string) (*ZhipuService, bool) {
	service, exists := m.services[request_id]
	return service, exists
}

// GetMessageModelByID 根据请求ID获取消息模型实例
func (m *ZhipuServiceManager) GetMessageModelByID(request_id string) (*models.MessageModel, bool) {
	service, exists := m.services[request_id]
	serviceMessageModel := service.message_model
	return serviceMessageModel, exists
}

// DeleteServiceByID 删除Zhipu服务实例和消息模型实例
func (m *ZhipuServiceManager) DeleteServiceByID(request_id string) {
	delete(m.services, request_id)
}
