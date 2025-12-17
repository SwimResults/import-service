package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetDurationFromTimeString(t *testing.T) {
	d2, _ := time.ParseDuration("2m34.69s")
	println(d2)
}

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

func TestEasyWkReactionToDuration_Simple(t *testing.T) {
	ti := 314

	d, _ := EasyWkReactionToDuration(ti)

	assert.Equal(t, "3.14s", d.String())
}

func TestEasyWkReactionToDuration_Zero(t *testing.T) {
	ti := 0
	d, _ := EasyWkReactionToDuration(ti)
	assert.Equal(t, "0s", d.String())
}

func TestEasyWkReactionToDuration_Max(t *testing.T) {
	ti := 999
	d, _ := EasyWkReactionToDuration(ti)
	assert.Equal(t, "9.99s", d.String())
}

func TestAlgeTimeToDuration(t *testing.T) {
	ti := 693060
	d, _ := AlgeTimeToDuration(ti)
	assert.Equal(t, "1m9.306s", d.String())
}
