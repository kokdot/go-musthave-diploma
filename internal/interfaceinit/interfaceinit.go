package interfaceinit

import (
	"fmt"
	// "time"
	"errors"
	"github.com/rs/zerolog"

	"github.com/kokdot/go-musthave-diploma/internal/repo"
	"github.com/kokdot/go-musthave-diploma/internal/store"
)

// var m  repo.Repo

func InterfaceInit(address string, accrualSysemAddress string, dataBaseURI string, logg zerolog.Logger) (repo.Repo, error) {
	if dataBaseURI == "" {
		return nil, errors.New("dataBaseURI is empty, failed to create DBStorage")
	}
	logg.Print("-------before--------NewDBStorage---------------------------------")
	d, err := store.NewDBStorage(address , accrualSysemAddress, dataBaseURI)
	logg.Print("----------d, err := -------", d, "------", err)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create DBStorage, err: %s", err)
	}
	return d, nil
}

