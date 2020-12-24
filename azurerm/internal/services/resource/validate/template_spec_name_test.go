package validate

import (
	"strings"
	"testing"
)

func TestTemplateSpecName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			Input:    "",
			Expected: false,
		},
		{
			Input:    "TestTemplateSpecName@",
			Expected: false,
		},
		{
			Input:    "Test Template Spec Name",
			Expected: false,
		},
		{
			Input:    strings.Repeat("s", 89),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 90),
			Expected: true,
		},
		{
			Input:    strings.Repeat("s", 91),
			Expected: false,
		},
		{
			Input:    "TestTemplateSpecName",
			Expected: true,
		},
	}

	for _, v := range testCases {
		_, errors := TemplateSpecName(v.Input, "template_spec_name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t but got %t (and %d errors)", v.Expected, result, len(errors))
		}
	}
}
