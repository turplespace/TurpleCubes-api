package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/turplespace/portos/internal/services/proxy"
)

type ProxyRequest struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	Subdomain string `json:"subdomain"`
}

func HandlePostProxy(w http.ResponseWriter, r *http.Request) {
	var req ProxyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.IP == "" || req.Port == 0 || req.Subdomain == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := proxy.GenerateNginxProxyConfig(req.IP, req.Port, req.Subdomain); err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to generate proxy config", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Proxy configuration generated successfully"))
}
