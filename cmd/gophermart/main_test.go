package main
 import (
	"testing"
 	"github.com/go-resty/resty/v2"
	"errors"
	"net/http"
	// "time"

	// "bytes"

	"github.com/stretchr/testify/assert"

 )
 const (
	serverAddress = "http://127.0.0.1"
 )

 type User struct {
	 Login string   `json:"login"`
	 Password string   `json:"password"`
 }
// type want struct {
// 	StatusCode  int
// 	contentType string
// 	result      string
// }
var tests = []struct {
	User User
	name   string
	url    string
	method string
	StatusCode  int
	contentType string
	result      string
}{
	{
		name: "**Регистрация пользователя корректный запрос**",
		url: "/api/user/register",
		method: "POST",
		StatusCode: 200,
		contentType: "application/json",
		result: "",
		User: User{
			Login: "Vasya",
			Password: "Vasya2023",
		},
	},
	{
		name: "**Регистрация пользователя логин занят**",
		url: "/api/user/register",
		method: "POST",
		StatusCode: 409,
		contentType: "application/json",
		result: "",
		User: User{
			Login: "Vasya",
			Password: "Vasya2022",
		},
	},
	{
		name: "**Регистрация пользователя не верный формат**",
		url: "/api/user/register",
		method: "POST",
		StatusCode: 400,
		contentType: "application/json",
		result: "",
		User: User{
			Login: "Vasya",
			Password: "Vasya2023",
	},
	},
}
// create HTTP client without redirects support
var errRedirectBlocked = errors.New("HTTP redirect blocked")
var redirPolicy = resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
	return errRedirectBlocked
})
var httpc = resty.New().
	SetBaseURL(serverAddress).
	SetRedirectPolicy(redirPolicy)

func TestRegister(t *testing.T) {
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// tt.Skip()
			req := httpc.R().
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("Content-Type", "application/json")
			// var result Metrics
			resp, err := req.
			SetBody(&User{
				Login:    tt.User.Login,
				Password: tt.User.Password,}).
			// SetResult(&result).
			Post("value/")


			assert.NoError(t, err)
			assert.Equal(t, tt.StatusCode, resp.StatusCode())
			// logg.Print(tt.name)
			if tt.name != "default" {
				assert.Equal(t, tt.contentType, resp.Header().Get("Content-Type"))
			}
			// if tt.method == http.MethodGet {
				// assert.Equal(t, tt.result, string(body))
			// }
		})
	}

}