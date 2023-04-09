package test

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"testing"
)

func Test_ChatGPT(t *testing.T) {

	config := openai.DefaultConfig("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")
	//proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	//transport := &http.Transport{
	//	Proxy: http.ProxyURL(proxyUrl),
	//}
	//config.HTTPClient = &http.Client{
	//	Transport: transport,
	//}

	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Please translate this txt to japanese : '{单身的要讲究刷存在感的方法。恋爱中的不宜太恋爱脑，需给彼此适度的空间一定耐心。'",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		t.Error(err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
	t.Log(resp.Choices[0].Message.Content)
}
