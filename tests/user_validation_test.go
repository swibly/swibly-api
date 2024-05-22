package tests

import (
	"strings"
	"testing"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
)

func TestUserModel_Register_Success(t *testing.T) {
	userModel := &dto.UserRegister{
		FirstName: "Test",
		LastName:  "Subject",
		Username:  "testsubject",
		Email:     "test.subject@example.com",
		Password:  "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) > 0 {
		var errorMessages strings.Builder
		for _, err := range errs {
			errorMessages.WriteString(err.Error() + "\n")
		}

		t.Errorf("Validation failed: %s", errorMessages.String())
	}
}

func TestUserModel_Login_Success(t *testing.T) {
	userModel := &dto.UserLogin{
		Username: "testsubject",
		Email:    "test.subject@example.com",
		Password: "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) > 0 {
		var errorMessages strings.Builder
		for _, err := range errs {
			errorMessages.WriteString(err.Error() + "\n")
		}

		t.Errorf("Validation failed: %s", errorMessages.String())
	}
}

func TestUserModel_Register_Failure_EmptyUsername(t *testing.T) {
	userModel := &dto.UserRegister{
		FirstName: "Test",
		LastName:  "Subject",
		Username:  "",
		Email:     "test.subject@example.com",
		Password:  "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}

func TestUserModel_Register_Failure_LongUsername(t *testing.T) {
	userModel := &dto.UserRegister{
		FirstName: "Test",
		LastName:  "Subject",
		Username:  "123456781234567812345678123456781234567899999123456789",
		Email:     "test.subject@example.com",
		Password:  "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}

func TestUserModel_Register_Failure_InvalidEmail(t *testing.T) {
	userModel := &dto.UserRegister{
		FirstName: "Test",
		LastName:  "Subject",
		Username:  "testsubject",
		Email:     "invalid.email", // Invalid email format
		Password:  "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}

func TestUserModel_Register_Failure_WeakPassword(t *testing.T) {
	userModel := &dto.UserRegister{
		FirstName: "Test",
		LastName:  "Subject",
		Username:  "testsubject",
		Email:     "test.subject@example.com",
		Password:  "weakpassword", // Does not meet the strength criteria
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}

func TestUserModel_Login_Failure_LongUsername(t *testing.T) {
	userModel := &dto.UserLogin{
		Username: "123456781234567812345678123456781234567899999123456789",
		Email:    "test.subject@example.com",
		Password: "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}

func TestUserModel_Login_Failure_InvalidEmail(t *testing.T) {
	userModel := &dto.UserLogin{
		Username: "testsubject",
		Email:    "invalid.email", // Invalid email format
		Password: "t&STngF3@tur3",
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}

func TestUserModel_Login_Failure_WeakPassword(t *testing.T) {
	userModel := &dto.UserLogin{
		Username: "testsubject",
		Email:    "test.subject@example.com",
		Password: "weakpassword", // Does not meet the strength criteria
	}

	errs := utils.ValidateStruct(userModel) // returns validator.ValidationErrors

	if len(errs) == 0 {
		t.Fail()
	}
}
