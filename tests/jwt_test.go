package tests

import (
	"testing"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
)

func TestJWT_Multi_Success(t *testing.T) {
	tokenString, err := utils.GenerateJWT(1)
	if err != nil {
		t.Fatal(err)
	}

	claims, err := utils.GetClaimsJWT(tokenString)
	if err != nil {
		t.Fatal(err)
	}

	if claims.Subject != "1" {
		t.Fail()
	}
}
