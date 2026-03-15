package translator

import (
	"context"
	"errors"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type Translator struct {
	llm llms.Model
}

func New(ctx context.Context, model, baseURL string) (*Translator, error) {
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(baseURL))
	if err != nil {
		return nil, err
	}
	if llm == nil {
		return nil, errors.New("ollama client is nil")
	}
	return &Translator{llm: llm}, nil
}

func (t *Translator) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	if t == nil || t.llm == nil {
		return "", errors.New("translator not initialized")
	}
	prompt := buildTranslatePrompt(text, sourceLang, targetLang)
	resp, err := llms.GenerateFromSinglePrompt(ctx, t.llm, prompt)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func buildTranslatePrompt(text, sourceLang, targetLang string) string {
	return "Translate the following text from " + sourceLang + " to " + targetLang + ". Reply with only the translation.\n\n" + text
}
