package hw09structvalidator

import (
	"errors"
	"fmt"
	"strconv"
)

type ValidationError struct {
	Field string
	Error error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	acc := "validation erros: " + strconv.Itoa(len(v))
	for i, e := range v {
		acc += fmt.Sprintf("\n%d. validation error in field %s: %v", i, e.Field, e.Error)
	}
	return acc
}

func (v ValidationError) IsError() bool {
	return v.Error != nil
}

type ProgramError struct {
	Msg string
}

func (e ProgramError) Error() string {
	return e.Msg
}

// 1. validaton errors -> STRINGS.
// 1.1 String must be exactly N characters.
type StringLengthError struct {
	Expected int
	Fact     int
}

func NewStrictStringLengthError(fact, expected int) error {
	return &StringLengthError{Expected: expected, Fact: fact}
}

func (e *StringLengthError) Error() string {
	return fmt.Sprintf("string must be exactly %d chars, got %d", e.Expected, e.Fact)
}

func (e *StringLengthError) Is(target error) bool {
	var sErr *StringLengthError
	ok := errors.As(target, &sErr)
	return ok && e.Expected == sErr.Expected && e.Fact == sErr.Fact
}

// 1.2. String must match regexp.
type StringRegExpError struct {
	Regexp string
}

func NewStringRegExpError(r string) error {
	return &StringRegExpError{Regexp: r}
}

func (e *StringRegExpError) Error() string {
	return fmt.Sprintf("string must match with regexp %s", e.Regexp)
}

func (e *StringRegExpError) Is(target error) bool {
	var sErr *StringRegExpError
	return errors.As(target, &sErr)
}

// 1.3. String must be one of the predefined values.
type StringNotAllowedError struct {
	Allowed string
	Fact    string
}

func NewStringNotAllowedError(s, vals string) error {
	return &StringNotAllowedError{Allowed: vals, Fact: s}
}

func (e *StringNotAllowedError) Error() string {
	return fmt.Sprintf("string must be one of: %v, got %v", e.Allowed, e.Fact)
}

func (e *StringNotAllowedError) Is(target error) bool {
	var sErr *StringNotAllowedError
	return errors.As(target, &sErr)
}

// 2. validaton errors -> INT
// 2.1. Lower bound for int field.
type IntMustBeLargerThanError struct {
	Expected int
	Fact     int
}

func NewIntMustBeLargerThanError(expected, fact int) error {
	return &IntMustBeLargerThanError{Expected: expected, Fact: fact}
}

func (e *IntMustBeLargerThanError) Error() string {
	return fmt.Sprintf("must be larger than: %d, got: %d", e.Expected, e.Fact)
}

func (e *IntMustBeLargerThanError) Is(target error) bool {
	var sErr *IntMustBeLargerThanError
	ok := errors.As(target, &sErr)
	return ok && e.Expected == sErr.Expected && e.Fact == sErr.Fact
}

// 2.2. Upper bound for int field.
type IntMustBeLowerThanError struct {
	Expected int
	Fact     int
}

func NewIntMustBeLowerThanError(expected, fact int) error {
	return &IntMustBeLowerThanError{Expected: expected, Fact: fact}
}

func (e *IntMustBeLowerThanError) Error() string {
	return fmt.Sprintf("must be larger than: %d, got: %d", e.Expected, e.Fact)
}

func (e *IntMustBeLowerThanError) Is(target error) bool {
	var sErr *IntMustBeLowerThanError
	ok := errors.As(target, &sErr)
	return ok && e.Expected == sErr.Expected && e.Fact == sErr.Fact
}

// 2.3. Int must be one of the predefined values.
type IntNotAllowedError struct {
	Allowed []string
	Fact    int
}

func NewIntNotAllowedError(fact int, allowed []string) error {
	return &IntNotAllowedError{Allowed: allowed, Fact: fact}
}

func (e *IntNotAllowedError) Error() string {
	return fmt.Sprintf("must be one of: %v, got %v", e.Allowed, e.Fact)
}

func (e *IntNotAllowedError) Is(target error) bool {
	var sErr *IntNotAllowedError
	return errors.As(target, &sErr)
}
