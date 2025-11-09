package Routes

import (
	Controllers "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/controllers"
	Middlewares "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// authentication routes
	auth := router.Group("/api/auth")
	auth.Use(Middlewares.RateLimitMiddleware(), Middlewares.BruteforceMiddleware())
	{
		// secure routes for login and register
		auth.POST("/login", Controllers.LoginHandler)
		auth.POST("/register/patient", Controllers.RegisterPatientHandler)
		auth.POST("/register/doctor", Controllers.RegisterDoctorsHandler)
	}

	// seure APIs with grouping and  protected with Authmiddlewares
	api := router.Group("/api")
	api.Use(Middlewares.AuthMiddleware)

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
