package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	authURL := "http://localhost:8081"
	catalogURL := "http://localhost:8082"
	orderURL := "http://localhost:8083"
	paymentURL := "http://localhost:8084"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// ---- CORS ----
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// ---- ПОЛУЧАЕМ USER_ID ЧЕРЕЗ AUTH SERVICE ----
		var userID string
		tokenStr := r.Header.Get("Authorization")
		if tokenStr != "" && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			userID = getUserIDFromAuthService(tokenStr)
		}

		// ---- МАРШРУТИЗАЦИЯ ----
		var targetURL string
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/v1/auth/"):
			targetURL = authURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/sellers") ||
			strings.HasPrefix(r.URL.Path, "/api/v1/admin/users"):
			targetURL = authURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/stats"):
			targetURL = authURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/categories"):
			targetURL = catalogURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/catalog/"):
			targetURL = catalogURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/orders"):
			targetURL = orderURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/payments"):
			targetURL = paymentURL + r.URL.Path
		case strings.HasPrefix(r.URL.Path, "/api/v1/analytics/"):
			targetURL = orderURL + r.URL.Path
		default:
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// ---- ДОБАВЛЯЕМ USER_ID В ЗАГОЛОВОК ----
		proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		proxyReq.Header = r.Header.Clone()
		if userID != "" {
			proxyReq.Header.Set("X-User-ID", userID)
		}
		proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

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
	log.Println("Payment Service: http://localhost:8084")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// getUserIDFromAuthService — проверяет токен через Auth Service и возвращает user_id
func getUserIDFromAuthService(token string) string {
	req, err := http.NewRequest("POST", "http://localhost:8081/api/v1/auth/validate", nil)
	if err != nil {
		log.Println("Failed to create validation request:", err)
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to validate token:", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Token validation failed with status:", resp.StatusCode)
		return ""
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode validation response:", err)
		return ""
	}

	log.Println("Token validated, user_id:", result["user_id"])
	return result["user_id"]
}