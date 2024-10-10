package usecase

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"

	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service/repository"
	"github.com/swibly/swibly-api/pkg/sender"
	"github.com/swibly/swibly-api/translations"
)

type PasswordResetUseCase struct {
	prr repository.PasswordResetRepository

	uuc UserUseCase
}

func NewPasswordResetUseCase() PasswordResetUseCase {
	return PasswordResetUseCase{prr: repository.NewPasswordResetRepository(), uuc: NewUserUseCase()}
}

func (pruc *PasswordResetUseCase) Request(dict translations.Translation, email string) error {
	data, err := pruc.prr.Request(email)
	if err != nil {
		return err
	}

	user, _ := pruc.uuc.GetByEmail(email)

	tmpl, err := template.New("email").Parse(dict.PasswordResetEmailTemplate)
	if err != nil {
		return err
	}

	mapping := map[string]string{
		"user": user.FirstName,
		"url":  fmt.Sprintf("https://www.swibly.com.br/reset/%s", url.QueryEscape(data.Key.String())),
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, mapping)
	if err != nil {
		return err
	}

	sender.SMTPSender.Send(email, dict.PasswordResetEmailSubject, body.String())

	return nil
}

func (pruc *PasswordResetUseCase) Reset(key, newPassword string) error {
	return pruc.prr.Reset(key, newPassword)
}

func (pruc *PasswordResetUseCase) IsKeyValid(key string) (*dto.PasswordResetInfo, bool, error) {
	return pruc.prr.IsKeyValid(key)
}
