package repo

// import "time"

// "time"
// "fmt"
// "sort" 
type Status int
const (
	NEW	Status			= iota + 1
	PROCESSING
	INVALID
	PROCESSED
)
var StatusSlice = map[Status]string{ 
	NEW: "NEW",
	PROCESSING: "PROCESSING",
	INVALID: "INVALID",
	PROCESSED: "PROCESSED",
}
type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}
type Order struct {
	Number string `json:"number"`
	Status string `json:"status"`
	Accrual int `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}
type Orders []Order

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
	ObtainNewOrder(userID, number int) error
	GetListOrders(userID int) *Orders
	UserIsPresentReturnUserID(name string) (int, bool)

	// GetOk()

}


