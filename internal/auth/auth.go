package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)
const Name string = "authentication"

func SetCookie(value string, secretKey []byte) *http.Cookie {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(value))
	signature := mac.Sum(nil)
	value = string(signature) + value
	value = base64.URLEncoding.EncodeToString([]byte(value))
	cookie := http.Cookie{
		Name: Name,
		Value: value,
		Path: "/",
		MaxAge: 3600,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}

	return &cookie
}
func ValidCookie(r *http.Request, secretKey []byte) (bool, error) {
	cookie, err := r.Cookie(Name)
	if err != nil {
		return false, err
	}
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return false, ErrInvalidValue
	}
	if len(value) < sha256.Size {
		return false, ErrInvalidValue
	}
	signature := value[:sha256.Size]
	value = value[sha256.Size:]
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)
	if !hmac.Equal([]byte(signature), expectedSignature) {
		return false, ErrInvalidValue
	}
	return true, nil
}