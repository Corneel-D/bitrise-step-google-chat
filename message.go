package main

// Message which can be send to Google Chat
// More info at https://developers.google.com/hangouts/chat/reference/message-formats/cards
type Message struct {
	Cards []Card `json:"cards"`
}

// Card property of a message. can contain a header and must have at least one section
type Card struct {
	// header object (optional)
	Header Header `json:"header,omitempty"`
	// sections object. At least one section is required
	Sections []Section `json:"sections"`
}

// Header property of a card
type Header struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	// imageStyle controls the shape of the header image, which may be square ("IMAGE") or circular ("AVATAR"). The default is square ("IMAGE").
	ImageStyle string `json:"imageStyle,omitempty"`
}

// Section of a card. can contain multiple widgets, but at least one is required. Sections are separated by a horizontal line
type Section struct {
	// Section header (optional)
	Header string `json:"header,omitempty"`
	// widgets object. At least one widget is required.
	Widgets []Widget `json:"widgets"`
}

// Widget of a section. Can contain only one type of UI element
type Widget struct {
	TextParagraph TextParagraph `json:"textParagraph,omitempty"`
	KeyValue      KeyValue      `json:"keyValue,omitempty"`
	Image         Image         `json:"image,omitempty"`
	// buttons object can contain one or more buttons. will be laid out horizontally
	Buttons []Button `json:"buttons,omitempty"`
}

// TextParagraph UI element
type TextParagraph struct {
	// The text to display inside the paragraph
	Text string `json:"text"`
}

// KeyValue UI element
type KeyValue struct {
	TopLabel         string  `json:"topLabel,omitempty"`
	Content          string  `json:"content,omitempty"`
	ContentMultiline string  `json:"contentMultiline,omitempty"`
	BottomLabel      string  `json:"bottomLabel,omitempty"`
	OnClick          OnClick `json:"onClick,omitempty,omitempty"`
	// either iconUrl of icon can be used
	IconURL string `json:"iconUrl,omitempty"`
	// either iconUrl of icon can be used
	Icon   string `json:"icon,omitempty"`
	Button Button `json:"button,omitempty"`
}

// Image UI element
type Image struct {
	ImageURL string  `json:"imageUrl"`
	OnClick  OnClick `json:"onClick,omitempty"`
}

// Button UI element. Can contain either a TextButton or an ImageButton
type Button struct {
	TextButton  TextButton  `json:"textButton,omitempty"`
	ImageButton ImageButton `json:"imageButton,omitempty"`
}

// TextButton UI element
type TextButton struct {
	Text    string  `json:"text"`
	OnClick OnClick `json:"onClick,omitempty"`
}

// ImageButton UI element
type ImageButton struct {
	// either iconUrl of icon can be used
	IconURL string `json:"iconUrl,omitempty"`
	// either iconUrl of icon can be used
	Icon    string  `json:"icon,omitempty"`
	OnClick OnClick `json:"onClick,omitempty"`
}

// OnClick handler object
type OnClick struct {
	OpenLink OpenLink `json:"openLink"`
}

// OpenLink object
type OpenLink struct {
	URL string `json:"url"`
}
