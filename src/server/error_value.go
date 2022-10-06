package main

import "errors"

var (
	ErrNilPointer = errors.New("null pointer data")
	ErrOldValue   = errors.New("old value")
	ErrNoValue    = errors.New("no value")
)
