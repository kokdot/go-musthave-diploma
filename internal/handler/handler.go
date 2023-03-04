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

func GetBalanceWithdrawals(w http.ResponseWriter, r *http.Request) {
	logg.Print("-----------------------------GetBalanceWithdrawals-------start-------------------------------------------")
	name, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		w.Header().Set("content-type", "application/json")
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	logg.Print("Получен запрос для пользователя: ", name, "Проверка cookie прошла успешно.")
	userID, ok := m.UserIsPresentReturnUserID(name)
	if !ok {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "такого пользователя. не существует вам необходимо пройти регистрацию или аутентификацию. login failed", http.StatusUnauthorized)
	}
	logg.Print("Данный пользователь присутствует в системе.")
	withdraws, err := m.GetBalanceWithdrawals(userID)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logg.Printf("Полученный список выведенных средств: %#v", withdraws)
	if len(*withdraws) == 0 {
		logg.Error().Err(repo.ErrNoDataForAnswer).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "нет ни одного списания", http.StatusNoContent )
	}
	bodyBytes, err := json.Marshal(&withdraws)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(bodyBytes)
	w.WriteHeader(http.StatusOK)
	
}
func Withdraw(w http.ResponseWriter, r *http.Request) {
	logg.Print("-----------------------------Withdraw-------start-------------------------------------------")
	name, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		w.Header().Set("content-type", "application/json")
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	logg.Print("Получен запрос для пользователя: ", name, "Проверка cookie прошла успешно.")
	userID, ok := m.UserIsPresentReturnUserID(name)
	if !ok {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "такого пользователя. не существует вам необходимо пройти регистрацию или аутентификацию. login failed", http.StatusUnauthorized)
	}
	logg.Print("Данный пользователь присутствует в системе.")
	withdraw := repo.Withdraw{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(bodyBytes, &withdraw)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logg.Printf("Getting of requets is: %#v\n", withdraw)
	number, err := strconv.Atoi(withdraw.Order)
	logg.Print("Получен заказ номер: ", number)
	if err != nil {
		logg.Error().Err(repo.ErrInvalidFormatNumberOfOrder).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if !luna.Valid(number) {
		logg.Error().Err(repo.ErrInvalidFormatNumberOfOrder).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	logg.Print("Проверка luna прошла успешно.")
	ok, err = m.PutWithdraw(userID, withdraw)
	
	if  !ok {
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrNoMoney):
				logg.Error().Err(err).Send()
				w.Header().Set("content-type", "application/json")	
				http.Error(w, "на счету недостаточно средств", http.StatusPaymentRequired)
			default:
				logg.Error().Err(err).Send()
				w.Header().Set("content-type", "application/json")	
				http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			}
		} else {
			logg.Error().Err(err).Send()
			w.Header().Set("content-type", "application/json")	
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		}

	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
func Balance(w http.ResponseWriter, r *http.Request) {
	logg.Print("-----------------------------Balance-------start-------------------------------------------")
	name, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		w.Header().Set("content-type", "application/json")
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	logg.Print("Получен запрос для пользователя: ", name, "Проверка cookie прошла успешно.")
	userID, ok := m.UserIsPresentReturnUserID(name)
	if !ok {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "такого пользователя. не существует вам необходимо пройти регистрацию или аутентификацию. login failed", http.StatusUnauthorized)
	}
	logg.Print("Данный пользователь присутствует в системе.")
	balance := m.GetBalance(userID)
	bodyBytes, err := json.Marshal(&balance)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(bodyBytes)
	w.WriteHeader(http.StatusOK)
}
func UploadOrders(w http.ResponseWriter, r *http.Request) {
	logg.Print("-----------------------------UploadOrders-------start-------------------------------------------")
	name, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		w.Header().Set("content-type", "application/json")
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	logg.Print("Получен запрос для пользователя: ", name, "Проверка cookie прошла успешно.")
	userID, ok := m.UserIsPresentReturnUserID(name)
	if !ok {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "такого пользователя. не существует вам необходимо пройти регистрацию или аутентификацию. login failed", http.StatusUnauthorized)
	}
	logg.Print("Данный пользователь присутствует в системе.")
	orders := m.GetListOrders(userID)
	logg.Printf("Полученный список заказов: %#v", orders)
	if len(*orders) == 0 {
		logg.Error().Err(repo.ErrNoDataForAnswer).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "нет данных для ответа", http.StatusNoContent )
	}
 	bodyBytes, err := json.Marshal(&orders)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(bodyBytes)
	w.WriteHeader(http.StatusOK)
}

