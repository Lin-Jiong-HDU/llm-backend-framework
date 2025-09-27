package main

import (
	"fmt"

	"github.com/Lin-Jiong-HDU/llm-backend-framework/internal/config"
	"github.com/Lin-Jiong-HDU/llm-backend-framework/internal/services"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// 初始化LLM服务
	client, err := services.NewZhipuClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("Failed to create Zhipu client: %v", err))
	}
	service := client.NewZhipuService()
	fmt.Println("Zhipu Service Request ID:", service.GetRequestID())

	response, err := service.ChatCompletion("你好，你叫什么名字?")
	if err != nil {
		panic(fmt.Sprintf("ChatCompletion error: %v", err))
	}
	fmt.Println("Response:", response.Choices[0].Message.Content)

}
