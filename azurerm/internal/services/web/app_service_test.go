package web

import (
	"testing"
)

func TestParseAppService(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *AppServiceResourceID
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing Sites Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/",
			Expected: nil,
		},
		{
			Name:  "App Service Resource ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/site1",
			Expected: &AppServiceResourceID{
				Name:          "site1",
				ResourceGroup: "mygroup1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/Sites/",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ParseAppServiceID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}

func TestValidateAppServiceID(t *testing.T) {
	cases := []struct {
		ID    string
		Valid bool
	}{
		{
			ID:    "",
			Valid: false,
		},
		{
			ID:    "nonsense",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web",
			Valid: false,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/sites/duckduckgo",
			Valid: true,
		},
		{
			ID:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Web/Sites/duckduckgo",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %q", tc.ID)
		_, errors := ValidateAppServiceID(tc.ID, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
