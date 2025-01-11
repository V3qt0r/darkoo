package services

import "darkoo/models"


type messageService struct {
	messageRepository models.IMessageRepository
}


func NewMessageService(messageRepository models.IMessageRepository) models.IMessageService {
	return &messageService {
		messageRepository : messageRepository,
	}
}


func (s *messageService) SendMessage(message *models.Message) (*models.Message, error) {
	return s.messageRepository.SendMessage(message)
}


func (s *messageService) GetMessagesInGroup(groupId, limit, page int) ([]models.Message, error) {
	return s.messageRepository.GetMessagesInGroup(groupId, limit, page)
}


func (s *messageService) GetUserMessagesInGroup(userId, groupId, limit, page int) ([]models.Message, error) {
	return s.messageRepository.GetUserMessagesInGroup(userId, groupId, limit, page)
}


func (s *messageService) DeleteMessage(id, userId, groupId int) error {
	return s.messageRepository.DeleteMessage(id, userId, groupId)
}


func (s *messageService) UpdateMessage(message models.Message) error {
	return s.messageRepository.UpdateMessage(message)
}


func (s *messageService) GetMessageById(id int) (*models.Message, error) {
	return s.messageRepository.GetMessageById(id)
}