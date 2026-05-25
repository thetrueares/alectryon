package vendor

import (
	"context"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
	"go.iain.rocks/alectryon/api/engine"
)

func NewOpenAI(apiKey string) *OpenAI {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAI{client: client}
}

type OpenAI struct {
	client openai.Client
}

func (oa OpenAI) Process(input engine.Input) engine.Output {
	resp, err := oa.client.Responses.New(context.TODO(), responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(input.Text)},
		Model: openai.ChatModelGPT5_4Mini,
	})

	if err != nil {
		return engine.Output{Text: err.Error()}
	}

	return engine.Output{Text: resp.OutputText()}
}
