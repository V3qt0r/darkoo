package api

import (
	"strings"
	"darkoo/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)


type RegisterPayload struct {
	Email 		string 	`json:"email"`
	Password 	string	`json:"password"`
	IsAbove18 	bool	`json:"isAbove18"`
}


func (p RegisterPayload) Sanitize() {
	p.Email = strings.TrimSpace(p.Email)
	p.Email = strings.ToLower(p.Email)
}


func (p RegisterPayload) Validate() error {
	return validation.ValidateStruct(&p, 
		validation.Field(&p.Email, validation.Required, is.EmailFormat),
		validation.Field(&p.Password, validation.Required, validation.Length(6, 150),),		
		validation.Field(&p.IsAbove18, validation.Required, validation.In(true).Error("You must be above 18 to register")),
	)
}


type UpdateUserPayload struct {
	Email 		string		`json:"email"`
	UserName	string		`json:"userName"`
}


func (p UpdateUserPayload) Validate() error {
	if p.Email != "" {
		return validation.Validate(p.Email, validation.Required, is.EmailFormat)
	}

	if p.UserName != "" {
		return validation.Validate(p.UserName, validation.Required)
	}

	return nil
}


func (p UpdateUserPayload) Sanitize() {
	p.Email = strings.TrimSpace(p.Email)
	p.Email = strings.ToLower(p.Email)
	p.UserName = strings.TrimSpace(p.UserName)
}


func (p UpdateUserPayload) ToEntity() models.User {
	var user models.User

	if p.UserName != "" {
		user.UserName = p.UserName
	}

	if p.Email != "" {
		user.Email = p.Email
	}

	return user
}