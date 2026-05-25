package engine

type Input struct {
	Text string
}

type Output struct {
	Text       string
	TokenCount int
}

type Engine struct {
	ai AiInterface
}

func (e Engine) Process(in Input) Output {
	return e.ai.Process(in)
}

func NewEngine(ai AiInterface) *Engine {
	return &Engine{ai: ai}
}

type AiInterface interface {
	Process(in Input) Output
}

type SimpleAi struct{}

func (s SimpleAi) Process(in Input) Output {
	return Output{Text: "AI Response: " + in.Text}
}
