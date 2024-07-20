package markov

import (
	"math/rand/v2"
	"slices"
	"strings"
)

func mapFunc[T any](
	s []T,
	f func(T) T,
) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

type Markov struct {
	words      map[string][]*string
	firstWords []*string
}

func New() *Markov {
	return &Markov{
		words:      map[string][]*string{},
		firstWords: []*string{},
	}
}

func (m *Markov) add(word, followup string) {
	// Follow up word list doesn't exist yet
	if _, ok := m.words[word]; !ok {
		m.words[word] = []*string{}
	}
	l := m.words[word]
	lowered := strings.ToLower(followup)
	l = append(l, &lowered)
	m.words[word] = l
}

func (m *Markov) AddSentence(fields []string) {
	lowered := mapFunc(fields, func(s string) string {
		return strings.ToLower(s)
	})
	m.firstWords = append(m.firstWords, &lowered[0])
	for i, word := range lowered {
		if i+1 > len(fields)-1 {
			continue
		}
		m.add(word, fields[i+1])
	}
}

func (m *Markov) randFirstWord() string {
	randIdx := rand.IntN(len(m.firstWords))
	return *m.firstWords[randIdx]
}

func (m *Markov) randFollowupFor(word string) string {
	keys := make([]string, len(m.words))
	i := 0
	for k := range m.words {
		keys[i] = k
		i++
	}

	lowered := strings.ToLower(word)
	if !slices.Contains(keys, lowered) {
		return ""
	}
	r := len(m.words[word]) - 1
	randIdx := 0
	if r > 0 {
		randIdx = rand.IntN(r + 1) // + 1 because it's exclusive
	}
	return *m.words[word][randIdx]
}

func (m *Markov) GenSentence() string {
	i := 0
	fullSentence := []string{}
	nextWord := m.randFirstWord()
	fullSentence = append(fullSentence, nextWord)
	for nextWord != "" {
		nextWord = m.randFollowupFor(nextWord)
		fullSentence = append(fullSentence, nextWord)

		i++
		// Avoid infinite loops
		if i > 64 {
			break
		}
	}
	return strings.Join(fullSentence, " ")
}
