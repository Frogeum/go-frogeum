// Copyright 2018 The go-frogeum Authors
// This file is part of the go-frogeum library.
//
// The go-frogeum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-frogeum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-frogeum library. If not, see <http://www.gnu.org/licenses/>.

package accounts

import (
	"testing"
)

func TestURLParsing(t *testing.T) {
	url, err := parseURL("https://frogeum.org")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "frogeum.org" {
		t.Errorf("expected: %v, got: %v", "frogeum.org", url.Path)
	}

	_, err = parseURL("frogeum.org")
	if err == nil {
		t.Error("expected err, got: nil")
	}
}

func TestURLString(t *testing.T) {
	url := URL{Scheme: "https", Path: "frogeum.org"}
	if url.String() != "https://frogeum.org" {
		t.Errorf("expected: %v, got: %v", "https://frogeum.org", url.String())
	}

	url = URL{Scheme: "", Path: "frogeum.org"}
	if url.String() != "frogeum.org" {
		t.Errorf("expected: %v, got: %v", "frogeum.org", url.String())
	}
}

func TestURLMarshalJSON(t *testing.T) {
	url := URL{Scheme: "https", Path: "frogeum.org"}
	json, err := url.MarshalJSON()
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if string(json) != "\"https://frogeum.org\"" {
		t.Errorf("expected: %v, got: %v", "\"https://frogeum.org\"", string(json))
	}
}

func TestURLUnmarshalJSON(t *testing.T) {
	url := &URL{}
	err := url.UnmarshalJSON([]byte("\"https://frogeum.org\""))
	if err != nil {
		t.Errorf("unexpcted error: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("expected: %v, got: %v", "https", url.Scheme)
	}
	if url.Path != "frogeum.org" {
		t.Errorf("expected: %v, got: %v", "https", url.Path)
	}
}

func TestURLComparison(t *testing.T) {
	tests := []struct {
		urlA   URL
		urlB   URL
		expect int
	}{
		{URL{"https", "frogeum.org"}, URL{"https", "frogeum.org"}, 0},
		{URL{"http", "frogeum.org"}, URL{"https", "frogeum.org"}, -1},
		{URL{"https", "frogeum.org/a"}, URL{"https", "frogeum.org"}, 1},
		{URL{"https", "abc.org"}, URL{"https", "frogeum.org"}, -1},
	}

	for i, tt := range tests {
		result := tt.urlA.Cmp(tt.urlB)
		if result != tt.expect {
			t.Errorf("test %d: cmp mismatch: expected: %d, got: %d", i, tt.expect, result)
		}
	}
}
