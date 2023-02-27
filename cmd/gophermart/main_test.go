package main

import (
	"strconv"
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

	"github.com/kokdot/go-musthave-diploma/internal/luna"

	// "time"

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
func NewJar() *cookiejar.Jar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return jar
}
var OrderNumbertests = []struct {
	name string
	number string
	StatusCode int
	contentType string
}{
	{
		name: "correct order",
		StatusCode: 202,
		contentType: "application/json",
	},
	{
		name: "reapit order",
		StatusCode: 200,
		contentType: "application/json",
	},
	// {
	// 	name: "another user order",
	// 	number: 100,
	// 	StatusCode: "",
	// },

}

func TestDownloadNumberOfOrder(t *testing.T) {
	fmt.Println("-----------------start-------------TestDownloadNumberOfOrder-----------------------------")
	var user = User{
			Login: "Vasya",
			Password: "Vasya2023",
		}
	var user1 = User{
			Login: "Misha",
			Password: "Misha2023",
		}
	var number = luna.GetOrderNumber()
	var sNumber = strconv.Itoa(number)
	client.Jar = NewJar()
	client.Jar.SetCookies(&urlReal, cookies)
	fmt.Println("client.Jar: ", client.Jar)
	bodyBytes, err := json.Marshal(&user)
	if err != nil {
		fmt.Println(err)
	}
	
	bodyReader := bytes.NewReader(bodyBytes)
	req, err := http.NewRequest(http.MethodPost, serverAddress + "/api/user/register", bodyReader)
	if err != nil {
		fmt.Println(err)
	}
	
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Println(err)
	} 
	resp.Body.Close()
	// -------------------------------------------Vasya-------------------------------------------------------------------
	
	for _, tt := range OrderNumbertests {
		t.Run(tt.name + "_Vasya", func(t *testing.T) {
			bodyBytes := []byte(sNumber)
			bodyReader := bytes.NewReader(bodyBytes)
			req, err = http.NewRequest(http.MethodPost, serverAddress + "/api/user/orders", bodyReader)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
			req.Header.Add("Accept", "application/json")
			resp2, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			_, err = io.Copy(io.Discard, resp2.Body)
			if err != nil {
				fmt.Println(err)
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.StatusCode, resp2.StatusCode)
			assert.Equal(t, tt.contentType, resp2.Header.Get("Content-Type"))
			resp2.Body.Close()
		})
	}
	for i := 0; i < 30; i++ {
		number = luna.GetOrderNumber()
		sNumber = strconv.Itoa(number)
		bodyBytes := []byte(sNumber)
		bodyReader := bytes.NewReader(bodyBytes)
		req, err = http.NewRequest(http.MethodPost, serverAddress + "/api/user/orders", bodyReader)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
		req.Header.Add("Accept", "application/json")
		resp2, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		_, err = io.Copy(io.Discard, resp2.Body)
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, 202, resp2.StatusCode)
		assert.Equal(t, "application/json", resp2.Header.Get("Content-Type"))
		resp2.Body.Close()
	}
	// -------------------------------------------Misha-------------------------------------------------------------------
	var number1 = luna.GetOrderNumber()
	var sNumber1 = strconv.Itoa(number1)
	client.Jar = NewJar()
	client.Jar.SetCookies(&urlReal, cookies)
	fmt.Println("client.Jar: ", client.Jar)
	bodyBytes1, err := json.Marshal(&user1)
	if err != nil {
		fmt.Println(err)
	}
	bodyReader1 := bytes.NewReader(bodyBytes1)
	req1, err := http.NewRequest(http.MethodPost, serverAddress + "/api/user/register", bodyReader1)
	if err != nil {
		fmt.Println(err)
	}
	req1.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req1.Header.Add("Accept", "application/json")
	resp1, err := client.Do(req1)
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(io.Discard, resp1.Body)
	if err != nil {
		fmt.Println(err)
	}
	resp1.Body.Close()
	
	for _, tt := range OrderNumbertests {
		t.Run(tt.name + "_Misha", func(t *testing.T) {
			bodyBytes := []byte(sNumber1)
			bodyReader := bytes.NewReader(bodyBytes)
			req, err = http.NewRequest(http.MethodPost, serverAddress + "/api/user/orders", bodyReader)
			if err != nil {
				fmt.Println(err)
			}
			req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
			req.Header.Add("Accept", "application/json")
			resp3, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			_, err = io.Copy(io.Discard, resp3.Body)
			if err != nil {
				fmt.Println(err)
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.StatusCode, resp3.StatusCode)
			assert.Equal(t, tt.contentType, resp3.Header.Get("Content-Type"))
			resp3.Body.Close()
		})
	}
	t.Run("incorrect order", func(t *testing.T) {
		bodyBytes := []byte("100")
		bodyReader := bytes.NewReader(bodyBytes)
		req, err = http.NewRequest(http.MethodPost, serverAddress + "/api/user/orders", bodyReader)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
		req.Header.Add("Accept", "application/json")
		resp4, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		_, err = io.Copy(io.Discard, resp4.Body)
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, 422, resp4.StatusCode)
		assert.Equal(t, "application/json", resp4.Header.Get("Content-Type"))
		resp4.Body.Close()
	})
	t.Run("order of enuther user", func(t *testing.T) {
		bodyBytes := []byte(sNumber)
		bodyReader := bytes.NewReader(bodyBytes)
		req, err = http.NewRequest(http.MethodPost, serverAddress + "/api/user/orders", bodyReader)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
		req.Header.Add("Accept", "application/json")
		resp5, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		_, err = io.Copy(io.Discard, resp5.Body)
		if err != nil {
			fmt.Println(err)
		}
		assert.NoError(t, err)
		assert.Equal(t, 409, resp5.StatusCode)
		assert.Equal(t, "application/json", resp5.Header.Get("Content-Type"))
		resp5.Body.Close()
	})
}

