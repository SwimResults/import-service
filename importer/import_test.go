package importer

import "testing"

func TestDoTheMagic(t *testing.T) {
	err := DoTheMagic()
	if err != nil {
		panic(err)
	}
}
