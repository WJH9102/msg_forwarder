package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/WJH9102/msg_forwarder/internal/config"
	"github.com/WJH9102/msg_forwarder/internal/mailer"
)

type sendRequest struct {
	Token   string `json:"token"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func main() {
	cfg := config.Load()
	m := mailer.New(cfg)

	http.HandleFunc("POST /api/send", func(w http.ResponseWriter, r *http.Request) {
		var req sendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, response{Success: false, Message: "invalid request body"})
			return
		}

		if req.Token != cfg.AuthToken {
			writeJSON(w, http.StatusUnauthorized, response{Success: false, Message: "invalid token"})
			return
		}

		if req.To == "" || req.Content == "" {
			writeJSON(w, http.StatusBadRequest, response{Success: false, Message: "fields 'to' and 'content' are required"})
			return
		}

		subject := req.Subject
		if subject == "" {
			subject = "New Message"
			if req.Name != "" {
				subject = "Message from " + req.Name
			}
		}

		var body strings.Builder
		if req.Name != "" {
			body.WriteString("Sender: " + req.Name + "\n\n")
		}
		body.WriteString(req.Content)
		body.WriteString("\n\n---\nSent via msg_forwarder at " + time.Now().Format(time.RFC3339))

		if err := m.Send(req.To, subject, body.String()); err != nil {
			writeJSON(w, http.StatusInternalServerError, response{Success: false, Message: "failed to send email: " + err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, response{Success: true, Message: "email sent to " + req.To})
	})

	addr := ":" + cfg.ServerPort
	fmt.Printf("msg_forwarder listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}