package models


type Group struct {
	Base
	Name 		string 		`json:"name" gorm:"unique"`
	Description string 		`json:"description"`
	Users 		[]User 		`gorm:"many2many:user_groups"`
	Messages    []Message   `gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE`
}


type UserGroup struct {
	Base
	UserId     uint    `gorm:"primaryKey" json:"userId"`
	GroupId    uint	   `gorm:"primaryKey" json:"groupId"`
	Banned     bool    `json:"ban" gorm:"type:bool;default:false"`
}


type IGroupRepository interface {
	CreateGroup(group *Group) (*Group, error)
	UpdateGroup(group Group) error
	GetGroupById(id int) (*Group, error)
	GetGroupsByUserId(userId, limit, page int) ([]Group, error)
	DeleteGroupById(id int) error
	BanUserFromGroup(groupId, userId int) error
	UnBanUserFromGroup(groupId, userId int) error
}


type IGroupService interface {
	CreateGroup(group *Group) (*Group, error)
	UpdateGroup(group Group) error
	GetGroupById(id int) (*Group, error)
	GetGroupsByUserId(userId, limit, page int) ([]Group, error)
	DeleteGroupById(id int) error
	BanUserFromGroup(groupId, userId int) error
	UnBanUserFromGroup(groupId, userId int) error
}