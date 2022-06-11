package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/adlio/trello"
	"github.com/go-pdf/fpdf"
)

// Struct to load configuration from config.json file
type Configuration struct {
	Url		string	`json:"Url"`
	ApiKey string	`json:"ApiKey"`
	Token	string	`json:"Token"`
	BoardId string	`json:"BoardId"`
	OutputLocation string	`json:"OutputLocation"`
}

// Struct to load Card body to run functions
type ParseCardBody struct {
	Card *trello.Card
	StartingTime string
}

// Struct to return information from card body
type ParseCardReturn struct {
	CardTitle string
	CardDesc string
	ColorR int
	ColorG int
	ColorB int
	StartTime string
	EndTime string
}

// Function to find the ID for list based on the name of the List
func QueryListId(client *trello.Client, boardId string, listName string) string{
	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {	log.Fatal(err)}
	lists, err := board.GetLists(trello.Defaults())
	var listId string
	for key, element := range lists{
		if element.Name == listName{
			listId = lists[key].ID
		}
	}
	if len(listId) < 3 {log.Fatal("Could not find list named `", listName, "`!!")}
	return listId
}

// Function to return all cards from a specified list
func QueryCards(client *trello.Client, listId string) []*trello.Card {
	list, err := client.GetList(listId, trello.Defaults()) 
	if err != nil {log.Fatal(err)}
	cards, err := list.GetCards(trello.Defaults())
	if err != nil {log.Fatal(err)}
	return cards
}

// Function to calculate the date for the upcoming Sunday to use in Titles
func FindSundayDate(cards []*trello.Card) *time.Time {
	for _, card := range cards {
		if card.Due != nil{
			return card.Due
		}
	}
	var weekday = int(time.Now().Weekday())
	nextSunday := time.Now().AddDate(0, 0, 7-weekday)
	return &nextSunday
}

// Function to parse a provided card and return information required for printing
func ParseCards(CardParseBody ParseCardBody) ParseCardReturn {
	
	// Finding the correct card title
	// Load Regex for finding the duration of the card
	cardTitleRegex := regexp.MustCompile(`^\s*(\d?\d:\d\d)`)
	// Generate the corrected card title
	var cardTitle string
	if !cardTitleRegex.MatchString(CardParseBody.Card.Name) {
		cardTitle = CardParseBody.Card.Name
	} else {
		cardTitle = cardTitleRegex.ReplaceAllString(CardParseBody.Card.Name, "")
	}
	
	// Cleaning up the card description
	// Load Regex for checking if card description has CCLI line
	cardDescRegex := regexp.MustCompile(`^-\s?CCLI\s?-\s?\d*`)
	// Generate the cleaned card description
	var cardDesc string
	if !cardDescRegex.MatchString(CardParseBody.Card.Desc) {
		cardDesc = CardParseBody.Card.Desc
	} else {
		cardDesc = cardDescRegex.ReplaceAllString(CardParseBody.Card.Desc, "")
	}

	// Converting the start time for the card into an array
	timeSliceStr := strings.Split(CardParseBody.StartingTime, ":")
	var timeSliceInt [2]int
	timeSliceInt[0], _ = strconv.Atoi(timeSliceStr[0])
	timeSliceInt[1], _ = strconv.Atoi(timeSliceStr[1])

	// Retrieving duration from card
	duration := "00:00"
	if !cardTitleRegex.MatchString(CardParseBody.Card.Name) {
		duration = "05:00"
	} else {
		duration = CardParseBody.Card.Name[0:5]
	}
	durationSliceStr := strings.Split(duration, ":")
	var durationSliceInt [2]int
	durationSliceInt[0], _ = strconv.Atoi(durationSliceStr[0])
	durationSliceInt[1], _ = strconv.Atoi(durationSliceStr[1])

	// Computing end time from the starting time + duration
	modHours := 0
	if timeSliceInt[1]>60 {modHours = timeSliceInt[1]%60}
	modHours += timeSliceInt[0]
	modMinutes := timeSliceInt[1] + durationSliceInt[0] + durationSliceInt[1]/60
	if modMinutes/60 > 0 {
		modHours += modMinutes/60
		modMinutes = modMinutes%60
	}

	// Converting hours and minutes into a string for printing
	var modHoursStr string
	if modHours < 10 {
		modHoursStr = "0" + strconv.Itoa(modHours)
	} else {
		modHoursStr = strconv.Itoa(modHours)
	}
	var modMinutesStr string
	if modMinutes < 10 {
		modMinutesStr = "0"+strconv.Itoa(modMinutes)
	} else{
		modMinutesStr = strconv.Itoa(modMinutes)
	}

	// Concatenating Hours and Minutes strings to obtain endTime
	endTime := modHoursStr+":"+modMinutesStr

	// Retrieve font colours based on card labels to print
	colorR := 0
	colorG := 0
	colorB := 0
	if len(CardParseBody.Card.Labels) > 0{
		cardTitle = cardTitle + " | " + CardParseBody.Card.Labels[0].Name
		switch cardLabel := CardParseBody.Card.Labels[0].Name; cardLabel {
		case "Song":
			colorR = 17
			colorG = 80
			colorB = 185
		case "DRIC songs":
			colorR = 17
			colorG = 80
			colorB = 185
		case "Announcements":
			colorR = 23
			colorG = 111
			colorB = 39
		default:
			colorR = 0
			colorG = 0
			colorB = 0
		}
	} 

	return ParseCardReturn{CardTitle: cardTitle, CardDesc: cardDesc, ColorR: colorR, ColorG: colorG, ColorB: colorB, StartTime: CardParseBody.StartingTime, EndTime: endTime}
}

