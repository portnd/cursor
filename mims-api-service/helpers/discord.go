package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gitlab.com/mims-api-service/models"
)

var RED int = 13370114
var GREEN int = 183299

var ERROR string = "ERROR : "
var SUCCESS string = "SUCCESS : "

var TIME string = "Time : "
var MESSAGE string = "Message : "
var Data string = "Data : "

func DiscordHrisLog(title, message string, isSuccess bool) error {
	dataTime := time.Now().Format(time.RFC822)
	url := os.Getenv("DISCORD_WEBHOOK_URL")
	method := "POST"

	var embeds []models.Embeds
	var fields []models.Fields

	var discord models.Discord
	discord.Username = "BOT (" + Env() + ")"
	discord.Content = "!!!!  HRIS LOG !!!!"

	var embed models.Embeds
	if isSuccess {
		embed.Color = GREEN
		embed.Title = SUCCESS + title

	} else {
		embed.Color = RED
		embed.Title = ERROR + title
	}

	var fieldTime models.Fields
	fieldTime.Inline = false
	fieldTime.Name = TIME
	fieldTime.Value = dataTime + " น."
	fields = append(fields, fieldTime)

	var fieldMessage models.Fields
	fieldMessage.Inline = false
	fieldMessage.Name = MESSAGE
	fieldMessage.Value = message
	fields = append(fields, fieldMessage)

	embed.Fields = fields

	embeds = append(embeds, embed)

	discord.Embeds = embeds

	payload, err := json.Marshal(discord)
	if err != nil {
		return err
	}

	_, err = RequestWithJsonBody(url, method, string(payload))
	if err != nil {
		return err
	}

	return nil
}

func DiscordHrisMatchLog(title, message string, data string, isSuccess bool) error {
	dataTime := time.Now().Format(time.RFC822)
	url := os.Getenv("DISCORD_WEBHOOK_URL")
	method := "POST"

	var embeds []models.Embeds
	var fields []models.Fields

	var discord models.Discord
	discord.Username = "BOT (" + Env() + ")"
	discord.Content = "!!!!  HRIS LOG !!!!"

	var embed models.Embeds
	if isSuccess {
		embed.Color = GREEN
		embed.Title = SUCCESS + title

	} else {
		embed.Color = RED
		embed.Title = ERROR + title
	}

	var fieldTime models.Fields
	fieldTime.Inline = false
	fieldTime.Name = TIME
	fieldTime.Value = dataTime + " น."
	fields = append(fields, fieldTime)

	var fieldMessage models.Fields
	fieldMessage.Inline = false
	fieldMessage.Name = MESSAGE
	fieldMessage.Value = message
	fields = append(fields, fieldMessage)

	var fieldData models.Fields
	fieldData.Inline = false
	fieldData.Name = Data
	fieldData.Value = data
	fields = append(fields, fieldData)

	embed.Fields = fields

	embeds = append(embeds, embed)

	discord.Embeds = embeds

	payload, err := json.Marshal(discord)
	if err != nil {
		return err
	}

	_, err = RequestWithJsonBody(url, method, string(payload))
	if err != nil {
		return err
	}

	return nil
}

func RequestWithJsonBody(url, method, data string) (interface{}, error) {

	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return string(body), nil
}

func Env() string {

	switch env := os.Getenv("ENV"); env {
	case "dev":
		return "Develop"
	case "stag":
		return "Stage"
	case "prod":
		return "Production"
	}
	return ""
}
