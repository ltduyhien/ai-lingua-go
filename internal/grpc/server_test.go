package grpc

import (
	"context"
	"errors"
	"testing"

	translationv1 "github.com/ltduyhien/ai-lingua-go/api/gen/translation/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeTranslator struct {
	result string
	err    error
	called bool
}

func (f *fakeTranslator) Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error) {
	f.called = true
	if f.err != nil {
		return "", f.err
	}
	return f.result, nil
}

type fakeCache struct {
	getResult string
	getOk     bool
	getErr    error
	setErr    error
	keyUsed   string
	setValue  string
}

func (f *fakeCache) Key(text, sourceLang, targetLang string) string {
	return "trans:" + sourceLang + ":" + targetLang + ":" + text
}

func (f *fakeCache) Get(ctx context.Context, key string) (string, bool, error) {
	f.keyUsed = key
	if f.getErr != nil {
		return "", false, f.getErr
	}
	return f.getResult, f.getOk, nil
}

func (f *fakeCache) Set(ctx context.Context, key, value string) error {
	f.keyUsed = key
	f.setValue = value
	return f.setErr
}

func TestTranslate_NoCache_Success(t *testing.T) {
	tr := &fakeTranslator{result: "translated"}
	srv := NewServer(tr, nil)
	ctx := context.Background()
	req := &translationv1.TranslateRequest{Text: "hello", SourceLang: "en", TargetLang: "es"}
	resp, err := srv.Translate(ctx, req)
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if resp.TranslatedText != "translated" {
		t.Errorf("got %q, want translated", resp.TranslatedText)
	}
	if !tr.called {
		t.Error("translator was not called")
	}
}

func TestTranslate_CacheHit_ReturnsCached(t *testing.T) {
	tr := &fakeTranslator{result: "from-llm"}
	c := &fakeCache{getResult: "cached", getOk: true}
	srv := NewServer(tr, c)
	ctx := context.Background()
	req := &translationv1.TranslateRequest{Text: "hi", SourceLang: "en", TargetLang: "fr"}
	resp, err := srv.Translate(ctx, req)
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if resp.TranslatedText != "cached" {
		t.Errorf("got %q, want cached", resp.TranslatedText)
	}
	if tr.called {
		t.Error("translator should not be called on cache hit")
	}
}

func TestTranslate_CacheMiss_CallsTranslatorAndSetsCache(t *testing.T) {
	tr := &fakeTranslator{result: "from-llm"}
	c := &fakeCache{getOk: false}
	srv := NewServer(tr, c)
	ctx := context.Background()
	req := &translationv1.TranslateRequest{Text: "hi", SourceLang: "en", TargetLang: "fr"}
	resp, err := srv.Translate(ctx, req)
	if err != nil {
		t.Fatalf("Translate: %v", err)
	}
	if resp.TranslatedText != "from-llm" {
		t.Errorf("got %q, want from-llm", resp.TranslatedText)
	}
	if !tr.called {
		t.Error("translator should be called on cache miss")
	}
	if c.setValue != "from-llm" {
		t.Errorf("cache Set value: got %q, want from-llm", c.setValue)
	}
}

func TestTranslate_EmptyText_InvalidArgument(t *testing.T) {
	tr := &fakeTranslator{result: "x"}
	srv := NewServer(tr, nil)
	ctx := context.Background()
	req := &translationv1.TranslateRequest{Text: "", SourceLang: "en", TargetLang: "es"}
	_, err := srv.Translate(ctx, req)
	if err == nil {
		t.Fatal("expected error")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		t.Errorf("got %v, want InvalidArgument", err)
	}
	if tr.called {
		t.Error("translator should not be called when text is empty")
	}
}

func TestTranslate_EmptySourceLang_InvalidArgument(t *testing.T) {
	tr := &fakeTranslator{result: "x"}
	srv := NewServer(tr, nil)
	ctx := context.Background()
	req := &translationv1.TranslateRequest{Text: "hi", SourceLang: "", TargetLang: "es"}
	_, err := srv.Translate(ctx, req)
	if err == nil {
		t.Fatal("expected error")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		t.Errorf("got %v, want InvalidArgument", err)
	}
}

func TestTranslate_TranslatorError_Internal(t *testing.T) {
	tr := &fakeTranslator{err: errors.New("ollama down")}
	srv := NewServer(tr, nil)
	ctx := context.Background()
	req := &translationv1.TranslateRequest{Text: "hi", SourceLang: "en", TargetLang: "es"}
	_, err := srv.Translate(ctx, req)
	if err == nil {
		t.Fatal("expected error")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Internal {
		t.Errorf("got %v, want Internal", err)
	}
}
