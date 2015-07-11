package main

import "testing"

func TestToSnake(t *testing.T) {
	cases := [][2]string{
		{"Name", "name"},
		{"OrderId", "order_id"},
		{"OrderID", "order_id"},
		{"URL", "url"},
		{"URLName", "url_name"},
		{"NameURL", "name_url"},
		{"TestO", "test_o"},
	}
	for _, c := range cases {
		if got, want := toSnake(c[0]), c[1]; got != want {
			t.Errorf("toSnake(%s) = %s; want %s", c[0], got, want)
		}
	}
}
