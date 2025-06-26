package main

import (
	"log"
	"net/http"
	"time"

	"hospital-management-system/api/routes"
	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/config"
	"hospital-management-system/internal/handlers"
	"hospital-management-system/internal/models"
	"hospital-management-system/internal/repository"
	"hospital-management-system/internal/services"
	"hospital-management-system/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.Server.GinMode)

	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	jwtService := auth.NewJWTService(cfg.JWT.Secret)

	userRepo := repository.NewUserRepository(database.GetDB())
	patientRepo := repository.NewPatientRepository(database.GetDB())

	authService := services.NewAuthService(userRepo, jwtService)
	patientService := services.NewPatientService(patientRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientService)

	createDefaultUsers(authService)

	router := routes.SetupRoutes(authHandler, patientHandler, jwtService)

	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Printf("Starting %s v%s", cfg.App.Name, cfg.App.Version)
	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Printf("Health check available at: http://localhost:%s/health", cfg.Server.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createDefaultUsers(authService *services.AuthService) {
	defaultUsers := []services.RegisterRequest{
		{
			Username:  "admin_receptionist",
			Email:     "receptionist@hospital.com",
			Password:  "password123",
			FirstName: "Admin",
			LastName:  "Receptionist",
			Role:      models.RoleReceptionist,
		},
		{
			Username:  "admin_doctor",
			Email:     "doctor@hospital.com",
			Password:  "password123",
			FirstName: "Admin",
			LastName:  "Doctor",
			Role:      models.RoleDoctor,
		},
	}

	for _, user := range defaultUsers {
		existingUser, err := authService.GetUserByID(1)
		if err != nil || existingUser == nil {
			_, err := authService.Register(user)
			if err != nil && err.Error() != "username already exists" && err.Error() != "email already exists" {
				log.Printf("Failed to create default user %s: %v", user.Username, err)
			} else if err == nil {
				log.Printf("Created default user: %s", user.Username)
			}
		}
	}

	log.Println("\n=== DEFAULT USERS FOR TESTING ===")
	log.Println("Receptionist:")
	log.Println("  Username: admin_receptionist")
	log.Println("  Password: password123")
	log.Println("Doctor:")
	log.Println("  Username: admin_doctor")
	log.Println("  Password: password123")
	log.Println("=================================")
}
