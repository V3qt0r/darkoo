package api

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type InitOnetimePasswordPayload struct {
	Email 		string 	`json:"email"`
}


func (r InitOnetimePasswordPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
	)
}