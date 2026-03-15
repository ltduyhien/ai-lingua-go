package grpc

import (
	"context"

	translationv1 "github.com/ltduyhien/ai-lingua-go/api/gen/translation/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Translator interface {
	Translate(ctx context.Context, text, sourceLang, targetLang string) (string, error)
}

type Cache interface {
	Key(text, sourceLang, targetLang string) string
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key, value string) error
}

type Server struct {
	translationv1.UnimplementedTranslationServiceServer
	translator Translator
	cache      Cache
}

func NewServer(tr Translator, c Cache) *Server {
	return &Server{translator: tr, cache: c}
}

func (s *Server) Translate(ctx context.Context, req *translationv1.TranslateRequest) (*translationv1.TranslateResponse, error) {
	if req.GetText() == "" {
		return nil, status.Error(codes.InvalidArgument, "text is required")
	}
	if req.GetSourceLang() == "" || req.GetTargetLang() == "" {
		return nil, status.Error(codes.InvalidArgument, "source_lang and target_lang are required")
	}
	if s.cache != nil {
		key := s.cache.Key(req.GetText(), req.GetSourceLang(), req.GetTargetLang())
		cached, ok, err := s.cache.Get(ctx, key)
		if err != nil {
		} else if ok {
			return &translationv1.TranslateResponse{TranslatedText: cached}, nil
		}
		translated, err := s.translator.Translate(ctx, req.GetText(), req.GetSourceLang(), req.GetTargetLang())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "translation failed: %v", err)
		}
		_ = s.cache.Set(ctx, key, translated)
		return &translationv1.TranslateResponse{TranslatedText: translated}, nil
	}
	translated, err := s.translator.Translate(ctx, req.GetText(), req.GetSourceLang(), req.GetTargetLang())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "translation failed: %v", err)
	}
	return &translationv1.TranslateResponse{TranslatedText: translated}, nil
}
