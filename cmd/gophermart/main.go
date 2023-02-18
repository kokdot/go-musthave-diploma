      package main

import (
	// "log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// "fmt"
	// "time"

	"github.com/kokdot/go-musthave-diploma/internal/handler"
	"github.com/kokdot/go-musthave-diploma/internal/interfaceinit"
	"github.com/kokdot/go-musthave-diploma/internal/onboardingserver"
	"github.com/kokdot/go-musthave-diploma/internal/store"
)


// var logg = log.Logger
func main() { 
    address, accrualSysemAddress, dataBaseURI, logg  := onboardingserver.OnboardingServer()
    logg.Print("--------------------main-------------------------------------------")
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
    // metricsserver.GetLogg(logg)0
    logg.Printf("interface m:   %#v", m)
    logg.Print("--------------------main--started-----------------------------------------")
    
    // определяем роутер chi
    r := chi.NewRouter()
    // зададим встроенные middleware, чтобы улучшить стабильность приложения
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)    
    r.Use(middleware.Recoverer)
    // r.Use(middleware.Compress(5, "gzip"))
    r.Use(middleware.Compress(5))
    // r.Get("/", handler.GetAll)
    // r.Get("/ping", handler.GetPing)

    r.Get("/", handler.GetOk)
    r.Post("/api/user/register", handler.Registration)
    r.Post("/api/user/test", handler.CheckUserLogin)
    r.Post("/api/user/login", handler.Authentication)
    r.Post("/api/user/orders", handler.DownloadOrderNumber)
    err = http.ListenAndServe(address, r)
    logg.Fatal().Err(err).Send()
    // log.Fatal(http.ListenAndServe(":8080", r))
}






















    // r.Post("/updates1/", handler.PostUpdateByBatch1)
    
    // r.Route("/update", func(r chi.Router) {
    //     r.Post("/", handler.PostUpdate)
    //     r.Route("/counter", func(r chi.Router) {
    //         r.Route("/{nameData}/{valueData}", func(r chi.Router) {
    //             r.Use(handler.PostCounterCtx)
    //             r.Post("/", handler.PostUpdateCounter)
    //         })
    //     })
    //     r.Route("/gauge", func(r chi.Router) {
    //         r.Route("/{nameData}/{valueData}", func(r chi.Router) {
    //             r.Use(handler.PostGaugeCtx)
    //             r.Post("/", handler.PostUpdateGauge)
    //         })
    //     })
    //     r.Route("/",func(r chi.Router) {
    //         r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
	// 	        w.Header().Set("content-type", "text/plain; charset=utf-8")
    //             w.WriteHeader(http.StatusNotImplemented)
    //             // fmt.Fprint(w, "line: 52; http.StatusNotImplemented")
	//         })
    //     })
    // })

    // r.Route("/value", func(r chi.Router) {
    //     r.Post("/", handler.GetValue)
	// 	r.Route("/counter", func(r chi.Router){
    //         r.Route("/{nameData}", func(r chi.Router) {
    //             r.Use(handler.GetCtx)
    //             r.Get("/", handler.GetCounter)
    //         })
    //     })
    //    	r.Route("/gauge", func(r chi.Router){
    //         r.Route("/{nameData}", func(r chi.Router) {
    //             r.Use(handler.GetCtx)
    //             r.Get("/", handler.GetGauge)
    //         })
    //     })
	// })



