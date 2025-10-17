package main

import (
	"ecom-backend/api/handler"
	"ecom-backend/api/router"
	"ecom-backend/application/service"
	"ecom-backend/infrastructure/database"
	"ecom-backend/infrastructure/persistence"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Load configuration from environment variables
	cfg := &database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "ecom"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Initialize database connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established")

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database migrations completed")

	// Initialize repositories (Infrastructure layer)
	productRepo := persistence.NewProductRepository(db)
	basketRepo := persistence.NewBasketRepository(db)
	orderRepo := persistence.NewOrderRepository(db)

	// Initialize services (Application layer)
	productService := service.NewProductService(productRepo)
	basketService := service.NewBasketService(basketRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, basketRepo, productRepo)

	// Initialize handlers (API layer)
	productHandler := handler.NewProductHandler(productService)
	basketHandler := handler.NewBasketHandler(basketService)
	orderHandler := handler.NewOrderHandler(orderService)

	// Setup router
	r := router.Setup(productHandler, basketHandler, orderHandler)

	// Start server
	port := getEnv("PORT", "8080")
	addr := ":" + port

	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as int or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
