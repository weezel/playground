package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"weezel/playground/cmd/tvprogs/programinfo"

	"github.com/PuerkitoBio/goquery"
)

const iltapuluURL = "https://www.iltapulu.fi/koko-paiva"

var (
	flagWholeDay  bool
	flagShowDates bool
	flagFromFile  bool
	flagToJSON    bool
)

var helsinkiTZ, _ = time.LoadLocation("Europe/Helsinki")

func parseIltapulu(r io.Reader, now time.Time) (programinfo.Channels, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return programinfo.Channels{}, fmt.Errorf("read document: %w", err)
	}

	channels := programinfo.Channels{}
	pat := `ul.daypart-block`
	doc.Find(pat).Each(func(_ int, s *goquery.Selection) {
		channelName, ok := s.Parent().Find("section > a.channel-header").Attr("title")
		if !ok {
			log.Panicf("Failed to get channel info")
		}

		programs := []programinfo.Program{}
		s.Find(`ul > li`).Each(func(_ int, ss *goquery.Selection) {
			startTime := ss.Find(`time`).Text()
			progTime, err := time.Parse("15.04", startTime)
			if err != nil {
				log.Fatalf("Failed to parse %q: %v\n", startTime, err)
			}
			progTime = time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				progTime.Hour(),
				progTime.Minute(),
				progTime.Second(),
				progTime.Nanosecond(),
				helsinkiTZ,
			)

			prog := ss.Find(`b > a.op`).Text()

			programs = append(programs, programinfo.Program{
				StartTime: progTime,
				Name:      prog,
			})
		})

		if _, ok := channels[channelName]; !ok {
			channels[channelName] = []programinfo.Program{}
		}
		progs := channels[channelName]
		progs = append(progs, programs...)
		channels[channelName] = progs
	})

	// Normalize times
	for _, programs := range channels {
		// If this is set to true, all programs will be marked for the next day from now on.
		nextDay := false
		for i := 0; i < len(programs); i++ {
			prog := &programs[i]

			curHour, _, _ := prog.StartTime.Clock()
			// If the iteration is before the midpoint of the program list,
			// it's assumed that programs having 20-23 as their hour,
			// are considered being started on the previous day.
			isHalfOfTheProgramsRemaining := i < len(programs)/2
			if isHalfOfTheProgramsRemaining && curHour >= 20 && curHour <= 23 {
				prog.StartTime = prog.StartTime.AddDate(0, 0, -1)
				continue
			}

			// Vice versa, after the midpoint and hour being 0-6 it's considered as
			// a next day's program.
			isHalfOfTheProgramsPassed := i > len(programs)/2
			if nextDay || isHalfOfTheProgramsPassed && curHour >= 0 && curHour <= 6 {
				prog.StartTime = prog.StartTime.AddDate(0, 0, 1)
				nextDay = true
			}
		}
	}

	return channels, nil
}

func fromInternet() []byte {
	res, err := http.Get(iltapuluURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Panicf("Status code error: %d %s\n", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panicf("Read HTTP body: %v\n", err)
	}

	return body
}

func fromFile() []byte {
	absPath, err := filepath.Abs("iltapulu.html")
	if err != nil {
		panic(err)
	}
	f, err := os.Open(absPath)
	if err != nil {
		log.Panicf("Open file %q: %v\n", absPath, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		log.Panicf("File all file %q: %v\n", absPath, err)
	}

	return data
}

func main() {
	flag.BoolVar(&flagFromFile, "f", false, "Read from iltapulu.html")
	flag.BoolVar(&flagWholeDay, "d", false, "Show the whole day's info")
	flag.BoolVar(&flagShowDates, "D", false, "Also print dates")
	flag.BoolVar(&flagToJSON, "j", false, "Dump to JSON")
	flag.Parse()

	var data []byte
	if flagFromFile {
		data = fromFile()
	} else {
		data = fromInternet()
	}

	now := time.Now().In(helsinkiTZ)
	var err error
	var channels programinfo.Channels
	if channels, err = parseIltapulu(bytes.NewReader(data), now); err != nil {
		log.Panicf("Parsing failed: %v\n", err)
	}

	switch {
	case flagWholeDay:
		channels.ShowWholeDay()
	case flagToJSON:
		j, err := channels.ToJSON()
		if err != nil {
			log.Panicf("JSON marshal failed: %v", err)
		}
		fmt.Printf("%s\n", j)
	default:
		channels.ShowUpcoming(now)
	}
}
