package handlers

import (
	"net/http"
	"strconv"

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

// GetDailyStats — получение дневной статистики для графиков
func (h *AdminStatsHandler) GetDailyStats(w http.ResponseWriter, r *http.Request) {
	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val > 0 {
			days = val
		}
	}
	
	stats, err := h.statsService.GetDailyStats(r.Context(), days)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, stats)
}