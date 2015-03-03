package flux

import (
	"fmt"
	"testing"
)

func TestCompile(t *testing.T) {
	regex := NewFlux()
	assertEquals(t, regex.Compile(), "//")
}

func TestAddPrefix(t *testing.T) {
	regex := NewFlux().StartOfLine()
	assertEquals(t, regex.Compile(), "/^/")
}


func TestAddSuffix(t *testing.T) {
	regex := NewFlux().EndOfLine()
	assertEquals(t, regex.Compile(), "/$/")
}

func TestModifiers(t *testing.T) {
	regex := NewFlux().Multiline()
	assertEquals(t, regex.Compile(), "//m")

	regex.IgnoreCase()
	assertEquals(t, regex.Compile(), "//mi")

	regex = NewFlux().Multiline().IgnoreCase().MatchNewLine()
	assertEquals(t, regex.Compile(), "//mis")
}

func TestThen(t *testing.T) {
	regex := NewFlux().Then("required")
	assertEquals(t, regex.Compile(), "/(required)/")
}

func TestMaybe(t *testing.T) {
	regex := NewFlux().Maybe("optional")
	assertEquals(t, regex.Compile(), "/(optional)?/")
}

func TestAny(t *testing.T) {
	regex := NewFlux().Any("abc")
	assertEquals(t, regex.Compile(), "/([abc])/")
}

func TestAnything(t *testing.T) {
	regex := NewFlux().Anything()
	assertEquals(t, regex.Compile(), "/(.*)/")
}

func TestAnythingBut(t *testing.T) {
	regex := NewFlux().AnythingBut("BUT")
	assertEquals(t, regex.Compile(), "/([^BUT]*)/")
}

func TestEither(t *testing.T) {
	regex := NewFlux().Either("one", "two", "three")
	assertEquals(t, regex.Compile(), "/(one|two|three)/")
}

func TestOrTry(t *testing.T) {
	// TODO implement me
	t.SkipNow()
}

func TestRange(t *testing.T) {
	regex := NewFlux().Range("a", "z", "0", "9")
	assertEquals(t, regex.Compile(), "/([a-z0-9])/")
}


func TestLength(t *testing.T) {
	regex := NewFlux().Word().Length(1, 5)
	assertEquals(t, regex.Compile(), "/(\\w{1,5})/")
}

func TestMin(t *testing.T) {
	regex := NewFlux().Word().Min(1)
	assertEquals(t, regex.Compile(), "/(\\w{1})/")
}

func assertEquals(t *testing.T, actual, expected string) {
	if expected == actual {
		t.Log("Empty Compile Test passed")
	} else {
		t.Error(fmt.Sprintf("Assertion failed! Expected: %s, Actual: %s", expected, actual))
	}
}
