package usecases

import (
	"encoding/json"
	"os"
	"sort"
	"time"

	"gitlab.com/mims-api-service/helpers"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/src/discord/domains"
	"go.mongodb.org/mongo-driver/bson"
)

type Usecase struct {
	Repository domains.Repository
}

func NewUsecase(repo domains.Repository) domains.UseCase {
	return &Usecase{Repository: repo}
}

var RED int = 13370114
var GREEN int = 183299

var ERROR string = "ERROR : "
var SUCCESS string = "SUCCESS : "

var TIME string = "Time : "
var MESSAGE string = "Message : "
var Data string = "Data : "

func (u *Usecase) GetHrisAndHsmsLog() error {

	hrisType := []string{"Section Geom", "Road Latest", "Road Section (Match Data)", "Road Group (Match Data)"}
	hsmsType := []string{"Hsms Bridge", "Hsms Guard", "Hsms Interchange", "Hsms Intersection", "Hsms Streetlight", "Hsms Railwaycrossing", "Hsms Trafficlight", "Hsms Uturnbridge"}

	var discord models.Discord
	discord.Username = "BOT (" + Env() + ")"

	var embeds []models.Embeds
	var fieldHriss []models.Fields
	var embedHris models.Embeds
	for _, item := range hrisType {

		var level1 bson.A

		level1 = append(level1, bson.M{"title": item})

		filter := bson.D{
			{"$and", level1}}

		data, err := u.Repository.GetHrisAndHsmsLogWithFilter(filter, "MONGODB_HRIS_ERROR_LOG")
		if err != nil {
			return err
		}

		if len(data) > 0 {

			sort.Slice(data, func(i, j int) bool {
				return data[i].DateTime.After(data[j].DateTime)
			})

			value := data[0]

			textName := ""
			if value.IsSuccess {
				textName = "🟢🟢🟢   SUCCESS   🟢🟢🟢"
			} else {
				textName = "🔴🔴🔴   ERROR   🔴🔴🔴"
			}

			textName += "\nประเภท : " + item
			textName += "\nเวลา : " + value.DateTime.Add(7*time.Hour).Format(time.RFC822) + " น."

			message := value.Message
			message += "\n----------------------------------------------"

			var field models.Fields
			field.Inline = false
			field.Name = textName
			field.Value = message
			fieldHriss = append(fieldHriss, field)

		}
	}

	embedHris.Title = "HRIS"
	embedHris.Fields = fieldHriss

	embeds = append(embeds, embedHris)
	var fieldHsmss []models.Fields
	var embedHsms models.Embeds
	for _, item := range hsmsType {

		var level1 bson.A

		level1 = append(level1, bson.M{"title": item})

		filter := bson.D{
			{"$and", level1}}

		data, err := u.Repository.GetHrisAndHsmsLogWithFilter(filter, "MONGODB_HSMS_ERROR_LOG")
		if err != nil {
			return err
		}

		if len(data) > 0 {

			sort.Slice(data, func(i, j int) bool {
				return data[i].DateTime.After(data[j].DateTime)
			})

			value := data[0]

			textName := ""
			if value.IsSuccess {
				textName = "🟢🟢🟢   SUCCESS   🟢🟢🟢"
			} else {
				textName = "🔴🔴🔴   ERROR   🔴🔴🔴"
			}

			textName += "\nประเภท : " + item
			textName += "\nเวลา : " + value.DateTime.Add(7*time.Hour).Format(time.RFC822) + " น."

			message := value.Message
			message += "\n----------------------------------------------"

			var field models.Fields
			field.Inline = false
			field.Name = textName
			field.Value = message
			fieldHsmss = append(fieldHsmss, field)

		}
	}

	embedHsms.Title = "HSMS"
	embedHsms.Fields = fieldHsmss

	embeds = append(embeds, embedHsms)

	discord.Embeds = embeds

	payload, err := json.Marshal(discord)
	if err != nil {
		return err
	}

	url := os.Getenv("DISCORD_WEBHOOK_URL")
	method := "POST"
	_, err = helpers.RequestWithJsonBody(url, method, string(payload))
	if err != nil {
		return err
	}

	return nil
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