// The main execution function
func main() {
	startTimePtr := flag.Int("startTime", 1100, "Starting time of the Service as a 4 digit number, eg: 10:30 starting time should be entered as 1030. The default set is 1100")
	forceSizePtr := flag.Bool("forceSize", false, "Force script to develop the PDF with the provided magFactor and descMagFactor")
	magFactorPtr := flag.Float64("magFactor", 1.6, "Magnification factor of the text except descriptions. This will default to 1.6")
	descMagFactorPtr := flag.Float64("descMagFactor", 1.6, "Magnification factor for the description, this will only affect the description added under a card. This will default to 1.6")
	listNamePtr := flag.String("listName", "Sunday Service", "List name can be selected if there are multiple services, if not specified then the default value is Sunday Service")
	flag.Parse()

	for {
		// Initiate Fpdf writer
		Fpdf := fpdf.New("P", "mm", "A4", "")

		// Load configuration from config file
		configuration := Configuration{}
		file,err := os.Open("config.json") 
		if err != nil {
			fmt.Println("Please ensure config.json file exists in the same folder as the script in the original format documented")
			log.Fatal(err)
		}
		byteValue, _ := ioutil.ReadAll(file)
		json.Unmarshal(byteValue,&configuration)
		

		// List current printing options
		fmt.Println("\n Developing Sunday Schedule with the following parameters:")
		fmt.Println("\tstartTime:", *startTimePtr)
		fmt.Println("\tmagFactor:", *magFactorPtr)
		fmt.Println("\tdescMagFactor:", *descMagFactorPtr)
		fmt.Println("\tlistName:", *listNamePtr)

		client := trello.NewClient(configuration.ApiKey, configuration.Token)
		var listId = QueryListId(client, configuration.BoardId, *listNamePtr)
		fmt.Println("\nRetrieving cards from Trello board...")
		var cards = QueryCards(client, listId)
		nextSundayDate := FindSundayDate(cards)
		serviceName := *listNamePtr + ": " + nextSundayDate.Format("January 2, 2006")
		
		CardStartTime := strconv.Itoa(*startTimePtr)
		CardStartTime = CardStartTime[:len(CardStartTime)-2] + ":" + CardStartTime[len(CardStartTime)-2:]

		// Set FPDF initial parameters
		Fpdf.SetMargins(20, 10, 15)
		Fpdf.SetFont("Arial", "", 20)
		Fpdf.SetAutoPageBreak(true, 15)
		Fpdf.SetTextColor(37, 77, 145)
		Fpdf.AddPage()
		
		// Add Damascus Road logo to service order
		imgOptions := fpdf.ImageOptions{
			AllowNegativePosition: false,
			ReadDpi:   true,
			ImageType: "PNG",
		}
		Fpdf.ImageOptions("./DRIC Full Logo.png", 150, 10, 45, 10, true, imgOptions, 0, "")
		
		// Develop PDF based on card information
		Fpdf.Write(-5, serviceName)
		Fpdf.SetFont("Arial", "", 12*(*magFactorPtr))
		Fpdf.SetTextColor(0, 0, 0)
		for _, card := range cards {
			if card.Name != ""{
				fmt.Println("\t... Printing out card for", card.Name)
				CardReturnBody := ParseCards(ParseCardBody{Card: card, StartingTime: CardStartTime})
				Fpdf.SetTextColor(CardReturnBody.ColorR, CardReturnBody.ColorG, CardReturnBody.ColorB)
				Fpdf.SetLeftMargin(15)
				Fpdf.Write(8*(*magFactorPtr), "\n")
				Fpdf.Write(8*(*magFactorPtr), CardStartTime + " " + CardReturnBody.CardTitle)
				if len(CardReturnBody.CardDesc) > 0{
					Fpdf.SetLeftMargin(32)
					Fpdf.SetFont("Arial", "I", 10*(*descMagFactorPtr))
					Fpdf.Write(8*(*magFactorPtr),"\n")
					Fpdf.Write(5*(*descMagFactorPtr),CardReturnBody.CardDesc)
					Fpdf.SetFont("Arial", "", 12*(*magFactorPtr))
				}
				CardStartTime = CardReturnBody.EndTime
			}
		}

		// Finding OutputLocation to print out file to depending on the OS
		var OutputFileName string
		UserOS := runtime.GOOS
		UserHomeDirectory, _ := os.UserHomeDir()
		switch UserOS{
			case "windows":
				if configuration.OutputLocation == ""{
					OutputFileName = UserHomeDirectory + `\Desktop\` + `\Sunday Service Schedule - ` + nextSundayDate.Format("January 2, 2006") + ".pdf"
				} else{
					OutputFileName = configuration.OutputLocation + `\Sunday Service Schedule - ` + nextSundayDate.Format("January 2, 2006") + ".pdf"
				}
			case "darwin":
				if configuration.OutputLocation == ""{
					OutputFileName = UserHomeDirectory + `/Desktop/` + "/Sunday Service Schedule - " + nextSundayDate.Format("January 2, 2006") + ".pdf"
				} else{
					OutputFileName = configuration.OutputLocation + "/Sunday Service Schedule - " + nextSundayDate.Format("January 2, 2006") + ".pdf"
				}
			case "linux":
				if configuration.OutputLocation == ""{
					OutputFileName = UserHomeDirectory + `/Desktop/` + "/Sunday Service Schedule - " + nextSundayDate.Format("January 2, 2006") + ".pdf"
				} else{
					OutputFileName = configuration.OutputLocation + "/Sunday Service Schedule - " + nextSundayDate.Format("January 2, 2006") + ".pdf"
				}
		}

		// Compute page numbers for PDF to be developed
		fmt.Println(Fpdf.PageNo())
		if Fpdf.PageNo() == 1 || *forceSizePtr {
			// Developing PDF file
			err := Fpdf.OutputFileAndClose(OutputFileName)
			if err != nil {log.Fatal(err)}
			break
		} else {
			if *magFactorPtr < 0.7 {
				fmt.Println("Too many lines in the Trello list, or run the script with the 'forceSize' parameter set to 'true'")
			}
			fmt.Println("*** Developed service order has more than one page and so going to try again by reducing font size... ***")
			// Decrement magFactors until it is only one page
			*magFactorPtr -= 0.05
			*descMagFactorPtr -= 0.05
		}
	}
}