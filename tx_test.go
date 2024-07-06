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
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	pgx "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	testifyrequire "github.com/stretchr/testify/require"

	unwrap "github.com/dhermes/pgx-unwrap"
)

func TestExtractTx(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	pool, err := sql.Open("pgx", "postgres://...")
	assert.Nil(err)
	t.Cleanup(func() {
		err := pool.Close()
		assert.Nil(err)
	})

	ctx := context.Background()
	tx, err := pool.BeginTx(ctx, nil)
	assert.Nil(err)
	t.Cleanup(func() {
		err := tx.Rollback()
		assert.Nil(err)
	})

	unwrapped, err := unwrap.ExtractTx(tx)
	assert.Nil(unwrapped)
	assert.Equal("not implemented", fmt.Sprintf("%v", err))

	// 1.
	txValue := reflect.ValueOf(tx)
	assert.Equal(reflect.Pointer, txValue.Type().Kind())

	// 2.
	txValue = txValue.Elem()
	assert.Equal(reflect.Struct, txValue.Type().Kind())

	// 3.
	txiValue := txValue.FieldByName("txi")
	assert.Equal(reflect.Interface, txiValue.Type().Kind())
	assert.True(txiValue.CanAddr())
	assert.False(txiValue.CanSet())
	assert.False(txiValue.CanInterface())

	// 4.
	txiValue, copied := unwrap.CopyReflectData(txiValue)
	assert.True(copied)
	assert.Equal(reflect.Interface, txiValue.Type().Kind())
	assert.True(txiValue.CanAddr())
	assert.True(txiValue.CanSet())
	assert.True(txiValue.CanInterface())

	// 5.
	wrapTxValue := reflect.ValueOf(txiValue.Interface())
	assert.Equal(reflect.Struct, wrapTxValue.Type().Kind())
	assert.False(wrapTxValue.CanAddr())
	assert.False(wrapTxValue.CanSet())
	assert.True(wrapTxValue.CanInterface())

	// 6.
	wrapTxValue = unwrap.CopyReflectValue(wrapTxValue)
	assert.Equal(reflect.Struct, wrapTxValue.Type().Kind())
	assert.True(wrapTxValue.CanAddr())
	assert.True(wrapTxValue.CanSet())
	assert.True(wrapTxValue.CanInterface())

	// 7.
	wrapTxType := wrapTxValue.Type()
	assert.Equal("github.com/jackc/pgx/v5/stdlib", wrapTxType.PkgPath())
	assert.Equal("wrapTx", wrapTxType.Name())

	// 8.
	p := unsafe.Pointer(wrapTxValue.UnsafeAddr())
	wt := (*wrapTx)(p)
	assert.NotNil(wt)
	assert.NotNil(wt.tx)
}

// wrapTx is vendored in from the pgx source:
// https://github.com/jackc/pgx/blob/v5.6.0/stdlib/sql.go#L874-L877
type wrapTx struct {
	ctx context.Context
	tx  pgx.Tx
}
