package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPollItemsRand(t *testing.T) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(99))
	s, err := newTestStorage()
	require.NoError(t, err)
	err = s.init()
	require.NoError(t, err)
	su, e := s.Obtain("0")
	require.NoError(t, e)
	poll := generateOurPoll()
	items := poll.items
	for i := 0; i < poll.size; i++ {
		answers := items[i].possibleAnswers
		if answers != nil {
			a := int(r.Float32() * float32(len(answers)))
			if items[i].validateAnswer != nil {
				items[i].validateAnswer(answers[a])
			}
			if items[i].persistAnswer != nil {
				e := items[i].persistAnswer(answers[a], su)
				require.NoError(t, e)
			}

		}

	}
}
