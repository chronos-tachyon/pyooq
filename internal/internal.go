package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chronos-tachyon/assert"
)

var reName = regexp.MustCompile(`^[A-Za-z][0-9A-Za-z]*(?:_[0-9A-Za-z]+)*$`)

type Void = struct{}

func IsValidName(str string) bool {
	return reName.MatchString(str)
}

func AssertValidName(kind string, str string) {
	if IsValidName(str) {
		return
	}
	assert.Raisef("%q is not a valid %s name", str, kind)
}

type IsValider interface{ IsValid() bool }

func AssertValid(obj IsValider) {
	if obj != nil && obj.IsValid() {
		return
	}
	assert.Raisef("%#v is not valid", obj)
}

func AssertIs[T comparable](preamble string, actual T, expect T) {
	if actual == expect {
		return
	}
	if preamble == "" {
		assert.Raisef("wrong %T; expected %v, but it's actually %v", actual, expect, actual)
		return
	}
	assert.Raisef("%s; expected %v, but it's actually %v", preamble, expect, actual)
}

func AssertOneOf[T comparable](preamble string, actual T, list ...T) {
	if IsOneOf(actual, list...) {
		return
	}
	assert.Raisef("%s; expected one of %v, but it's actually %v", preamble, list, actual)
}

func IsOneOf[T comparable](actual T, list ...T) bool {
	for _, expect := range list {
		if actual == expect {
			return true
		}
	}
	return false
}

func NotValidError(obj any) error {
	return fmt.Errorf("%#v is not valid", obj)
}

func Contains[K comparable, V any](in map[K]V, key K) bool {
	_, found := in[key]
	return found
}

func Lookup[K comparable, V any](in map[K]V, key K, format string, args ...any) V {
	if value, found := in[key]; found {
		return value
	}
	panic(fmt.Errorf(format, args...))
}

func CopyList[T any](in []T) []T {
	size := len(in)
	if size <= 0 {
		return nil
	}
	out := make([]T, size)
	copy(out, in)
	return out
}

func CopyMap[K comparable, V any](in map[K]V) map[K]V {
	size := len(in)
	if size <= 0 {
		return nil
	}
	out := make(map[K]V, size)
	for k, v := range in {
		out[k] = v
	}
	return out
}

func RecycleList[T any](ptr *[]T, initCapacity uint) {
	in := *ptr
	if in == nil {
		*ptr = make([]T, 0, initCapacity)
		return
	}
	var zero T
	for index := range in {
		in[index] = zero
	}
	*ptr = in[:0]
}

func RecycleMap[K comparable, V any](ptr *map[K]V, initCapacity uint) {
	in := *ptr
	if in == nil {
		*ptr = make(map[K]V, initCapacity)
		return
	}
	clear(in)
	*ptr = in
}

func BuildEnumMap[T ~byte](size T, goNames []string, names []string) map[string]T {
	m := make(map[string]T, 64)
	for x := T(0); x < size; x++ {
		goName, name := goNames[x], names[x]
		m[goName] = x
		m[name] = x
	}
	return m
}

func FinishEnumMap[T ~byte](m map[string]T) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for _, k := range keys {
		v := m[k]
		lc := strings.ToLower(k)
		uc := strings.ToUpper(k)
		m[lc] = v
		m[uc] = v
	}
}

func RangeErr(index uint, size uint) error {
	return fmt.Errorf("index out of bounds: %d >= %d", index, size)
}
