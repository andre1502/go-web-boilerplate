package token

import (
	cerror "boilerplate/utils/error"
	"boilerplate/utils/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtMapClaims struct {
	UserId    uint64 `json:"user_id"`
	AccountId uint64 `json:"account_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uint64, tokenHourLifeSpan int, secretKey []byte) (string, error) {
	now := time.Now()

	claims := JwtMapClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Unix(now.Unix(), 0)),
			ExpiresAt: jwt.NewNumericDate(time.Unix(now.Add(time.Hour*time.Duration(tokenHourLifeSpan)).Unix(), 0)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := token.SignedString(secretKey)

	if err != nil {
		logger.Sugar.Error(err)
		return "", cerror.Fail(cerror.FuncName(), "failed_generate_token", nil, err)
	}

	return res, nil
}

func TokenValid(tokenString string, secretKey []byte) (*JwtMapClaims, error) {
	claims := &JwtMapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Sugar.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, cerror.Fail(cerror.FuncName(), "invalid_token", nil, nil)
		}

		return secretKey, nil
	})

	if (err != nil) || (!token.Valid) {
		if err != nil {
			logger.Sugar.Error(err)
		}

		return nil, cerror.Fail(cerror.FuncName(), "invalid_token", nil, err)
	}

	return claims, nil
}
