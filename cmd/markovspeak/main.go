package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"
	"strings"

	"weezel/playground/cmd/markovspeak/markov"
)

var omittable = []string{
	"--",
	"-->",
	"<--",
	"▬▬▶",
	"◀▬▬",
}

type UserComment struct {
	Nick    string
	Message string
}

var ErrLine = errors.New("malfunctioned line")

func parseLineToObject(fields []string) (UserComment, error) {
	if len(fields) < 4 {
		return UserComment{}, errors.Join(ErrLine)
	}

	return UserComment{
		Nick:    fields[2],
		Message: strings.Join(fields[3:], " "),
	}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Missing filename")
		os.Exit(1)
	}

	mf, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer mf.Close()

	cf, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer cf.Close()
	if err = pprof.StartCPUProfile(cf); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m := markov.New()

	// Expect these to be weechat logs in the following format:
	// 2012-09-08 20:05:58     WeeZeL  very message such wow
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		switch {
		case len(fields) < 4:
			continue
		case strings.HasPrefix(fields[3], "@"):
			continue
		case slices.Contains(omittable, fields[2]):
			continue
		}

		var uc UserComment
		uc, err = parseLineToObject(fields)
		if err != nil {
			log.Panicf("Failed on line %q: %v\n", line, err)
		}

		m.AddSentence(strings.Fields(uc.Message))
	}

	for range 10 {
		fmt.Println(m.GenSentence())
	}

	if err = pprof.WriteHeapProfile(mf); err != nil {
		panic(err)
	}
}
