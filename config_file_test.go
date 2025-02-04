package setlist

import (
	"testing"
)

func TestStartUrl(t *testing.T) {
	tt := []struct {
		name            string
		identityStoreID string
		ssoFriendlyName string
		expected        string
	}{
		{
			"has friendly name",
			"d-012345",
			"foo",
			"https://foo.awsapps.com/start",
		},
		{
			"missing friendly name",
			"d-012345",
			"",
			"https://d-012345.awsapps.com/start",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			c := ConfigFile{
				IdentityStoreId: tc.identityStoreID,
				FriendlyName:    tc.ssoFriendlyName,
			}
			actual := c.StartURL()

			if actual != tc.expected {
				t.Errorf("unexpected output: got %v, want %v", actual, tc.expected)
			}
		})
	}
}

func TestHasFriendlyName(t *testing.T) {
	tests := []struct {
		name  string
		given string
		want  bool
	}{
		{"has_friendly_name", "foo", true},
		{"has_no_friendly_name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ConfigFile{
				FriendlyName: tt.given,
			}
			got := c.hasFriendlyName()
			if got != tt.want {
				t.Fatalf("Expected %v but got %v", tt.want, got)
			}
		})
	}
}
