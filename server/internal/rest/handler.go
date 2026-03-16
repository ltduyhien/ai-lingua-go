package rest

import (
	"encoding/json"
	"net/http"

	translationv1 "github.com/ltduyhien/ai-lingua-go/api/gen/translation/v1"
	grpchandler "github.com/ltduyhien/ai-lingua-go/internal/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
	CustomPrompt string `json:"custom_prompt"`
}

type TranslateResponse struct {
	TranslatedText string `json:"translated_text"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Handler struct {
	grpc *grpchandler.Server
}

func NewHandler(grpc *grpchandler.Server) *Handler {
	return &Handler{grpc: grpc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost || r.URL.Path != "/api/translate" {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
		return
	}
	var req TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}
	protoReq := &translationv1.TranslateRequest{
		Text:         req.Text,
		SourceLang:   req.SourceLang,
		TargetLang:   req.TargetLang,
		CustomPrompt: req.CustomPrompt,
	}
	resp, err := h.grpc.Translate(r.Context(), protoReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.InvalidArgument {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: st.Message()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, TranslateResponse{TranslatedText: resp.TranslatedText})
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
