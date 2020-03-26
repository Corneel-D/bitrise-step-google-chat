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

func Test_simpleToAdvancedFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Replace at start of string",
			input:  "*test test* test",
			output: "<b>test test</b> test",
		},
		{
			name:   "Replace in middle of string",
			input:  "test *test* test",
			output: "test <b>test</b> test",
		},
		{
			name:   "Replace at end of string",
			input:  "test *test test*",
			output: "test <b>test test</b>",
		},
		{
			name:   "Replace around whole string",
			input:  "*test test test*",
			output: "<b>test test test</b>",
		},
		{
			name:   "Replace multiple",
			input:  "*test* test *test*",
			output: "<b>test</b> test <b>test</b>",
		},
		{
			name:   "Replace nothing with other format",
			input:  "_test test test_",
			output: "_test test test_",
		},
		{
			name:   "Replace nothing with reverse format",
			input:  "<b>test test test</b>",
			output: "<b>test test test</b>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := simpleToAdvancedFormat("*", "b", tc.input)

			if formatted != tc.output {
				t.Errorf("Substitution failed.\nExpected: %s\nActual: %s", tc.output, formatted)
			}
		})
	}
}

func Test_advancedToSimpleFormat(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Replace at start of string",
			input:  "<b>test test</b> test",
			output: "*test test* test",
		},
		{
			name:   "Replace in middle of string",
			input:  "test <b>test</b> test",
			output: "test *test* test",
		},
		{
			name:   "Replace at end of string",
			input:  "test <b>test test</b>",
			output: "test *test test*",
		},
		{
			name:   "Replace around whole string",
			input:  "<b>test test test</b>",
			output: "*test test test*",
		},
		{
			name:   "Replace multiple",
			input:  "<b>test</b> test <b>test</b>",
			output: "*test* test *test*",
		},
		{
			name:   "Replace nothing with other format",
			input:  "<i>test test test</i>",
			output: "<i>test test test</i>",
		},
		{
			name:   "Replace nothing with reverse format",
			input:  "*test test test*",
			output: "*test test test*",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := advancedToSimpleFormat("b", "*", tc.input)

			if formatted != tc.output {
				t.Errorf("Substitution failed.\nExpected: %s\nActual: %s", tc.output, formatted)
			}
		})
	}
}

func Test_SimpleToAdvancedFormatting(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Different formats",
			input:  "*test* _test_ ~test~\n<https://example.org|test>",
			output: "<b>test</b> <i>test</i> <strike>test</strike>\n<a href=\"https://example.org\">test</a>",
		},
		{
			name:   "Ignore backtics",
			input:  "*test* _test_ `test` ~test~\n```test```",
			output: "<b>test</b> <i>test</i> `test` <strike>test</strike>\n```test```",
		},
		{
			name:   "Mixed formats",
			input:  "_*test* test ~test_ <https://example.org|test>~",
			output: "<i><b>test</b> test <strike>test</i> <a href=\"https://example.org\">test</a></strike>",
		},
		{
			name:   "Mixed formats, not working over a newline",
			input:  "_*test* test ~test_\n<https://example.org|test>~",
			output: "<i><b>test</b> test ~test</i>\n<a href=\"https://example.org\">test</a>~",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := SimpleToAdvancedFormatting(tc.input)

			if formatted != tc.output {
				t.Errorf("Substitution failed.\nExpected: %s\nActual: %s", tc.output, formatted)
			}
		})
	}
}

func Test_AdvancedToSimpleFormatting(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Different formats",
			input:  "<b>test</b> <i>test</i> <strike>test</strike>\n<a href=\"https://example.org\">test</a>",
			output: "*test* _test_ ~test~\n<https://example.org|test>",
		},
		{
			name:   "Mixed formats",
			input:  "<i><b>test</b> test <strike>test</i> <a href=\"https://example.org\">test</a></strike>",
			output: "_*test* test ~test_ <https://example.org|test>~",
		},
		{
			name:   "Mixed formats, not working over a newline",
			input:  "<i><b>test</b> test <strike>test</i>\n<a href=\"https://example.org\">test</a></strike>",
			output: "_*test* test <strike>test_\n<https://example.org|test></strike>",
		},
		{
			name:   "Replace line breaks",
			input:  "<i><b>test</b> test <strike>test</i><br><a href=\"https://example.org\">test</a></strike>",
			output: "_*test* test <strike>test_\n<https://example.org|test></strike>",
		},
		{
			name:   "Strip underline",
			input:  "<b>test</b> <i>test</i> <strike>test</strike> <u>test</u>",
			output: "*test* _test_ ~test~ test",
		},
		{
			name:   "Strip font color",
			input:  "<b>test</b> <i>test</i> <strike>test</strike> <font color=\"green\">test</font>",
			output: "*test* _test_ ~test~ test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formatted := AdvancedToSimpleFormatting(tc.input)

			if formatted != tc.output {
				t.Errorf("Substitution failed.\nExpected: %s\nActual: %s", tc.output, formatted)
			}
		})
	}
}
