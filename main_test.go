package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_selectValue(t *testing.T) {
	tests := []struct {
		name      string
		success   bool
		ifSuccess string
		ifFailed  string
		output    string
	}{
		{
			name:      "Success",
			success:   true,
			ifSuccess: "Successfull",
			ifFailed:  "Failed",
			output:    "Successfull",
		}, {
			name:      "Failed with fail message",
			success:   false,
			ifSuccess: "Successfull",
			ifFailed:  "Failed",
			output:    "Failed",
		}, {
			name:      "Failed with no fail message",
			success:   false,
			ifSuccess: "Successfull",
			ifFailed:  "",
			output:    "Successfull",
		}, {
			name:      "Success with empty message and non-empty fail message",
			success:   true,
			ifSuccess: "",
			ifFailed:  "Failed",
			output:    "",
		}, {
			name:      "Fail with empty message and non-empty fail message",
			success:   false,
			ifSuccess: "",
			ifFailed:  "Failed",
			output:    "Failed",
		}, {
			name:      "Fail with empty message and empty fail message",
			success:   false,
			ifSuccess: "",
			ifFailed:  "",
			output:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			success = tc.success

			selected := selectValue(tc.ifSuccess, tc.ifFailed)

			if tc.output != selected {
				t.Errorf("Returned string is not correct: expected %+v, got %+v", tc.output, selected)
			}
		})
	}
}

func Test_selectAvancedFormatValue(t *testing.T) {
	tests := []struct {
		name      string
		success   bool
		ifSuccess string
		ifFailed  string
		transform bool
		output    string
	}{
		{
			name:      "Success without transforming message",
			success:   true,
			ifSuccess: "*Successfull*",
			ifFailed:  "*Failed*",
			transform: false,
			output:    "*Successfull*",
		}, {
			name:      "Success with transforming message",
			success:   true,
			ifSuccess: "*Successfull*",
			ifFailed:  "*Failed*",
			transform: true,
			output:    "<b>Successfull</b>",
		}, {
			name:      "Failed without transforming message",
			success:   false,
			ifSuccess: "*Successfull*",
			ifFailed:  "*Failed*",
			transform: false,
			output:    "*Failed*",
		}, {
			name:      "Failed with transforming message",
			success:   false,
			ifSuccess: "*Successfull*",
			ifFailed:  "*Failed*",
			transform: true,
			output:    "<b>Failed</b>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			success = tc.success

			selected := selectAvancedFormatValue(tc.ifSuccess, tc.ifFailed, tc.transform)

			if tc.output != selected {
				t.Errorf("Returned string is not correct: expected %+v, got %+v", tc.output, selected)
			}
		})
	}
}

func Test_selectSimpleFormatValue(t *testing.T) {
	tests := []struct {
		name      string
		success   bool
		ifSuccess string
		ifFailed  string
		transform bool
		output    string
	}{
		{
			name:      "Success without transforming message",
			success:   true,
			ifSuccess: "<b>Successfull</b>",
			ifFailed:  "<b>Failed</b>",
			transform: false,
			output:    "<b>Successfull</b>",
		}, {
			name:      "Success with transforming message",
			success:   true,
			ifSuccess: "<b>Successfull</b>",
			ifFailed:  "<b>Failed</b>",
			transform: true,
			output:    "*Successfull*",
		}, {
			name:      "Failed without transforming message",
			success:   false,
			ifSuccess: "<b>Successfull</b>",
			ifFailed:  "<b>Failed</b>",
			transform: false,
			output:    "<b>Failed</b>",
		}, {
			name:      "Failed with transforming message",
			success:   false,
			ifSuccess: "<b>Successfull</b>",
			ifFailed:  "<b>Failed</b>",
			transform: true,
			output:    "*Failed*",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			success = tc.success

			selected := selectSimpleFormatValue(tc.ifSuccess, tc.ifFailed, tc.transform)

			if tc.output != selected {
				t.Errorf("Returned string is not correct: expected %+v, got %+v", tc.output, selected)
			}
		})
	}
}

