package toking

import (
	"fmt"
	"crypto/sha256"
	"crypto/rand"
)
func RandBytesKeyString(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return []byte{}, err
    }
    return b, nil
}
func Sha256 (b []byte) string {
	h := sha256.New()
    h.Write(b)
    dst := h.Sum(nil)

    result := fmt.Sprintf("%x", dst)
	return result
} 
func Toking() string{
	x, err := RandBytesKeyString(16)
		if err != nil {
        fmt.Println(err)
    	}
	return Sha256(x)
}	
 