package Controllers

import (
	"net/http"
	"strings"
	"time"

	Configs "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/config"
	Models "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/models"
	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request struct for registration
type RegisterRequest struct {
	Fullname        string `json:"fullname" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	PhoneNumber     string `json:"phonenumber" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirmpassword" binding:"required"`
}

// Request struct for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// Register patient
func RegisterPatientHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check duplicate email
	var existing Models.User
	if err := Configs.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	hashPassword, err := Utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := Models.User{
		UUID:        uuid.New().String(),
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashPassword,
		Role:        "patient",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := Configs.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register patient"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Patient registered successfully"})
}

// Register doctor
func RegisterDoctorsHandler(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check duplicate email
	var existing Models.User
	if err := Configs.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	hashPassword, err := Utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := Models.User{
		UUID:        uuid.New().String(),
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashPassword,
		Role:        "doctor",
		IsActive:    false, // doctors require admin approval
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := Configs.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register doctor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Doctor registered successfully. Await admin approval."})
}
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user by email + role
	var user Models.User
	if err := Configs.DB.Where("email = ? AND role = ?", req.Email, req.Role).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	isMatch, err := Utils.VerifyPassword(req.Password, user.Password)
	if err != nil || !isMatch {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	if req.Role == "doctor" && !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account pending admin approval"})
		return
	}

	jwtManager, _ := Utils.NewJWTManager("./certs/jwt_private.pem", "./certs/jwt_public.pem", 24*time.Hour)
	token, _ := jwtManager.GenerateToken(user.ID, user.Role)

	// Set token in secure HTTP-only cookie
	c.SetCookie("auth_token", token, 3600*24, "/", "localhost", true, true) // adjust domain for production

	// Redirect to dashboard
	if user.Role == "patient" {
		c.Redirect(http.StatusSeeOther, "/patient/dashboard")
		return
	} else if user.Role == "doctor" {
		c.Redirect(http.StatusSeeOther, "/doctor/dashboard")
	} else {
		c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}
}
