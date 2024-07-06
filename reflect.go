// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unwrap

import (
	"reflect"
	"unsafe"
)

// CopyReflectData unsafely copies a pointer so it can be modified. This copies
// the underlying data (by address) and just changes the flags on the
// `reflect.Value` wrapping that data.
//
// By default, `reflect` sets the read-only flag on values that are internal,
// e.g. those produced by `FieldByName()` on unexported fields. This function
// makes a new `reflect.Value` without the read-only flag set.
//
// If the value is not addressable, the value cannot and will not be copied
// and this function will return a `copied=false` value.
//
// H/T: https://stackoverflow.com/a/43918797/1068170
func CopyReflectData(v reflect.Value) (reflect.Value, bool) {
	if !v.CanAddr() {
		return v, false
	}

	v2 := reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	return v2, true
}

// CopyReflectValue copies a value so it can be modified.
//
// H/T: https://stackoverflow.com/a/43918797/1068170
func CopyReflectValue(v reflect.Value) reflect.Value {
	v2 := reflect.New(v.Type()).Elem()
	v2.Set(v)
	return v2
}
