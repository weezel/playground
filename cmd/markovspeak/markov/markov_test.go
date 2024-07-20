package markov

import (
	"slices"
	"strings"
	"testing"
)

func stringDeref(s string) *string {
	return &s
}

func TestMarkov_randFirstWord(t *testing.T) {
	type fields struct {
		words      map[string][]*string
		firstWords []*string
	}
	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{
			name: "One item",
			fields: fields{
				words: map[string][]*string{
					"jorma": {
						stringDeref("first"),
					},
				},
				firstWords: []*string{
					stringDeref("first"),
				},
			},
			want: "first",
		},
		{
			name: "Two items",
			fields: fields{
				words: map[string][]*string{
					"jorma": {
						stringDeref("first"),
						stringDeref("star"),
					},
				},
				firstWords: []*string{
					stringDeref("first"),
				},
			},
			want: "first",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Markov{
				words:      tt.fields.words,
				firstWords: tt.fields.firstWords,
			}
			if got := m.randFirstWord(); got != tt.want {
				t.Errorf("Markov.randFirstWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarkov_AddSentence(t *testing.T) {
	type fields struct {
		words      map[string][]*string
		firstWords []*string
	}
	type args struct {
		fields []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				words:      map[string][]*string{},
				firstWords: []*string{},
			},
			args: args{
				fields: strings.Fields("Oh no you again"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Markov{
				words:      tt.fields.words,
				firstWords: tt.fields.firstWords,
			}
			m.AddSentence(tt.args.fields)
		})
	}
}

func TestMarkov_randFollowupFor(t *testing.T) {
	type fields struct {
		words      map[string][]*string
		firstWords []*string
	}
	type args struct {
		word string
	}
	tests := []struct {
		name   string
		args   args
		want   []string
		fields fields
	}{
		{
			name: "Simple case",
			fields: fields{
				words: map[string][]*string{
					"first": {
						stringDeref("ball"),
					},
				},
				firstWords: []*string{},
			},
			args: args{
				word: "first",
			},
			want: []string{"ball"},
		},
		{
			name: "Two cases",
			fields: fields{
				words: map[string][]*string{
					"first": {
						stringDeref("star"),
						stringDeref("ball"),
					},
				},
				firstWords: []*string{},
			},
			args: args{
				word: "first",
			},
			want: []string{"star", "ball"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Markov{
				words:      tt.fields.words,
				firstWords: tt.fields.firstWords,
			}
			if got := m.randFollowupFor(tt.args.word); !slices.Contains(tt.want, got) {
				t.Errorf("Markov.randFollowupFor() = %v, want %v", got, tt.want)
			}
		})
	}
}
