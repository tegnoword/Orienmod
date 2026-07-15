package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tegnoword/orienmod/internal/adapters/input/http/router"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func main() {
	// 1. Cargar Variables de Entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No se pudo cargar el archivo .env, usando variables del sistema")
	}

	// 2. Validar ruta de credenciales de Google
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		log.Fatal("❌ Error: GOOGLE_APPLICATION_CREDENTIALS no definida en el entorno")
	}

	// 3. Validar existencia física del JSON de credenciales
	fileInfo, err := os.Stat(credPath)
	if os.IsNotExist(err) {
		log.Fatalf("❌ Error: El archivo de credenciales no existe en: %s", credPath)
	} else if err != nil {
		log.Fatalf("❌ Error al verificar credenciales: %v", err)
	}
	if fileInfo.Size() == 0 {
		log.Fatalf("❌ Error: El archivo '%s' está vacío", credPath)
	}

	log.Printf("✅ Archivo de credenciales detectado en: %s (%d bytes)", credPath, fileInfo.Size())

	// 4. Inicializar adaptador de Google Sheets (Fase de Herramientas)
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithCredentialsFile(credPath))
	if err != nil {
		log.Fatalf("❌ Error crítico al inicializar Google Sheets: %v", err)
	}

	log.Println("🚀 ¡Servicio de Google Sheets inicializado con éxito!")
	_ = sheetsService // Este 'sheetsService' lo inyectarás luego en tus adaptadores de salida (output)

	// 5. Inicializar tu Router nativo
	// Pasamos "nil" o un mock temporalmente como 'classroomAdapter' hasta que implementes el caso de uso/puerto correspondiente
	appRouter := router.NewRouter(nil)

	// 6. Configurar puerto y arrancar el Servidor con TU router
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("📡 Servidor levantado usando tu Router Hexagonal en http://localhost:%s", port)

	// Levantamos el servidor HTTP pasándole tu estructura appRouter (que implementa ServeHTTP)
	err = http.ListenAndServe(":"+port, appRouter)
	if err != nil {
		log.Fatalf("❌ Error crítico al iniciar el servidor: %v", err)
	}
}
