package scene

import (
	"net/http"

	"game-server/internal/httputil"
	"game-server/internal/model"
	"game-server/internal/store"

	"go.uber.org/zap"
)

// Handler scene 相关 HTTP 接口
type Handler struct {
	pg    *store.PostgresStore
	redis *store.RedisStore
}

func NewHandler(pg *store.PostgresStore, redis *store.RedisStore) *Handler {
	return &Handler{pg: pg, redis: redis}
}

// HandleDetail GET/POST /scene/detail
func (h *Handler) HandleDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sceneName, cred, ok := h.parseCommon(w, r)
	if !ok {
		return
	}

	openID, ok := h.resolveOpenID(w, cred)
	if !ok {
		return
	}

	profile, err := h.pg.GetOrCreateGameProfile(openID, sceneName)
	if err != nil {
		zap.L().Error("查询 game_profile 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	settings, err := model.ParseSceneSettings(profile.SceneSettings)
	if err != nil {
		zap.L().Error("解析 scene_settings 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"scene":  sceneName,
		"layout": settings.LayoutResponse(),
	})
}

// HandleReveal GET/POST /scene/slot/reveal
func (h *Handler) HandleReveal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		httputil.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	sceneName, cred, ok := h.parseCommon(w, r)
	if !ok {
		return
	}

	slotName := r.URL.Query().Get("slot")
	if !httputil.ValidName(slotName) {
		httputil.WriteError(w, http.StatusBadRequest, "invalid slot")
		return
	}

	openID, ok := h.resolveOpenID(w, cred)
	if !ok {
		return
	}

	profile, err := h.pg.GetOrCreateGameProfile(openID, sceneName)
	if err != nil {
		zap.L().Error("查询 game_profile 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	settings, err := model.ParseSceneSettings(profile.SceneSettings)
	if err != nil {
		zap.L().Error("解析 scene_settings 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if settings.Credit <= 0 {
		httputil.WriteError(w, http.StatusPaymentRequired, "insufficient credit")
		return
	}

	slot, exists := settings.Slots[slotName]
	if !exists {
		httputil.WriteError(w, http.StatusBadRequest, "unknown slot")
		return
	}
	if slot.Visible == "Y" {
		httputil.WriteError(w, http.StatusConflict, "slot already revealed")
		return
	}

	slot.Visible = "Y"
	settings.Slots[slotName] = slot
	settings.Credit--

	raw, err := settings.JSON()
	if err != nil {
		zap.L().Error("序列化 scene_settings 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	if err := h.pg.SaveSceneSettings(profile.ID, raw); err != nil {
		zap.L().Error("保存 scene_settings 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"scene":  sceneName,
		"layout": settings.LayoutResponse(),
	})
}

func (h *Handler) parseCommon(w http.ResponseWriter, r *http.Request) (sceneName, cred string, ok bool) {
	sceneName = r.URL.Query().Get("scene")
	if !httputil.ValidName(sceneName) {
		httputil.WriteError(w, http.StatusBadRequest, "invalid scene")
		return "", "", false
	}
	cred = r.URL.Query().Get("auth_credential")
	if cred == "" {
		httputil.WriteError(w, http.StatusUnauthorized, "auth_credential is required")
		return "", "", false
	}
	return sceneName, cred, true
}

func (h *Handler) resolveOpenID(w http.ResponseWriter, cred string) (string, bool) {
	openID, err := h.redis.GetAuthOpenID(cred)
	if err != nil {
		zap.L().Error("校验 auth_credential 失败", zap.Error(err))
		httputil.WriteError(w, http.StatusInternalServerError, "internal error")
		return "", false
	}
	if openID == "" {
		httputil.WriteError(w, http.StatusUnauthorized, "invalid auth_credential")
		return "", false
	}
	return openID, true
}
