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
	Title             string          `env:"title"`
	TitleOnError      string          `env:"title_on_error"`
	Subtitle          string          `env:"subtitle"`
	SubtitleOnError   string          `env:"subtitle_on_error"`
	ImageURL          string          `env:"image"`
	ImageURLOnError   string          `env:"image_on_error"`
	ImageStyle        string          `env:"image_style,opt[IMAGE,AVATAR]"`
	ImageStyleOnError string          `env:"image_style_on_error,opt[IMAGE,AVATAR]"`
	Text              string          `env:"text"`
	TextOnError       string          `env:"text_on_error"`
}

// success is true if the build is successful, false otherwise.
var success = os.Getenv("BITRISE_BUILD_STATUS") == "0"

func selectValue(ifSuccess, ifFailed string) string {
	if success || ifFailed == "" {
		return ifSuccess
	}
	return ifFailed
}

func newMessage(c Config) Message {
	msg := Message{
		Cards: []Card{{
			Header: Header{
				Title:      selectValue(c.Title, c.TitleOnError),
				Subtitle:   selectValue(c.Subtitle, c.SubtitleOnError),
				ImageURL:   selectValue(c.ImageURL, c.ImageURLOnError),
				ImageStyle: selectValue(c.ImageStyle, c.ImageStyleOnError),
			},
			Sections: []Section{{
				Widgets: []Widget{{
					TextParagraph: TextParagraph{
						Text: selectValue(c.Text, c.TextOnError),
					},
				}, {
					Buttons: []Button{{
						TextButton: TextButton{
							Text: "example website",
							OnClick: OnClick{
								OpenLink: OpenLink{
									URL: "https://example.com",
								},
							},
						},
					}, {
						ImageButton: ImageButton{
							Icon: "CLOCK",
							OnClick: OnClick{
								OpenLink: OpenLink{
									URL: "https://example.org",
								},
							},
						},
					}, {
						ImageButton: ImageButton{
							IconURL: "https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg",
							OnClick: OnClick{
								OpenLink: OpenLink{
									URL: "https://bitrise.io/",
								},
							},
						},
					}},
				}},
			}},
		}},
	}

	return msg
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

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)
	log.SetEnableDebugLog(conf.Debug)

	msg := newMessage(conf)
	if err := postMessage(conf, msg); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nGoogle Chat message successfully sent! 🚀\n")

	os.Exit(0)
}
