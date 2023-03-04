package interfaceinit

import (
	"fmt"
	// "time"
	"errors"
	"github.com/rs/zerolog"

	"github.com/kokdot/go-musthave-diploma/internal/repo"
	"github.com/kokdot/go-musthave-diploma/internal/store"
)


func InterfaceInit(address string, accrualSysemAddress string, dataBaseURI string, logg zerolog.Logger) (repo.Repo, error) {
	logg.Print("-------------------InterfaceInit---------start--------------")
	if dataBaseURI == "" {
		return nil, errors.New("dataBaseURI is empty, failed to create DBStorage")
	}
	d, err := store.NewDBStorage(address , accrualSysemAddress, dataBaseURI)
	logg.Print("----------d: ", d)
	logg.Print("--------err: ", err)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create DBStorage, err: %s", err)
	}
	return d, nil
}

