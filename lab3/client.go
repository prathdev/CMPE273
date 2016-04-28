package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/routes"
)

type DataSet struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var InstanceC DataSet

var dataInstances map[string]map[string]DataSet

func main() {
	dataInstances = make(map[string]map[string]DataSet)
	mux1 := routes.New()
	input := os.Args[1:]

	ports := strings.Split(string(input[0]), "-")

	startport, _ := strconv.Atoi(ports[0])
	endport, _ := strconv.Atoi(ports[1])
	mux1.Get("/", GetAllKeys)
	mux1.Get("/:key", GetKey)
	mux1.Put("/:key/:value", PutKey)
	http.Handle("/", mux1)

	for index := startport; index <= endport; index++ {
		go func() {


			dataInstances[strconv.Itoa(index)] = make(map[string]DataSet)

			log.Print("Listening to" + strconv.Itoa(index))
			http.ListenAndServe(":"+strconv.Itoa(index), nil)
		}()
		time.Sleep(1 * time.Second)
	}
	for 1 < 2 {

	}
}



func GetAllKeys(w http.ResponseWriter, r *http.Request) {
	portNum := strings.Split(r.Host, ":")[1]
	log.Print("Inside GetAll")
				log.Print(portNum)
	valuesarr := make([]DataSet, 0)
				for _, value := range dataInstances[portNum] {
		valuesarr = append(valuesarr, value)
	}
	returnjson, _ := json.Marshal(valuesarr)
	w.Write([]byte(returnjson))
}

func GetKey(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	keydata := params.Get(":key")
				portNum := strings.Split(r.Host, ":")[1]
	log.Print("Inside Get")
	log.Print(portNum)


	keyvalstruct := dataInstances[portNum][keydata]
		if keyvalstruct.Key != "" {
		returnjson, _ := json.Marshal(keyvalstruct)
		w.Write([]byte(returnjson))
	} else {
		w.WriteHeader(201)
	}

}
func PutKey(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
		keydata := params.Get(":key")
	valuedata := params.Get(":value")

				portNum := strings.Split(r.Host, ":")[1]

		keyvalstruct := DataSet{Key: keydata, Value: valuedata}
			dataInstances[portNum][keydata] = keyvalstruct
	if valuedata == dataInstances[portNum][keydata].Value {
		w.WriteHeader(204)
	} else {
		w.WriteHeader(404)
	}
}