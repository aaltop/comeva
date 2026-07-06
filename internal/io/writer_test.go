package io

import (
	"strings"
	"testing"

	testingUtils "github.com/aaltop/comeva/internal/testing"
)

func FixtureStringModifier(counter *stringLengthCounter) (modifier Modifier[string]) {
	modifier = NewDefaultModifier[string]()

	// add slash to start, calculate and store the length without modification,
	// add another slash to start.
	modifier.Add(addSlash, func(in string) (out string) { counter.Count(in); return in }, addSlash)
	return modifier
}

func addSlash(in string) (out string) {
	return "/" + in
}

type stringLengthCounter struct {
	stringLength uint
}

func (counter *stringLengthCounter) Count(value string) {
	counter.stringLength = uint(len(value))
}

func textModifierTest(t *testing.T, fun modifierFunc[string], counter *stringLengthCounter) {
	var in = "string"
	var outFirst = addSlash(in)
	var expectedLength = uint(len(outFirst))
	var expected = addSlash(outFirst)

	var received = fun(in)

	if expected != received {
		t.Error(testingUtils.ValueMismatch("string", expected, received))
	}

	var receivedLength = counter.stringLength
	if receivedLength != expectedLength {
		t.Error(testingUtils.ValueMismatch("uint", expectedLength, receivedLength))
	}
}

// Default TextModifier works.
func TestDefaultTextModifier(t *testing.T) {

	var counter = &stringLengthCounter{}
	var modifier Modifier[string] = FixtureStringModifier(counter)
	textModifierTest(t, modifier.Apply, counter)

}

// TextModifierWriter works.
func TestTextModifierWriter(t *testing.T) {

	var counter = &stringLengthCounter{}
	var modifier = FixtureStringModifier(counter)
	var stringBuilder = &strings.Builder{}
	var writer = NewStringModifierWriter(modifier, stringBuilder)

	var modify = func(in string) (out string) {
		writer.Write(append([]byte{}, in...))
		return stringBuilder.String()
	}

	textModifierTest(t, modify, counter)
}
