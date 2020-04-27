package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Message which can be send to Google Chat
// More info at https://developers.google.com/hangouts/chat/reference/message-formats/cards
type Message struct {
	Text  string `json:"text,omitempty"`
	Cards []Card `json:"cards"`
}

// Card property of a message. can contain a header and must have at least one section
type Card struct {
	// header object (optional)
	Header *Header `json:"header,omitempty"`
	// sections object. At least one section is required
	Sections []Section `json:"sections,omitempty"`
}

// Header property of a card
type Header struct {
	Title    string
	Subtitle string
	ImageURL string
	// imageStyle controls the shape of the header image, which may be "square" ("IMAGE") or "circular" ("AVATAR"). The default is "square" ("IMAGE").
	ImageStyle string
}

// MarshalJSON implements json.Marshaler.MarshalJSON.
func (h Header) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	if h.Title != "" {
		m["title"] = h.Title
	}

	if h.Subtitle != "" {
		m["subtitle"] = h.Subtitle
	}

	if h.ImageURL != "" {
		m["imageUrl"] = h.ImageURL

		if h.ImageStyle != "" {
			m["imageStyle"] = h.ImageStyle

			switch h.ImageStyle {
			case "circular":
				m["imageStyle"] = "AVATAR"
			case "square":
				fallthrough
			default:
				m["imageStyle"] = "IMAGE"
			}
		}
	}

	return json.Marshal(m)
}

// CreateHeader creates a Header struct if at least one field is present, return nil pointer otherwise
func CreateHeader(title, subtitle, imageURL, imageStyle string) *Header {
	if title == "" && subtitle == "" && imageURL == "" && imageStyle == "" {
		return nil
	}

	return &Header{
		Title:      title,
		Subtitle:   subtitle,
		ImageURL:   imageURL,
		ImageStyle: imageStyle,
	}
}

// Section of a card. can contain multiple widgets, but at least one is required. Sections are separated by a horizontal line
type Section struct {
	// Section header (optional)
	Header string `json:"header,omitempty"`
	// widgets object. At least one widget is required.
	Widgets []*Widget `json:"widgets,omitempty"`
}

// Widget of a section. Can contain only one type of UI element
type Widget struct {
	TextParagraph *TextParagraph `json:"textParagraph,omitempty"`
	KeyValue      *KeyValue      `json:"keyValue,omitempty"`
	Image         *Image         `json:"image,omitempty"`
	// buttons object can contain one or more buttons. will be laid out horizontally
	Buttons []*Button `json:"buttons,omitempty"`
}

// TextParagraph UI element
type TextParagraph struct {
	// The text to display inside the paragraph
	Text string `json:"text,omitempty"`
}

// KeyValue UI element
type KeyValue struct {
	TopLabel         string   `json:"topLabel,omitempty"`
	Content          string   `json:"content,omitempty"`
	ContentMultiline string   `json:"contentMultiline,omitempty"`
	BottomLabel      string   `json:"bottomLabel,omitempty"`
	OnClick          *OnClick `json:"onClick,omitempty"`
	// either iconUrl of icon can be used
	IconURL string `json:"iconUrl,omitempty"`
	// either iconUrl of icon can be used
	Icon   string  `json:"icon,omitempty"`
	Button *Button `json:"button,omitempty"`
}

