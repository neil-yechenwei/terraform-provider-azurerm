package hanaonazure

import "testing"

func TestValidateHanaSapMonitorName(t *testing.T) {
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
			input:    "hello1",
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

		_, errors := ValidateHanaSapMonitorName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateHanaDBName(t *testing.T) {
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
			input:    "HELLO",
			expected: true,
		},
		{
			// can't be lowercase letters
			input:    "hello",
			expected: false,
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
			// 63 chars
			input:    "QWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPASD",
			expected: true,
		},
		{
			// 64 chars
			input:    "QWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPASDY",
			expected: true,
		},
		{
			// 65 chars
			input:    "QWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPQWERTYUIOPASDYH",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateHanaDBName(v.input, "hana_db_name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateHanaDBUserName(t *testing.T) {
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
			// can't be lowercase letters
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
			// 31 chars
			input:    "qwertyuiopqwertyuiopqwertyuiopa",
			expected: true,
		},
		{
			// 32 chars
			input:    "qwertyuiopqwertyuiopqwertyuiopaa",
			expected: true,
		},
		{
			// 33 chars
			input:    "qwertyuiopqwertyuiopqwertyuiopaaa",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateHanaDBUserName(v.input, "hana_db_username")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}

func TestValidateHanaDBPassword(t *testing.T) {
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
			// can't be lowercase letters
			input:    "hello",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: true,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: true,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: true,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: true,
		},
		{
			// 62 chars
			input:    "qwertyuiopqwertyuiopqwertyuiopaqwertyuiopqwertyuiopqwertyuiopab",
			expected: true,
		},
		{
			// 63 chars
			input:    "qwertyuiopqwertyuiopqwertyuiopaaqwertyuiopqwertyuiopqwertyuiopab",
			expected: true,
		},
		{
			// 64 chars
			input:    "qwertyuiopqwertyuiopqwertyuiopaaaqwertyuiopqwertyuiopqwertyuiopab",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ValidateHanaDBPassword(v.input, "hana_db_password")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
