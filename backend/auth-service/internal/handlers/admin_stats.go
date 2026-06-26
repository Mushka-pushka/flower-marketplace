package handlers

import (
	"net/http"

	"github.com/Mushka-pushka/flower-marketplace/backend/auth-service/internal/service"
)

type AdminStatsHandler struct {
	statsService *service.AdminStatsService
}

func NewAdminStatsHandler(statsService *service.AdminStatsService) *AdminStatsHandler {
	return &AdminStatsHandler{statsService: statsService}
}

// GetAdminStats — получение общей статистики для администратора
func (h *AdminStatsHandler) GetAdminStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsService.GetAdminStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}