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

	"github.com/tegnoword/orienmod/configs"
	"github.com/tegnoword/orienmod/internal/adapters/input/http/router"
	"github.com/tegnoword/orienmod/internal/adapters/output/google"
	"github.com/tegnoword/orienmod/internal/adapters/output/storage"
	"golang.org/x/oauth2"
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

		// Fallback a variables de entorno
		oauthConfig = configs.NewOAuthConfig()
	}

	// ✅ VERIFICACIÓN CRÍTICA: asegurar que oauthConfig no sea nil
	if oauthConfig == nil {
		log.Fatal("❌ No se pudo cargar la configuración OAuth2. Revisa credentials.json o variables de entorno.")
	}

	// ✅ VERIFICACIÓN CRÍTICA: asegurar que ClientID no esté vacío
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

	// Verificar que tokenStore no sea nil
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

	// ✅ VERIFICACIÓN: asegurar que googleAdapter no sea nil
	if googleAdapter == nil {
		log.Fatal("❌ Error al crear GoogleClientAdapter")
	}
	log.Println("✅ Adaptador de Google inicializado")

	// ============================================
	// 4. ROUTER
	// ============================================
	log.Println("📋 Configurando router...")

	r := router.NewRouter(oauthConfig, tokenStore, googleAdapter)

	// ✅ VERIFICACIÓN: asegurar que router no sea nil
	if r == nil {
		log.Fatal("❌ Error al crear Router")
	}
	log.Println("✅ Router configurado")

	// ============================================
	// 5. SERVIDOR HTTP
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
	// 6. INICIAR SERVIDOR
	// ============================================
	log.Printf("🚀 Servidor corriendo en http://localhost:%s", port)
	log.Println("📚 Endpoints disponibles:")
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

	// Iniciar servidor en goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Error al iniciar servidor: %v", err)
		}
	}()

	// ============================================
	// 7. GRACEFUL SHUTDOWN
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
