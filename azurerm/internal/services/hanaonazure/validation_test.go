package hanaonazure

import "testing"

func TestValidateHanaOnAzureSapMonitorName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: false,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 29 chars
			input:    "qwertyuiopasdfghjklzxcvbnmqwe",
			expected: true,
		},
		{
			// 30 chars
			input:    "qwertyuiopasdfghjklzxcvbnmqwea",
			expected: true,
		},
		{
			// 31 chars
			input:    "qwertyuiopasdfghjklzxcvbnmqweaf",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateHanaOnAzureSapMonitorName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
