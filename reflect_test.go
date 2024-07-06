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

package unwrap_test

import (
	"reflect"
	"testing"

	testifyrequire "github.com/stretchr/testify/require"

	unwrap "github.com/dhermes/pgx-unwrap"
)

func TestCopyReflectData(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	i := 1337

	v := reflect.ValueOf(i)
	v2, copied := unwrap.CopyReflectData(v)
	assert.False(copied)
	assert.Equal(v, v2)

	v = reflect.ValueOf(&i).Elem()
	v2, copied = unwrap.CopyReflectData(v)
	assert.True(copied)
	assert.Equal(v, v2)
}

func TestCopyReflectValue(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	i := 1337
	v := reflect.ValueOf(i)
	v2 := unwrap.CopyReflectValue(v)
	assert.NotEqual(v, v2)
	assert.False(v.CanAddr())
	assert.True(v2.CanAddr())
}