// --------------------------------------------------------------------------------
var Registertests = []struct {
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

func TestRegister(t *testing.T) {
	t.Skip()
	fmt.Println("-----------------start-------------TestRegister-----------------------------")
	client.Jar = NewJar()
	// куки можно устанавливать клиенту для всех запросов по определённому URL
	client.Jar.SetCookies(&urlReal, cookies)
	fmt.Println("client.Jar: ", client.Jar)
	// а можно добавлять к конкретному запросу
	// request.AddCookie(cookie)
	for _, tt := range Registertests {
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
// func TestOrderCreate(t *testing.T) {
// 	fmt.Println("-----------------start-------------TestOrderCreatew-----------------------------")
// 	jar, err := cookiejar.New(nil)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		client.Jar = jar
// 	}
// 	// куки можно устанавливать клиенту для всех запросов по определённому URL
// 	client.Jar.SetCookies(&urlReal, cookies)
// 	fmt.Println("client.Jar: ", client.Jar)
// 	// а можно добавлять к конкретному запросу
// 	// request.AddCookie(cookie)
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			bodyBytes, err := json.Marshal(&tt.User)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
//  			bodyReader := bytes.NewReader(bodyBytes)
// 			req, err := http.NewRequest(http.MethodPost, serverAddress + tt.url, bodyReader)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 			req.Header.Add("Accept", "application/json")
// 			resp, err := client.Do(req)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			// fmt.Printf("resp Cookies: %#v\n", resp.Cookies())
// 			// fmt.Println("client.Jar: ", client.Jar)
// 			defer resp.Body.Close()
// 			_, err = io.Copy(io.Discard, resp.Body)
// 			if err != nil {
// 				fmt.Println(err)
// 			} 
// 			// for _, c := range resp.Cookies(){
// 			// 	fmt.Println("c: ", c)
// 			// }
// 			// tt.Skip()
// 			// req := httpc.R().
// 			// SetHeader("Accept-Encoding", "gzip").
// 			// SetHeader("Content-Type", "application/json")
// 			// var result Metrics
// 			// resp, err := req.
// 			// SetBody(&User{
// 			// 	Login:    tt.User.Login,
// 			// 	Password: tt.User.Password,}).
// 			// // SetResult(&result).
// 			// Post(tt.url)


// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.StatusCode, resp.StatusCode)
// 			// logg.Print(tt.name)
// 			// if tt.name != "default" {
// 			// 	assert.Equal(t, tt.contentType, resp.Header.Get("Content-Type"))
// 			// }
// 			// if tt.method == http.MethodGet {
// 				// assert.Equal(t, tt.result, string(body))
// 			// }
// 			time.Sleep(1* time.Second)
// 		})
// 	}

// }