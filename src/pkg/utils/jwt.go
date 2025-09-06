package utils

import (
    "time"

    jwt_token "github.com/golang-jwt/jwt/v5"
)

func GenerateSuperAdminToken(superadminEmail, superadminName, superadminId string) (string, error) {
    claims := jwt_token.MapClaims{
        "super_admin_id":    superadminId,
        "super_admin_email": superadminEmail,
        "super_admin_name":  superadminName,
        "exp":               time.Now().Add(time.Hour * 24 * 7).Unix(),
    }
    return jwt_token.NewWithClaims(jwt_token.SigningMethodHS256, claims).SignedString([]byte("srujankm1234"))
}
