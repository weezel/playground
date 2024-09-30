package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"weezel/playground/cmd/tvprogs/programinfo"
)

func Test_parseIltapuluScheduleOrder(t *testing.T) {
	absPath, err := filepath.Abs("iltapulu.html")
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(absPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	someTime := time.Date(2024, 7, 21, 0, 0, 0, 0, helsinkiTZ)
	var channels programinfo.Channels
	if channels, err = parseIltapulu(f, someTime); err != nil {
		t.Fatalf("parse Iltapulu: %v", err)
	}

	for _, progs := range channels {
		for i := 1; i < len(progs); i++ {
			prev := progs[i-1]
			current := progs[i]
			if prev.StartTime.After(current.StartTime) {
				t.Errorf("Previous progam has start time after the current program: \n%#v\n%#v\n",
					prev,
					current,
				)
			}
		}
	}
}
