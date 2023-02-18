package main

import (
	"testing"
	"time"
	// "github.com/go-resty/resty/v2"
	// "errors"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	// "time"

	// "bytes"

	"github.com/stretchr/testify/assert"
)
 const (
	serverAddress = "http://127.0.0.1:8080"
	Name string = "authentication"
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
		name: "Регистрация пользователя корректный запрос",
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
		name: "Проверка авторизации пользователя",
		url: "/api/user/test",
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
		name: "Аутентификация пользователя корректный запрос",
		url: "/api/user/login",
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
		name: "Проверка авторизации пользователя",
		url: "/api/user/test",
		method: "POST",
		StatusCode: 200,
		contentType: "application/json",
		result: "",
		User: User{
			Login: "Vasya",
			Password: "Vasya2023",
		},
	},
	// {
	// 	name: "Регистрация пользователя логин занят",
	// 	url: "/api/user/register",
	// 	method: "POST",
	// 	StatusCode: 409,
	// 	contentType: "application/json",
	// 	result: "",
	// 	User: User{
	// 		Login: "Vasya",
	// 		Password: "Vasya2022",
	// 	},
	// },
	// {
	// 	name: "**Регистрация пользователя не верный формат**",
	// 	url: "/api/user/register",
	// 	method: "POST",
	// 	StatusCode: 400,
	// 	contentType: "application/json",
	// 	result: "",
	// 	User: User{
	// 		Login: "Vasya",
	// 		Password: "Vasya2023",
	// },
	// },
}
// create HTTP client without redirects support
// var errRedirectBlocked = errors.New("HTTP redirect blocked")
// var redirPolicy = resty.RedirectPolicyFunc(func(_ *http.Request, _ []*http.Request) error {
// 	return errRedirectBlocked
// })
// var httpc = resty.New().
// 	SetBaseURL(serverAddress).
// 	SetRedirectPolicy(redirPolicy)
var client = http.Client{}
var urlReal = url.URL{
    Scheme:     "http",  
    Host:       "localhost",
    Path:       "/",
}
var cookies = []*http.Cookie{
    {
		Name:   "some_token",
    	Value:  "some_token",
    	MaxAge: 300,
	},
}
func TestRegister(t *testing.T) {
	fmt.Println("-----------------start-------------TestRegister-----------------------------")
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
	} else {
		client.Jar = jar
	}
	// куки можно устанавливать клиенту для всех запросов по определённому URL
	client.Jar.SetCookies(&urlReal, cookies)
	fmt.Println("client.Jar: ", client.Jar)
	// а можно добавлять к конкретному запросу
	// request.AddCookie(cookie)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, err := json.Marshal(&tt.User)
			if err != nil {
				fmt.Println(err)
			}
 			bodyReader := bytes.NewReader(bodyBytes)
			req, err := http.NewRequest(http.MethodPost, serverAddress + tt.url, bodyReader)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			req.Header.Add("Accept", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			// fmt.Printf("resp Cookies: %#v\n", resp.Cookies())
			// fmt.Println("client.Jar: ", client.Jar)
			defer resp.Body.Close()
			_, err = io.Copy(io.Discard, resp.Body)
			if err != nil {
				fmt.Println(err)
			} 
			// for _, c := range resp.Cookies(){
			// 	fmt.Println("c: ", c)
			// }
			// tt.Skip()
			// req := httpc.R().
			// SetHeader("Accept-Encoding", "gzip").
			// SetHeader("Content-Type", "application/json")
			// var result Metrics
			// resp, err := req.
			// SetBody(&User{
			// 	Login:    tt.User.Login,
			// 	Password: tt.User.Password,}).
			// // SetResult(&result).
			// Post(tt.url)


			assert.NoError(t, err)
			assert.Equal(t, tt.StatusCode, resp.StatusCode)
			// logg.Print(tt.name)
			// if tt.name != "default" {
			// 	assert.Equal(t, tt.contentType, resp.Header.Get("Content-Type"))
			// }
			// if tt.method == http.MethodGet {
				// assert.Equal(t, tt.result, string(body))
			// }
			time.Sleep(1* time.Second)
		})
	}

}