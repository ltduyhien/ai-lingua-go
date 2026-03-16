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

func (t *Translator) Translate(ctx context.Context, text, sourceLang, targetLang, customPrompt string) (string, error) {
	if t == nil || t.llm == nil {
		return "", errors.New("translator not initialized")
	}
	prompt := buildTranslatePrompt(text, sourceLang, targetLang, customPrompt)
	resp, err := llms.GenerateFromSinglePrompt(ctx, t.llm, prompt,
		llms.WithMaxTokens(2048),
		llms.WithTemperature(0.1),
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func buildTranslatePrompt(text, sourceLang, targetLang, customPrompt string) string {
	base := "Translate the following text from " + sourceLang + " to " + targetLang + ".\n" +
		"Rules:\n" +
		"- Output ONLY the translation in the target language, nothing else.\n" +
		"- Translate every word, phrase, and sentence; do not leave ANY character in the source language script (for example, no Chinese characters may remain in the output).\n" +
		"- Do not add explanations, summaries, or duplicate rephrasings; translate the text exactly once.\n" +
		"- Preserve the same paragraph and line-break structure as the source.\n\n" +
		"- If the translated paragraph structure is broken, or is not properly structured, fix it with a proper restructuring of the paragraphs."

	if customPrompt != "" {
		base += "\n\nAdditional instructions:\n" + customPrompt
	}

	return base + "\n\nText to translate:\n" + text
}
