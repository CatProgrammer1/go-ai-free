package chat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Models struct {
	Gpt4o     string
	GeminiPro string
	
	//ClaudeSonnet-3.5
	ClaudeSonnet string
	Basic        any
}

type Modes struct {
	ImageGeneration AgentMode
}

const (
	api = "https://www.blackbox.ai/api/chat"
)

var (
	TO_CLEAN = []string{
		"Generated by BLACKBOX.AI, try unlimited chat https://www.blackbox.ai",
	}

	CLIENT = http.Client{}

	MODELS = Models{
		Gpt4o:        "gpt-4o",
		GeminiPro:    "gemini-pro",
		ClaudeSonnet: "claude-sonnet-3.5",
		Basic:        nil,
	}

	MODES = Modes{
		ImageGeneration: AgentMode{
			Mode: true,
			ID:   "ImageGenerationLV45LJp",
			Name: "Image Generation",
		},
	}
)

type AgentMode struct {
	ID   string `json:"id"`
	Mode bool   `json:"mode"`
	Name string `json:"name"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Chat struct {
	Messages []Message `json:"messages"`

	ClickedAnswer2 bool `json:"clickedAnswer2"`
	ClickedAnswer3 bool `json:"clickedAnswer3"`

	CodeModelMode         bool `json:"codeModelMode"`
	ClickedForceWebSearch bool `json:"clickedForceWebSearch"`

	AgentMode         AgentMode `json:"agentMode"`
	TrendingAgentMode AgentMode `json:"trendingAgentMode"`

	MaxTokens uint `json:"maxTokens"`

	UserSystemPrompt string `json:"userSystemPrompt"`

	PlaygroundTemperature float32 `json:"playgroundTemperature"`
	UserSelectedModel     *string `json:"userSelectedModel"`
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func PostAPI(data any) string {
	jsonEncoded, err := json.Marshal(data)
	handle(err)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonEncoded))
	handle(err)

	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")

	resp, err := CLIENT.Do(req)
	handle(err)

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	handle(err)

	response := string(respData)

	for _, cleanFrom := range TO_CLEAN {
		response = strings.ReplaceAll(response, cleanFrom, "")
	}

	return response
}

func (chat *Chat) SendMessage(m Message) string {
	chat.Messages = append(chat.Messages, m)

	resp := PostAPI(chat)

	chat.Messages = append(chat.Messages, Message{
		Role:    "assistant",
		Content: resp,
	})

	return resp
}

// Change AgentModes manualy
func New(model *string, codeModelMode bool, temperature float32) *Chat {
	return &Chat{
		Messages:              []Message{},
		CodeModelMode:         codeModelMode,
		PlaygroundTemperature: temperature,
		MaxTokens:             1024,
		UserSelectedModel:     model,
	}
}
