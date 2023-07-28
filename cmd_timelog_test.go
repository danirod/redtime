package main

import "testing"

func TestDecimal(t *testing.T) {
	if time := formatAsTime(0.25); time != "0:15" {
		t.Errorf("Expected 0.25 to yield 0:15, but yielded %s", time)
	}
	if time := formatAsTime(1); time != "1:00" {
		t.Errorf("Expected 1 to yield 1:00, but yielded %s", time)
	}
	if time := formatAsTime(1.5); time != "1:30" {
		t.Errorf("Expected 1.5 to yield 1:30, but yielded %s", time)
	}
	if time := formatAsTime(0.017); time != "0:01" {
		t.Errorf("Expected 0.016 to yield 0:01, but yielded %s", time)
	}
	if time := formatAsTime(0.01); time != "0:00" {
		t.Errorf("Expected .01 to yield 0:01, but yielded %s", time)
	}
}
