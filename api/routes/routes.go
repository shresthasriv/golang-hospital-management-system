package routes

import (
	"hospital-management-system/internal/auth"
	"hospital-management-system/internal/handlers"
	"hospital-management-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authHandler *handlers.AuthHandler,
	patientHandler *handlers.PatientHandler,
	jwtService *auth.JWTService,
) *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Hospital Management System API is running",
		})
	})

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
		}

		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.AuthMiddleware(jwtService))
		{
			authProtected.POST("/refresh", authHandler.RefreshToken)
			authProtected.GET("/profile", authHandler.GetProfile)
		}

		patients := v1.Group("/patients")
		patients.Use(middleware.AuthMiddleware(jwtService))
		{
			patients.GET("", middleware.RequireReceptionistOrDoctor(), patientHandler.ListPatients)
			patients.GET("/search", middleware.RequireReceptionistOrDoctor(), patientHandler.SearchPatients)
			patients.GET("/:id", middleware.RequireReceptionistOrDoctor(), patientHandler.GetPatient)
			patients.GET("/by-patient-id/:patient_id", middleware.RequireReceptionistOrDoctor(), patientHandler.GetPatientByPatientID)
			patients.PUT("/:id", middleware.RequireReceptionistOrDoctor(), patientHandler.UpdatePatient)

			patients.POST("", middleware.RequireReceptionist(), patientHandler.CreatePatient)
			patients.DELETE("/:id", middleware.RequireReceptionist(), patientHandler.DeletePatient)
		}
	}

	return router
}
