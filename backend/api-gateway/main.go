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

		// ---- ПОЛУЧАЕМ USER_ID И ROLE ЧЕРЕЗ AUTH SERVICE ----
		var userID string
		var role string
		tokenStr := r.Header.Get("Authorization")
		if tokenStr != "" && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			userID, role = getUserInfoFromAuthService(tokenStr)
		}

		// ---- МАРШРУТИЗАЦИЯ ----
		var targetURL string
		
		// СОХРАНЯЕМ ПОЛНЫЙ URL С ПАРАМЕТРАМИ
		fullPath := r.URL.RequestURI()
		
		switch {
		case strings.HasPrefix(r.URL.Path, "/api/v1/auth/"):
			targetURL = authURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/sellers") ||
			strings.HasPrefix(r.URL.Path, "/api/v1/admin/users"):
			targetURL = authURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/stats"):
			targetURL = authURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/admin/categories"):
			targetURL = catalogURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/catalog/"):
			targetURL = catalogURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/orders"):
			targetURL = orderURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/payments"):
			targetURL = paymentURL + fullPath
		case strings.HasPrefix(r.URL.Path, "/api/v1/analytics/"):
			targetURL = orderURL + fullPath
		default:
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		log.Printf("Proxying: %s -> %s", r.URL.String(), targetURL)

		// ---- ДОБАВЛЯЕМ USER_ID И ROLE В ЗАГОЛОВОК ----
		proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		proxyReq.Header = r.Header.Clone()
		if userID != "" {
			proxyReq.Header.Set("X-User-ID", userID)
		}
		if role != "" {
			proxyReq.Header.Set("X-User-Role", role)
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

// getUserInfoFromAuthService — проверяет токен через Auth Service и возвращает user_id и role
func getUserInfoFromAuthService(token string) (string, string) {
	req, err := http.NewRequest("POST", "http://localhost:8081/api/v1/auth/validate", nil)
	if err != nil {
		log.Println("Failed to create validation request:", err)
		return "", ""
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to validate token:", err)
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Token validation failed with status:", resp.StatusCode)
		return "", ""
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode validation response:", err)
		return "", ""
	}

	userID := result["user_id"]
	role := result["role"]
	if role == "" {
		role = "customer" // по умолчанию
	}

	log.Printf("Token validated: user_id=%s, role=%s", userID, role)
	return userID, role
}