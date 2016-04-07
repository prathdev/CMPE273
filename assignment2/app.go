package main

import (
	"encoding/json"
	"fmt"
	"github.com/drone/routes"
	"github.com/mkilling/goejdb"
	"github.com/naoina/toml"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"strconv"
    "net/rpc"
    "net/url"
)

type Listnr int
var hashDB map[string]*rpc.Client //replicate systems
var Config TomlConfig
var Jb *goejdb.Ejdb //remove and try
var Collct *goejdb.EjColl

type User struct {
	Email          string                       `json:"email"`
	Zip            string                       `json:"zip"`
	Country        string                       `json:"country"`
	Profession     string                       `json:"profession"`
	Favorite_color string                       `json:"favorite_color"`
	Is_smoking     string                       `json:"is_smoking"`
	Favorite_sport string                       `json:"favorite_sport"`
	Food           map[string]string            `json:"food"`
	Music          map[string]string            `json:"music"`
	Movie          map[string][]string          `json:"movie"`
	Travel         map[string]map[string]string `json:"travel"`
}

type TomlConfig struct {
	Database struct {
		filename string
		PortNo  int
	}

	Replication struct {
		ServerPort int
		ReplicaSys  []string
	}
}

func main() {
	arg := os.Args[1]  //arg-toml
	f, err := os.Open(arg)
	if err != nil {
		fmt.Print(err)
	}
	defer f.Close()
	buf, errIO := ioutil.ReadAll(f)
	if errIO != nil {
		fmt.Print(errIO)
	}
	if errMar := toml.Unmarshal(buf, &Config); err != nil {
		fmt.Print(errMar)
	}
	Jb, err = goejdb.Open(Config.Database.filename, goejdb.JBOWRITER|goejdb.JBOCREAT|goejdb.JBOREADER)
	if err != nil {
		fmt.Print("")
	}
	/*if err != nil {
		os.Exit(1)
	} */
	
	//Jb.Wait()
	defer Jb.Close() 

	Collct, err = Jb.GetColl("customerbook")
	if err != nil {
		Collct, err = Jb.CreateColl("customerbook", nil)
	}
	mux := routes.New()
	mux.Get("/profile/:email", GetProfile)
	mux.Post("/profile", PostProfile)
	mux.Put("/profile/:email", UpdateProfile)
	mux.Del("/profile/:email", DeleteProfile)
	log.Println("Http Started on", Config.Database.PortNo)
    go RPCListen()
	port := ":" + strconv.Itoa(Config.Database.PortNo)
	http.ListenAndServe(port, mux)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get(":email")
	res, _ := Collct.Find(`{"email" :` + string(email) + `}`)
	if len(res) > 0 {
		var m map[string]interface{}
		bson.Unmarshal(res[0], &m)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(m)
	} else {
		w.WriteHeader(404)
	}
}

