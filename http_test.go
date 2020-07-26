package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeyboardForOptions(t *testing.T) {
	kb := keyboardFromOptions(nil, []string{"one", "two"})
	require.Len(t, kb.Buttons, 2)
}
