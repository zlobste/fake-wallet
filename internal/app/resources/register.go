package resources

import (
	"encoding/json"
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegister(r *http.Request) (RegisterRequest, error) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal register request")
	}

	return request, validateRegister(request)
}

func validateRegister(r RegisterRequest) error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Name, validation.Required, validation.Length(3, 50),
			validation.Match(regexp.MustCompile("^[A-Za-z]+")),
		),
		validation.Field(
			&r.Surname, validation.Required, validation.Length(3, 50),
			validation.Match(regexp.MustCompile("^[A-Za-z]+")),
		),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 50)),
	)
}
