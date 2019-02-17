package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type User struct {
	gorm.Model
	//Id           uint    `gorm:"primary_key"`
	FirstName string `gorm:"varchar(255);not null"`
	LastName  string `gorm:"varchar(255);not null"`
	Username  string `gorm:"column:username"`
	Email     string `gorm:"column:email;unique_index"`
	Password  string `gorm:"column:password;not null"`

	Comments []Comment `gorm:"foreignkey:UserId"`

	Roles     []Role     `gorm:"many2many:users_roles;"`
	UserRoles []UserRole `gorm:"foreignkey:UserId"`
}

// What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
// Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
// You can change the value in bcrypt.DefaultCost to adjust the security index.
// 	err := userModel.setPassword("password0")
func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}

// Database will only save the hashed string, you should check it by util function.
// 	if err := serModel.checkPassword("password0"); err != nil { password error }
func (u *User) IsValidPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func (user *User) BeforeSave(db *gorm.DB) (err error) {
	if len(user.Roles) == 0 {
		// role := Role{}
		userRole := Role{}
		// db.Model(&role).Where("name = ?", "ROLE_USER").First(&userRole)
		db.Model(&Role{}).Where("name = ?", "ROLE_USER").First(&userRole)
		//db.Where(&models.Role{Name: "ROLE_USER"}).Attrs(models.Role{Description: "For standard Users"}).FirstOrCreate(&userRole)
		user.Roles = append(user.Roles, userRole)
	}
	return
}

// Generate JWT token associated to this user
func (user *User) GenerateJwtToken() string {
	// jwt.New(jwt.GetSigningMethod("HS512"))
	jwt_token := jwt.New(jwt.SigningMethodHS512)

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	jwt_token.Claims = jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return token
}

func (user *User) IsAdmin() bool {
	for _, role := range user.Roles {
		if role.Name == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}
func (user *User) IsNotAdmin() bool {
	return !user.IsAdmin()
}
