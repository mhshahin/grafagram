package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	TelegramBaseURL = "https://api.telegram.org"

	SendMessageCmd = "sendMessage"
)

type AlertPayload struct {
	Title       string `json:"title,omitempty"`
	RuleID      int64  `json:"ruleId,omitempty"`
	RuleName    string `json:"ruleName,omitempty"`
	State       string `json:"state,omitempty"`
	EvalMatches []struct {
		Value  float64     `json:"value,omitempty"`
		Metric string      `json:"metric,omitempty"`
		Tags   interface{} `json:"tags,omitempty"`
	} `json:"evalMatches,omitempty"`
	OrgID       int `json:"orgId,omitempty"`
	DashboardID int `json:"dashboardId,omitempty"`
	PanelID     int `json:"panelId,omitempty"`
	Tags        struct {
	} `json:"tags,omitempty"`
	RuleURL  string `json:"ruleUrl,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	Message  string `json:"message,omitempty"`
}

type Message struct {
	ChatId    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type TelegramPayload struct {
	RuleName string `json:"ruleName,omitempty"`
	State    string `json:"state,omitempty"`
	Message  string `json:"message,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	HasImage bool
}

func main() {
	e := echo.New()
	e.POST("/alert", alertHandler())

	e.Logger.Fatal(e.Start(":1323"))
}

func alertHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		var alertPayload = new(AlertPayload)

		err := c.Bind(alertPayload)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		caser := cases.Title(language.AmericanEnglish)

		tgPayload := &TelegramPayload{
			RuleName: alertPayload.RuleName,
			State:    caser.String(alertPayload.State),
			Message:  alertPayload.Message,
		}

		if len(alertPayload.ImageURL) > 0 {
			tgPayload.ImageURL = alertPayload.ImageURL
			tgPayload.HasImage = true
		}

		buf := new(bytes.Buffer)

		tmpl := template.Must(template.ParseFiles("alert-layout.html"))
		err = tmpl.Execute(buf, tgPayload)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		err = sendTelegramMessage(buf.String())
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "OK!"})
	}
}

func sendTelegramMessage(message string) error {
	chatID := os.Getenv("CHAT_ID")
	url := generateTelegramUrl(SendMessageCmd)

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(Message{
		ChatId:    chatID,
		Text:      message,
		ParseMode: "html",
	})
	if err != nil {
		return err
	}

	_, err = newRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	return nil
}

func generateTelegramUrl(cmd string) string {
	botToken := os.Getenv("BOT_TOKEN")

	url := fmt.Sprintf("%s/bot%s/%s", TelegramBaseURL, botToken, cmd)

	return url
}

func newRequest(method, url string, body io.Reader) (io.ReadCloser, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.Println(res.StatusCode)
		return nil, err
	}

	return res.Body, nil
}
