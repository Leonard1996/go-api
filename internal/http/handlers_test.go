package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pack-calculator/internal/packutil"
	"pack-calculator/internal/repo"
)

type memRepo struct {
	sizes []int
}

func (m *memRepo) ListPackSizes(ctx context.Context) ([]int, error) {
	out := append([]int(nil), m.sizes...)
	sort.Ints(out)
	return out, nil
}

func (m *memRepo) ReplacePackSizes(ctx context.Context, sizes []int) error {
	clean, err := packutil.NormalizePackSizes(sizes, repo.ErrInvalidPackSizes)
	if err != nil {
		return err
	}
	m.sizes = clean
	return nil
}

func TestHandleListPackSizes(t *testing.T) {
	repo := &memRepo{sizes: []int{250, 500, 1000}}
	mux := NewRouter(repo)

	req := httptest.NewRequest(http.MethodGet, "/v1/pack-sizes", nil)
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	var body packSizesResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &body))
	assert.Equal(t, []int{250, 500, 1000}, body.PackSizes)
}

func TestHandleReplacePackSizes(t *testing.T) {
	repo := &memRepo{sizes: []int{250}}
	mux := NewRouter(repo)

	payload := []byte(`{"pack_sizes":[500,250,1000]}`)
	req := httptest.NewRequest(http.MethodPut, "/v1/pack-sizes", bytes.NewReader(payload))
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	var body packSizesResponse
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &body))
	assert.Equal(t, []int{250, 500, 1000}, body.PackSizes)
}

func TestHandleCalculate(t *testing.T) {
	repo := &memRepo{sizes: []int{250, 500, 1000}}
	mux := NewRouter(repo)

	payload := []byte(`{"amount":501}`)
	req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewReader(payload))
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)

	var body struct {
		Packs map[int]int `json:"packs"`
	}
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &body))
	assert.Equal(t, map[int]int{250: 1, 500: 1}, body.Packs)
}

func TestHandleCalculateInvalidJSON(t *testing.T) {
	repo := &memRepo{sizes: []int{250, 500}}
	mux := NewRouter(repo)

	req := httptest.NewRequest(http.MethodPost, "/v1/calculate", bytes.NewReader([]byte("{")))
	res := httptest.NewRecorder()
	mux.ServeHTTP(res, req)

	require.Equal(t, http.StatusBadRequest, res.Code)
}
