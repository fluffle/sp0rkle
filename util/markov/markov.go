package markov

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/fluffle/golog/logging"
)

const (
	ACTION_START   = "!ACTION_START"
	SENTENCE_START = "!SENTENCE_START"
	SENTENCE_END   = "!SENTENCE_END"
)

var NOT_ENOUGH_DATA = errors.New("Not enough data")
var TOO_MUCH_DATA = errors.New("Too much data")

type Link struct {
	Value string
	Uses  int
}

type Source interface {
	GetLinks(value string) ([]Link, error)
}

func Action(data Source) (string, error) {
	s, err := generate(data, ACTION_START, 50)
	return strings.Join(s, " "), err
}

func Sentence(data Source) (string, error) {
	s, err := generate(data, SENTENCE_START, 50)
	return strings.Join(s, " "), err
}

func generate(data Source, start string, length int) ([]string, error) {
	current, output := start, make([]string, 0, length)

	for len(output) < length {
		children, err := data.GetLinks(current)
		if err != nil {
			logging.Error("Error getting markov links: %v", err)
			return output, err
		}

		sum := 0
		for _, child := range children {
			sum += child.Uses
			if len(output) > 4*length/5 && child.Value == SENTENCE_END {
				// start to limit at 80% of length, prefer SENTENCE_END if valid
				return output, nil
			}
		}
		if sum == 0 {
			return output, NOT_ENOUGH_DATA
		}

		r := rand.Intn(sum)

		for _, child := range children {
			r -= child.Uses
			if r >= 0 {
				continue
			}
			if child.Value == SENTENCE_END {
				return output, nil
			}
			output = append(output, child.Value)
			current = child.Value
			break
		}
	}
	return output, TOO_MUCH_DATA
}
