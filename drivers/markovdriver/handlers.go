package markovdriver

import (
	"github.com/fluffle/sp0rkle/base"
	//	"github.com/fluffle/sp0rkle/bot"
	"fmt"
	"strings"
)

func isStorable(line *base.Line) bool {
	return !line.Addressed
}

func processWord(word string) string {
	// TODO: Discard urls, dots, spaces
	word = strings.TrimSpace(word)
	word = strings.ToLower(word)
	return word
}

func recordMarkov(line *base.Line) {
	nick, _ := line.Storable()
	sentence := line.Args[1]

	if !isStorable(line) {
		return
	}

	words := strings.Split(sentence, " ")
	output := make([]string, 0, len(words))
	for _, word := range words {
		clean_word := processWord(word)
		if len(clean_word) > 0 {
			output = append(output, clean_word)
		}
	}

	err := mc.AddSentence(output, "user:"+string(nick))
	fmt.Printf("%v", err)
}
