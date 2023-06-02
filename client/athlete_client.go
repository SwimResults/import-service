package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/swimresults/athlete-service/model"
	"net/http"
)

var apiUrl = "http://localhost:8086/"

func addAthlete(athlete model.Athlete) {
	b, err := json.Marshal(athlete)
	if err != nil {
		return
	}
	r, err := http.NewRequest("POST", apiUrl+"athlete", bytes.NewBuffer(b))
	if err != nil {
		return
	}
	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	newAthlete := &model.Athlete{}
	err = json.NewDecoder(res.Body).Decode(newAthlete)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusCreated {
		panic(res.Status)
	}

	fmt.Println("Id:", newAthlete.Identifier)
	fmt.Println("Name:", newAthlete.Name)
	fmt.Println("Year:", newAthlete.Year)
	fmt.Println("DSV:", newAthlete.DsvId)
}
