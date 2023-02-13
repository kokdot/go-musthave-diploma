package repo

import (
	// "time"
	// "fmt"
	// "sort"
) 

type Repo interface {

// 	Save(mtx *Metrics) (*Metrics, error)
// 	Get(id string) (*Metrics, error)
// 	GetAll() (StoreMap, error)
// 	SaveCounterValue(name string, counter Counter) (Counter, error)
// 	SaveGaugeValue(name string, gauge Gauge) error
// 	GetCounterValue(name string) (Counter, error)
// 	GetGaugeValue(name string) (Gauge, error)
// 	GetAllValues() string
// 	ReadStorage() error
// 	WriteStorage() error
// 	GetURL() string
// 	GetKey() string
// 	GetStoreFile() string
// 	GetRestore() bool
// 	GetStoreInterval() time.Duration 
// 	GetDataBaseDSN() string
// 	GetPing() (bool, error)
// 	SaveByBatch1(*StoreMap) (*StoreMap, error)
// 	SaveByBatch([]Metrics) (*[]Metrics, error)
}

type User struct {
	Name string `json:"login"`
	Password string `json:"password"`
}

type BD interface {
	UserRegistrate(u User) error
	GetSecretKey() string

}

// type Consumer interface {
//     ReadStorage() (*StoreMap, error) // для чтения события
//     Close() error               // для закрытия ресурса (файла)
// }
// type Producer interface {
//     WriteStorage() error // для записи события
//     Close() error            // для закрытия ресурса (файла)
// }
// type Metrics struct {
// 	ID    string   `json:"id"`              // имя метрики
// 	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
// 	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
// 	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
// 	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
// 	// Hash  []byte   `json:"hash,omitempty"`  // значение хеш-функции
// }
// type Counter int64
// type Gauge float64
// type StoreMap map[string]Metrics
// func StoreMapToString(smPtr *StoreMap) string {
// 	if smPtr == nil {
// 		return ""
// 	}
// 	var str string
// 	var v Gauge
// 	var d Counter
// 	var i int
// 	keys := make([]string, 0, len(*smPtr))
// 	for k := range *smPtr {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)
// 	for _, key := range keys {
// 		i++
// 		if (*smPtr)[key].Delta != nil {
// 			d = *(*smPtr)[key].Delta
// 		}
// 		if (*smPtr)[key].Value != nil {
// 			v = *(*smPtr)[key].Value
// 		}
// 		str += fmt.Sprintf("%d; %s: %v %v\n",i , key, v, d)
// 	}
// 	return str
// }
