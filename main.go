package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
)

// Config object containing all configurations set in bitrise
type Config struct {
	Debug bool `env:"is_debug_mode,opt[yes,no]"`

	// Message
	WebhookURL        stepconf.Secret `env:"webhook_url"`
	Message           string          `env:"message"`
	MessageOnError    string          `env:"message_on_error"`
	Title             string          `env:"title"`
	TitleOnError      string          `env:"title_on_error"`
	Subtitle          string          `env:"subtitle"`
	SubtitleOnError   string          `env:"subtitle_on_error"`
	ImageURL          string          `env:"image"`
	ImageURLOnError   string          `env:"image_on_error"`
	ImageStyle        string          `env:"image_style,opt[square,circular]"`
	ImageStyleOnError string          `env:"image_style_on_error,opt[square,circular]"`
	Text              string          `env:"text"`
	TextOnError       string          `env:"text_on_error"`
	KeyValue          string          `env:"key_value"`
	KeyValueOnError   string          `env:"key_value_on_error"`
	Buttons           string          `env:"buttons"`
	ButtonsOnError    string          `env:"buttons_on_error"`

	ConvertSimpleToAvancedFormat bool `env:"convert_simple_to_advanced_format,opt[yes,no]"`
	ConvertAvancedToSimpleFormat bool `env:"convert_advanced_to_simple_format,opt[yes,no]"`
}

// success is true if the build is successful, false otherwise.
var success = os.Getenv("BITRISE_BUILD_STATUS") == "0"

func selectValue(ifSuccess, ifFailed string) string {
	if success || ifFailed == "" {
		return ifSuccess
	}
	return ifFailed
}

func selectAvancedFormatValue(ifSuccess, ifFailed string, simpleToAvancedFormat bool) string {
	selected := selectValue(ifSuccess, ifFailed)

	if simpleToAvancedFormat {
		selected = SimpleToAdvancedFormatting(selected)
	}

	return selected
}

func selectSimpleFormatValue(ifSuccess, ifFailed string, advancedToSimpleFormat bool) string {
	selected := selectValue(ifSuccess, ifFailed)

	if advancedToSimpleFormat {
		selected = AdvancedToSimpleFormatting(selected)
	}

	return selected
}

func newMessage(c Config) (msg Message, err error) {
	sections := []Section{}

	text := selectAvancedFormatValue(c.Text, c.TextOnError, c.ConvertSimpleToAvancedFormat)
	if text != "" {
		sections = append(sections, Section{
			Widgets: []*Widget{{
				TextParagraph: &TextParagraph{
					Text: text,
				},
			}},
		})
	}

	keyValueConfig := selectValue(c.KeyValue, c.KeyValueOnError)
	if keyValueConfig != "" {
		var keyValueWidgets []*Widget
		keyValueWidgets, err = ParseKeyValues(keyValueConfig)

		if err != nil {
			return
		}

		if keyValueWidgets != nil {
			sections = append(sections, Section{
				Widgets: keyValueWidgets,
			})
		}
	}

	buttonConfig := selectValue(c.Buttons, c.ButtonsOnError)
	if buttonConfig != "" {
		var buttons []*Button
		buttons, err = parseButtons(buttonConfig)

		if err != nil {
			return
		}

		sections = append(sections, Section{
			Widgets: []*Widget{{
				Buttons: buttons,
			}},
		})
	}

	message := selectSimpleFormatValue(c.Message, c.MessageOnError, c.ConvertAvancedToSimpleFormat)
	if message == "" {
		message = selectSimpleFormatValue(c.Title, c.TitleOnError, c.ConvertAvancedToSimpleFormat)
	}
	if message == "" {
		message = selectSimpleFormatValue(c.Text, c.TextOnError, c.ConvertAvancedToSimpleFormat)
	}

	msg = Message{
		Text: message,
		Cards: []Card{{
			Header: CreateHeader(
				selectAvancedFormatValue(c.Title, c.TitleOnError, c.ConvertSimpleToAvancedFormat),
				selectAvancedFormatValue(c.Subtitle, c.SubtitleOnError, c.ConvertSimpleToAvancedFormat),
				selectValue(c.ImageURL, c.ImageURLOnError),
				selectValue(c.ImageStyle, c.ImageStyleOnError),
			),
			Sections: sections,
		}},
	}

	return
}

// postMessage sends a message to a channel.
func postMessage(conf Config, msg Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Debugf("Request to Google Chat: %s\n", b)

	url := string(conf.WebhookURL)

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("server error: %s, failed to read response: %s", resp.Status, err)
		}
		return fmt.Errorf("server error: %s, response: %s", resp.Status, body)
	}

	return nil
}

func validate(conf *Config) error {
	if conf.WebhookURL == "" {
		return fmt.Errorf("WebhookURL is empty. You need to provide one")
	}

	if conf.Text == "" && conf.Buttons == "" && conf.KeyValue == "" {
		return fmt.Errorf("Text, keyValue and buttons are empty. You need to provide at least one")
	}

	return nil
}

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)
	log.SetEnableDebugLog(conf.Debug)

	if err := validate(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}

	msg, err := newMessage(conf)
	if err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	if err := postMessage(conf, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nGoogle Chat message successfully sent! 🚀\n")

	os.Exit(0)
}
