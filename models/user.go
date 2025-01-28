package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	OneTimeLoginOTPType     OTPType = "OneTimeLoginOTPType"
)


type OTPType string


type User struct {
	Base
	UserName 		string		 `json:"userName"`
	Email			string		 `json:"email" gorm:"unique"`
	Password 		string   	 `json:"-"`
	Gender          string       `json:"gender"`
	Groups 			[]Group		 `gorm:"many2many:user_groups"`
	Messages  		[]Message	 `gorm:"constraint:OnUpdate:CASCADE, OnDelete:SET NULL"`
	ImageNum		int 		 `json:"image"`
	Ban 			bool  		 `json:"ban" gorm:"type:bool;default:false"`
	TotpEnabled     bool  		 `json:"totpEnabled" gorm:"type:bool;default:false"`
	TotpSecret      string 		 `json:"_"`
	OneTimePassword string   	 `json:"-"`
	OneTimePasswordExpiry time.Time `json:"oneTimePasswordExpiry"`
	OneTimePasswordValid  bool      `json:"oneTimePasswordValid" gorm:"type:bool;default:false"`
}


type IUserRepository interface {
	RegisterUser(user *User) (*User, error)
	JoinGroup(userId, groupId int) error
	LeaveGroup(userId, groupId int) error
	GetUserById(id int) (*User, error)
	GetUserByUUID(uuid string) (*User, error)
	GetUserByEmailOrUserName(email string) (*User, error)
	GetUsersByGroupId(groupId, limit, page int) ([]User, error)
	CreateOneTimePassword(user *User, password string, expiry time.Time) error
	InvalidateOneTimePassword(user *User) error
	UpdateUserTOTP(user User, totpSecret string, totpEnabled bool) error
	UpdateUser(user User) error
	UpdatePassword(userId int, password string) error
	UpdateUserImageNum(userId, num int) (int, error)
}


type IUserService interface {
	RegisterUser(user *User) (*User, error)
	Login(email, password string) (*User, error)
	GetUserById(id int) (*User, error)
	GetUserByUUID(uuid string) (*User, error)
	GetUserByEmailOrUserName(email string) (*User, error)
	GetUsersByGroupId(groupId, limit, page int) ([]User, error)
	UpdateUser(user User) error
	UpdatePassword(userId int, password string) error
	ConfirmPassword(userId int, password string) error
	InitLoginWithOneTimePassword(email string) (error)
	LoginWithOneTimePassword(email, code string) (*User, error)
	GenerateOneTimePasswordForUser(user *User, otpType OTPType, duration time.Duration) (string, error)
	InvalidateOneTimePassword(user *User) error
	EnrollTOTP(userId int) ([]byte, error)
	VerifyTOTP(userId int, verifyTOTP VerifyTOTPRequest) error
	DisableTOTP(userId int) error
	UpdateUserImageNum(userId, num int) (int, error)
	JoinGroup(userId, groupId int) error
	// LeaveGroup(userId, groupId int) error
}


func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

type VerifyTOTPRequest struct {
	Totp string `json: "totp"`
}

func (r VerifyTOTPRequest) Validate() error {
	return nil
}

// func (r VerifyTOTPRequest) Sanitize() {

// }