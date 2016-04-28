package main

import (
	"hash/crc32"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type HashingCir struct {
	portnum string
	hash    int
}

type Circle []HashingCir //circle creation

func (a Circle) Len() int           { return len(a) }
func (a Circle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Circle) Less(i, j int) bool { return a[i].hash < a[j].hash }
func HashCode(value string) int {
	returnHash := crc32.ChecksumIEEE([]byte(value))
	return int(float32(returnHash))
}
func main() {
	ServerPortNumber := make([]HashingCir, 0)
	inputs := os.Args[1:]
//port openings
	ports := strings.Split(inputs[0], "-")
	startport, _ := strconv.Atoi(ports[0])
	endport, _ := strconv.Atoi(ports[1])

	keyvaluedata := strings.Split(inputs[1], ",")

	numberofports := endport - startport + 1

	for index := startport; index <= endport; index++ {
		currentport := strconv.Itoa(index)
		tempnode := HashingCir{portnum: currentport, hash: ConsistentHash(currentport, numberofports)}
		ServerPortNumber = append(ServerPortNumber, tempnode)
	}
	sort.Sort(Circle(ServerPortNumber))
//arrange in circel
	x := len(keyvaluedata)
	for index := 0; index < x; index++ {

		values := strings.Split(keyvaluedata[index], "->")
		keyvalue := values[0]
		datavalue := values[1]

		//Call Function
		keyhash := ConsistentHash(keyvalue, numberofports)
		callport := ""
		for index := 0; index < len(ServerPortNumber); index++ {
			if keyhash <= ServerPortNumber[index].hash {
				callport = ServerPortNumber[index].portnum
				break
			}
		}
		if callport == "" {
			callport = ServerPortNumber[0].portnum
		}

		//Call Ends
		url := "http://localhost:" + callport + "/" + keyvalue + "/" + datavalue
		log.Print(url)
		client := &http.Client{}
		req, err := http.NewRequest("PUT", url, nil)
		response, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		if response != nil {
			log.Print(" yeah!!!")
		} else {
			log.Print(" unsucessful")
		}
	}
}



func ConsistentHash(value string, numberofports int) int {
	consistenthash := (HashCode(value)) / (numberofports * 5000)
	return consistenthash
}