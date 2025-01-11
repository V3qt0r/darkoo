package repository

import (
	"darkoo/models"
	"darkoo/apperrors"

	"log"

	"gorm.io/gorm"
)


type groupRepository struct {
	DB *gorm.DB
}


func NewGroupRepository(db *gorm.DB) models.IGroupRepository {
	return &groupRepository{ DB: db, }
}


func (r *groupRepository) CreateGroup(group *models.Group) (*models.Group, error) {
	if err := r.DB.Create(&group).Error; err != nil {
		log.Print("Could not create group")
		return nil, apperrors.NewBadRequest("Could not create group")
	}

	return group, nil
}


func (r *groupRepository) GetGroupById(id int) (*models.Group, error) {
	group := &models.Group{}

	if err := r.DB.Where("id = ?", id).First(&group).Error; err != nil {
		log.Print("Could not get group with ID: %d\n", id)
		return nil, apperrors.NewBadRequest("Could not get group with provided ID")
	}

	return group, nil
}


func (r *groupRepository) UpdateGroup(group models.Group) error {
	groupId := group.ID
	foundGroup, _ := r.GetGroupById(int(groupId))

	if foundGroup == nil {
		log.Printf("Could not find group with ID: %d\n", groupId)
		return apperrors.NewBadRequest("Could not find group with provided ID")
	}

	updatedDetails := map[string] interface{}{}
	if group.Name != "" {
		updatedDetails["Name"] = group.Name
	}

	if group.Description != "" {
		updatedDetails["Description"] = group.Description
	}

	if err := r.DB.Model(&foundGroup).Updates(updatedDetails).Error; err != nil {
		log.Print("Could not update group")
		return apperrors.NewInternal()
	}

	return nil
}


func (r *groupRepository) GetGroupsByUserId(id, limit, page int) ([]models.Group, error) {
	user := &models.User{}
	var groups []models.Group
	
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		log.Print("Could not fin user with ID: %d\n", id)
		return groups, apperrors.NewBadRequest("Could not find user with provided ID")
	}

	if err := r.DB.Model(&user).Association("Groups").Find(&groups).Error; err != nil {
		log.Print("Could not find user_groups association")
		return groups, apperrors.NewBadRequest("Could not find user_groups association")
	}

	return groups, nil
}


func (r *groupRepository) DeleteGroupById(id int) error {
	if err := r.DB.Where("id = ?", id).Delete(&models.Group{}).Error; err != nil {
		log.Printf("Could not delete group with ID: %d\n", id)
		return apperrors.NewBadRequest("Could not delete group")
	}

	return nil
}


func (r *groupRepository) BanUserFromGroup(groupId, userId int) error { 
	group := &models.Group{}
	user := &models.User{}
	userGroup := &models.UserGroup{}

	if err := r.DB.Where("id = ?", groupId).First(&group).Error; err != nil {
		log.Printf("Group with ID %d not found\n", groupId)
		return apperrors.NewBadRequest("Group with provided ID not found")
	}

	if err := r.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		log.Printf("User with ID %d not found\n", userId)
		return apperrors.NewBadRequest("User with provided ID not found")
	}

	if err := r.DB.Where("user_id = ? AND group_id = ? AND banned = ?", userId, groupId, false).First(&userGroup).Error; err != nil {
		log.Print("User is already banned or does not belong to this group")
		return apperrors.NewBadRequest("User is already banned or does not belong to this group")
	}
	
	if err := r.DB.Model(&userGroup).Updates(models.UserGroup{ Banned : true}).Error; err != nil {
		log.Print("Could not updated user banned status")
		return apperrors.NewBadRequest("Could not update user banned status")
	}

	return nil
}


func (r *groupRepository) UnBanUserFromGroup(groupId, userId int) error { 
	group := &models.Group{}
	user := &models.User{}
	userGroup := &models.UserGroup{}

	if err := r.DB.Where("id = ?", groupId).First(&group).Error; err != nil {
		log.Printf("Group with ID %d not found\n", groupId)
		return apperrors.NewBadRequest("Group with provided ID not found")
	}

	if err := r.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		log.Printf("User with ID %d not found\n", userId)
		return apperrors.NewBadRequest("User with provided ID not found")
	}

	if err := r.DB.Where("user_id = ? AND group_id = ? AND banned = ?", userId, groupId, true).First(&userGroup).Error; err != nil {
		log.Print("User is not banned or does not belong to this group")
		return apperrors.NewBadRequest("User is not banned or does not belong to this group")
	}

	if err := r.DB.Model(&userGroup).Updates(models.UserGroup{ Banned : false}).Error; err != nil {
		log.Print("Could not updated user banned status")
		return apperrors.NewBadRequest("Could not update user banned status")
	}

	return nil
}