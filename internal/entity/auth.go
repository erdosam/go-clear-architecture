package entity

import "github.com/dgrijalva/jwt-go"

type JwtPayload struct {
	AccountId         string `json:"account_id"`
	Username          string `json:"username"`
	Name              string `json:"name"`
	MobileCallingCode string `json:"mobile_calling_code"`
	MobileNo          string `json:"mobile_no"`
	Email             string `json:"email"`
	Gender            string `json:"gender"`
	ClientKey         string `json:"client_key"`
}

type CustomClaims struct {
	jwt.StandardClaims
	JwtPayload
}
