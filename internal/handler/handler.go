package handler

import (
	"encoding/json"
	"errors"
	"strconv"

	// "fmt"
	"io"

	// "time"
	"net/http"
	// "strconv"

	// "github.com/go-chi/chi/v5"
	"github.com/kokdot/go-musthave-diploma/internal/luna"
	"github.com/kokdot/go-musthave-diploma/internal/repo"

	// "github.com/kokdot/go-musthave-diploma/internal/store"
	"github.com/kokdot/go-musthave-diploma/internal/auth"
	"github.com/kokdot/go-musthave-diploma/internal/toking"
	"github.com/rs/zerolog"
)
 
var ErrUserIsPresent = errors.New("user is present")
var ErrPasswordIsEmpty = errors.New("password is empty")
var ErrUserNotPresent error = errors.New("user not present")
var ErrPasswordAndLoginMismatch = errors.New("password and login mismatch")
var ErrInternalServerError = errors.New("internal server error")
var ErrInvalidFormatNumberOfOrder = errors.New("invalid format number of order")//неверный формат номера заказа
var ErrOrderUsedUser = errors.New("this order being download yet")//номер заказа уже был загружен этим пользователем
var ErrOrderUsedUnotherUser = errors.New("this order being download yet by unother user")//номер заказа уже был загружен другим пользователем


var m  repo.Repo
var logg zerolog.Logger

func PutM(M repo.Repo) {
	m = M
}
func GetLogg(loggReal zerolog.Logger)  {
	logg = loggReal
}
func CheckCookieAutentication(r *http.Request) (string, bool, error) {
	name, ok, err := auth.ValidCookie(r, m.GetSeckretKey())
	return name, ok, err
}

func DownloadOrderNumber(w http.ResponseWriter, r *http.Request) {
	name, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	ok = m.UserIsPresent(name)
	if !ok {
		logg.Error().Err(err).Send()	
		http.Error(w, "такой пользователь не существует. login failed", http.StatusUnauthorized)
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	number, err := strconv.Atoi(string(bodyBytes))
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !luna.Valid(number) {
		logg.Error().Err(ErrInvalidFormatNumberOfOrder).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	ok = m.CheckExistOrderNumber(number)
	var userId int
	if ok {
		userId = m.GetIdOrderOwner(number)
		userName := m.GetUserNameById(userId)
		if userName == name {
			logg.Error().Err(ErrOrderUsedUser).Send()	
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			logg.Error().Err(ErrOrderUsedUnotherUser).Send()	
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusConflict)
			return
		}
	}
	err = m.ObtainNewOrder(userId, number)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logg.Print("новый номер заказа принят в обработку")	
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	return
}

func GetOk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	logg.Print("m: ", m)
	_, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
func Authentication(w http.ResponseWriter, r *http.Request) {
	logg.Print("--------------------Registration------------1-------------start-------------------------------")
	u := repo.User{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logg.Printf("Getting of requets is: %#v\n", u)
	logg.Print("m: ", m)
	ok, err := m.UserAuthenticate(u)
	if !ok {
		if err != nil {
			switch {
			case errors.Is(err, ErrPasswordAndLoginMismatch):
				logg.Error().Err(err).Send()	
				http.Error(w, "неверная пара логин/пароль", http.StatusUnauthorized)
			case errors.Is(err, ErrPasswordIsEmpty):
				logg.Error().Err(err).Send()	
				http.Error(w, "пароль не может быть пустым", http.StatusBadRequest)
			default:
				logg.Error().Err(err).Send()	
				http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			}
		}
	} 
	cookie := auth.SetCookie(toking.Toking(), m.GetSeckretKey())
	logg.Print("cookie: ", cookie)
	http.SetCookie(w, cookie)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	logg.Print("--------------------Registration------------1-------------start-------------------------------")
	u := repo.User{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &u)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logg.Printf("Getting of requets is: %#v\n", u)
	logg.Print("m: ", m)
	err = m.UserRegistrate(u)

	if err != nil {
		switch {
		case errors.Is(err, ErrUserIsPresent):
			http.Error(w, "логин уже занят", http.StatusConflict)//PasswordIsEmpty
		case errors.Is(err, ErrPasswordIsEmpty):
			http.Error(w, "пароль не может быть пустым", http.StatusBadRequest)//StatusBadRequest
		default:
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)//StatusBadRequest 
		}
	}
	
	cookie := auth.SetCookie(toking.Toking(), m.GetSeckretKey())
	logg.Print("cookie: ", cookie)

	http.SetCookie(w, cookie)
	
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
// func GetPing(w http.ResponseWriter, r *http.Request) {
// 	ok, err := m.GetPing()
//  	if !ok {
// 		w.Header().Set("content-type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		logg.Printf("%s", err)
// 		return
// 	} else {
// 		w.Header().Set("content-type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}
// }
