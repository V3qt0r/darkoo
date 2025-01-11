package models


type Message struct {
	Base
	Content  		string 		`json:"content"`
	ContentType     string      `gorm:"not null" json"contentType"`
	AttachmentUrl  *string		`gorm:"type:text" json"attachmentUrl`
	GroupId         uint 		`json:"_"`
	Group			Group       `gorm:"foreignKey:GroupId; constraint:OnUpdate:CASCADE, OnDelete:CASCADE" json"-"`
	UserId          uint        `json:"-"`
	User  			User 		`gorm:"foreignKey:UserId; constraint:OnUpdate:CASCADE, OnDelete:SET NULL"`
}


type IMessageRepository interface {
	SendMessage(message *Message) (*Message, error)
	GetMessagesInGroup(groupId, limit, page int) ([]Message, error)
	GetUserMessagesInGroup(userId, groupId, limit, page int) ([]Message, error)
	DeleteMessage(id, userId, groupId int) error
	UpdateMessage(message Message) error
	GetMessageById(id int) (*Message, error)
}


type IMessageService interface {
	SendMessage(message *Message) (*Message, error)
	GetMessagesInGroup(groupId, limit, page int) ([]Message, error)
	GetUserMessagesInGroup(userId, groupId, limit, page int) ([]Message, error)
	DeleteMessage(id, userId, groupId int) error
	UpdateMessage(message Message) error
	GetMessageById(id int) (*Message, error)
}