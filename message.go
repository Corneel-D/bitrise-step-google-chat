package main

import (
	"encoding/json"
	"fmt"
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
	OnClick          *OnClick `json:"onClick,omitempty,omitempty"`
	// either iconUrl of icon can be used
	IconURL string `json:"iconUrl,omitempty"`
	// either iconUrl of icon can be used
	Icon   string  `json:"icon,omitempty"`
	Button *Button `json:"button,omitempty"`
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
