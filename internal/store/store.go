package store

import (
	"context"
	// "errors"
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
func (d DBStorage) UpdateAccrual(allOrdersMap *repo.AllOrdersMap) {
	for numberStr, order := range *allOrdersMap {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			logg.Error().Err(err).Send()
		}
		if number == 0 {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel() 
		query := `UPDATE Orders
			SET
			Accrual = $1
			WHERE Id = $2
		`
		accrual := int(order.Accrual * 100)
		_, err = d.dbconn.ExecContext(ctx, query, accrual, order.ID)
		if err != nil {
			logg.Error().Err(err).Send()
		}
	}
}
func (d DBStorage) GetNotDoneOrders(allOrdersMap *repo.AllOrdersMap) error {
	logg.Print("----------------------GetNotDoneOrders----start--------------------------------")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `
	select 
	Id,
	Number
	from Orders where Status < 4;
	`
    rows, err := d.dbconn.QueryContext(ctx, query)
	if err != nil {
		logg.Error().Err(err).Send()
        return err
    }
    // обязательно закрываем перед возвратом функции
    defer rows.Close()
    // пробегаем по всем записям
    for rows.Next() {
		var order repo.Order
        var number int
		var id int
		err = rows.Scan(&id, &number)
		// err = rows.Scan(&order.ID, &number)
        if err != nil {
			logg.Error().Err(err).Send()
            return err
        }
		numberStr:= strconv.Itoa(number)
		order.Number = numberStr
		order.ID = id
		logg.Printf("order: %#v", order)
		logg.Printf("allOrdersMap: %#v", *allOrdersMap)
        (*allOrdersMap)[numberStr] = order
    }
    // проверяем на ошибки
    err = rows.Err()
    if err != nil {
		logg.Error().Err(err).Send()
        return err
    }
    return nil
}
func (d DBStorage) GetBalanceWithdrawals(userID int) (*repo.Withdraws, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `
	select 
	number
	(select Number from Orders where OrderId = OrderId as number), 
	Withdrawn, 
	PrcessedAt
	from Withdrawns where UserId=$1;
	`
    rows, err := d.dbconn.QueryContext(ctx, query, userID)
	if err != nil {
        return nil, err
    }
    // обязательно закрываем перед возвратом функции
    defer rows.Close()
    var withdraws repo.Withdraws
    // пробегаем по всем записям
    for rows.Next() {
		var withdraw repo.Withdraw
        var withdrawn int
		var number int
		var processedAt time.Time
		err = rows.Scan(&number, &withdrawn, &processedAt)
        if err != nil {
            return nil, err
        }
		withdraw.Order = strconv.Itoa(number)
		withdraw.Sum = float64(withdrawn) / 100
		withdraw.ProcessedAt = processedAt.Format(time.RFC3339)
        withdraws = append(withdraws, withdraw)
    }
    // проверяем на ошибки
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return &withdraws, nil
}

func (d DBStorage) PutWithdraw(userID int, withdraw repo.Withdraw) (bool, error) {
	accrual := d.GetAccrualForUser(userID)
	if withdraw.Sum > accrual {
		return false, repo.ErrNoMoney
	}
	orderID, err := d.ObtainNewOrder(userID, int(withdraw.Sum * 100))
	if err != nil {
		logg.Printf("не удалось загрузить новый заказ: %v", err)
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel() 
	query := `INSERT INTO Withdraws
    (
        UserId, 
        OrderId, 
        Withdraw, 
		ProcessedAt
    ) values($1, $2, $3, $4);
    `
    withdrawn := int(withdraw.Sum * 100)
	_, err = d.dbconn.ExecContext(ctx, query, userID, orderID, withdrawn, time.Now())
    if err != nil {
		logg.Printf("не удалось выполнить запрос на списание: %v", err)
		return false, err
	}
	
	return true, nil
}
func (d DBStorage) GetAccrualForUser(userID int) float64 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select SUM(Accrual) from Orders where UserId=$1;
	`
    row := d.dbconn.QueryRowContext(ctx, query, userID)
	var sum int
	_ = row.Scan(&sum)
   
    return float64(sum) / 100
}
func (d DBStorage) GetBalance(userID int) *repo.Balance {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
    query := `
	select SUM(Withdrawn) from Withdrawns where UserId=$1;
	`
    row := d.dbconn.QueryRowContext(ctx, query, userID)
	var withdrawn int
	_ = row.Scan(&withdrawn)
	current := d.GetAccrualForUser(userID)
   var balance = repo.Balance{
		Current: float64(current),
		Withdrawn: float64(withdrawn) / 100,
   }
    return &balance
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
			order.Accrual = float64(accrual.Int64) / 100
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
func (d DBStorage) ObtainNewOrder(userID, number int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `INSERT INTO Orders
    (
        UserId, 
        Number, 
		Status,
		UploadedAt
    ) values($1, $2, $3, $4) RETURNING Id;
    `
    row := d.dbconn.QueryRowContext(ctx, query, userID, number, repo.NEW, time.Now())
	var orderID int
	err := row.Scan(&orderID)
	if err != nil {
		logg.Printf("не удалось выполнить запрос создания заказа: %v", err)
		return 0, err
	}
	return orderID, nil
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
		return nil, repo.ErrUserNotPresent
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
	logg.Print("--------------------UserAuthenticate------------1-------------start--------------")
	u1ptr, err := d.UserGet(u.Name)
	if err != nil {
		logg.Error().Err(err).Send()
		return false, repo.ErrUserNotPresent
	}
	if u.Password == "" {
		logg.Error().Err(repo.ErrPasswordIsEmpty).Send()
		return false, repo.ErrPasswordIsEmpty
	}
	u.Password = toking.Sha256([]byte(u.Password))
	logg.Print("after hash u.Password: ", u.Password)
	ok := u.Password == u1ptr.Password
	if !ok {
		logg.Error().Err(repo.ErrPasswordAndLoginMismatch).Send()
		return false, repo.ErrPasswordAndLoginMismatch
	} else {
		logg.Print("Аутентификация прошла успешно.")
		return true, nil
	}
}
func (d DBStorage) UserRegistrate(u repo.User) error {
	logg.Print("--------------------UserRegistrate------------1-------------start---------------")
	ok := d.UserIsPresent(u.Name)
	if  ok {
		return repo.ErrUserIsPresent
	}
	if u.Password == "" {
		return repo.ErrPasswordIsEmpty
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
		DROP TABLE IF EXISTS Withdrawns;
		DROP TABLE IF EXISTS Orders;
		DROP TABLE IF EXISTS Users;
		CREATE TABLE Users
        (
			Id SERIAL PRIMARY KEY,
            Name VARCHAR(255) NOT NULL UNIQUE,
            Password VARCHAR(255),
			Withdrawn INTEGER
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
		CREATE TABLE Withdrawns
		(
			Id SERIAL PRIMARY KEY,
			UserId INTEGER,
			OrderId INTEGER,
			Withdrawn INTEGER,
			ProcessedAt timestamptz,
			FOREIGN KEY (UserId) REFERENCES Users (Id) ON DELETE CASCADE,
			FOREIGN KEY (OrderId) REFERENCES Orders (Id) ON DELETE CASCADE
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
