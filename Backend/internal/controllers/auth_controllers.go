package Controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	Configs "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/config"
	Models "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/models"
	Utils "github.com/Aks98068/Secure_Health_record_managemnt_system-SHRMS/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterPatientHandler(c *gin.Context) {
	//  take data form the html form
	fullname := c.PostForm("fullname")
	email := strings.ToLower(strings.TrimSpace(c.PostForm("name")))
	phonenumber := c.PostForm("phonenumber")
	password := c.PostForm("password")
	confirmpassword := c.PostForm("confirmpassword")
	role := "patient"

	// validation
	if fullname == "" || email == "" || phonenumber == "" || password == "" || confirmpassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please fill all required fields and try again."})
		return
	}

	if password != confirmpassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password do not match."})
	}

	// check duplicate email or usename
	var existing Models.User
	if err := Configs.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email alrady exists"})
		return
	}

	hashPassword, err := Utils.HashPassword(password)
	if err != nil {
		log.Println("password hashing faield :", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user."})
	}

	user := Models.User{
		UUID:        uuid.New().String(),
		Fullname:    fullname,
		Email:       email,
		PhoneNumber: phonenumber,
		Password:    hashPassword,
		Role:        role,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := Configs.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to register user."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Patient registered successfully."})
}

func RegisterDoctorsHandler(c *gin.Context) {

	fullname := c.PostForm("fullname")
	email := strings.TrimSpace(c.PostForm("email"))
	phonenumber := strings.TrimSpace(c.PostForm("phonenumber"))
	password := strings.TrimSpace(c.PostForm("password"))
	confirmpassword := strings.TrimSpace(c.PostForm("confirmpassword"))
	role := "doctor"

	if fullname == "" || email == "" || phonenumber == "" || password == "" || confirmpassword == "" || role == "" {
		fmt.Println("all field required")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "all fields are required"})
	}

	if password != confirmpassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password does n"})
	}

	var existing Models.User
	if err := Configs.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exist"})
		return
	}

	hashedpassword, err := Utils.HashPassword(password)
	if err != nil {
		fmt.Printf("password hashing failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register User"})
	}

	user := Models.User{
		UUID:        uuid.New().String(),
		Fullname:    fullname,
		Email:       email,
		PhoneNumber: phonenumber,
		Password:    hashedpassword,
		Role:        role,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := Configs.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register Doctor"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Doctor created successfully"})

}

var jwtSecret = []byte("")

func LoginHandler(c *gin.Context) {

	email := strings.TrimSpace(c.PostForm("email"))
	password := strings.TrimSpace(c.PostForm("password"))
	role := strings.ToLower(strings.TrimSpace(c.PostForm("role")))

	if email == "" || password == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "all fields are required"})
	}
	validRoles := map[string]bool{
		"admin":   true,
		"doctor":  true,
		"patient": true,
	}

	if !validRoles[role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	//fetch user by email and role
	var user Models.User
	if err := Configs.DB.Where("email = ? AND  role = ?", email, role).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// check password (Argon2)
	if !Utils.VerifyPassword(password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect Password"})
		return
	}
	// Doctors must be approved by admin
	if role == "doctor" && user.IsActive == false {
		c.JSON(http.StatusForbidden, gin.H{"error": "your account is pending admin approval"})
		return
	}

	// generate JWT token
	token, err := Utils.GenerateToken(user.UUID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"role":    role,
	})

}
