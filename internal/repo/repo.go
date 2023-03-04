package repo

import (
	"errors"

	// "github.com/kokdot/go-musthave-diploma/internal/accrual"
)

// "time"
// "fmt"
// "sort"
type Status int
const (
	NEW	Status	= iota + 1
	REGISTERED		
	PROCESSING
	INVALID
	PROCESSED
)
var ErrUserIsPresent = errors.New("user is present")
var ErrPasswordIsEmpty = errors.New("password is empty")
var ErrUserNotPresent error = errors.New("user not present")
var ErrPasswordAndLoginMismatch = errors.New("password and login mismatch")
var ErrInternalServerError = errors.New("internal server error")
var ErrInvalidFormatNumberOfOrder = errors.New("invalid format number of order")//неверный формат номера заказа
var ErrOrderUsedUser = errors.New("this order being download yet")//номер заказа уже был загружен этим пользователем
var ErrOrderUsedUnotherUser = errors.New("this order being download yet by unother user")//номер заказа уже был загружен другим пользователем
var ErrNoDataForAnswer = errors.New("there is no data for answer")//нет данных для ответа
var ErrNoMoney = errors.New("there is no enuf money for this order")//на счету недостаточно средств
var Err429 = errors.New("429")//429


type Withdraw struct {
	Order string `json:"order"`
	Sum int `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}
type Balance struct {
	Current float64 `json:"current"`
	Withdrawn int `json:"withdrawn"`
}
var StatusSlice = map[Status]string{ 
	NEW: "NEW",
	REGISTERED: "REGISTERED",
	PROCESSING: "PROCESSING",
	INVALID: "INVALID",
	PROCESSED: "PROCESSED",
}
var StatusSliceFeedBack = map[string]Status{ 
	"NEW": NEW,
	"REGISTERED": REGISTERED,
	"PROCESSING": PROCESSING,
	"INVALID": INVALID,
	"PROCESSED": PROCESSED,
}
type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}
type Order struct {
	ID int		`json:"-"`
	Number string `json:"number"`
	Status string `json:"status"`
	Accrual int `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}
type Orders []Order
type Withdraws []Withdraw
type AllOrdersMap map[string]Order //  string - order

type Repo interface {
	UserRegistrate(u User) error
	GetSeckretKey() []byte
	UserIsPresent(name string) bool
	UserAuthenticate(u User) (bool, error)
	UserGet(name string) (*User, error)
	CheckExistOrderNumber(number int) bool
	GetIDOrderOwner(naumber int) int
	GetUserNameByID(userID int) string
	GetUserIDByName(name string) int
	ObtainNewOrder(userID, number int) (int, error)
	GetListOrders(userID int) *Orders
	UserIsPresentReturnUserID(name string) (int, bool)
	GetBalance(userID int) *Balance
	GetAccrualForUser(userID int) int
	PutWithdraw(userID int, withdraw Withdraw) (bool, error)
	GetBalanceWithdrawals(userID int) (*Withdraws, error)
	UpdateAccrual(allOrdersMap *AllOrdersMap)
	GetNotDoneOrders(allOrdersMap *AllOrdersMap) error
	// GetOk()

}


