package flux


import (
	"fmt"
	"strings"
	"regexp"
)

type Flux struct {
	pattern   []string
	prefixes  []string
	suffixes  []string
	modifiers []string
}

// Compiles prefixes/pattern/suffixes/modifiers into a regular expression
func (f *Flux) Compile() {
	pattern   := strings.Join(f.pattern, "")
	prefixes  := strings.Join(f.prefixes, "")
	suffixes  := strings.Join(f.suffixes, "")
	modifiers := strings.Join(f.modifiers, "")

	return fmt.Sprintf("/%s%s%s/%s", prefixes, pattern, suffixes, modifiers)
}

// Clears all pattern components to create a fresh expression
func (f *Flux) Clear() (*Flux) {
	f.pattern   = f.pattern[0:0]
	f.prefixes  = f.prefixes[0:0]
	f.suffixes  = f.suffixes[0:0]
	f.modifiers  = f.modifiers[0:0]

	return f
}

func (f *Flux) RawGroup(value string) (*Flux) {
	return add(f, value, "(%s)")
}

func (f *Flux) Once() (*Flux) {
	return f.Length(1, 0)
}
func (f *Flux) Min(min int) (*Flux) {
	return f.Length(min, min - 1)
}

func (f *Flux) Max(max int) (*Flux) {
	return f.Length(max, max)
}

func (f *Flux) Length(min, max int) (*Flux) {
	lastSegmentKey := getLastSegmentKey(f)
	var lengthPattern string

	if max > min {
		lengthPattern = fmt.Sprintf("{%d,%d}", min, max)
	} else {
		lengthPattern = fmt.Sprintf("{%d}", min)
	}

	return replaceQuantifierByKey(f, lastSegmentKey, lengthPattern)

}


//--------------------------------------------------------------------------------
// MODIFIERS
//--------------------------------------------------------------------------------
func (f *Flux) StartOfLine() (*Flux) {
	return addPrefix(f, "")
}

func (f *Flux) EndOfLine() (*Flux) {
	return addSuffix(f, "$")
}

func (f *Flux) IgnoreCase() (*Flux) {
	return addModifier(f, "i")
}

func (f *Flux) OneLine() (*Flux) {
	return removeModifier(f, "m")
}

func (f *Flux) Multiline() (*Flux) {
	return addModifier(f, "m")
}

func (f *Flux) MatchNewLine() (*Flux) {
	return addModifier(f, "s")
}

//--------------------------------------------------------------------------------
// @=LANGUAGE
//--------------------------------------------------------------------------------
func (f *Flux) Find(value string) (*Flux) {
	return f.Then(value)
}

func (f *Flux) Then(value string) (*Flux) {
	return add(f, value, "(%s)")
}

func (f *Flux) Maybe(value string) (*Flux) {
	return add(f, value, "(%s)?")
}

func (f *Flux) Either(values ...string) (*Flux) {
	return raw(f, strings.Join(values, "|"), "(%s)")
}

func (f *Flux) Any(value string) (*Flux) {
	return add(f, value, "([%s])")
}

func (f *Flux) Anything() (*Flux) {
	return raw(f, ".*", "(%s)")
}

func (f *Flux) AnythingBut(value string) (*Flux) {
	return add(f, value, "([^%s]*)")
}

func (f *Flux) LineBreak() (*Flux) {
	return raw(f, "(\\n|\\r\\n)", "%s")
}

func (f *Flux) Tab() (*Flux) {
	return raw(f, "(\\t)", "%s")
}

func (f *Flux) Word() (*Flux) {
	return raw(f, "(\\w+)", "%s")
}

func (f *Flux) Letters() (*Flux) {
	return raw(f, "([a-zA-Z]+)", "%s")
}

func (f *Flux) Digits() (*Flux) {
	return raw(f, "(\\d+)", "%s")
}

func (f *Flux) OrTry(value string) (*Flux) {
	addPrefix(f, "(")
	addSuffix(f, ")")
	return raw(f, value, ")|((%s)")
}

// TODO add range function

func NewFlux() Flux {
	return Flux{}
}

func add(f *Flux, value, format string) (*Flux) {
	f.pattern = append(f.pattern, fmt.Sprintf(format, regexp.QuoteMeta(value)))
	return f
}

func raw(f *Flux, value, format string) (*Flux) {
	f.pattern = append(f.pattern, fmt.Sprintf(format, value))
	return f
}

func addPrefix(f *Flux, prefix string) (*Flux) {
	if !stringInSlice(prefix, f.prefixes) {
		f.prefixes = append(f.prefixes, strings.TrimSpace(prefix))
	}
	return f
}

func addSuffix(f *Flux, suffix string) (*Flux) {
	if !stringInSlice(suffix, f.suffixes) {
		f.suffixes = append(f.suffixes, strings.TrimSpace(suffix))
	}
	return f
}

func addModifier(f *Flux, modifier string) (*Flux) {
	if !stringInSlice(modifier, f.modifiers) {
		f.modifiers = append(f.modifiers, strings.TrimSpace(modifier))
	}
	return f
}

func removeModifier(f *Flux, modifier string) (*Flux) {
	for i, mod := range f.modifiers {
		if mod == modifier {
			f.modifiers = append(f.modifiers[:i], f.modifiers[i+1:]...)
		}
	}
	return f
}

func getLastSegmentKey(f *Flux) (int) {
	// TODO implement me
	return 0
}

func replaceQuantifierByKey(f *Flux, key int, replacement string) (*Flux) {
	// TODO implement me
	return f
}

// checks if the given string is in the given slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
