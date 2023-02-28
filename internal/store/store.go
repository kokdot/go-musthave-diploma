package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"database/sql"

	"github.com/rs/zerolog"

	_ "github.com/jackc/pgx/v5/stdlib"
	// "github.com/kokdot/go-musthave-diploma/internal/auth"
	"github.com/kokdot/go-musthave-diploma/internal/repo"
	"github.com/kokdot/go-musthave-diploma/internal/toking"
)
var ErrUserIsPresent error = errors.New("user is present")
var ErrUserNotPresent error = errors.New("user not present")
var ErrPasswordIsEmpty = errors.New("password is empty")
var ErrPasswordAndLoginMismatch = errors.New("password and login mismatch")
var logg zerolog.Logger
type DBStorage struct {
	// StoreMap      *StoreMap
	accrualSysemAddress    string
	address       string
	dataBaseURI   string
	secretKey []byte
	dbconn        *sql.DB
}
func GetLogg(loggReal zerolog.Logger)  {
	logg = loggReal
}

func (d DBStorage) GetSeckretKey() []byte {
	return d.secretKey
}
func (d DBStorage) GetListOrders(userID int) *repo.Orders {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `
	select Number, Status, Accrual, UploadedAt from Orders where UserId=$1;
	`
	rows, err := d.dbconn.QueryContext(ctx, query, userID)
    if err != nil {
        return nil
    }
    // обязательно закрываем перед возвратом функции
    defer rows.Close()
    var orders repo.Orders
    // пробегаем по всем записям
    for rows.Next() {
        var order repo.Order
		var number int
		var uploadedAt time.Time
		var status int
		var accrual sql.NullInt64
		err = rows.Scan(&number, &status, &accrual, &uploadedAt)
        if err != nil {
            return nil
        }
		order.Number = strconv.Itoa(number)
		order.UploadedAt = uploadedAt.Format(time.RFC3339)

		order.Status = repo.StatusSlice[repo.Status(status)]
		if accrual.Valid {
			order.Accrual = int(accrual.Int64)
		} //else {
			//order.Accrual = 0
		//}
		// fmt.Println("")
        orders = append(orders, order)
    }
    // проверяем на ошибки
    err = rows.Err()
    if err != nil {
        return nil
    }
    return &orders
}
func (d DBStorage) ObtainNewOrder(userID, number int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `INSERT INTO Orders
    (
        UserId, 
        Number, 
		Status,
		UploadedAt
    ) values($1, $2, $3, $4)
    `
    _, err := d.dbconn.ExecContext(ctx, query, userID, number, repo.NEW, time.Now())
    if err != nil {

		logg.Printf("не удалось выполнить запрос создания заказа: %v", err)
		return err
	}
	return nil
}

