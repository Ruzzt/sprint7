package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	city := "moscow"
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList[city])},
	}
	for _, v := range requests {
		url := "/cafe?city=" + city + "&count=" + strconv.Itoa(v.count)
		res := httptest.NewRecorder() //обращаюсь к серверу и в ответ от серва получаю ответ
		req := httptest.NewRequest("GET", url, nil)

		handler.ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code)

		body := res.Body.String()
		if body == "" {
			assert.Equal(t, 0, v.want, "ожидается пустой ответ")
			continue
		}

		cafes := strings.Split(body, ",")
		assert.Equal(t, v.want, len(cafes), "неверное количество кафе в ответе")

	}

}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		url := "/cafe?city=moscow&search=" + v.search
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, url, nil)

		handler.ServeHTTP(res, req)
		require.Equal(t, http.StatusOK, res.Code)

		body := res.Body.String()

		if body == "" {
			assert.Equal(t, 0, v.wantCount, "ожидается пустой ответ")
			continue
		}

		cafes := strings.Split(body, ",")

		assert.Equal(t, v.wantCount, len(cafes), "неверное количество кафе")

		searchLower := strings.ToLower(v.search)
		for _, name := range cafes {
			nameLower := strings.ToLower(name)
			assert.Contains(t, nameLower, searchLower, "название должно содержать подстроку поиска")
		}
	}
}
