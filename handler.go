package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type SendRequest struct {
	Token   string `json:"token"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func NewHandler(mailer *Mailer, cfg *Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/send", sendHandler(mailer, cfg))
	return mux
}

func sendHandler(mailer *Mailer, cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "invalid request body",
			})
			return
		}

		if req.Token != cfg.AuthToken {
			writeJSON(w, http.StatusUnauthorized, Response{
				Success: false,
				Message: "invalid token",
			})
			return
		}

		if req.To == "" || req.Content == "" {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "fields 'to' and 'content' are required",
			})
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

		if err := mailer.Send(req.To, subject, body.String()); err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "failed to send email: " + err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "email sent to " + req.To,
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}