package models

import "time"

type Log struct {
	IDLog        int       `json:"id_log"`
	Method       string    `json:"method"`
	Path         string    `json:"path"`
	Protocol     string    `json:"protocol"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int       `json:"response_time"` // en milisegundos
	UserAgent    *string   `json:"user_agent"`
	IP           string    `json:"ip"`
	Hostname     string    `json:"hostname"`
	Body         *string   `json:"body"`
	Params       *string   `json:"params"`
	Query        *string   `json:"query"`
	Email        *string   `json:"email"`        // del usuario autenticado
	Username     *string   `json:"username"`     // del usuario autenticado
	Role         *string   `json:"role"`         // rol del usuario
	LogLevel     string    `json:"log_level"`    // info, warning, error
	Environment  string    `json:"environment"`  // development, production
	NodeVersion  string    `json:"node_version"` // versión del sistema
	PID          int       `json:"pid"`          // process ID
	Timestamp    time.Time `json:"timestamp"`
	URL          string    `json:"url"` // URL completa
	CreatedAt    time.Time `json:"created_at"`
}
