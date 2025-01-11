package services

import "darkoo/models"


type groupService struct {
	groupRepository models.IGroupRepository
}


func NewGroupService(GroupRepository models.IGroupRepository) models.IGroupService {
	return &groupService{
		groupRepository: GroupRepository,
	}
}



func (s *groupService) CreateGroup(group *models.Group) (*models.Group, error) {
	return s.groupRepository.CreateGroup(group)
}


func (s *groupService) UpdateGroup(group models.Group)error {
	return s.groupRepository.UpdateGroup(group)
}


func (s *groupService) GetGroupById(id int) (*models.Group, error) {
	return s.groupRepository.GetGroupById(id)
}


func (s *groupService) GetGroupsByUserId(userId, limit, page int) ([]models.Group, error) {
	return s.groupRepository.GetGroupsByUserId(userId, limit, page)
}



func (s *groupService) DeleteGroupById(id int) error {
	return s.groupRepository.DeleteGroupById(id)
}


func (s *groupService) BanUserFromGroup(groupId, userId int) error {
	return s.groupRepository.BanUserFromGroup(groupId, userId)
}


func (s *groupService) UnBanUserFromGroup(groupId, userId int) error {
	return s.groupRepository.UnBanUserFromGroup(groupId, userId)
}