package repository

import (
	"darkoo/models"
	"darkoo/apperrors"

	"errors"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)


type userRepository struct {
	DB *gorm.DB
}


func NewUserRepository(db *gorm.DB) models.IUserRepository {
	return &userRepository{ DB: db,}
}


func (r *userRepository) RegisterUser(user *models.User) (*models.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		log.Printf("Could not register user with email: %s. %v\n", user.Email, result.Error)
		return nil, apperrors.NewAuthorization("Could not register user")
	}

	return user, nil
}


func (r *userRepository) JoinGroup(userId, groupId int) error {
	group := &models.Group{}
	user := &models.User{}
	userGroup := &models.UserGroup{}

	if err := r.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		log.Print("User does not exist")
		return apperrors.NewBadRequest("User does not exist")
	}

	if err := r.DB.Where("id = ?", groupId).First(&group).Error; err != nil {
		log.Print("Group does not exist")
		return apperrors.NewBadRequest("Group does not exist")
	}

	err := r.DB.Where("user_id = ? AND group_id = ?", userId, groupId).First(&userGroup) 
	
	if err.Error != nil {
		if errors.Is(err.Error, gorm.ErrRecordNotFound){
			log.Print("User is not already a member of group")
		}else {
			return apperrors.NewInternal()
		}
	}

	if err == nil {
		log.Print("User is already a member of group")
		return apperrors.NewBadRequest("User is already a member of group")
	}

	userGroup.UserId = user.ID
	userGroup.GroupId = group.ID

	if err := r.DB.Create(&userGroup).Error; err != nil {
		log.Print("Could not join group")
		return apperrors.NewBadRequest("Could not join group")
	}

	return nil
}


func (r *userRepository) LeaveGroup(userId, groupId int) error {
	group := &models.Group{}
	user := &models.User{}
	userGroup := &models.UserGroup{}

	if err := r.DB.Where("id = ?", groupId).First(&group).Error; err != nil {
		log.Printf("Could not find group with provided ID")
		return apperrors.NewNotFound("Group", strconv.Itoa(groupId))
	}

	if err := r.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		log.Printf("Could not find user with ID: %v\n", userId)
		return apperrors.NewNotFound("User", strconv.Itoa(userId))
	}


	if err := r.DB.Where("user_id = ? AND group_id = ?", userId, groupId).First(&userGroup).Error; err != nil {
		log.Print("User is not a part of this group")
		return apperrors.NewBadRequest("User is not a part of this group")
	}

	if err := r.DB.Where("id = ?", userGroup.ID).Delete(&models.UserGroup{}).Error; err != nil {
		log.Print("Could not delete user group association")
		return apperrors.NewBadRequest("Could not delete user group association")
	}

	return nil
}


func (r *userRepository) GetUserById(id int) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		log.Printf("Could not find user with ID: %d\n", id)
		return nil, apperrors.NewBadRequest("Could not find user with provided ID")
	}

	return user, nil
}


func (r *userRepository) GetUserByUUID(uuid string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		log.Printf("Could not find user with uuid %s\n", uuid)
		return nil, apperrors.NewBadRequest("Could not find user with provided uuid")
	}

	return user, nil
}

func (r *userRepository) GetUserByEmailOrUserName(email string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("email = ?", email).Or("user_name = ?", email).First(&user).Error; err != nil {
		log.Printf("Could not find user with eamil or username: %s\n", email)
		return nil, apperrors.NewBadRequest("Could not find user with provided email or username")
	}

	return user, nil
}


func (r *userRepository) GetUsersByGroupId(id, limit, page int) ([]models.User, error) {
	group := &models.Group{}
	var users []models.User

	if err := r.DB.Where("id = ?", id).First(&group).Error; err != nil {
		return users, apperrors.NewBadRequest("Could not find group with provided ID")
	}

	if err := r.DB.Joins("JOIN user_groups ON user_groups.user_id = users.id").
					Where("user_groups.group_id = ?", id).
					Find(&users).Error; err != nil {
						log.Print("Failed to get users by group ID")
						return users, apperrors.NewBadRequest("Failed to get users by group ID")
					}

	return users, nil
}


func (r *userRepository) CreateOneTimePassword(user *models.User, password string, expiry time.Time) error {
	result := r.DB.Model(&user).Updates(models.User{OneTimePassword: password, OneTimePasswordValid: true, OneTimePasswordExpiry: expiry})

	if result.Error != nil {
		log.Print("Could not create one time password")
		return result.Error
	}
	return nil
}


func (r *userRepository) InvalidateOneTimePassword(user *models.User) error {
	if err := r.DB.Model(&user).Updates(models.User{OneTimePasswordValid: false}).Error; err != nil {
		log.Print("Could not invalidate one time password")
		return apperrors.NewBadRequest("Could not update one time password")
	}

	return  nil
}


func (r *userRepository) UpdateUserTOTP(user models.User, totpSecret string, totpEnabled bool) error {
	if err := r.DB.Model(&user).Updates(map[string] interface{}{"totp_secret": totpSecret, "totp_enabled": totpEnabled}).Error; err != nil {
		log.Print("Could not update user TOTP")
		return apperrors.NewBadRequest("Could not update user TOTP")
	}

	return nil
}


func (r *userRepository) UpdateUser(user models.User) error {
	userId := user.ID
	foundUser, _ := r.GetUserById(int(userId))

	if foundUser == nil {
		log.Printf("Could not find user with ID: %d", int(userId))
		return apperrors.NewBadRequest("Could not find user with provided ID")
	}

	updatedDetails := map[string] interface{}{}

	if user.UserName != "" {
		updatedDetails["UserName"] = user.UserName
	}

	if user.Email != "" {
		updatedDetails["Email"] = user.Email
	}

	if err := r.DB.Model(&foundUser).Updates(updatedDetails).Error; err != nil {
		log.Print("Could not update user")
		return apperrors.NewBadRequest("Could not update user")
	}

	return nil
}


func (r *userRepository) UpdatePassword(userId int, password string) error {
	foundUser, _ := r.GetUserById(userId)

	if foundUser == nil {
		log.Printf("Could not find user with ID: %d", userId)
		return apperrors.NewBadRequest("Could not find user with provided ID")
	}

	if err := r.DB.Model(&foundUser).Updates(models.User{Password: password}).Error; err != nil {
		log.Print("Could not update user password")
		return apperrors.NewBadRequest("Could not update user password")
	}

	return nil
}