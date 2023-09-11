package main

import (
	"errors"
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

func (k *remappedEscKey) reset() {
	k.pressed = 0
}

func (m *model) setEscRemapKey(escRemap string) error {
	if escRemap == "" {
		return errors.New("invalid remapped escape key: empty string provided")
	}

	k := escRemap[0]
	for i := 0; i < len(escRemap); i++ {
		kr := rune(escRemap[i])
		if unicode.IsLetter(kr) || unicode.IsDigit(kr) {
			return errors.New("remapped escape key must not be alphanumeric")
		}
		if escRemap[i] != k {
			return errors.New("remapped escape key must not contain different characters")
		}
	}

	m.esc = &remappedEscKey{
		key:     key.NewBinding(key.WithKeys(string(k))),
		presses: len(escRemap),
	}
	return nil
}

func defaultEscRemapKey() *remappedEscKey {
	return &remappedEscKey{
		key: key.NewBinding(key.WithKeys("")),
	}
}
