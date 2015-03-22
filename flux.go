package flux

import (
	"fmt"
	"regexp"
	"strings"
)

type Flux struct {
	pattern   []string
	prefixes  []string
	suffixes  []string
	modifiers []string

	regex *regexp.Regexp
}

// regexp.Compile
func (f *Flux) Compile() (*regexp.Regexp, error) {
	if f.regex == nil {
		if reg, err := regexp.Compile(f.String()); err != nil {
			f.regex = nil
			return nil, err
		} else {
			f.regex = reg
		}
	}
	return f.regex, nil
}

// To reuse the compiled regexp
// regexp.MustCompile
func (f *Flux) MustCompile() *regexp.Regexp {
	return regexp.MustCompile(f.String())
}

// concatenate prefixes/pattern/suffixes/modifiers into a regular expression
// prints the current String representation from the regex
func (f *Flux) String() string {
	pattern := strings.Join(f.pattern, "")
	prefixes := strings.Join(f.prefixes, "")
	suffixes := strings.Join(f.suffixes, "")
	modifiers := strings.Join(f.modifiers, "")

	plainRegex := fmt.Sprintf("%s%s%s", prefixes, pattern, suffixes)
	if len(modifiers) > 0 {
		plainRegex = fmt.Sprintf("(?%s)%s", modifiers, plainRegex)
	}
	return plainRegex
}

func (f *Flux) Replace(src, repl string) string {
	return regexp.MustCompile(f.String()).ReplaceAllString(src, repl)
}

func (f *Flux) Match(value string) (bool, error) {
	r, err := regexp.Compile(f.String())

	if err == nil {
		return r.MatchString(value), nil
	}
	return false, err
}

// returns a map with all named matches
// experimental.... (Currently Only for groups added with #NamedGroup(name, rawGroup))
func (f *Flux) NamedMatches(s string) (map[string]string, error) {
	captures := make(map[string]string)

	r, err := f.Compile()
	if err != nil {
		return captures, err
	}

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures, nil
	}

	// get all named captures
	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}

	return captures, nil
}

// Clears all pattern components to create a fresh expression
func (f *Flux) Clear() *Flux {
	f.pattern = f.pattern[0:0]
	f.prefixes = f.prefixes[0:0]
	f.suffixes = f.suffixes[0:0]
	f.modifiers = f.modifiers[0:0]

	return resetCachedRegex(f)
}

func (f *Flux) RawGroup(value string) *Flux {
	return add(f, value, "(%s)")
}

//very experimental
func (f *Flux) NamedGroup(name, value string) *Flux {
	return raw(f, value, fmt.Sprintf("(?P<%s>%s)", name, "%s"))
}

// alias for Min(1)
func (f *Flux) Once() *Flux {
	return f.Length(1, 0)
}

func (f *Flux) Min(min int) *Flux {
	return f.Length(min, min-1)
}

