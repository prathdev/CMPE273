package main

import (
	"os"
	"sort"
	"log"
	"hash/fnv"
	"strings"
	"strconv"
	"net/http"
	"fmt"
)

var ServerInstances []int
var ServerHostMap= make(map[int]string)
var SortedServerHostMap= make(map[int]string)
var HashKeytoServer= make(map[int]string)

//find the server that has max score of key to store it
func findMaxScoreServer(keyVal string)(server int){

	var maxScore int
	var maxNode int
	var score int

	for i:=0;i<len(ServerInstances);i++{
		score = storeKeyValtoServer(strconv.Itoa(ServerInstances[i]), keyVal)
		if score > maxScore {
			maxScore = score
			maxNode = ServerInstances[i]
		}
	}
	return maxNode;
	/*for i:=0;i<len(ServerInstances);i++{
		if keyVal<=ServerInstances[i]{
			return SortedServerHostMap[ServerInstances[i]]
		}
	}
	return SortedServerHostMap[ServerInstances[1]]*/
}

//Get the hash value for the string
func storeKeyValtoServer(ip string,strKey string) int {
	keyValuetoStore := fnv.New32a()
	keyValuetoStore.Write([]byte(strKey))
	keyValuetoStore.Write([]byte(ip))
	return int(keyValuetoStore.Sum32())
}

func main() {
	ports := strings.Split(os.Args[1],"-") //instantiate to all server ports
	fromPort,_ := strconv.Atoi(ports[0])
	toPort,_ := strconv.Atoi(ports[1])
	for i:=fromPort; i<=toPort ;i++ {
		ServerInstances = append(ServerInstances, i)
	}
	inputValues := strings.Split(os.Args[2],",")
	
	//get  keyvalue pairs from arguments
	for i:=0; i<len(inputValues) ;i++ {
		inputItem := strings.Split(inputValues[i],"->")
		inputItemKey,_ := strconv.Atoi(inputItem[0])
		HashKeytoServer[inputItemKey]=inputItem[1]
	}
	
	/*for ser := range ServerHostMap {
		ServerInstances = append(ServerInstances, ser)
	}	
	sort.Ints(ServerInstances)
	for _, ser := range ServerInstances {
		SortedServerHostMap[ser]=ServerHostMap[ser]
	}*/
	
	//put all keys to server
	var inputKeyList []int
	for key := range HashKeytoServer {
		inputKeyList = append(inputKeyList, key)
	}
	
	//sort keys and find the server with max score
	sort.Ints(inputKeyList)
	for _, key := range inputKeyList {
		keyStr:=strconv.Itoa(key)
		serverName:=strconv.Itoa(findMaxScoreServer(keyStr))
		serverPath :="http://localhost:"+serverName+"/"+keyStr+"/"+HashKeytoServer[key]
		fmt.Print(serverPath)
		clientObj := &http.Client{}
		requestObj, err := http.NewRequest("PUT", serverPath, nil)
		_,err= clientObj.Do(requestObj)
		if err != nil {
			log.Fatal(err)
		}
	}
}
