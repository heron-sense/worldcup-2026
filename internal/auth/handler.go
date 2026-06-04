package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"game-server/internal/httputil"
	"game-server/internal/store"

	"go.uber.org/zap"
)

// Handler auth 接口：根据 open_id 颁发 auth_credential
type Handler struct {
	redis *store.RedisStore
}

func NewHandler(redis *store.RedisStore) *Handler {
	return &Handler{redis: redis}
}

// ServeHTTP 处理 GET/POST /auth?open_id=...
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	openID := r.URL.Query().Get("open_id")
	if openID == "" {
		httputil.WriteError(w, http.StatusBadRequest, "open_id is required")
		return
	}

	credential, err := newCredential()
	if err != nil {
		zap.L().Error("生成 auth_credential 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.redis.SetAuthCredential(credential, openID); err != nil {
		zap.L().Error("保存 auth_credential 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{
		"auth_credential": credential,
	})
}

func newCredential() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
