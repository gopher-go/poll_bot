package poll_bot

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyboardForOptions(t *testing.T) {
	kb := keyboardFromOptions(nil, []string{"one", "two"})
	require.Len(t, kb.Buttons, 2)
}
