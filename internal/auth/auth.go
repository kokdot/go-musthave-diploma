package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"errors"
	"fmt"
)
const Name string = "authentication"
const (
	sizeOfKeyString = 32
)
var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

func SetCookie(value []byte, name string, seckretKey []byte) *http.Cookie {
	mac := hmac.New(sha256.New, seckretKey)
	mac.Write(value)//value is toking
	signature := mac.Sum(nil)

	var value1 = make([]byte, 0)
	fmt.Println("signature: ", signature)
	fmt.Println("value: ", value)
	fmt.Println("seckretKey: ", seckretKey)
	fmt.Println("bytes of name: ", []byte(name))
	fmt.Printf("len of string value: %v\n", len(value))
	value1 = append(value1, signature...)
	fmt.Println("value1: ", value1)
	value1 = append(value1, value...)
	fmt.Println("value1: ", value1)
	value1 = append(value1, []byte(name)...)
	fmt.Println("value1: ", value1)
	fmt.Printf("len of []bytes value1: %v\n", len(value1))
	valueStr := base64.URLEncoding.EncodeToString(value1)
	fmt.Printf("len of string valueStr: %v\n", len(valueStr))
	cookie := http.Cookie{
		Name: Name,
		Value: valueStr,
		Path: "/",
		MaxAge: 300,
		// SameSite: http.SameSiteNoneMode,
		// HttpOnly: true,
		// Secure: false,
		// SameSite: http.SameSiteLaxMode,
	}

	return &cookie
}
func ValidCookie(r *http.Request, secretKey []byte) (string, bool, error) {
	cookie, err := r.Cookie(Name)
	if err != nil {
		return "", false, err
	}
	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", false, ErrInvalidValue
	}
	if len(value) < sha256.Size {
		return "", false, ErrInvalidValue
	}

	fmt.Println("len of value is: ", len(value))
	signature := value[:sha256.Size]
	fmt.Printf("signature is: %v, size is: %v\n", signature, len(signature))
	valueTaking := value[sha256.Size:sha256.Size + sizeOfKeyString]
	fmt.Printf("valueTaking is: %v, size is: %v\n", valueTaking, len(valueTaking))
	name := value[sha256.Size + sizeOfKeyString:]
	fmt.Printf("name is: %v, size is: %v\n", name, len(name))
	nameStr := string(name)
	fmt.Printf("nameStr is: %v, size is: %v\n", nameStr, len(nameStr))
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(valueTaking))
	expectedSignature := mac.Sum(nil)
	if !hmac.Equal(signature, expectedSignature) {
		return "", false, ErrInvalidValue
	}
	return nameStr, true, nil
}