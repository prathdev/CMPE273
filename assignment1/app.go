package main

import (
	"github.com/drone/routes"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Profile struct {
	Email string `json:"email"`
	Zip string `json:"zip"`
	Country string `json:"country"`
	Profession string	`json:"profession"`
	Favorite_color string `json:"favorite_color"`
	Is_smoking string	`json:"is_smoking"`
	Favorite_sport string `json:"favorite_sport"`
	Food Foodtype `json:"food"`
	Music Musictype `json:"music"`
	Movie Movietype `json:"movie"`
	Travel Traveltype `json:"travel"`
}

type Foodtype struct {
	Type string `json:"type"`
	Drink_alcohol string `json:"drink_alcohol"`
} 

type Musictype struct {
	Spotify_user_id string `json:"spotify_user_id"`
}

type Movietype struct{
	Tv_shows []string `json:"tv_shows"`
	Movies []string `json:"movies"`
}

type Traveltype struct{
	Flight Flighttype `json:"flight"`
}

type Flighttype struct{
	Seat string `json:"seat"`
}

var profile []Profile

func main() {
	mux := routes.New()

	mux.Get("/profile/:email", GetProfile)
	mux.Post("/profile", PostProfile)
	mux.Del("/profile/:email",RemoveProfile)
	mux.Put("/profile/:email",UpdateProfile)
	http.Handle("/", mux)
	log.Println("Listening...")
	http.ListenAndServe(":3000", nil)
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	email := params.Get(":email")
	position := isProfilePresent(email)
	if position < 0 {
		w.WriteHeader(404)
		return
	}
	resp, err := json.Marshal(profile[position])
	if err != nil {
		w.WriteHeader(505)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(resp))


}

func isProfilePresent(email string) int{

	for index , value := range profile {
		if value.Email == email{
			return index
		}
	}
	return -1
}

func PostProfile(w http.ResponseWriter, r *http.Request){
	
	var temp Profile
	receivedJson, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(receivedJson, &temp)
	if err != nil{
		http.Error(w,"Error Unmarshalling",http.StatusBadRequest)
		return
	}
	if isProfilePresent(temp.Email) > 0 {
		w.WriteHeader(422)
		return
	}
	profile = append(profile,temp)
	w.WriteHeader(http.StatusCreated)
}

func RemoveProfile(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	email := params.Get(":email")
	position := isProfilePresent(email)
	if position == -1 {
		w.WriteHeader(204)
		return
	}
	profile = append(profile[:position],profile[position+1:]...)
	w.WriteHeader(204)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	email := params.Get(":email")
	position := isProfilePresent(email)
	if position == -1 {
				w.WriteHeader(204)
		return
	}
	var temp map[string]interface{}
	receivedJson, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(receivedJson, &temp)

	fmt.Println(temp)

	if err != nil{
		http.Error(w,"Error Unmarshalling",http.StatusBadRequest)
		return
	}

	json_temp_profile,_ := json.Marshal(profile[position])

	var map_profile map[string]interface{}
	_ = json.Unmarshal(json_temp_profile, &map_profile)

	map_profile = editProf(temp, map_profile)

	json_temp_profile,_ = json.Marshal(map_profile)
	_ = json.Unmarshal(json_temp_profile, &profile[position])
		

	w.WriteHeader(204)
}

func editProf(updateDate map[string]interface{}, profileData map[string]interface{}) map[string]interface{}{
	for i := range updateDate {
		profileData[i] = updateDate[i]
	}
	return profileData
}
