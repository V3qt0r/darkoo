package api

import (
	"strings"
	"darkoo/models"

	validation "github.com/go-ozzo/ozzo-validation"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)


type SendMessagePayload struct {
	Content 		string 		`json:"content"`
	ContentType 	string 		`json:"contentType"`
	AttachmentUrl  *string		`json:"attachmentUrl"`
}


func (p SendMessagePayload) Sanitize() {
	p.ContentType = strings.TrimSpace(p.ContentType)
}


func (p SendMessagePayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Content, validation.Length(1, 500)),
		validation.Field(&p.ContentType, validation.Required),
	)
}


type UpdateMessagePayload struct {
	Content string `json:"content"`
}

func (p UpdateMessagePayload) Validate() error {
	if p.Content != "" {
		return validation.Validate(p.Content, validation.Required, validation.Length(1, 500))
	}
	return nil
}


func (p UpdateMessagePayload) ToEntity() models.Message {
	var message models.Message

	if p.Content != "" {
		message.Content = p.Content
	}

	return message
}