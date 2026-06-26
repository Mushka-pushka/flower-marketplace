package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Настройки сервисов
	authURL := "http://localhost:8081"
	catalogURL := "http://localhost:8082"
	orderURL := "http://localhost:8083"

	// Главный обработчик всех запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var targetURL string

		// Определяем, куда перенаправить запрос
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/v1/auth/"):
			targetURL = authURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/catalog/"):
			targetURL = catalogURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/orders"):
			targetURL = orderURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/"):
			targetURL = authURL + r.URL.Path
		default:
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		// Создаём запрос к нужному сервису
		proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Копируем заголовки
		proxyReq.Header = r.Header.Clone()
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

		// Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Копируем ответ обратно клиенту
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	log.Println("API Gateway starting on port 8080")
	log.Println("Auth Service: http://localhost:8081")
	log.Println("Catalog Service: http://localhost:8082")
	log.Println("Order Service: http://localhost:8083")
	log.Fatal(http.ListenAndServe(":8080", nil))
}