// ParseKeyValues parses a simpler KeyValue json array to the KeyValue object used by the api
func ParseKeyValues(raw string) (widgets []*Widget, err error) {
	// ignore empty string
	if raw == "" {
		return
	}

	var keyValueInput []KeyValueInput

	err = json.Unmarshal([]byte(raw), &keyValueInput)
	if err != nil {
		return
	}

	// Throw error for empty array
	if len(keyValueInput) == 0 {
		// Ignore an empty array: It is a way to remove the keyValues in case of an error
		return
	}

	widgets = []*Widget{}

	for _, keyValue := range keyValueInput {
		if keyValue.Content == "" {
			err = fmt.Errorf("KeyValue Should have content")
			return
		}

		// either iconUrl of icon can be used
		if keyValue.IconURL != "" && keyValue.Icon != "" {
			err = fmt.Errorf("KeyValue object should have either an iconUrl, an icon, or neither, but not both")
			return
		}

		var onClick *OnClick
		if keyValue.OnClick != "" {
			onClick = &OnClick{
				OpenLink: &OpenLink{
					URL: keyValue.OnClick,
				},
			}
		}

		var button *Button
		if keyValue.Button != nil {
			if keyValue.Button.OnClick == "" {
				err = fmt.Errorf("KeyValue button should have an onClick value")
				return
			}

			if (keyValue.Button.Text != "" && (keyValue.Button.IconURL != "" || keyValue.Button.Icon != "")) || (keyValue.Button.IconURL != "" && (keyValue.Button.Text != "" || keyValue.Button.Icon != "")) {
				err = fmt.Errorf("KeyValue button should have either a text, an iconUrl, or an icon field, not multiple")
				return
			}

			var buttonOnClick = &OnClick{
				OpenLink: &OpenLink{
					URL: keyValue.Button.OnClick,
				},
			}

			var textButton *TextButton
			if keyValue.Button.Text != "" {
				textButton = &TextButton{
					Text:    keyValue.Button.Text,
					OnClick: buttonOnClick,
				}
			}

			var imageButton *ImageButton
			if keyValue.Button.Icon != "" || keyValue.Button.IconURL != "" {
				imageButton = &ImageButton{
					IconURL: keyValue.Button.IconURL,
					Icon:    keyValue.Button.Icon,
					OnClick: buttonOnClick,
				}
			}

			if textButton == nil && imageButton == nil {
				err = fmt.Errorf("KeyValue button should have either an iconUrl, an icon or text field")
				return
			}

			button = &Button{
				TextButton:  textButton,
				ImageButton: imageButton,
			}
		}

		widgets = append(widgets, &Widget{
			KeyValue: &KeyValue{
				TopLabel:         keyValue.TopLabel,
				Content:          keyValue.Content,
				ContentMultiline: strconv.FormatBool(keyValue.ContentMultiline),
				BottomLabel:      keyValue.BottomLabel,
				OnClick:          onClick,
				IconURL:          keyValue.IconURL,
				Icon:             keyValue.Icon,
				Button:           button,
			},
		})
	}

	return
}

// Image UI element
type Image struct {
	ImageURL string   `json:"imageUrl,omitempty"`
	OnClick  *OnClick `json:"onClick,omitempty"`
}

// Button UI element. Can contain either a TextButton or an ImageButton
type Button struct {
	TextButton  *TextButton  `json:"textButton,omitempty"`
	ImageButton *ImageButton `json:"imageButton,omitempty"`
}

// TextButton UI element
type TextButton struct {
	Text    string   `json:"text,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

// ImageButton UI element
type ImageButton struct {
	// either iconUrl of icon can be used
	IconURL string `json:"iconUrl,omitempty"`
	// either iconUrl of icon can be used
	Icon    string   `json:"icon,omitempty"`
	OnClick *OnClick `json:"onClick,omitempty"`
}

func parseButtons(s string) (buttons []*Button, err error) {
	var buttonConf [][3]string
	buttonConf, err = triples(s)
	if err != nil {
		return
	}

	for _, triple := range buttonConf {
		var button *Button

		button, err = parseButton(triple)
		if err != nil {
			return
		}

		buttons = append(buttons, button)
	}

	return
}

func parseButton(triple [3]string) (button *Button, err error) {
	onClick := &OnClick{
		OpenLink: &OpenLink{
			URL: triple[2],
		},
	}

	if triple[0] == "text" {
		button = &Button{
			TextButton: &TextButton{
				Text:    triple[1],
				OnClick: onClick,
			},
		}

	} else if triple[0] == "builtin-icon" {
		button = &Button{
			ImageButton: &ImageButton{
				Icon:    triple[1],
				OnClick: onClick,
			},
		}

	} else if triple[0] == "custom-icon" {
		button = &Button{
			ImageButton: &ImageButton{
				IconURL: triple[1],
				OnClick: onClick,
			},
		}

	} else {
		err = fmt.Errorf("Unknown button type %s", triple[0])
	}

	return
}

// pairs slices every lines in s into two substrings separated by the first pipe
// character and returns a slice of those pairs.
func triples(s string) (ps [][3]string, err error) {
	s = strings.TrimSpace(s)

	for _, line := range strings.Split(s, "\n") {
		var tripple [3]string

		tripple, err = splitTriple(line)

		if err != nil {
			return
		}

		ps = append(ps, tripple)
	}
	return
}

func splitTriple(s string) (triple [3]string, err error) {
	splitString := strings.SplitN(s, "|", 3)

	if len(splitString) == 3 && splitString[0] != "" && splitString[1] != "" && splitString[2] != "" {
		triple = [3]string{splitString[0], splitString[1], splitString[2]}

	} else {
		err = fmt.Errorf("Could not parse button with declaration %s", s)
	}

	return
}

// OnClick handler object
type OnClick struct {
	OpenLink *OpenLink `json:"openLink,omitempty"`
}

// OpenLink object
type OpenLink struct {
	URL string `json:"url,omitempty"`
}

// ButtonInput defines the input format of the Button objects
type ButtonInput struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"iconUrl,omitempty"` // either iconUrl of icon can be used
	Icon    string `json:"icon,omitempty"`    // either iconUrl of icon can be used
	OnClick string `json:"onClick,omitempty"`
}

