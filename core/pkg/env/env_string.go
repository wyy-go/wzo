// Code generated by "stringer -type=Env -linecomment"; DO NOT EDIT.

package env

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[None-0]
	_ = x[Dev-1]
	_ = x[Test-2]
	_ = x[Uat-3]
	_ = x[Prod-4]
}

const _Env_name = "nonedevtestuatprod"

var _Env_index = [...]uint8{0, 4, 7, 11, 14, 18}

func (i Env) String() string {
	if i < 0 || i >= Env(len(_Env_index)-1) {
		return "Env(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Env_name[_Env_index[i]:_Env_index[i+1]]
}
