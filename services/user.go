package services

import (
	"darkoo/models"
	"darkoo/apperrors"
	"darkoo/utils"

	"log"
	"os"
	"time"

	"github.com/sethvargo/go-password/password"
)

type userService struct {
	UserRepository models.IUserRepository
}


func NewUserService(UserRepository models.IUserRepository) models.IUserService {
	return &userService{
		UserRepository: UserRepository,
	}
}


func (s *userService) RegisterUser(user *models.User) (*models.User, error) {
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to register user")
		return nil, apperrors.NewInternal()
	}

	user.Password = hashedPassword
	return s.UserRepository.RegisterUser(user)
}


func (s *userService) Login(email, password string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByEmailOrUserName(email)

	if err != nil {
		log.Printf("Could not find user with email: %s\n", email)
		return nil, apperrors.NewBadRequest("Could not find user with provided email")
	}

	match, err := comparePassword(user.Password, password)

	if err != nil {
		log.Print("Error checking password\n")
		return nil, apperrors.NewInternal()
	}

	if !match {
		log.Print("Incorrect password provided")
		return nil, apperrors.NewBadRequest("Incorrect password provided")
	}

	return user, nil
}


func (s *userService) GetUserById(id int) (*models.User, error) {
	return s.UserRepository.GetUserById(id)
}


func (s *userService) GetUserByUUID(uuid string) (*models.User, error) {
	return s.UserRepository.GetUserByUUID(uuid)
}


func (s *userService) GetUserByEmailOrUserName(email string) (*models.User, error) {
	return s.UserRepository.GetUserByEmailOrUserName(email)
}


func (s *userService) GetUsersByGroupId(id, limit, page int) ([]models.User, error) {
	return s.UserRepository.GetUsersByGroupId(id, limit, page)
}


func (s *userService) UpdateUser(user models.User) error {
	return s.UserRepository.UpdateUser(user)
}


func (s *userService) UpdatePassword(userId int, password string) error {
	hashedPassword, err := hashPassword(password)

	if err != nil {
		log.Print("Error processing password")
		return apperrors.NewInternal()
	}

	return s.UserRepository.UpdatePassword(userId, hashedPassword)
}


func (s *userService) ConfirmPassword(userId int, password string) error {
	user, err := s.UserRepository.GetUserById(userId)

	if err != nil {
		log.Printf("Could not find user with ID: %d", userId)
		return apperrors.NewBadRequest("Could not find user with provided ID")
	}

	match, err := comparePassword(user.Password, password)

	if err != nil {
		log.Print("Error processing password")
		return apperrors.NewInternal()
	}

	if !match {
		log.Print("Invalid password provided")
		return apperrors.NewAuthorization("Invalid password provided")
	}

	return nil
}


func (s *userService) InitLoginWithOneTimePassword(email string) error {
	user, err := s.UserRepository.GetUserByEmailOrUserName(email)

	if err != nil {
		log.Print("Could not find user by provided email")
		return apperrors.NewBadRequest("Could not find user by provided email")
	}

	otp, err := s.GenerateOneTimePasswordForUser(user, models.OneTimeLoginOTPType, time.Hour)

	if err != nil {
		log.Print("Could not generate one time password for user")
		return apperrors.NewInternal()
	}

	err = utils.SendEmail(os.Getenv("EMAIL_SENDER_EMAIL"), user.Email, "One-Time-Password", otp)
	if err != nil {
		log.Print("Could not send mail to user")
		return apperrors.NewInternal()
	}

	return nil
}


func (s *userService) GenerateOneTimePasswordForUser(user *models.User, otpType models.OTPType, duration time.Duration) (string, error) {

	otp, err := s.generateOTPPasskey(otpType)
	if err != nil {
		log.Print(err)
		return "", err
	}

	expiry := time.Now().Add(duration)

	hashedPassword, err := hashPassword(otp)
	if err != nil {
		log.Print(err)
		return "", err
	}

	err = s.UserRepository.CreateOneTimePassword(user, hashedPassword, expiry)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return otp, nil
}


func (s *userService) generateOTPPasskey(otpType models.OTPType) (string, error) {
	var passKey string
	var err error

	if otpType == models.OneTimeLoginOTPType {
		passKey, err = password.Generate(6, 2, 0, false, false)
	}

	if err != nil {
		log.Print(err)
		return "", err
	}
	return passKey, nil
}


