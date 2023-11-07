package importer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEventImportable_exclude(t *testing.T) {
	exclude := []int{1, 2, 3, 4}
	include := []int{}
	events := [3]int{0, 4, 5}
	asserts := [3]bool{true, false, true}

	for i, event := range events {
		assert.Equal(t, asserts[i], IsEventImportable(event, exclude, include))
	}
}

func TestIsEventImportable_include(t *testing.T) {
	var exclude []int
	include := []int{1, 2, 3, 4}
	events := [3]int{0, 4, 5}
	asserts := [3]bool{false, true, false}

	for i, event := range events {
		assert.Equal(t, asserts[i], IsEventImportable(event, exclude, include))
	}
}

func TestIsEventImportable_all(t *testing.T) {
	var exclude []int
	var include []int
	events := [3]int{0, 4, 5}

	for _, event := range events {
		assert.Equal(t, true, IsEventImportable(event, exclude, include))
	}
}

func TestIsEventImportable_nothing(t *testing.T) {
	exclude := []int{0, 4, 5}
	events := [3]int{0, 4, 5}

	for _, event := range events {
		assert.Equal(t, false, IsEventImportable(event, exclude, nil))
	}
}
