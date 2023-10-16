package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEasyWkTimeToDuration_Simple(t *testing.T) {
	ti := 31415

	d, _ := EasyWkTimeToDuration(ti)

	assert.Equal(t, "3m14.15s", d.String())
}

func TestEasyWkTimeToDuration_Zero(t *testing.T) {
	ti := 0
	d, _ := EasyWkTimeToDuration(ti)
	assert.Equal(t, "0s", d.String())
}

func TestEasyWkTimeToDuration_MaxDay(t *testing.T) {
	ti := 595999
	d, _ := EasyWkTimeToDuration(ti)
	assert.Equal(t, "59m59.99s", d.String())
}
