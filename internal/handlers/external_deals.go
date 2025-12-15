package handlers

import (
	"encoding/json"
	"net/http"

	"serverio_darbas/internal/services"
)

type ExternalHandler struct {
	deals *services.DealsService
}

func NewExternalHandler(deals *services.DealsService) *ExternalHandler {
	return &ExternalHandler{deals: deals}
}

func (h *ExternalHandler) Deals(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "missing query param: title", http.StatusBadRequest)
		return
	}

	out, err := h.deals.GetDeals(r.Context(), title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
