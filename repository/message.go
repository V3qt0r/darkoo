package repository

import (
	"darkoo/models"
	"darkoo/apperrors"
	"darkoo/websocket"

	"log"

	"gorm.io/gorm"
)


type messageRepository struct {
	DB *gorm.DB
}


func NewMessageRepository(db *gorm.DB) models.IMessageRepository {
	return &messageRepository{ DB: db,}
}


func (r *messageRepository) SendMessage(message *models.Message) (*models.Message, error) {
	group := &models.Group{}
	user := &models.User{}
	userGroup := &models.UserGroup{}

	if err := r.DB.Where("id = ?", message.GroupId).First(&group).Error; err != nil {
		log.Printf("Could not find group with ID: %v\n", message.GroupId)
		return nil, apperrors.NewBadRequest("Could not find group with provided ID")
	}

	if err := r.DB.Where("id = ?", message.UserId).First(&user).Error; err != nil {
		log.Printf("Could not find user with ID: %v\n", message.UserId)
		return nil, apperrors.NewBadRequest("Could not find user with provided ID")
	}

	if err := r.DB.Where("user_id = ? AND group_id = ? AND banned = ?", user.ID, group.ID, false).First(&userGroup).Error; err != nil {
		log.Print("Cannot send messages. User is not a member of group")
		return nil, apperrors.NewBadRequest("Cannot send messages. User is not a member of group")
	}

	if message.ContentType != "text" && message.AttachmentUrl == nil {
		log.Print("Attachment URL is required for this type of message")
		return nil, apperrors.NewBadRequest("Attachment URL is required for this type of message")
	}

	if err := r.DB.Create(&message).Error; err != nil {
		log.Print("Could not send message")
		return nil, apperrors.NewInternal()
	}

	websocket.Hub.Broadcast <- message

	return message, nil
}


func (r *messageRepository) GetMessagesInGroup(groupId, limit, page int) ([]models.Message, error) {
	var messages []models.Message

	if err := r.DB.Where("group_id = ?", groupId).Find(&messages).Error; err != nil {
		log.Print("Could not get messages")
		return messages, apperrors.NewInternal()
	}

	return messages, nil
}


func (r *messageRepository) GetMessageById(id int) (*models.Message, error) {
	message := &models.Message{}

	if err := r.DB.Where("id = ?", id).First(&message).Error; err != nil {
		log.Printf("Could not find message with ID: %d\n", id)
		return nil, apperrors.NewBadRequest("Could not find message with provided ID")
	}

	return message, nil
}


func (r *messageRepository) UpdateMessage(message models.Message) error {
	id := message.ID
	foundMessage, _ := r.GetMessageById(int(id))
	userGroup := &models.UserGroup{}

	if foundMessage == nil {
		log.Printf("Could not find message with ID: %d\n", int(id))
		return apperrors.NewBadRequest("Could not find message with ID")
	}

	if err := r.DB.Where("user_id = ? AND group_id = ? AND banned = ?", foundMessage.UserId, foundMessage.GroupId, false).First(&userGroup).Error; err != nil {
		log.Print("Unable to edit message")
		return apperrors.NewBadRequest("Unable to edit message")
	}


	updatedContent := map[string] interface{}{}

	if message.ContentType == "text"{
		if message.Content != "" {
			updatedContent["Content"] = message.Content
		}
	} else {
		log.Print("You can only edit text messages")
		return apperrors.NewBadRequest("You can only edit text messages")
	}

	if err := r.DB.Model(&foundMessage).Updates(updatedContent).Error; err != nil {
		log.Print("Could not edit message")
		return apperrors.NewInternal()
	}

	return nil
}


func (r *messageRepository) DeleteMessage(id, userId, groupId int) error {
	userGroup := &models.UserGroup{}

	if err := r.DB.Where("user_id = ? AND group_id = ? AND banned = ?", userId, groupId, false).First(&userGroup).Error; err != nil {
		log.Print("Unauthorized to delete user message from this group")
		return apperrors.NewBadRequest("Unauthorized to delete user message from this group")
	}


	if err := r.DB.Where("id = ? AND user_id = ? AND group_id = ?", id, userId, groupId).Delete(&models.Message{}).Error; err != nil {
		log.Print("Error deleting message")
		return apperrors.NewInternal()
	}

	return nil
}


func (r *messageRepository) GetUserMessagesInGroup(userId, groupId, limit, page int) ([]models.Message, error){
	var messages []models.Message


	if err := r.DB.Where("user_id = ? AND group_id = ?", userId, groupId).Find(&messages).Error; err != nil {
		log.Print("Could not get user messages from group")
		return messages, apperrors.NewBadRequest("Could not get user messages from group")
	}

	return messages, nil
}