func (d DBStorage) CheckExistOrderNumber(number int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select exists(select 1 from Orders where Number=$1);
	`
    row := d.dbconn.QueryRowContext(ctx, query, number)
	var ok bool
	_ = row.Scan(&ok)
   
    return ok
}
func (d DBStorage) GetIDOrderOwner(number int) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select UserId from Orders where Number=$1;
	`
    row := d.dbconn.QueryRowContext(ctx, query, number)
	var userID int
	_ = row.Scan(&userID)
   
    return userID
}
func (d DBStorage) GetUserNameByID(userID int) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select Name from Users where Id=$1;
	`
    row := d.dbconn.QueryRowContext(ctx, query, userID)
	var name string
	_ = row.Scan(&name)
   
    return name
}
func (d DBStorage) GetUserIDByName(name string) int {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select Id from Users where Name=$1;
	`
    row := d.dbconn.QueryRowContext(ctx, query, name)
	var userID int
	_ = row.Scan(&userID)
   
    return userID
}
func (d DBStorage) UserGet(name string) (*repo.User, error) {
	ok := d.UserIsPresent(name)
	if !ok {
		return nil, ErrUserNotPresent
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select Name, Password from Users where Name=$1;
	`
    row := d.dbconn.QueryRowContext(ctx, query, name)
	var u repo.User
	_ = row.Scan(&u.Name, &u.Password)
   
    return &u, nil
}

func (d DBStorage) UserIsPresent(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//select exists(select 1 from contact where id=12)
	
    query := `
	select exists(select 1 from Users where Name=$1);
	`
    row := d.dbconn.QueryRowContext(ctx, query, name)
	var ok bool
	_ = row.Scan(&ok)
   
    return ok
}

func (d DBStorage) UserIsPresentReturnUserID(name string) (int, bool) {
	ok := d.UserIsPresent(name)
	if !ok {
		return 0, false
	}
	return d.GetUserIDByName(name), true
}
func (d DBStorage) UserAuthenticate(u repo.User) (bool, error) {
	logg.Print("--------------------UserAuthenticate------------1-------------start-------------------------------")
	u1ptr, err := d.UserGet(u.Name)
	if err != nil {
		logg.Error().Err(err).Send()
		return false, ErrUserNotPresent
	}
	if u.Password == "" {
		logg.Error().Err(ErrPasswordIsEmpty).Send()
		return false, ErrPasswordIsEmpty
	}
	u.Password = toking.Sha256([]byte(u.Password))
	logg.Print("after hash u.Password: ", u.Password)
	ok := u.Password == u1ptr.Password
	if !ok {
		logg.Error().Err(ErrPasswordAndLoginMismatch).Send()
		return false, ErrPasswordAndLoginMismatch
	} else {
		logg.Print("Аутентификация прошла успешно.")
		return true, nil
	}
}
func (d DBStorage) UserRegistrate(u repo.User) error {
	logg.Print("--------------------UserRegistrate------------1-------------start-------------------------------")
	ok := d.UserIsPresent(u.Name)
	if  ok {
		return ErrUserIsPresent
	}
	if u.Password == "" {
		return ErrPasswordIsEmpty
	}
	u.Password = toking.Sha256([]byte(u.Password))
	logg.Print("after hash u.Password: ", u.Password)
	err := d.UserCreate(u)
	if err != nil {
		return err
	}
	return nil
}
func (d DBStorage) UserCreate(u repo.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `INSERT INTO Users
    (
        Name, 
        Password 
    ) values($1, $2)
    `
    _, err := d.dbconn.ExecContext(ctx, query, u.Name, u.Password)
    if err != nil {

		logg.Printf("не удалось выполнить запрос создания пользователя: %v", err)
		return err
	}
	return nil
}

func (d DBStorage) GetSecretKey() []byte {
	return d.secretKey
}

func NewDBStorage(address, accrualSysemAddress, dataBaseURI string) (*DBStorage, error){
	logg.Print("-------------------------NewDBStorage-----------------1---")
    dbconn, err := sql.Open("pgx", dataBaseURI)
	if err != nil {
		logg.Print("-------------------------NewDBStorage-----------------2---")
		return nil, err
	}
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = dbconn.PingContext(ctx); err != nil {
		logg.Print("-------------------------NewDBStorage-----------------3---")
		return nil, err
	}
    secretKey, err := toking.RandBytesKeyString(32)
	if err != nil {
		logg.Print("-------------------------NewDBStorage-----------------4---")
		return nil, err
	}
    var dbStorage =   DBStorage{
        // StoreMap: &sm,
		address: address,
		accrualSysemAddress: accrualSysemAddress,
		dataBaseURI: dataBaseURI,
		secretKey: secretKey,
        dbconn: dbconn,
    }
    if err = dbStorage.CreateTableUsers(); err != nil {
		logg.Print("-------------------------NewDBStorage-----------------5---")
        return nil, err
    }

    return &dbStorage , nil
}

func (d DBStorage) CreateTableUsers() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
		DROP TABLE IF EXISTS Orders;
		DROP TABLE IF EXISTS Users;
		CREATE TABLE Users
        (
			Id SERIAL PRIMARY KEY,
            Name VARCHAR(255) NOT NULL UNIQUE,
            Password VARCHAR(255)
        );
		CREATE TABLE Orders
		(
			Id SERIAL PRIMARY KEY,
			UserId INTEGER,
			Number BIGINT NOT NULL,
			Accrual INTEGER,
			Status INTEGER NOT NULL,
			UploadedAt timestamptz,
			FOREIGN KEY (UserId) REFERENCES Users (Id) ON DELETE CASCADE
		);
	`
    _, err := d.dbconn.ExecContext(ctx, query)
    if err != nil {
		return fmt.Errorf("не удалось выполнить запрос создания таблицы Users: %v", err)
	}
    return nil
}

func (d DBStorage) GetPing() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := d.dbconn.PingContext(ctx); err != nil {
		return false, err
	}
	logg.Print("Ping Ok")
	return true, nil
}
