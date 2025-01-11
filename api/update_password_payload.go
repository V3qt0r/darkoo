package api

import "github.com/go-ozzo/ozzo-validation"

type UpdatePasswordPayload struct {
	Password 		string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}


func (p UpdatePasswordPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Password, validation.Required),
		validation.Field(&p.ConfirmPassword, validation.Required),	
	)
}


type ConfirmPasswordPayload struct {
	Password	string	`json:"password"`
}


func (p ConfirmPasswordPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Password, validation.Required),
	)
}