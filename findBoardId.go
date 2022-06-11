package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/adlio/trello"
)

type Configuration struct {
	Url		string	`json:"url"`
	ApiKey string	`json:"api_key"`
	Token	string	`json:"token"`
}


func main() {
	configuration := Configuration{}
	file,err := os.Open("config.json") 
	if err != nil {log.Fatal(err)}
	byteValue, _ := ioutil.ReadAll(file)
	json.Unmarshal(byteValue,&configuration)

	client := trello.NewClient(configuration.ApiKey, configuration.Token)

	member, err := client.GetMember("damascusroadic1", trello.Defaults())
	if err != nil {log.Fatal(err)}

	fmt.Println("Logged in as", member.FullName)

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {log.Fatal(err)}
	for _, board := range boards {
		if board.Name != ""{
			fmt.Printf("Board Name: %s\tBoard ID: %s", board.Name, board.ID)
		}
		fmt.Println(board.Name)
	}
}