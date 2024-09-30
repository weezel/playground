package programinfo

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var ChannelOrder = []string{
	"Yle TV1",
	"Yle TV2",
	"MTV3",
	"Nelonen",
	"Yle Teema & Fem",
	"MTV Sub",
	"TV5",
	"Liv",
	"Jim",
	"Kutonen",
	"TLC",
	"STAR Channel",
	"MTV Ava",
	"Hero",
	"Frii",
	"National Geographic",
}

type Program struct {
	StartTime time.Time `json:"start_time,omitempty"`
	Name      string    `json:"name,omitempty"`
}

func (p Program) String() string {
	// if flagShowDates {
	// 	return fmt.Sprintf("[%s] %s", p.StartTime.Format("2006-01-02 15:04"), p.Name)
	// }
	return fmt.Sprintf("%s  %s", p.StartTime.Format("15:04"), p.Name)
}

type Channels map[string][]Program

func (c Channels) ShowUpcoming(offset time.Time) {
	log.Println("Offset start time:", offset)
	for _, channel := range ChannelOrder {
		fmt.Printf("%s\n", channel)

		for i := 0; i < len(c[channel]); i++ {
			prog := c[channel][i]
			if i+1 < len(c[channel]) {
				isOnAir := offset.After(prog.StartTime) &&
					offset.Before(c[channel][i+1].StartTime)
				isAfter := offset.After(prog.StartTime)
				if !isOnAir && isAfter {
					continue
				}
			}

			fmt.Printf("\t%s\n", prog.String())
		}
	}
}

func (c Channels) ShowWholeDay() {
	for _, channel := range ChannelOrder {
		fmt.Printf("%s\n", channel)

		for _, prog := range c[channel] {
			fmt.Printf("\t%s\n", prog.String())
		}
	}
}

func (c Channels) GetChannelWholeDay(name string) string {
	buf := new(strings.Builder)
	for _, prog := range c[name] {
		buf.WriteString(prog.String() + "\n")
	}
	return buf.String()
}

func (c Channels) ToJSON() ([]byte, error) {
	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return j, nil
}
