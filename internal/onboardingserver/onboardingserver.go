package onboardingserver

import (
    // "time"
	"strconv"
	"os"
	"flag"
    "github.com/caarlos0/env/v6"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)
const (
    Address = "127.0.0.1:8081"
    AccrualSysemAddress = ""
    DataBaseURI = "postgres://postgres:postgrespw@localhost:49153"
    Debug = false
)

type Config struct {
    Address  string 		`env:"RUN_ADDRESS"`// envDefault:"127.0.0.1:8080"`
    AccrualSysemAddress string 			`env:"ACCRUAL_SYSTEM_ADDRESS"`
    DataBaseURI string    `env:"DATABASE_URI"`
}
var (
    addressReal = Address
    dataBaseURIReal = DataBaseURI
    accrualSysemAddressReal = AccrualSysemAddress
    cfg Config
    logg zerolog.Logger
)

func GetLogg() zerolog.Logger {
	return logg
}

func OnboardingServer() (string, string, string, zerolog.Logger) {
	logg.Print("---------onboarding-------------------")
    err := env.Parse(&cfg)
    if err != nil {
        logg.Print(err)
    }

    addressPtr := flag.String("a", "127.0.0.1:8081", "ip adddress of server")
    accrualSysemAddressPtr := flag.String("r", "", "address system of accrual")
    DataBaseURIPtr := flag.String("d", "", "Data Base URI")
    debug := flag.Bool("debug", false, "sets log level to debug")

    flag.Parse()
    zerolog.SetGlobalLevel(zerolog.DebugLevel)
    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// log.Logger = log.With().Caller().Logger()
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
	logg = log.Logger

    addressReal = *addressPtr
    accrualSysemAddressReal = *accrualSysemAddressPtr
    if *DataBaseURIPtr != "" {
        dataBaseURIReal = *DataBaseURIPtr
    }

    if cfg.Address != "" {
        addressReal	= cfg.Address
    }
     
    if cfg.AccrualSysemAddress != "" {
        accrualSysemAddressReal	= cfg.AccrualSysemAddress
    }
    if cfg.DataBaseURI != "" {
        dataBaseURIReal	= cfg.DataBaseURI
    }
    logg.Print("----------OnboardingServer-------------has-finished-----------")
    logg.Print("Address:  ", Address)
    logg.Print("AccrualSysemAddress:  ", AccrualSysemAddress)
    logg.Print("DataBaseURI:  ", DataBaseURI)
    logg.Print("---------------------------flag------------------------------")
    logg.Print("addressPtr:", *addressPtr)
    logg.Print("accrualSysemAddressPtr:", *accrualSysemAddressPtr)
    logg.Print("DataBaseURIPtr:", *DataBaseURIPtr)
    logg.Print("debug:     ", *debug)
    logg.Print("---------------------------cfg------------------------------")
    logg.Print("cfg.Address:", cfg.Address)
    logg.Print("cfg.AccrualSysemAddress:", cfg.AccrualSysemAddress)
    logg.Print("cfg.DataBaseURI:", cfg.DataBaseURI)
    logg.Print("------------------------real---------------------------------")
    logg.Print("addressReal:", addressReal)
    logg.Print("accrualSysemAddressReal:", accrualSysemAddressReal)
    logg.Print("dataBaseURIReal:", dataBaseURIReal)
    return addressReal, accrualSysemAddressReal, dataBaseURIReal, log.Logger
}



func GetAddress() string {
    return addressReal
}

func GetAccrualSysemAddress() string {
return accrualSysemAddressReal
}