func PostProfile(w http.ResponseWriter, r *http.Request) {
	var person User
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		fmt.Println(err)
	}
	found := false
	res, _ := Collct.Find(`{"email" :` + string(person.Email) + `}`)
	if len(res) > 0 {
		found = true
	}
	valid := true
	if found == false {
		if person.Is_smoking != "yes" && person.Is_smoking != "no" {
			w.WriteHeader(417)
			valid = false
		}
		if person.Food["type"] != "vegetarian" && person.Food["type"] != "meat_eater" && person.Food["type"] != "eat_everything" {
			w.WriteHeader(417)
			valid = false
		}
		if person.Food["drink_alcohol"] != "yes" && person.Food["drink_alcohol"] != "no" {
			w.WriteHeader(417)
			valid = false
		}
		if person.Travel["flight"]["seat"] != "aisle" && person.Travel["flight"]["seat"] != "window" {
			w.WriteHeader(417)
			valid = false
		}
		if valid == true {
			bsrec, _ := bson.Marshal(person)
			Collct.SaveBson(bsrec)
            AddtoReplica(person)
			w.WriteHeader(201)
		}
	} else {
		w.WriteHeader(404)
	}
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get(":email")
	found := false
	res, _ := Collct.Find(`{"email" :` + string(email) + `}`)
	var v User
	if len(res) > 0 {
		found = true
		bson.Unmarshal(res[0], &v)
	}
	valid := true
	var update User
	err_json := json.NewDecoder(r.Body).Decode(&update)
	if err_json != nil {
		fmt.Println(err_json)
	}
	if found == true {
		if update.Zip != "" {
			v.Zip = update.Zip
		}
		if update.Country != "" {
			v.Country = update.Country
		}
		if update.Profession != "" {
			v.Profession = update.Profession
		}
		if update.Favorite_color != "" {
			v.Favorite_color = update.Favorite_color
		}
		if update.Is_smoking != "" {
			v.Is_smoking = update.Is_smoking
			if v.Is_smoking != "yes" && v.Is_smoking != "no" {
				w.WriteHeader(417)
				valid = false
			}
		}
		if update.Favorite_sport != "" {
			v.Favorite_sport = update.Favorite_sport
		}
		if update.Food != nil {
			v.Food = update.Food
			if v.Food["type"] != "vegetarian" && v.Food["type"] != "meat_eater" && v.Food["type"] != "eat_everything" {
				w.WriteHeader(417)
				valid = false
			}
			if v.Food["drink_alcohol"] != "yes" && v.Food["drink_alcohol"] != "no" {
				w.WriteHeader(417)
				valid = false
			}
		}
		if update.Music != nil {
			v.Music = update.Music
		}
		if update.Movie != nil {
			v.Movie = update.Movie
		}
		if update.Travel != nil {
			v.Travel = update.Travel
			if v.Travel["flight"]["seat"] != "aisle" && v.Travel["flight"]["seat"] != "window" {
				w.WriteHeader(417)
				valid = false
			}
		}
		if valid == true {
			q := map[string]interface{}{
				"email": email,
				"$set":  v}
			qjson, _ := json.Marshal(q)
			Collct.Update(string(qjson))
            DBSysUpdate(v)
			w.WriteHeader(204)
		} else {
			w.WriteHeader(404)
		}
	} else {
		w.WriteHeader(404)
	}
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get(":email")
	q := map[string]interface{}{
		"email":    email,
		"$dropall": true}
	qjson, _ := json.Marshal(q)
	res, _ := Collct.Find(`{"email" :` + string(email) + `}`)
	var v User
	if len(res) > 0 {
		bson.Unmarshal(res[0], &v)
	}
    Collct.Update(string(qjson))
    CopySystem(v)
	w.WriteHeader(204)
}

func (l *Listnr) Update(pref User, reply *bool) error {
	q := map[string]interface{}{
		"email": pref.Email,
		"$set":  pref}
	qjson, _ := json.Marshal(q)
	Collct.Update(string(qjson))
	return nil
}

func (l *Listnr) Delete(pref User, reply *bool) error {
	q := map[string]interface{}{
		"email":    pref.Email,
		"$dropall": true}
	qjson, _ := json.Marshal(q)
    Collct.Update(string(qjson))
	return nil
}

func (l *Listnr) Insert(pref User, reply *bool) error {
	bsrec, _ := bson.Marshal(pref)
	Collct.SaveBson(bsrec)
	return nil
}

func RPCListen() {
	address := fmt.Sprintf(":%d", Config.Replication.ServerPort)
	listener := new(Listnr)
	rpc.Register(listener)
	log.Println("RPC Started on", Config.Replication.ServerPort)
	rpc.HandleHTTP()
	http.ListenAndServe(address, nil)
}

func DBSysUpdate(pref User) {
	setDBReplica()
	var reply bool
	for _, singleDB := range hashDB {
		err := singleDB.Call("Listnr.Update", pref, &reply)
		if err != nil {
			log.Println(err)
		}
	}

}

func CloseRPConnection() {
	for _, singleDB := range hashDB {
		singleDB.Close()
	}
}

func AddtoReplica(pref User) {
	setDBReplica()
	var reply bool
	for _, singleDB := range hashDB {
		err := singleDB.Call("Listnr.Insert", pref, &reply)
		if err != nil {
			log.Println(err)
		}
	}
}

func CopySystem(pref User) {
	setDBReplica()
	var reply bool
	for _, singleDB := range hashDB {
		err := singleDB.Call("Listnr.Delete", pref, &reply)
		if err != nil {
			log.Println(err)
		}
	}
}

func setDBReplica() {
	if len(hashDB) == 0 {
		hashDB = make(map[string]*rpc.Client)
		for _, singleDB := range Config.Replication.ReplicaSys {
			url, err := url.Parse(singleDB)
			if err != nil {
				log.Fatal("Cannot parse URL", err)
			}
			hashDB[url.Host], err = rpc.DialHTTP("tcp", url.Host)
		}
	}
}