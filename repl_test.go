package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "\tCharmander\nBulbasaur\r\nPIKACHU\t",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
		{
			input:    "PIKACHU",
			expected: []string{"pikachu"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    " \t\n\r ",
			expected: []string{},
		},
	}

	for _, testCase := range cases {
		actual := cleanInput(testCase.input)
		if len(actual) != len(testCase.expected) {
			t.Errorf("cleanInput(%q): got %d words, want %d", testCase.input, len(actual), len(testCase.expected))
			continue
		}
		for index, word := range actual {
			if word != testCase.expected[index] {
				t.Errorf("cleanInput(%q): word %d = %q, want %q", testCase.input, index, word, testCase.expected[index])
			}
		}
	}
}

func TestCommandExploreArguments(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "missing area"},
		{name: "blank area", args: []string{" "}},
		{name: "extra arguments", args: []string{"canalave-city-area", "extra"}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if err := commandExplore(&config{}, testCase.args...); err == nil {
				t.Fatal("expected an argument error before making a request")
			}
		})
	}
}
