package main

import (
	"reflect"
	"testing"
)

func TestParseNicknameMapping(t *testing.T) {
	tt := []struct {
		name     string
		mapping  string
		expected map[string]string
	}{
		{
			"knowngood",
			"01234=foo,9876=bar",
			map[string]string{
				"01234": "foo",
				"9876":  "bar",
			},
		},
		{
			"empty mapping",
			"",
			map[string]string{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseNicknameMapping(tc.mapping)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("unexpected output: got %v, want %v", actual, tc.expected)
			}
		})
	}
}

func TestNicknameFor(t *testing.T) {
	tt := []struct {
		name      string
		accountId string
		mapping   map[string]string
		expected  string
	}{
		{
			"has nickname",
			"01234",
			map[string]string{
				"01234": "foo",
				"9876":  "bar",
			},
			"foo",
		},
		{
			"missing account id",
			"01234",
			map[string]string{
				"9876": "bar",
			},
			"NoNickname-01234",
		},
		{
			"missing nickname",
			"01234",
			map[string]string{
				"01234": "",
				"9876":  "bar",
			},
			"NoNickname-01234",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := nicknameFor(tc.accountId, tc.mapping)

			if actual != tc.expected {
				t.Errorf("unexpected output: got %v, want %v", actual, tc.expected)
			}
		})
	}
}
