package main

import (
	"bytes"
	"fmt"
	"flag"
	"os"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type config struct {
	Username string
	Password string
	Front string
	Back string
}

const (
	DIRECTMAILERS_URL = "https://print.directmailers.com/api/v1/postcard/"
)

func main() {
	config := setConfig()

	requestBody, err := json.Marshal(map[string]interface{}{
		"Front": config.Front,
		"Back": config.Back,
		"WaitForRender": "true",
		"Description": "Sample Request",
		"Size": "4.25x6",
		"DryRun": "true",
		"VariablePayload": map[string]string{
			"name": "Gizmo",
		},
		"To": map[string]string{
			"Name": "Jane Doe",
			"AddressLine1": "123 N Test",
			"City": "Apopka",
			"State": "FL",
			"Zip": "32703",
		},
		"From": map[string]string{
			"Name": "Me DUDE",
			"AddressLine1": "123 N Test",
			"City": "Apopka",
			"State": "FL",
			"Zip": "32703",
		},
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", DIRECTMAILERS_URL, bytes.NewBuffer(requestBody))
	req.Header.Add("Authorization", "Basic " + config.encodeAuth())
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var parsedResp map[string]interface{}

	err = json.Unmarshal(body, &parsedResp)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(parsedResp["RenderedPdf"])
}

func setConfig() config {
	username := flag.String("username", "", "Enter in username")
	password := flag.String("password", "", "Enter in password")
	front := flag.String("front", "", "Enter in front")
	back := flag.String("back", "", "Enter in back")


	flag.Parse()

	if *username == "" || username == nil {
		fmt.Println("Must pass in username flag")
		os.Exit(1)
	}

	if *password == "" || password == nil {
		fmt.Println("Must pass in password flag")
		os.Exit(1)
	}

	if *front == "" || front == nil {
		fmt.Println("Must pass in front flag")
		os.Exit(1)
	}

	if *back == "" || back == nil {
		fmt.Println("Must pass in back flag")
		os.Exit(1)
	}

	return config{
		Username: *username,
		Password: *password,
		Front: *front,
		Back: *back,
	}
}

func (c config) encodeAuth() string {
	auth := c.Username + ":" + c.Password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
