package api

import "github.com/dgrijalva/jwt-go"

const (
	jwtSigningKey = "10fa4f27-6a69-45c1-9a88-dfcecdbdc3d8"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS512
)