// KeyValueInput defines the input format of the KeyValue object
type KeyValueInput struct {
	TopLabel         string       `json:"topLabel,omitempty"`
	Content          string       `json:"content,omitempty"`
	ContentMultiline bool         `json:"contentMultiline,omitempty"`
	BottomLabel      string       `json:"bottomLabel,omitempty"`
	OnClick          string       `json:"onClick,omitempty"`
	IconURL          string       `json:"iconUrl,omitempty"` // either iconUrl of icon can be used
	Icon             string       `json:"icon,omitempty"`    // either iconUrl of icon can be used
	Button           *ButtonInput `json:"button,omitempty"`
}

// simpleToAdvancedFormat converts one of chats markdown-like simple formats to its corresponding html based advanced format
func simpleToAdvancedFormat(simpleFormat string, advancedFormat string, stringToFormat string) string {
	var re = regexp.MustCompile("\\" + simpleFormat + "(.+?)\\" + simpleFormat + "")
	return re.ReplaceAllString(stringToFormat, "<"+advancedFormat+">${1}</"+advancedFormat+">")
}

// advancedToSimpleFormat converts one of chats html based advanced formats to its corresponding markdown-like simple format
func advancedToSimpleFormat(advancedFormat string, simpleFormat string, stringToFormat string) string {
	var re = regexp.MustCompile("<" + advancedFormat + ">(.+?)</" + advancedFormat + ">")
	return re.ReplaceAllString(stringToFormat, simpleFormat+"${1}"+simpleFormat)
}

// SimpleToAdvancedFormatting converts google chats simple formatting to the advanced formatting (as far as possible). Code strings and code blocks are not stripped.
func SimpleToAdvancedFormatting(simple string) (formatted string) {
	// Replacing links -- has to be done first because the other replacements will mess up this one
	var linkRegexp = regexp.MustCompile(`<(.+?)\|(.+?)>`)
	formatted = linkRegexp.ReplaceAllString(simple, `<a href="$1">$2</a>`)

	// General format replacing
	formatted = simpleToAdvancedFormat("*", "b", formatted)
	formatted = simpleToAdvancedFormat("_", "i", formatted)
	formatted = simpleToAdvancedFormat("~", "strike", formatted)

	return
}

// AdvancedToSimpleFormatting converts google chats advanced formatting to the simple formatting (as far as possible)
func AdvancedToSimpleFormatting(simple string) (formatted string) {
	// Replacing line break
	formatted = strings.ReplaceAll(simple, "<br>", "\n")

	// General format replacing
	formatted = advancedToSimpleFormat("b", "*", formatted)
	formatted = advancedToSimpleFormat("i", "_", formatted)
	formatted = advancedToSimpleFormat("strike", "~", formatted)

	// Replacing links
	var linkRegexp = regexp.MustCompile(`<a href="(.+?)">(.+?)</a>`)
	formatted = linkRegexp.ReplaceAllString(formatted, `<$1|$2>`)

	// Stripping underline by replacing it with nothing
	formatted = advancedToSimpleFormat("u", "", formatted)

	// Stripping font color
	var fontColorRegexp = regexp.MustCompile(`<font color=".+?">(.+?)</font>`)
	formatted = fontColorRegexp.ReplaceAllString(formatted, `$1`)

	return
}
