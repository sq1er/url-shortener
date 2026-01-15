package jwt_test

import (
	"testing"

	"github.com/sq1er/url-shortener/pkg/jwt"
)

func TestJWTCreate(t *testing.T) {
	const email = "a@a.ru"
	jwtService := jwt.NewJWT("Of4G[KEy$kvN@#]Ft[_XSp_i,$)+IycB>7}Wjy>#Vd9")

	token, err := jwtService.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("token is invalid")
	}
	if data.Email != email {
		t.Fatalf("email %s not equal %s", data.Email, email)
	}
}
