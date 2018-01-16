package main

import (
	"testing"
)

func TestCfg(t *testing.T) {
	var i int = 1
	if i != 1 {
		t.Errorf("i must eq 1")
	}
}
