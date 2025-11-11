package Routes

import (
	Controllers "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/controllers"
	Middlewares "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/middlewares"
	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes with proper middleware chain
func SetupRoutes(router *gin.Engine, jwtManager *Utils.JWTManager) {
	// Apply global middlewares in optimal order
	router.Use(

		Middlewares.LoggingMiddleware(), // Logging for all requests
		gin.Recovery(),                  // Recovery middleware for panic handling
	)

	// Public authentication routes with security middlewares
	auth := router.Group("/api/auth")
	auth.Use(
		Middlewares.RateLimitMiddleware(),  // Rate limiting for auth endpoints
		Middlewares.BruteforceMiddleware(), // Brute force protection
	)
	{
		// Authentication endpoints
		auth.POST("/login", Controllers.LoginHandler)
		auth.POST("/register/patient", Controllers.RegisterPatientHandler)
		auth.POST("/register/doctor", Controllers.RegisterDoctorsHandler)
	}

	// Protected API routes with authentication and additional security
	api := router.Group("/api")
	api.Use(
		Middlewares.CORSMiddleware(),
		Middlewares.RateLimitMiddleware(),      // Rate limiting for API calls
		Middlewares.AuthMiddleware(jwtManager), // JWT authentication with proper dependency injection
	)

	// APIs for patient role based access control with grouping
	patient := api.Group("/patient")
	patient.Use(Middlewares.RoleGuard("patient"))
	{
		patient.GET("/profile", Controllers.PatientProfile)
		patient.GET("/book-appointment", Controllers.BookAppointment)
		patient.GET("/appointments", Controllers.GetAppointments)
	}

	// APIs for doctor role based acees control with grouping routes
	doctor := api.Group("/doctor")
	doctor.Use(Middlewares.RoleGuard("doctor"))
	{
		doctor.GET("/dashboard", Controllers.DoctorDashboard)
		doctor.POST("/updates-status", Controllers.UpdateStatus)
		doctor.GET("/patients", Controllers.GetDoctorPatients)
	}

	admin := api.Group("/admin")
	admin.Use(Middlewares.RoleGuard("admin"))
	{
		admin.GET("/dashboard", Controllers.AdminDashboard)
		admin.POST("/add-doctor", Controllers.AddDoctor)
		admin.DELETE("/delete-user/:id", Controllers.DeleteUser)
	}

}
