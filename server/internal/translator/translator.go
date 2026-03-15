// Package translator provides AI-powered translation using LangChainGo and Ollama.
// It is used at runtime by the gRPC handler; it does not run during make or Docker build.
package translator

import (
	// context is used to pass cancellation and timeouts into the LLM call so we can abort if needed.
	"context"
	"errors"
	"strings"
	// llms is the LangChainGo core type for "language model"; we use it for GenerateFromSinglePrompt.
	"github.com/tmc/langchaingo/llms"
	// ollama is the LangChainGo driver for the Ollama local API; we use it to create the LLM client.
	"github.com/tmc/langchaingo/llms/ollama"
)

// Translator holds the Ollama LLM client used to generate translations.
// We keep it as a struct so we can reuse one connection for many requests instead of creating a new client per call.
type Translator struct {
	// llm is the LangChainGo LLM interface; the concrete type is *ollama.LLM from ollama.New().
	llm llms.Model
}

// New creates a Translator that talks to Ollama with the given model and base URL.
// model is the Ollama model name (e.g. "llama2"); baseURL is the Ollama API root (e.g. "http://localhost:11434").
// The caller (e.g. main or a wire function) passes these from config so translator stays independent of the config package.
func New(ctx context.Context, model, baseURL string) (*Translator, error) {
	// ollama.New builds an *ollama.LLM; we pass options via variadic arguments.
	// WithModel sets which model to use; WithServerURL sets the API base (optional; default is localhost:11434).
	llm, err := ollama.New(ollama.WithModel(model), ollama.WithServerURL(baseURL))
	if err != nil {
		// If Ollama is unreachable or options are invalid, we return the error so the server can fail fast at startup.
		return nil, err
	}
	// Some ollama.New implementations can return (nil, nil); guard so we never store a nil llm and panic in Translate.
	if llm == nil {
		return nil, errors.New("ollama client is nil")
	}
	// Wrap the concrete *ollama.LLM in our struct; we use llms.Model so we could swap to another LLM later if needed.
	return &Translator{llm: llm}, nil
}

// Translate sends the text and language pair to the LLM and returns the translated text.
// sourceLang and targetLang are short codes (e.g. "en", "es"); they are embedded in the prompt so the model knows the direction.
func (t *Translator) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	// Guard against nil llm so we return an error instead of panicking (e.g. if New stored nil).
	if t == nil || t.llm == nil {
		return "", errors.New("translator not initialized")
	}
	// Build a single prompt that instructs the model to translate; we keep it simple so the model returns only the translation.
	prompt := buildTranslatePrompt(text, sourceLang, targetLang)
	// GenerateFromSinglePrompt calls the LLM with one string; it handles encoding and the HTTP request to Ollama.
	// We pass ctx so the call can be cancelled or timed out; we use the default temperature (no extra option).
	resp, err := llms.GenerateFromSinglePrompt(ctx, t.llm, prompt)
	if err != nil {
		// Network errors, timeouts, or Ollama errors are returned to the caller (the gRPC handler) to map to gRPC status.
		return "", err
	}
	// resp is the raw string from the model; we trim whitespace so small formatting quirks don't leak into the response.
	return strings.TrimSpace(resp), nil
}

// buildTranslatePrompt returns a prompt string that tells the model to translate from sourceLang to targetLang.
// We use a separate function so the main Translate logic stays readable and we can test or change the prompt shape in one place.
func buildTranslatePrompt(text, sourceLang, targetLang string) string {
	// Simple instruction plus the text; the model is expected to reply with only the translation.
	return "Translate the following text from " + sourceLang + " to " + targetLang + ". Reply with only the translation.\n\n" + text
}
