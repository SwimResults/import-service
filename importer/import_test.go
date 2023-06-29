package importer

import "testing"

func TestDoTheMagic(t *testing.T) {
	err := DoTheMagic()
	if err != nil {
		panic(err)
	}
}

func TestDoTheResultsMagic(t *testing.T) {
	err := DoTheResultsMagic()
	if err != nil {
		panic(err)
	}
}
