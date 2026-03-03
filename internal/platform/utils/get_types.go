package utils

import "github.com/golang-jwt/jwt/v5"

func GetIntClaim(claims jwt.MapClaims, key string) int64 {
	val, ok := claims[key].(float64)
	if ok {
		return int64(val)
	}
	return -1
}