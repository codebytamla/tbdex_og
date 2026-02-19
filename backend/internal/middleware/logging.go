package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.wroteHeader = true
	}
	return rw.ResponseWriter.Write(b)
}

// RequestLogger logs incoming HTTP requests with method, path, status, duration, and IP.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := wrapResponseWriter(w)
		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		ip := realIP(r)

		// Color-code status for terminal readability
		statusColor := colorForStatus(wrapped.status)
		methodColor := colorForMethod(r.Method)

		log.Printf("%s %s %s | %s%d%s | %10v | %s",
			methodColor, r.Method, colorReset,
			statusColor, wrapped.status, colorReset,
			duration.Round(time.Microsecond),
			r.URL.Path,
		)

		// Log additional details for errors (4xx, 5xx)
		if wrapped.status >= 400 {
			log.Printf("       ↳ IP: %s | UA: %s", ip, truncate(r.UserAgent(), 60))
		}
	})
}

// realIP extracts the client's real IP address.
func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	// Remove port from RemoteAddr
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ANSI color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func colorForStatus(status int) string {
	switch {
	case status >= 500:
		return colorRed // Server error
	case status >= 400:
		return colorYellow // Client error
	case status >= 300:
		return colorCyan // Redirect
	case status >= 200:
		return colorGreen // Success
	default:
		return colorWhite
	}
}

func colorForMethod(method string) string {
	switch method {
	case http.MethodGet:
		return colorGreen
	case http.MethodPost:
		return colorBlue
	case http.MethodPut:
		return colorPurple
	case http.MethodPatch:
		return colorYellow
	case http.MethodDelete:
		return colorRed
	default:
		return colorWhite
	}
}