func (f *Flux) Length(min, max int) *Flux {
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
// Helper for the ^ prefix
func (f *Flux) StartOfLine() *Flux {
	return addPrefix(f, "^")
}

// Helper for the $ suffix
func (f *Flux) EndOfLine() *Flux {
	return addSuffix(f, "$")
}

// Adds a modifier to ignore cases
func (f *Flux) IgnoreCase() *Flux {
	return addModifier(f, "i")
}

// Removes the 'm' modifier if it exists
func (f *Flux) OneLine() *Flux {
	return removeModifier(f, "m")
}

// Adds the 'm' modifier
func (f *Flux) Multiline() *Flux {
	return addModifier(f, "m")
}

func (f *Flux) MatchNewLine() *Flux {
	return addModifier(f, "s")
}

//--------------------------------------------------------------------------------
// @=LANGUAGE
//--------------------------------------------------------------------------------
// Alias to Flux#then
func (f *Flux) Find(value string) *Flux {
	return f.Then(value)
}

// Adds a search parameter
func (f *Flux) Then(value string) *Flux {
	return add(f, value, "(%s)")
}

// Optional search parameter
func (f *Flux) Maybe(value string) *Flux {
	return add(f, value, "(%s)?")
}

// Takes multiple arguments are creates an OR list.
// Output would be one|two|three etc
func (f *Flux) Either(values ...string) *Flux {
	return raw(f, strings.Join(values, "|"), "(%s)")
}

// Creates a [%s] search param
func (f *Flux) Any(value string) *Flux {
	return add(f, value, "([%s])")
}

// Adds a wildcard parameter
func (f *Flux) Anything() *Flux {
	return raw(f, ".*", "(%s)")
}

// Matches anything but the given arguments
func (f *Flux) AnythingBut(value string) *Flux {
	return add(f, value, "([^%s]*)")
}

func (f *Flux) LineBreak() *Flux {
	return raw(f, "(\\n|\\r\\n)", "%s")
}

func (f *Flux) Tab() *Flux {
	return raw(f, "(\\t)", "%s")
}

func (f *Flux) Word() *Flux {
	return raw(f, "(\\w+)", "%s")
}

func (f *Flux) Letters() *Flux {
	return raw(f, "([a-zA-Z]+)", "%s")
}

func (f *Flux) Digits() *Flux {
	return raw(f, "(\\d+)", "%s")
}

// experimental...
// This is bound to change
func (f *Flux) OrTry() *Flux {
	addPrefix(f, "(")
	addSuffix(f, ")")
	return raw(f, "", ")|(%s")
}

// Creates a range character class
// You can create a-z0-9 by calling Flux.range("a", "z", "0", "9")
func (f *Flux) Range(values ...string) *Flux {
	// validate pramas
	if len(values)%2 != 0 {
		return f
	}

	ranges := []string{}
	for i := 1; i < len(values); i += 2 {
		ranges = append(ranges, fmt.Sprintf("%s-%s", values[i-1], values[i]))
	}
	return raw(f, strings.Join(ranges, ""), "([%s])")
}

// creates a new Flux instance
func NewFlux() *Flux {
	newFlux := Flux{}
	return &newFlux
}

//--------------------------------------------------------------------------------
// HELPER FUNCTIONS
//--------------------------------------------------------------------------------

func add(f *Flux, value, format string) *Flux {
	f.pattern = append(f.pattern, fmt.Sprintf(format, regexp.QuoteMeta(value)))
	return resetCachedRegex(f)
}

func raw(f *Flux, value, format string) *Flux {
	f.pattern = append(f.pattern, fmt.Sprintf(format, value))
	return resetCachedRegex(f)
}

func addPrefix(f *Flux, prefix string) *Flux {
	if !stringInSlice(prefix, f.prefixes) {
		f.prefixes = append(f.prefixes, strings.TrimSpace(prefix))
	}
	return resetCachedRegex(f)
}

func addSuffix(f *Flux, suffix string) *Flux {
	if !stringInSlice(suffix, f.suffixes) {
		f.suffixes = append(f.suffixes, strings.TrimSpace(suffix))
	}
	return resetCachedRegex(f)
}

func addModifier(f *Flux, modifier string) *Flux {
	if !stringInSlice(modifier, f.modifiers) {
		f.modifiers = append(f.modifiers, strings.TrimSpace(modifier))
	}
	return resetCachedRegex(f)
}

func removeModifier(f *Flux, modifier string) *Flux {
	for i, mod := range f.modifiers {
		if mod == modifier {
			f.modifiers = append(f.modifiers[:i], f.modifiers[i+1:]...)
		}
	}
	return resetCachedRegex(f)
}

func getLastSegmentKey(f *Flux) int {
	return len(f.pattern) - 1
}

func replaceQuantifierByKey(f *Flux, key int, replacement string) *Flux {
	subject := f.pattern[key]
	replacementPattern := "%s%s"

	if strings.LastIndex(subject, ")") != -1 {
		subject = strings.TrimRight(subject, ")")
		replacementPattern += ")"
	}

	subject = removeQuantifier(f, subject)
	f.pattern[key] = fmt.Sprintf(replacementPattern, subject, replacement)
	return resetCachedRegex(f)

}

func removeQuantifier(f *Flux, pattern string) string {
	if strings.LastIndex(pattern, "+") != -1 && strings.LastIndex(pattern, "\\+") == -1 {
		return strings.TrimRight(pattern, "+")
	} else if strings.LastIndex(pattern, "*") != -1 && strings.LastIndex(pattern, "\\*") == -1 {
		return strings.TrimRight(pattern, "*")
	} else if strings.LastIndex(pattern, "?") != -1 && strings.LastIndex(pattern, "\\?") == -1 {
		return strings.TrimRight(pattern, "?")
	}
	return pattern
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

func resetCachedRegex(f *Flux) *Flux {
	f.regex = nil
	return f
}
