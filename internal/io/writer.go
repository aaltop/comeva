package io

import (
	"io"
	"os"
)

// modifierFunc represents a function that is given a value, may modify
// it with side-effects, and returns the result of the modification,
// which is generally just a value of the same type, and does not necessarily
// have to be based on the passed value.
type modifierFunc[T any] func(in T) (out T)

type modifier[T any] []modifierFunc[T]

// Modifier composes a number of modifier functions.
type Modifier[T any] interface {
	// Add adds the modifier functions in the given order to the modifier.
	Add(modifiers ...modifierFunc[T])
	// Apply applies the modifiers to the passed value. If no
	// modifiers have been added using [Modifier.Add], this
	// should be identity.
	Apply(in T) (out T)
}

func NewDefaultModifier[T any]() (modif Modifier[T]) {
	return &modifier[T]{}
}

func (modif *modifier[T]) Add(modifiers ...modifierFunc[T]) {
	var temp = modifier[T]{}
	temp = append(*modif, modifiers...)
	*modif = temp
}

func (modifier *modifier[T]) Apply(in T) (out T) {
	if len(*modifier) == 0 {
		return in
	}
	out = in
	for _, modi := range *modifier {
		out = modi(out)
	}
	return
}

// modifierWriter implements [io.Writer] and has a [Modifier] that intercepts and modifies
// any written bytes.
type modifierWriter struct {
	Modifier Modifier[[]byte]
	out      io.Writer
}

func (writer *modifierWriter) Write(p []byte) (n int, e error) {

	var temp = writer.Modifier.Apply(p)
	_, e = writer.out.Write(temp)
	// doesn't say that couldn't return non-nil error even if n is len(p).
	return len(p), e
}

// StringModifierWriter implements [io.Writer] and [io.StringWriter] and has a
// [Modifier] that intercepts and modifies any written string.
type StringModifierWriter struct {
	Modifier Modifier[string]
	out      io.Writer
}

func (writer *StringModifierWriter) Write(p []byte) (n int, e error) {
	return writer.WriteString(string(p))
}

func (writer *StringModifierWriter) WriteString(s string) (n int, e error) {
	var temp = writer.Modifier.Apply(s)
	n, e = writer.out.Write(append([]byte{}, temp...))
	return
}

func NewDefaultStringModifierWriter() (writer *StringModifierWriter) {
	var modifier = NewDefaultModifier[string]()
	return NewStringModifierWriter(modifier, os.Stdout)
}

func NewStringModifierWriter(modifier Modifier[string], out io.Writer) (writer *StringModifierWriter) {
	return &StringModifierWriter{
		Modifier: modifier,
		out:      out,
	}
}
