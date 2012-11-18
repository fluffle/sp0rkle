package markov

import (
	"errors"
	"fmt"
	"math/rand"
)

var SENTENCE_START string = "!SENTENCE_START"
var SENTENCE_END string = "!SENTENCE_END"

var ERROR_OVERFLOW = errors.New("Overflowed output bytes")
var NOT_ENOUGH_DATA = errors.New("Not enough data")

type Link struct {
	Value string
	Uses  int
}

type Source interface {
	GetLinks(value string) ([]Link, error)
}

func Generate(data Source, start_word string, seed int64, max_bytes int) (string, error) {
	current_node := start_word
	output := ""

	random := rand.New(rand.NewSource(seed))

	for len(output) < max_bytes {
		fmt.Printf("%s - %s\n", output, current_node)
		children, err := data.GetLinks(current_node)
		if err != nil {
			return output, err
		}

		sum := 0

		for _, child := range children {
			sum += child.Uses
		}
		if sum == 0 {
			return output, NOT_ENOUGH_DATA
		}

		r := random.Intn(sum)

		for _, child := range children {
			r -= child.Uses
			if r >= 0 {
				continue
			}
			if child.Value == SENTENCE_END {
				return output, nil
			}
			output += " " + child.Value
			current_node = child.Value
			break
		}
	}
	return output, ERROR_OVERFLOW

}
