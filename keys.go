package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
)

var (
	keyQuitForce         = key.NewBinding(key.WithKeys("ctrl+c"))
	keyQuitWithDirectory = key.NewBinding(key.WithKeys("ctrl+d"))
	keyQuitWithSelected  = key.NewBinding(key.WithKeys("ctrl+x"))

	keyEsc           = key.NewBinding(key.WithKeys("esc"))
	keySelect        = key.NewBinding(key.WithKeys("enter"))
	keyBack          = key.NewBinding(key.WithKeys("backspace"))
	keyTab           = key.NewBinding(key.WithKeys("tab"))
	keyFileSeparator = key.NewBinding(key.WithKeys("/"))

	keyUp    = key.NewBinding(key.WithKeys("up", "k"))
	keyDown  = key.NewBinding(key.WithKeys("down", "j"))
	keyLeft  = key.NewBinding(key.WithKeys("left", "h"))
	keyRight = key.NewBinding(key.WithKeys("right", "l"))

	keyDebugMode  = key.NewBinding(key.WithKeys("d"))
	keyHelpMode   = key.NewBinding(key.WithKeys("H"))
	keySearchMode = key.NewBinding(key.WithKeys("i"))

	keyToggleFollowSymlink = key.NewBinding(key.WithKeys("f"))
	keyToggleHidden        = key.NewBinding(key.WithKeys("a"))
	keyToggleList          = key.NewBinding(key.WithKeys("L"))

	keyDismissError = key.NewBinding(key.WithKeys("e"))
)

type remappedEscKey struct {
	key     key.Binding
	presses int
	pressed int
}

func (k *remappedEscKey) triggered() bool {
	k.pressed++
	if k.pressed == k.presses {
		k.pressed = 0
		return true
	}
	return false
}

func (m *model) setEscRemapKey() {
	m.esc = &remappedEscKey{
		key:     key.NewBinding(key.WithKeys("")), // No-op key binding.
		presses: 1,
	}

	escRemap := os.Getenv(envEscRemap)
	if escRemap == "" {
		return
	}

	keyRune := rune(escRemap[0])
	if unicode.IsLetter(keyRune) || unicode.IsDigit(keyRune) {
		m.setError(
			fmt.Errorf("remapped escape key [%s] must not be alphanumeric", string(escRemap[0])),
			"invalid remapped esc key",
		)
		return
	}

	m.esc.key = key.NewBinding(key.WithKeys(string(escRemap[0])))
	if len(escRemap) == 1 {
		return
	}
	presses, err := strconv.Atoi(escRemap[1:])
	if err != nil {
		m.setError(
			fmt.Errorf("remapped escape key [%s] must contain integer digits after the first character", escRemap),
			"invalid remapped esc key",
		)
		return
	}

	m.esc.presses = presses
}