func (s *userService) LoginWithOneTimePassword(email, password string) (*models.User, error) {
	user, err := s.UserRepository.GetUserByEmailOrUserName(email)

	if err != nil {
		log.Printf("Could not find user with email %s\n", email)
		return nil, apperrors.NewBadRequest("Could not find user with provided email")
	}

	if user.OneTimePasswordExpiry.Before(time.Now()) {
		log.Print("One time password already expired")
		return nil, apperrors.NewBadRequest("One time password already expired")
	}

	if !user.OneTimePasswordValid {
		log.Print("Invalid user credentials provided")
		return nil, apperrors.NewBadRequest(apperrors.InvalidCredentials)
	}

	match, err := comparePassword(user.OneTimePassword, password)

	if err != nil {
		log.Print("Error comparing passwords")
		return nil, apperrors.NewInternal()
	}

	if !match {
		log.Print("In correct password provided")
		return nil, apperrors.NewBadRequest(apperrors.InvalidCredentials)
	}

	err = s.UserRepository.InvalidateOneTimePassword(user)

	if err != nil {
		log.Print("Could not invalidate one-time-password")
		return nil, apperrors.NewBadRequest("Could not invalidate one-time-password")
	}

	return user, nil
}


func (s *userService) InvalidateOneTimePassword(user *models.User) error {
	err := s.UserRepository.InvalidateOneTimePassword(user)

	if err != nil {
		log.Print("Could not invalidate one-time-password")
		return apperrors.NewBadRequest("Could not invalidate one-time-password")
	}

	return nil
}


func (s *userService) EnrollTOTP(userId int) ([]byte, error) {
	user, err := s.UserRepository.GetUserById(userId)

	if err != nil {
		log.Print("Can not enroll totp. User does not exist!")
		return nil, apperrors.NewBadRequest("Could not enroll totp. User does not exist!")
	}

	totpSecret := utils.GenerateTOTPSecret()
	hashedTotpSecret, err := utils.Encrypt(totpSecret)

	if err != nil {
		log.Printf("Error on totp enrollment: %v\n", err)
		return nil, apperrors.NewInternalWithMessage("Could not enroll totp. Please try again")
	}

	err = s.UserRepository.UpdateUserTOTP(*user, hashedTotpSecret, false)

	if err != nil {
		log.Print("Could not update totp.")
		return nil, apperrors.NewInternalWithMessage("Could not update totp. Please try again")
	}

	return utils.GenerateTOTPQRCode(totpSecret, user.Email)
}


func (s *userService) VerifyTOTP(userId int, totp models.VerifyTOTPRequest) error {
	user, err := s.UserRepository.GetUserById(userId)

	if err != nil {
		log.Print("Could not verify totp. User with ID %d not found\n", userId)
		return apperrors.NewBadRequest("Could not verify totp. User not found")
	}

	match, err := s.verifyUserTotp(user, totp.Totp)
	if err != nil {
		log.Printf("Error verifing user totp %v\n", err)
		return apperrors.NewInternal()
	}

	if !match {
		log.Print("One time password is not valid")
		return apperrors.NewBadRequest("One time password is not valid")
	}

	var updatedUser models.User
	updatedUser.ID = user.ID
	updatedUser.TotpEnabled = true

	err = s.UpdateUser(updatedUser)

	if err != nil {
		log.Print("Error updating user totp")
		return apperrors.NewInternal()
	}

	return nil
}


func (s *userService) verifyUserTotp(user *models.User, totp string) (bool, error) {
	if totp == "" {
		log.Printf("User totp is empty. Totp should be be enrolled!")
		return false, apperrors.NewBadRequest("User totp is invalid. Please enroll your totp.")
	}

	totpSecret, err := utils.Decrypt(totp)
	if err != nil {
		return false, apperrors.NewInternalWithMessage("Unable to verify totp. Please try again.")
	}

	return utils.VerifyTOTP(totpSecret, totp), nil
}


func (s *userService) DisableTOTP(userId int) error {
	user, err := s.UserRepository.GetUserById(userId)

	if err != nil {
		log.Printf("Unable to disable totp because user with ID %d is does not exist.\n", userId)
		return apperrors.NewBadRequest("Could not disable totp because user does not exist")
	}

	err = s.UserRepository.UpdateUserTOTP(*user, "", false)

	if err != nil {
		log.Print("Unable to disable user totp. Please try again.")
		return apperrors.NewBadRequest("Unable to disable user totp. Please try again.")
	}

	return nil
}