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
	resp, err := llms.GenerateFromSinglePrompt(ctx, t.llm, prompt,
		llms.WithMaxTokens(1024),
		llms.WithTemperature(0.2),
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

func buildTranslatePrompt(text, sourceLang, targetLang string) string {
	return "Translate the following text from " + sourceLang + " to " + targetLang + ".\n" +
		"Rules: Output ONLY the translation in the target language. Translate every word—do not leave any source-language words untranslated (use the target language equivalent even for technical or rare terms). Do not include the original text, any part of it, or any explanation. Do not mix source and target language in the output. Preserve the same paragraph structure and line breaks as the source: use newlines in the same places so that lists, verses, or multi-paragraph text keep their layout.\n\n" +
		"Text to translate:\n" + text
}
