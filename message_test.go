package main

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tests the triples function
func Test_triples(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output [][3]string
		err    string
	}{
		{
			name:  "Text button",
			input: "text|example.org website|https://example.org",
			output: [][3]string{
				[3]string{"text", "example.org website", "https://example.org"},
			},
			err: "",
		},
		{
			name:  "Built-in icon button",
			input: "builtin-icon|EMAIL|mailto:user@example.org",
			output: [][3]string{
				[3]string{"builtin-icon", "EMAIL", "mailto:user@example.org"},
			},
			err: "",
		},
		{
			name:  "Custom icon button",
			input: "custom-icon|https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg|https://bitrise.io",
			output: [][3]string{
				[3]string{"custom-icon", "https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg", "https://bitrise.io"},
			},
			err: "",
		},
		{
			name:  "Invalid button",
			input: "text|example.org website",
			err:   "Could not parse button with declaration text|example.org website",
		},
		{
			name:  "Multiple buttons",
			input: "text|example.org website|https://example.org\nbuiltin-icon|EMAIL|mailto:user@example.org\ncustom-icon|https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg|https://bitrise.io",
			output: [][3]string{
				[3]string{"text", "example.org website", "https://example.org"},
				[3]string{"builtin-icon", "EMAIL", "mailto:user@example.org"},
				[3]string{"custom-icon", "https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg", "https://bitrise.io"},
			},
			err: "",
		},
		{
			name:  "Multiple buttons with one invalid",
			input: "text|example.org website|https://example.org\nbuiltin-icon|mailto:user@example.org\ncustom-icon|https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg|https://bitrise.io",
			output: [][3]string{
				[3]string{"text", "example.org website", "https://example.org"},
			},
			err: "Could not parse button with declaration builtin-icon|mailto:user@example.org",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			splitTriples, err := triples(tc.input)

			if (err == nil && tc.err != "") || (err != nil && err.Error() != tc.err) {
				t.Errorf("Unexpected error: %s", err)
				return
			}

			if !reflect.DeepEqual(splitTriples, tc.output) {
				t.Errorf("Returned triples are not correct: expected %+v, got %+v", tc.output, splitTriples)
			}
		})
	}
}

func Test_parseButtons(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []*Button
		err    string
	}{
		{
			name:  "Text button",
			input: "text|example.org website|https://example.org",
			output: []*Button{{
				TextButton: &TextButton{
					Text: "example.org website",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "https://example.org",
						},
					},
				},
			}},
			err: "",
		},
		{
			name:  "Built-in icon button",
			input: "builtin-icon|EMAIL|mailto:user@example.org",
			output: []*Button{{
				ImageButton: &ImageButton{
					Icon: "EMAIL",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "mailto:user@example.org",
						},
					},
				},
			}},
			err: "",
		},
		{
			name:  "Custom icon button",
			input: "custom-icon|https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg|https://bitrise.io",
			output: []*Button{{
				ImageButton: &ImageButton{
					IconURL: "https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "https://bitrise.io",
						},
					},
				},
			}},
			err: "",
		},
		{
			name:  "Invalid button",
			input: "text|example.org website",
			err:   "Could not parse button with declaration text|example.org website",
		},
		{
			name:  "Multiple buttons",
			input: "text|example.org website|https://example.org\nbuiltin-icon|EMAIL|mailto:user@example.org\ncustom-icon|https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg|https://bitrise.io",
			output: []*Button{{
				TextButton: &TextButton{
					Text: "example.org website",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "https://example.org",
						},
					},
				},
			}, {
				ImageButton: &ImageButton{
					Icon: "EMAIL",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "mailto:user@example.org",
						},
					},
				},
			}, {
				ImageButton: &ImageButton{
					IconURL: "https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "https://bitrise.io",
						},
					},
				},
			}},
			err: "",
		},
		{
			name:   "Multiple buttons with one invalid",
			input:  "text|example.org website|https://example.org\nbuiltin-icon|mailto:user@example.org\ncustom-icon|https://pbs.twimg.com/profile_images/1039432724120051712/wFlFGsF3_400x400.jpg|https://bitrise.io",
			output: nil,
			err:    "Could not parse button with declaration builtin-icon|mailto:user@example.org",
		},
		{
			name:  "Parse button with trailing newline",
			input: "text|example.org website|https://example.org\n",
			output: []*Button{{
				TextButton: &TextButton{
					Text: "example.org website",
					OnClick: &OnClick{
						OpenLink: &OpenLink{
							URL: "https://example.org",
						},
					},
				},
			}},
			err: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buttons, err := parseButtons(tc.input)

			if (err == nil && tc.err != "") || (err != nil && err.Error() != tc.err) {
				t.Errorf("Unexpected error: %s", err)
				return
			}

			bts, btsErr := json.Marshal(buttons)
			out, outErr := json.Marshal(tc.output)
			if btsErr != nil || outErr != nil {
				t.Errorf("Could not marshal json!\n%s\n%s", btsErr, outErr)
			}

			if !cmp.Equal(buttons, tc.output) {
				t.Errorf("Returned buttons are not correct:\nexpected: %s\ngot: %s", out, bts)
			}
		})
	}
}
