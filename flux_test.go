package flux

import (
	"fmt"
	"testing"
)

func TestCompile(t *testing.T) {
	regex := NewFlux()
	assertEquals(t, regex.Compile(), "")
}

func TestAddPrefix(t *testing.T) {
	regex := NewFlux().StartOfLine()
	assertEquals(t, regex.Compile(), "^")
}


func TestAddSuffix(t *testing.T) {
	regex := NewFlux().EndOfLine()
	assertEquals(t, regex.Compile(), "$")
}

func TestModifiers(t *testing.T) {
	regex := NewFlux().Multiline()
	assertEquals(t, regex.Compile(), "(?m)")

	regex.IgnoreCase()
	assertEquals(t, regex.Compile(), "(?mi)")

	regex = NewFlux().Multiline().IgnoreCase().MatchNewLine()
	assertEquals(t, regex.Compile(), "(?mis)")
}

func TestThen(t *testing.T) {
	regex := NewFlux().Then("required")
	assertEquals(t, regex.Compile(), "(required)")
}

func TestMaybe(t *testing.T) {
	regex := NewFlux().Maybe("optional")
	assertEquals(t, regex.Compile(), "(optional)?")
}

func TestAny(t *testing.T) {
	regex := NewFlux().Any("abc")
	assertEquals(t, regex.Compile(), "([abc])")
}

func TestAnything(t *testing.T) {
	regex := NewFlux().Anything()
	assertEquals(t, regex.Compile(), "(.*)")
}

func TestAnythingBut(t *testing.T) {
	regex := NewFlux().AnythingBut("BUT")
	assertEquals(t, regex.Compile(), "([^BUT]*)")
}

func TestEither(t *testing.T) {
	regex := NewFlux().Either("one", "two", "three")
	assertEquals(t, regex.Compile(), "(one|two|three)")
}

func TestOrTry(t *testing.T) {
	regex := NewFlux().StartOfLine().Find("dev.").Any("abc").OrTry().Maybe("live.").EndOfLine()
	assertEquals(t, regex.Compile(), "^((dev\\.)([abc]))|((live\\.)?)$")

	regex = regex.IgnoreCase()
	assertEquals(t, regex.Compile(), "(?i)^((dev\\.)([abc]))|((live\\.)?)$")
}

func TestRange(t *testing.T) {
	regex := NewFlux().Range("a", "z", "0", "9")
	assertEquals(t, regex.Compile(), "([a-z0-9])")
}


func TestLength(t *testing.T) {
	regex := NewFlux().Word().Length(1, 5)
	assertEquals(t, regex.Compile(), "(\\w{1,5})")

	regex = NewFlux().Letters()
	assertEquals(t, regex.Compile(), "([a-zA-Z]+)")
	regex = regex.Length(1, 5)
	assertEquals(t, regex.Compile(), "([a-zA-Z]{1,5})")
}

func TestMin(t *testing.T) {
	regex := NewFlux().Word().Min(1)
	assertEquals(t, regex.Compile(), "(\\w{1})")
}

// alias for Min(1)
func TestOnce(t *testing.T) {
	regex := NewFlux().Word().Once()
	assertEquals(t, regex.Compile(), "(\\w{1})")
}

func TestSimpleMatch(t *testing.T) {
	match, err := NewFlux().Find("Hallo").Match("Hallo")
	assertNoError(t, err)
	assertTrue(t, match)
}

func TestMatch(t *testing.T) {
	regex := NewFlux().StartOfLine().Find("dev").EndOfLine()

	match, err := regex.Match("dev")
	assertNoError(t, err)
	assertTrue(t, match)

	match, err = regex.Match("DEV")
	assertNoError(t, err)
	assertTrue(t, !match) // match should be false

	regex.IgnoreCase()
	match, err = regex.Match("DEV")
	assertNoError(t, err)
	assertTrue(t, match) // now it should be true
}

func assertEquals(t *testing.T, actual, expected string) {
	if expected == actual {
		t.Log("Assertion passed.")
	} else {
		t.Error(fmt.Sprintf("Assertion failed! Expected: %s, Actual: %s", expected, actual))
	}
}

func assertTrue(t *testing.T, actual bool) {
	if actual {
		t.Log("Assertion passed.")
	} else {
		t.Error(fmt.Sprintf("Assertion failed! Expected to be true, but is: %v", actual))
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(fmt.Sprintf("Assertion failed! Error: %v", err))
	} else {
		t.Log("Assertion passed.")
	}
}
