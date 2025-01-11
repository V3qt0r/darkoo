package api

import (
	"strings"
	"darkoo/models"

	validation "github.com/go-ozzo/ozzo-validation"
	// "github.com/go-ozzo/ozzo-validation/v4/is"
)


type CreateGroupPayload struct {
	Name 		string 	`json:"name"`
	Description string	`json:"description"`
}


func (p CreateGroupPayload) Sanitize() {
	p.Name = strings.TrimSpace(p.Name)
}


func (p CreateGroupPayload) Validate() error {
	return validation.ValidateStruct(&p, 
		validation.Field(&p.Name, validation.Required, validation.Length(3, 30)),
		validation.Field(&p.Description, validation.Length(3, 160)),
	)
}


type UpdateGroupPayload struct {
	Name 		string 		`json:"name"`
	Description string 		`json:"description"`
}


func (p UpdateGroupPayload) Validate() error {
	if p.Name != "" {
		return validation.Validate(p.Name, validation.Required, validation.Length(3, 30))
	}

	if p.Description != "" {
		return validation.Validate(p.Description, validation.Length(3, 160))
	}
	return nil
}


func (p UpdateGroupPayload) Sanitize() {
	p.Name = strings.TrimSpace(p.Name)
}


func (p UpdateGroupPayload) ToEntity() models.Group {
	var group models.Group

	if p.Name != "" {
		group.Name = p.Name
	}

	if p.Description != "" {
		group.Description = p.Description
	}

	return group
}