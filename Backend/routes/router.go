package Routes

import (
	Controllers "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/controllers"
	Middlewares "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/middlewares"
	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures protected HTML pages and API routes
func SetupRoutes(router *gin.Engine, jwtManager *Utils.JWTManager) {

	// ------------------------
	// Auth APIs (JSON) only
	// ------------------------
	auth := router.Group("/api/auth")
	auth.Use(
		Middlewares.RateLimitMiddleware(),
		Middlewares.BruteforceMiddleware(),
	)
	{
		auth.POST("/login", Controllers.LoginHandler)
		auth.POST("/register/patient", Controllers.RegisterPatientHandler)
		auth.POST("/register/doctor", Controllers.RegisterDoctorsHandler)
	}

	// ------------------------
	// Protected HTML pages (JWT + Role)
	// ------------------------
	// Patient pages
	patient := router.Group("/patient")
	patient.Use(
		Middlewares.AuthMiddleware(jwtManager),
		Middlewares.RoleGuard("patient"),
	)
	{
		patient.GET("/dashboard", Controllers.PatientDashboard)
		patient.GET("/book-appointment", Controllers.BookAppointmentPage) // optional HTML page
	}

	// Doctor pages
	doctor := router.Group("/doctor")
	doctor.Use(
		Middlewares.AuthMiddleware(jwtManager),
		Middlewares.RoleGuard("doctor"),
	)
	{
		doctor.GET("/dashboard", Controllers.DoctorDashboard)
	}

	// Admin pages
	admin := router.Group("/admin")
	admin.Use(
		Middlewares.AuthMiddleware(jwtManager),
		Middlewares.RoleGuard("admin"),
	)
	{
		admin.GET("/dashboard", Controllers.AdminDashboard)
	}

	// ------------------------
	// Protected API routes (JSON)
	// ------------------------
	api := router.Group("/api")
	api.Use(
		Middlewares.CORSMiddleware(),
		Middlewares.RateLimitMiddleware(),
		Middlewares.AuthMiddleware(jwtManager),
	)

	// Patient APIs
	patientAPI := api.Group("/patient")
	patientAPI.Use(Middlewares.RoleGuard("patient"))
	{
		patientAPI.GET("/appointments", Controllers.GetAppointments)
		patientAPI.POST("/book-appointment", Controllers.BookAppointment)
		patientAPI.GET("/profile", Controllers.PatientProfile)
	}

	// Doctor APIs
	doctorAPI := api.Group("/doctor")
	doctorAPI.Use(Middlewares.RoleGuard("doctor"))
	{
		doctorAPI.GET("/patients", Controllers.GetDoctorPatients)
		doctorAPI.POST("/update-status", Controllers.UpdateStatus)
	}

	// Admin APIs
	adminAPI := api.Group("/admin")
	adminAPI.Use(Middlewares.RoleGuard("admin"))
	{
		adminAPI.POST("/add-doctor", Controllers.AddDoctor)
		adminAPI.DELETE("/delete-user/:id", Controllers.DeleteUser)
	}
}