func DownloadOrderNumber(w http.ResponseWriter, r *http.Request) {
	logg.Print("-----------------------------DownloadOrderNumber-------start-------------------------------------------")
	name, ok, err := CheckCookieAutentication(r)
	if !ok {
		logg.Error().Err(err).Send()
		w.Header().Set("content-type", "application/json")
		http.Error(w, "логин или пароль не совпадают. login failed", http.StatusUnauthorized)
	}
	logg.Print("Получен запрос для пользователя: ", name, "Проверка cookie прошла успешно.")
	ok = m.UserIsPresent(name)
	if !ok {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		http.Error(w, "такого пользователя. не существует вам необходимо пройти регистрацию или аутентификацию. login failed", http.StatusUnauthorized)
	}
	logg.Print("Данный пользователь присутствует в системе.")
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
	logg.Print("Получен заказ номер: ", number)
	if !luna.Valid(number) {
		logg.Error().Err(repo.ErrInvalidFormatNumberOfOrder).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	logg.Print("Проверка luna прошла успешно.")
	ok = m.CheckExistOrderNumber(number)
	logg.Print("Провека сущеустаования данного номера заказа: ", ok)
	var userID int
	if ok {
		logg.Print("Данный заказ уже сущетвует.")
		userID = m.GetIDOrderOwner(number)
		
		logg.Print("Id ползователя, чей это заказ: ", userID)
		userName := m.GetUserNameByID(userID)
		if userName == name {
			logg.Print("заказ принадлежит пользователю с Id: ", userID)
			logg.Error().Err(repo.ErrOrderUsedUser).Send()	
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			logg.Print("заказ не принадлежит пользователю с Id: ", userID)
			logg.Error().Err(repo.ErrOrderUsedUnotherUser).Send()	
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusConflict)
			return
		}
	} else {
		userID = m.GetUserIDByName(name)
	}
	logg.Print("создаем заказ для пользователя с Id: ", userID, "; и номером заказа :  ", number)
	_, err = m.ObtainNewOrder(userID, number)
	if err != nil {
		logg.Error().Err(err).Send()	
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logg.Print("новый номер заказа принят в обработку")	
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	// return
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
			case errors.Is(err, repo.ErrPasswordAndLoginMismatch):
				logg.Error().Err(err).Send()	
				http.Error(w, "неверная пара логин/пароль", http.StatusUnauthorized)
			case errors.Is(err, repo.ErrPasswordIsEmpty):
				logg.Error().Err(err).Send()	
				http.Error(w, "пароль не может быть пустым", http.StatusBadRequest)
			default:
				logg.Error().Err(err).Send()	
				http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			}
		}
	} 
	cookie := auth.SetCookie(toking.Toking(), u.Name, m.GetSeckretKey())
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
		case errors.Is(err, repo.ErrUserIsPresent):
			http.Error(w, "логин уже занят", http.StatusConflict)//PasswordIsEmpty
		case errors.Is(err, repo.ErrPasswordIsEmpty):
			http.Error(w, "пароль не может быть пустым", http.StatusBadRequest)//StatusBadRequest
		default:
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)//StatusBadRequest 
		}
	}
	
	cookie := auth.SetCookie(toking.Toking(), u.Name, m.GetSeckretKey())
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
