package toking

import (
	"fmt"
	"crypto/sha256"
	"crypto/rand"
)
const (
	sizeOfKeyString = 32
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
func Toking() []byte{
	x, err := RandBytesKeyString(sizeOfKeyString)
		if err != nil {
        fmt.Println(err)
    	}
	return x
}	
 