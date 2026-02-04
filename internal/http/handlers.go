package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"pack-calculator/internal/pack"
	"pack-calculator/internal/packutil"
	"pack-calculator/internal/repo"
)

type Handler struct {
	repo repo.PackSizeRepository
}

type packSizesRequest struct {
	PackSizes []int `json:"pack_sizes"`
}

type packSizesResponse struct {
	PackSizes []int `json:"pack_sizes"`
}

type calculateRequest struct {
	Amount int `json:"amount"`
}

func NewRouter(repo repo.PackSizeRepository) *http.ServeMux {
	h := &Handler{repo: repo}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.handleHealth)
	mux.HandleFunc("GET /v1/pack-sizes", h.handleListPackSizes)
	mux.HandleFunc("PUT /v1/pack-sizes", h.handleReplacePackSizes)
	mux.HandleFunc("POST /v1/calculate", h.handleCalculate)
	return mux
}

func NewHandler(repo repo.PackSizeRepository) http.Handler {
	return NewRouter(repo)
}

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) handleListPackSizes(w http.ResponseWriter, r *http.Request) {
	sizes, err := h.repo.ListPackSizes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list pack sizes")
		return
	}

	writeJSON(w, http.StatusOK, packSizesResponse{PackSizes: sizes})
}

func (h *Handler) handleReplacePackSizes(w http.ResponseWriter, r *http.Request) {
	var req packSizesRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	clean, err := packutil.NormalizePackSizes(req.PackSizes, repo.ErrInvalidPackSizes)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid pack sizes")
		return
	}

	if err := h.repo.ReplacePackSizes(r.Context(), clean); err != nil {
		if errors.Is(err, repo.ErrInvalidPackSizes) {
			writeError(w, http.StatusBadRequest, "invalid pack sizes")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update pack sizes")
		return
	}

	writeJSON(w, http.StatusOK, packSizesResponse{PackSizes: clean})
}

func (h *Handler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	var req calculateRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Amount < 0 {
		writeError(w, http.StatusBadRequest, pack.ErrInvalidAmount.Error())
		return
	}

	sizes, err := h.repo.ListPackSizes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list pack sizes")
		return
	}

	res, err := pack.Solve(req.Amount, sizes)
	if err != nil {
		if errors.Is(err, pack.ErrInvalidAmount) || errors.Is(err, pack.ErrInvalidPackSizes) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "calculation failed")
		return
	}

	writeJSON(w, http.StatusOK, res)
}

func decodeJSON(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
