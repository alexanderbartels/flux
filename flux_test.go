package flux

import (
	"fmt"
	"testing"
)

func TestCompile(t *testing.T) {
	regex := NewFlux()
	assertEquals(t, regex.String(), "")

	res, err := regex.Compile()
	assertNoError(t, err)
	assertEquals(t, regex.String(), res.String())
	assertEquals(t, regex.String(), regex.MustCompile().String())

	regex = regex.Any("abcdef")
	assertEquals(t, regex.String(), "([abcdef])")

	res, err = regex.Compile()
	assertNoError(t, err)
	assertEquals(t, regex.String(), res.String())
	assertEquals(t, regex.String(), regex.MustCompile().String())
}

func TestAddPrefix(t *testing.T) {
	regex := NewFlux().StartOfLine()
	assertEquals(t, regex.String(), "^")
}


func TestAddSuffix(t *testing.T) {
	regex := NewFlux().EndOfLine()
	assertEquals(t, regex.String(), "$")
}

func TestModifiers(t *testing.T) {
	regex := NewFlux().Multiline()
	assertEquals(t, regex.String(), "(?m)")

	regex.IgnoreCase()
	assertEquals(t, regex.String(), "(?mi)")

	regex = NewFlux().Multiline().IgnoreCase().MatchNewLine()
	assertEquals(t, regex.String(), "(?mis)")
}

func TestNamedGroup(t *testing.T) {
	regex := NewFlux().NamedGroup("myGroup", "[A-zA-Z]+")
	assertEquals(t, regex.String(), "(?P<myGroup>[A-zA-Z]+)")
}

func TestThen(t *testing.T) {
	regex := NewFlux().Then("required")
	assertEquals(t, regex.String(), "(required)")
}

func TestMaybe(t *testing.T) {
	regex := NewFlux().Maybe("optional")
	assertEquals(t, regex.String(), "(optional)?")
}

func TestAny(t *testing.T) {
	regex := NewFlux().Any("abc")
	assertEquals(t, regex.String(), "([abc])")
}

func TestAnything(t *testing.T) {
	regex := NewFlux().Anything()
	assertEquals(t, regex.String(), "(.*)")
}

func TestAnythingBut(t *testing.T) {
	regex := NewFlux().AnythingBut("BUT")
	assertEquals(t, regex.String(), "([^BUT]*)")
}

func TestEither(t *testing.T) {
	regex := NewFlux().Either("one", "two", "three")
	assertEquals(t, regex.String(), "(one|two|three)")
}

func TestOrTry(t *testing.T) {
	regex := NewFlux().StartOfLine().Find("dev.").Any("abc").OrTry().Maybe("live.").EndOfLine()
	assertEquals(t, regex.String(), "^((dev\\.)([abc]))|((live\\.)?)$")

	regex = regex.IgnoreCase()
	assertEquals(t, regex.String(), "(?i)^((dev\\.)([abc]))|((live\\.)?)$")
}

func TestRange(t *testing.T) {
	regex := NewFlux().Range("a", "z", "0", "9")
	assertEquals(t, regex.String(), "([a-z0-9])")
}


func TestLength(t *testing.T) {
	regex := NewFlux().Word().Length(1, 5)
	assertEquals(t, regex.String(), "(\\w{1,5})")

	regex = NewFlux().Letters()
	assertEquals(t, regex.String(), "([a-zA-Z]+)")
	regex = regex.Length(1, 5)
	assertEquals(t, regex.String(), "([a-zA-Z]{1,5})")
}

func TestMin(t *testing.T) {
	regex := NewFlux().Word().Min(1)
	assertEquals(t, regex.String(), "(\\w{1})")
}

// alias for Min(1)
func TestOnce(t *testing.T) {
	regex := NewFlux().Word().Once()
	assertEquals(t, regex.String(), "(\\w{1})")
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

func TestPhoneMatchReplace(t *testing.T) {
	phone := "6124240013"
	regex := NewFlux().StartOfLine().Maybe("(").Digits().Length(3,3).Maybe(")").Maybe(" ").Digits().Length(3,3).Maybe("-").Digits().Length(4,4).EndOfLine()

	match, err := regex.Match(phone)
	assertNoError(t, err)
	assertTrue(t, match)

	repl :=  regex.Replace(phone, "($2) $5-$7" ) // $2 -> Steht für die 2. Gruppe Digits().Length(3,3), $5 für die 5. Grupper, usw.
	assertEquals(t, repl, "(612) 424-0013" );
}

func TestUrlMatchAndReplace(t *testing.T) {
	regex := NewFlux().StartOfLine().Find("http").Maybe("s").Then("://").Maybe("www.").AnythingBut(".").Either(".co", ".com", ".de").IgnoreCase().EndOfLine();
	assertEquals(t, regex.String(), "(?i)^(http)(s)?(://)(www\\.)?([^\\.]*)(.co|.com|.de)$")

	match, err := regex.Match("http://selvinortiz.com")
	assertNoError(t, err)
	assertTrue(t, match)

	repl := regex.Replace("http://selvinortiz.com", "$5$6")
	assertEquals(t, repl, "selvinortiz.com")
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