func Test_newMessage(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		output Message
		err    string
	}{
		{
			name: "Create message with header, text, keyValues and buttons",
			config: Config{
				WebhookURL: "URL",
				Title:      "title",
				Subtitle:   "subtitle",
				ImageURL:   "image url",
				ImageStyle: "image style",
				Text:       "text",
				KeyValue:   `[{"topLabel":"top label","content":"content\non <i>multiple<\/i> lines","contentMultiline":true,"bottomLabel":"bottom label","onClick":"https://example.org","button":{"text":"button text","onClick":"http://example.org"}},{"content":"<b>content<\/b> :)","iconUrl":"https://example.com","button":{"iconUrl":"http://example.com","onClick":"http://example.org"}},{"content":"content.","icon":"CLOCK","button":{"icon":"MULTIPLE_PEOPLE","onClick":"http://example.org"}}]`,
				Buttons:    "text|example.org website|https://example.org",
			},
			output: Message{
				Text: "title",
				Cards: []Card{{
					Header: &Header{
						Title:      "title",
						Subtitle:   "subtitle",
						ImageURL:   "image url",
						ImageStyle: "image style",
					},
					Sections: []Section{{
						Widgets: []*Widget{{
							TextParagraph: &TextParagraph{
								Text: "text",
							},
						}},
					}, {
						Widgets: []*Widget{{
							KeyValue: &KeyValue{
								TopLabel:         "top label",
								Content:          "content\non <i>multiple</i> lines",
								ContentMultiline: "true",
								BottomLabel:      "bottom label",
								OnClick: &OnClick{
									OpenLink: &OpenLink{
										URL: "https://example.org",
									},
								},
								Button: &Button{
									TextButton: &TextButton{
										Text: "button text",
										OnClick: &OnClick{
											OpenLink: &OpenLink{
												URL: "http://example.org",
											},
										},
									},
								},
							},
						}, {
							KeyValue: &KeyValue{
								Content:          "<b>content</b> :)",
								ContentMultiline: "false",
								IconURL:          "https://example.com",
								Button: &Button{
									ImageButton: &ImageButton{
										IconURL: "http://example.com",
										OnClick: &OnClick{
											OpenLink: &OpenLink{
												URL: "http://example.org",
											},
										},
									},
								},
							},
						}, {
							KeyValue: &KeyValue{
								Content:          "content.",
								ContentMultiline: "false",
								Icon:             "CLOCK",
								Button: &Button{
									ImageButton: &ImageButton{
										Icon: "MULTIPLE_PEOPLE",
										OnClick: &OnClick{
											OpenLink: &OpenLink{
												URL: "http://example.org",
											},
										},
									},
								},
							},
						}},
					}, {
						Widgets: []*Widget{{
							Buttons: []*Button{{
								TextButton: &TextButton{
									Text: "example.org website",
									OnClick: &OnClick{
										OpenLink: &OpenLink{
											URL: "https://example.org",
										},
									},
								},
							}},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with header and text",
			config: Config{
				WebhookURL: "URL",
				Title:      "title",
				Subtitle:   "subtitle",
				ImageURL:   "image url",
				ImageStyle: "image style",
				Text:       "text",
			},
			output: Message{
				Text: "title",
				Cards: []Card{{
					Header: &Header{
						Title:      "title",
						Subtitle:   "subtitle",
						ImageURL:   "image url",
						ImageStyle: "image style",
					},
					Sections: []Section{{
						Widgets: []*Widget{{
							TextParagraph: &TextParagraph{
								Text: "text",
							},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with header, and buttons",
			config: Config{
				WebhookURL: "URL",
				Title:      "title",
				Subtitle:   "subtitle",
				ImageURL:   "image url",
				ImageStyle: "image style",
				Buttons:    "text|example.org website|https://example.org",
			},
			output: Message{
				Text: "title",
				Cards: []Card{{
					Header: &Header{
						Title:      "title",
						Subtitle:   "subtitle",
						ImageURL:   "image url",
						ImageStyle: "image style",
					},
					Sections: []Section{{
						Widgets: []*Widget{{
							Buttons: []*Button{{
								TextButton: &TextButton{
									Text: "example.org website",
									OnClick: &OnClick{
										OpenLink: &OpenLink{
											URL: "https://example.org",
										},
									},
								},
							}},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with header, text, keyValues and buttons",
			config: Config{
				WebhookURL: "URL",
				KeyValue:   `[{"topLabel":"top label","content":"content\non <i>multiple<\/i> lines","contentMultiline":true,"bottomLabel":"bottom label","onClick":"https://example.org","button":{"text":"button text","onClick":"http://example.org"}}]`,
			},
			output: Message{
				Cards: []Card{{
					Sections: []Section{{
						Widgets: []*Widget{{
							KeyValue: &KeyValue{
								TopLabel:         "top label",
								Content:          "content\non <i>multiple</i> lines",
								ContentMultiline: "true",
								BottomLabel:      "bottom label",
								OnClick: &OnClick{
									OpenLink: &OpenLink{
										URL: "https://example.org",
									},
								},
								Button: &Button{
									TextButton: &TextButton{
										Text: "button text",
										OnClick: &OnClick{
											OpenLink: &OpenLink{
												URL: "http://example.org",
											},
										},
									},
								},
							},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with buttons",
			config: Config{
				WebhookURL: "URL",
				Buttons:    "text|example.org website|https://example.org",
			},
			output: Message{
				Cards: []Card{{
					Sections: []Section{{
						Widgets: []*Widget{{
							Buttons: []*Button{{
								TextButton: &TextButton{
									Text: "example.org website",
									OnClick: &OnClick{
										OpenLink: &OpenLink{
											URL: "https://example.org",
										},
									},
								},
							}},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with text",
			config: Config{
				WebhookURL: "URL",
				Text:       "text",
			},
			output: Message{
				Text: "text",
				Cards: []Card{{
					Sections: []Section{{
						Widgets: []*Widget{{
							TextParagraph: &TextParagraph{
								Text: "text",
							},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with message, header, and buttons",
			config: Config{
				WebhookURL: "URL",
				Message:    "message",
				Title:      "title",
				Subtitle:   "subtitle",
				ImageURL:   "image url",
				ImageStyle: "image style",
				Buttons:    "text|example.org website|https://example.org",
			},
			output: Message{
				Text: "message",
				Cards: []Card{{
					Header: &Header{
						Title:      "title",
						Subtitle:   "subtitle",
						ImageURL:   "image url",
						ImageStyle: "image style",
					},
					Sections: []Section{{
						Widgets: []*Widget{{
							Buttons: []*Button{{
								TextButton: &TextButton{
									Text: "example.org website",
									OnClick: &OnClick{
										OpenLink: &OpenLink{
											URL: "https://example.org",
										},
									},
								},
							}},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with message, header, buttons, and empty keyvalue",
			config: Config{
				WebhookURL: "URL",
				Message:    "message",
				Title:      "title",
				Subtitle:   "subtitle",
				ImageURL:   "image url",
				ImageStyle: "image style",
				Buttons:    "text|example.org website|https://example.org",
				KeyValue:   "[]",
			},
			output: Message{
				Text: "message",
				Cards: []Card{{
					Header: &Header{
						Title:      "title",
						Subtitle:   "subtitle",
						ImageURL:   "image url",
						ImageStyle: "image style",
					},
					Sections: []Section{{
						Widgets: []*Widget{{
							Buttons: []*Button{{
								TextButton: &TextButton{
									Text: "example.org website",
									OnClick: &OnClick{
										OpenLink: &OpenLink{
											URL: "https://example.org",
										},
									},
								},
							}},
						}},
					}},
				}},
			},
			err: "",
		},
		{
			name: "Create message with button error",
			config: Config{
				WebhookURL: "URL",
				Title:      "title",
				Subtitle:   "subtitle",
				ImageURL:   "image url",
				ImageStyle: "image style",
				Text:       "text",
				Buttons:    "invalid",
			},
			output: Message{
				Cards: nil,
			},
			err: "Could not parse button with declaration invalid",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Make sure success is set to true, so less config is needed. selectValue is tested in another test
			success = true

			message, err := newMessage(tc.config)

			if (err == nil && tc.err != "") || (err != nil && err.Error() != tc.err) {
				t.Errorf("Unexpected error: %s", err)
				return
			}

			msg, msgErr := json.Marshal(message)
			out, outErr := json.Marshal(tc.output)
			if msgErr != nil || outErr != nil {
				t.Errorf("Could not marshal json!\n%s\n%s", msgErr, outErr)
			}

			if !cmp.Equal(msg, out) {
				t.Errorf("Returned message is not correct:\nexpected:   %s\ngot:        %s", out, msg)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		err    string
	}{
		{
			name: "No webhook",
			config: &Config{
				WebhookURL: "",
				Text:       "Text",
				Buttons:    "Buttons",
			},
			err: "WebhookURL is empty. You need to provide one",
		}, {
			name: "No Text or buttons",
			config: &Config{
				WebhookURL: "URL",
			},
			err: "Text, keyValue and buttons are empty. You need to provide at least one",
		}, {
			name: "Text",
			config: &Config{
				WebhookURL: "URL",
				Text:       "Text",
			},
			err: "",
		}, {
			name: "KeyValue",
			config: &Config{
				WebhookURL: "URL",
				KeyValue:   `[{"content": "content"}]`,
			},
			err: "",
		}, {
			name: "Buttons",
			config: &Config{
				WebhookURL: "URL",
				Buttons:    "Buttons",
			},
			err: "",
		}, {
			name: "Text and KeyValue",
			config: &Config{
				WebhookURL: "URL",
				Text:       "Text",
				KeyValue:   `[{"content": "content"}]`,
			},
			err: "",
		}, {
			name: "Text and buttons",
			config: &Config{
				WebhookURL: "URL",
				Text:       "Text",
				Buttons:    "Buttons",
			},
			err: "",
		}, {
			name: "KeyValue and buttons",
			config: &Config{
				WebhookURL: "URL",
				KeyValue:   `[{"content": "content"}]`,
				Buttons:    "Buttons",
			},
			err: "",
		}, {
			name: "Text, KeyValue and buttons",
			config: &Config{
				WebhookURL: "URL",
				Text:       "Text",
				KeyValue:   `{"content": "content"}`,
				Buttons:    "Buttons",
			},
			err: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.config)

			if (err == nil && tc.err != "") || (err != nil && err.Error() != tc.err) {
				t.Errorf("Unexpected error: %s", err)
				return
			}
		})
	}
}
