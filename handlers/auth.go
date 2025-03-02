package handlers

import(
	"encoding/json"
	"net/http"
	"time"
	
	"digi-notice-board/db"
	"digi-notice-board/models"
	
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)


var jwtKey = []byte("default#44")

//have to change this hardcoded test key to environment variable before pushing the code to production

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	//hashing user password
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}
	// default user role is 'user'
	
	user.Password = string(hashedPassword)
	user.Role = "user" 
	user.CreatedAt = time.Now()
	
	//user record creation and error message
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error creating user",http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

//login payload

type LoginRequest struct {
	Email		string	`json:"email"`
	Password		string	`json:"password"`
}

//login handles user auth and returns a jwt

func Login(w http.ResponseWriter, r *http.Request){
	var creds LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	
	//finding user by email
	
	var user models.User
	if err := db.DB.Where("email = ?", creds.Email).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	
	//password verification
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		//log.Println("Bcrypt error:", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	
	//jwt creation with user ID and role
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":	user.ID,
		"role":		user.Role,
		"exp":		time.Now().Add(72 * time.Hour).Unix(),
	})
	
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
