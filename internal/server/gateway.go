package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	applicantsv1 "github.com/Thrun12/golang-assignment/api/proto/v1"
)

// swaggerSpec holds the loaded OpenAPI spec
var swaggerSpec []byte

// NewGatewayServer creates a new HTTP gateway server for the gRPC service
func NewGatewayServer(ctx context.Context, grpcAddress string, corsOrigins []string, db *sql.DB, logger *zap.Logger) (http.Handler, error) {
	// Load swagger spec if not already loaded
	if len(swaggerSpec) == 0 {
		data, err := os.ReadFile("api/proto/v1/applicants.swagger.json")
		if err != nil {
			logger.Warn("failed to load swagger spec", zap.Error(err))
		} else {
			swaggerSpec = data
			logger.Info("loaded swagger specification")
		}
	}

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customMatcher),
		runtime.WithErrorHandler(customErrorHandler),
	)

	// Setup connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Register gRPC-Gateway handlers
	if err := applicantsv1.RegisterApplicantsServiceHandlerFromEndpoint(ctx, mux, grpcAddress, opts); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	// Create HTTP handler with CORS
	handler := corsMiddleware(mux, corsOrigins, logger)

	// Add health check and swagger endpoints
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", healthCheckHandler(db, logger))
	healthMux.HandleFunc("/healthz", healthCheckHandler(db, logger))
	healthMux.HandleFunc("/ready", readinessCheckHandler(db, logger))
	healthMux.HandleFunc("/swagger.json", swaggerJSONHandler)
	healthMux.HandleFunc("/docs", docsRedirectHandler)
	healthMux.HandleFunc("/docs/", swaggerUIHandler)

	// Combine handlers
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/health") || strings.HasPrefix(r.URL.Path, "/ready") {
			healthMux.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/docs") || strings.HasPrefix(r.URL.Path, "/swagger") {
			healthMux.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	}), nil
}

// customMatcher matches incoming HTTP headers to gRPC metadata
func customMatcher(key string) (string, bool) {
	switch strings.ToLower(key) {
	case "x-request-id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// customErrorHandler handles errors from gRPC-Gateway
func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler, allowedOrigins []string, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(allowedOrigins) > 0 && allowedOrigins[0] == "*" {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// healthCheckHandler handles health check requests
func healthCheckHandler(db *sql.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check database connectivity
		ctx := r.Context()
		if err := db.PingContext(ctx); err != nil {
			logger.Error("health check failed: database unreachable", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "unhealthy",
				"error":  "database unreachable",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "healthy",
			"database": "connected",
		})
	}
}

// readinessCheckHandler handles readiness check requests
func readinessCheckHandler(db *sql.DB, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check if database is ready
		ctx := r.Context()
		if err := db.PingContext(ctx); err != nil {
			logger.Warn("readiness check failed: database not ready", zap.Error(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"status": "not_ready",
				"reason": "database not ready",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   "ready",
			"database": "ready",
		})
	}
}

// swaggerJSONHandler serves the OpenAPI specification
func swaggerJSONHandler(w http.ResponseWriter, r *http.Request) {
	if len(swaggerSpec) == 0 {
		http.Error(w, "Swagger specification not available", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(swaggerSpec)
}

// docsRedirectHandler redirects /docs to /docs/
func docsRedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
}

// swaggerUIHandler serves the Swagger UI
func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	// Serve a simple Swagger UI that loads from CDN
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Job Applicants API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>
        body { margin: 0; padding: 0; }
        .topbar { display: none; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}
