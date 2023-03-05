package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kokdot/go-musthave-diploma/internal/accrual"
	"github.com/kokdot/go-musthave-diploma/internal/handler"
	"github.com/kokdot/go-musthave-diploma/internal/interfaceinit"
	"github.com/kokdot/go-musthave-diploma/internal/onboardingserver"
	"github.com/kokdot/go-musthave-diploma/internal/repo"
	"github.com/kokdot/go-musthave-diploma/internal/store"
)


func main() { 
    address, accrualSysemAddress, dataBaseURI, logg  := onboardingserver.OnboardingServer()
    logg.Print("----------------------------main------------start-------------------------")
    logg.Print("address:  ", address)
    logg.Print("accrualSysemAddress:  ", accrualSysemAddress)
    logg.Print("dataBaseURI:  ", dataBaseURI)

    m, err := interfaceinit.InterfaceInit(address, accrualSysemAddress, dataBaseURI, logg)
    if err != nil {
        logg.Printf("\nthere in error in starting interface and restore data: %s", err)
    }
    handler.PutM(m)
    handler.GetLogg(logg)
    store.GetLogg(logg)
    accrual.GetLogg(logg, accrualSysemAddress)
    logg.Printf("---------interface m:   %#v", m)
    logg.Print("---------------------------main-----------------is-going--------------")
    
    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)    
    r.Use(middleware.Recoverer)
    r.Use(middleware.Compress(5))
    // r.Get("/ping", handler.GetPing)

    r.Get("/", handler.GetOk)
    r.Post("/api/user/register", handler.Registration)
    r.Post("/api/user/test", handler.CheckUserLogin)
    r.Post("/api/user/login", handler.Authentication)
    r.Post("/api/user/orders", handler.DownloadOrderNumber)
    r.Get("/api/user/orders", handler.UploadOrders)
    r.Get("/api/user/balance", handler.Balance)
    r.Get("/api/user/withdraw", handler.Withdraw)
    r.Get("/api/user/balance/withdrawals", handler.GetBalanceWithdrawals)
    var allOrdersMap = make(repo.AllOrdersMap, 0)
    go func(allOrdersMap *repo.AllOrdersMap) {
		for {
			<-time.After(time.Second * 5)
            list := m.GetListOrders(1)

            logg.Print("----------------------------list Vasya: -----------------------")
            logg.Printf("list Vasya: %#v", *list)
            err := m.GetNotDoneOrders(allOrdersMap)
            if err != nil {
                logg.Error().Err(err).Send()
                continue
            }
            err = accrual.GetAccrual(allOrdersMap)
            if errors.Is(err, repo.Err429) {
                <-time.After(time.Second * 60)
                continue
            }
            if err != nil {
                logg.Error().Err(err).Send()
                continue
            }
			
		}
	}(&allOrdersMap)
    err = http.ListenAndServe(address, r)
    logg.Fatal().Err(err).Send()
}
