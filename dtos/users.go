package dtos

import (
	"github.com/melardev/GoGonicEcommerceApi/models"
)

type RegisterRequestDto struct {
	Username             string `form:"username" json:"username" xml:"username"  binding:"required"`
	FirstName            string `form:"first_name" json:"first_name" xml:"first_name" binding:"required"`
	LastName             string `form:"last_name" json:"last_name" xml:"last_name" binding:"required"`
	Email                string `form:"email" json:"email" xml:"email" binding:"required"`
	Password             string `form:"password" json:"password" xml:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation" xml:"password-confirmation" binding:"required"`
}

type LoginRequestDto struct {
	// Username string `form:"username" json:"username" xml:"username" binding:"exists,username"`
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password"json:"password" binding:"exists,min=8,max=255"`

	userModel models.User `json:"-"`
}

func CreateLoginSuccessful(user *models.User) map[string]interface{} {
	var roles = make([]string, len(user.Roles))

	for i := 0; i < len(user.Roles); i++ {
		roles[i] = user.Roles[i].Name
	}

	return map[string]interface{}{
		"success": true,
		"token":   user.GenerateJwtToken(),
		"user": map[string]interface{}{
			"username": user.Username,
			"id":       user.ID,
			"roles":    roles,
		},
	}
}

func GetUserBasicInfo(user models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	}
}
