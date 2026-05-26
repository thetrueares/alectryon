package engine

import "go.iain.rocks/alectryon/api/entities"

type Input struct {
	Text    string
	History []entities.HistoryEntity
}

type Output struct {
	Text       string
	TokenCount int
}

type EngineInterface interface {
	Process(in Input) Output
}

type Engine struct {
	ai                AiInterface
	historyRepository entities.HistoryRepository
}

func (e Engine) Process(in Input) Output {
	return e.ai.Process(in)
}

func NewEngine(ai AiInterface, historyRepository entities.HistoryRepository) EngineInterface {
	return &Engine{ai: ai, historyRepository: historyRepository}
}

type AiInterface interface {
	Process(in Input) Output
}

type SimpleAi struct{}

func (s SimpleAi) Process(in Input) Output {
	return Output{Text: "AI Response: " + in.Text}
}
