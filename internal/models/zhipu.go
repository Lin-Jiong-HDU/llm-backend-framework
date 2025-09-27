package models

import (
	"github.com/Lin-Jiong-HDU/llm-backend-framework/internal/config"
	"github.com/yankeguo/zhipu"
)

// MessageModel 消息模型
type MessageModel struct {
	messages    []zhipu.ChatCompletionMessage
	token_usage int64
	request_id  string
}

// NewMessageModel 创建新的消息模型实例
func NewMessageModel(request_id string) *MessageModel {
	return &MessageModel{
		messages:    []zhipu.ChatCompletionMessage{},
		token_usage: 0,
		request_id:  request_id,
	}
}

// AddTokenUsage 增加token使用量
func (m *MessageModel) AddTokenUsage(tokens int64) {
	m.token_usage += tokens
}

func (m *MessageModel) AddMessage(msg string, role string) *MessageModel {
	m.messages = append(m.messages, zhipu.ChatCompletionMessage{
		Role:    role,
		Content: msg,
	})
	return m
}

// NewMessageModelWithInitialPrompt 创建带有初始提示的消息模型实例
func NewMessageModelWithInitialPrompt(cfg *config.Config, request_id string) *MessageModel {

	messageModel := NewMessageModel(request_id)
	messageModel.messages = append(messageModel.messages, zhipu.ChatCompletionMessage{
		Role:    "system",
		Content: cfg.Prompt.Prompt,
	})
	return messageModel
}

// GetMessages 获取消息列表
func (m *MessageModel) GetMessages() *[]zhipu.ChatCompletionMessage {
	return &m.messages
}

// GetRequestID 获取请求ID
func (m *MessageModel) GetRequestID() string {
	return m.request_id
}

// ClearMessages 清除消息列表
func (m *MessageModel) ClearMessages() {
	m.messages = []zhipu.ChatCompletionMessage{}
}

// SetMessages 设置消息列表
func (m *MessageModel) SetMessages(messages []zhipu.ChatCompletionMessage) {
	m.messages = messages
}

// GetTokenUsage 获取token使用情况
func (m *MessageModel) GetTokenUsage() int64 {
	return m.token_usage
}
