//go:generate stringer -type=Env -linecomment

package env

import "os"

const Mode = "_wzo_"

type Env int8

const (
	None Env = iota // none
	Dev             // dev
	Test            // test
	Uat             // uat
	Prod            // prod
)

var current = Dev

func Set(s string) {
	switch s {
	case Dev.String():
		current = Dev
	case Test.String():
		current = Test
	case Uat.String():
		current = Uat
	case Prod.String():
		current = Prod
	default:
		current = None
	}
}

func Get() Env {
	if current == None {
		env := os.Getenv(Mode)
		if env == "" {
			current = Dev
		} else {
			Set(env)
			if current == None {
				current = Dev
			}
		}
	}

	return current
}

// SetDev set dev.
func SetDev() {
	current = Dev
}

// SetTest set test.
func SetTest() {
	current = Test
}

// SetUat set uat.
func SetUat() {
	current = Uat
}

// SetProd set prod.
func SetProd() {
	current = Prod
}

// Is mode equal target.
func Is(target Env) bool { return current == target }

// Valid return true if one of dev, test, uat or prod.
func Valid() bool { return current > None && current <= Prod }

// IsDev is dev or not.
func IsDev() bool {
	return Get() == Dev
}

// IsTest is test or not.
func IsTest() bool {
	return Get() == Test
}

// IsUat is uat or not.
func IsUat() bool {
	return Get() == Uat
}

// IsProd is prod or not.
func IsProd() bool {
	return Get() == Prod
}

// IsTesting dev or test
func IsTesting() bool { return IsDev() || IsTest() }

// IsRelease uat or prod
func IsRelease() bool { return IsUat() || IsProd() }
