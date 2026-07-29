// cmd/api/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/tegnoword/orienmod/configs"
	"github.com/tegnoword/orienmod/internal/adapters/input/http/router"
	"github.com/tegnoword/orienmod/internal/adapters/output/google"
	"github.com/tegnoword/orienmod/internal/adapters/output/storage"
	graphql "github.com/tegnoword/orienmod/internal/graph"
)

func main() {
	log.Println("🚀 Iniciando Orienmod...")
	log.Println("📋 Cargando configuración OAuth2...")

	var oauthConfig *oauth2.Config
	var err error

	// Intentar cargar desde credentials.json
	oauthConfig, err = configs.NewOAuthConfigFromFile()
	if err != nil {
		log.Printf("⚠️ Error al cargar credentials.json: %v", err)
		log.Println("⚠️ Intentando con variables de entorno...")
		oauthConfig = configs.NewOAuthConfig()
	}

	if oauthConfig == nil {
		log.Fatal("❌ No se pudo cargar la configuración OAuth2. Revisa credentials.json o variables de entorno.")
	}

	if oauthConfig.ClientID == "" {
		log.Fatal("❌ GOOGLE_CLIENT_ID no configurado. Revisa credentials.json o variables de entorno.")
	}

	log.Printf("✅ OAuth2 configurado correctamente")
	log.Printf("   Client ID: %s...", oauthConfig.ClientID[:15])
	log.Printf("   Redirect URI: %s", oauthConfig.RedirectURL)

	// ============================================
	// 2. ALMACENAMIENTO DE TOKENS
	// ============================================
	log.Println("📋 Inicializando almacenamiento de tokens...")
	tokenStore := storage.NewMemoryTokenStore()

	if tokenStore == nil {
		log.Fatal("❌ Error al crear MemoryTokenStore")
	}
	log.Println("✅ Almacenamiento de tokens inicializado (en memoria)")

	// ============================================
	// 3. ADAPTADOR DE GOOGLE
	// ============================================
	log.Println("📋 Inicializando adaptador de Google...")

	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	if spreadsheetID == "" {
		log.Println("⚠️ GOOGLE_SPREADSHEET_ID no configurada (opcional)")
		spreadsheetID = "placeholder"
	}

	googleAdapter := google.NewGoogleClientAdapter(oauthConfig, tokenStore, spreadsheetID)

	if googleAdapter == nil {
		log.Fatal("❌ Error al crear GoogleClientAdapter")
	}
	log.Println("✅ Adaptador de Google inicializado")

	// ============================================
	// 4. GRAPHQL HANDLER
	// ============================================
	log.Println("📋 Configurando GraphQL...")

	// ✅ Crear handler de GraphQL con graphql-go
	graphQLHandler := graphql.NewGraphQLHandler(googleAdapter, tokenStore, oauthConfig)

	if graphQLHandler == nil {
		log.Fatal("❌ Error al crear GraphQL Handler")
	}
	log.Println("✅ GraphQL configurado")

	// ============================================
	// 5. ROUTER
	// ============================================
	log.Println("📋 Configurando router...")

	r := router.NewRouter(oauthConfig, tokenStore, googleAdapter)

	if r == nil {
		log.Fatal("❌ Error al crear Router")
	}

	// ✅ Agregar rutas de GraphQL al router
	r.AddGraphQLRoutes(graphQLHandler)

	log.Println("✅ Router configurado")

	// ============================================
	// 6. SERVIDOR HTTP
	// ============================================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ============================================
	// 7. INICIAR SERVIDOR
	// ============================================
	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	log.Println("📚 Endpoints disponibles:")
	log.Println("")
	log.Println("🔗 REST API:")
	log.Println("  GET  /health")
	log.Println("  GET  /api/v1/auth/google")
	log.Println("  GET  /api/v1/auth/google/callback")
	log.Println("  POST /api/v1/auth/logout")
	log.Println("  GET  /api/v1/courses")
	log.Println("  POST /api/v1/courses")
	log.Println("  GET  /api/v1/courses/{id}/students")
	log.Println("  POST /api/v1/courses/{id}/sync")
	log.Println("  QUERY /api/v1/courses/search")
	log.Println("  QUERY /api/v1/students/search")
	log.Println("  POST /api/v1/classroom/webhook")
	log.Println("")
	log.Println("🔗 GRAPHQL:")
	log.Println("  GET  /graphql    - Playground de GraphQL")
	log.Println("  POST /query      - Endpoint GraphQL")
	log.Println("")

	// Iniciar servidor en goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error al iniciar servidor: %v", err)
		}
	}()

	// ============================================
	// 8. GRACEFUL SHUTDOWN
	// ============================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Apagando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Error al apagar servidor: %v", err)
	}

	log.Println("✅ Servidor apagado correctamente")
}
