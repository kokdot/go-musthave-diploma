package accrual

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/rs/zerolog"
	"github.com/kokdot/go-musthave-diploma/internal/repo"
)
var logg zerolog.Logger
var accrualSysemAddress string
var B  = make(repo.AllOrdersMap, 0)
func GetLogg(loggReal zerolog.Logger, accrualSysemAddressReal string)  {
	logg = loggReal
	accrualSysemAddress = accrualSysemAddressReal 
}
func GetAccrual(orders *repo.AllOrdersMap) error {
	logg.Print("-----------------------GetAccrual------------start----------")
	// func (c *Client) Do(req *Request) (*Response, error)
	client := http.Client{}
	for number, order := range *orders {
		logg.Printf("number: %v\n", number)
		logg.Printf("order: %v\n", order)
		var order1 repo.Order
		
		url := fmt.Sprintf("%s/api/orders/%s", accrualSysemAddress, number)
		// url := fmt.Sprintf("http://localhost:8080/api/orders/%s", number)
		logg.Printf("---------------------------url: %v\n", url)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			logg.Error().Err(err).Send()
			logg.Print("-----------GetAccrual-------------return-----------------------")
			return err
		}
		request.Header.Add("Accept", "application/json") 
		requestDump, err := httputil.DumpRequest(request, true)
		if err != nil {
			fmt.Println(err.Error())

		}
		logg.Print("request Dump: ", string(requestDump)) 
		response, err := client.Do(request)
		if err != nil {
			logg.Error().Err(err).Send()
			logg.Print("-----------GetAccrual-------------return-----------------------")
			return err
		}
		logg.Printf("---------------------!!!!!!!!!!!!!!!!!!!!----response: %#v", response)
		switch response.StatusCode {
		case 200:
			logg.Print("-----------GetAccrual-------------status-OK-----------------------")
			
			bodyBytes, err := io.ReadAll(response.Body)
			if err != nil {
				logg.Error().Err(err).Send()
				logg.Print("-----------GetAccrual-------------return-----------------------")
				return err
			}
			err = json.Unmarshal(bodyBytes, &order1)
			if err != nil {
				logg.Error().Err(err).Send()
				logg.Print("-----------GetAccrual-------------return-----------------------")
				return err
			}
			order.Accrual = order1.Accrual
			order.Status = order1.Status
			logg.Printf("order: %#v", order)
			logg.Printf("order1: %#v", order1)
			logg.Printf("orders: %#v", orders)
		case 429:
			logg.Print("-----------GetAccrual-------------return-----------------------")
			return repo.Err429
		}

	}
	logg.Print("-----------GetAccrual-------------return---------OK--------------")
	return nil
}
