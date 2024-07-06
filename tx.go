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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	pgx "github.com/jackc/pgx/v5"
)

// wrapTx is vendored in from the pgx source:
// https://github.com/jackc/pgx/blob/v5.6.0/stdlib/sql.go#L874-L877
type wrapTx struct {
	ctx context.Context
	tx  pgx.Tx
}

// ExtractTx unwraps a `pgx.Tx` from a standard library transaction
// created from the `github.com/jackc/pgx/v5/stdlib` driver.
func ExtractTx(tx *sql.Tx) (pgx.Tx, error) {
	// 1. reflect.Pointer (*sql.Tx)
	txValue := reflect.ValueOf(tx)

	// 2. reflect.Struct (sql.Tx)
	txValue = txValue.Elem()

	// 3. reflect.Interface (driver.Tx via sql.Tx{}.txi; read-only)
	txiValue := txValue.FieldByName("txi")

	// 4. reflect.Interface (driver.Tx via sql.Tx{}.txi; read-write)
	txiValue, copied := CopyReflectData(txiValue)
	if !copied {
		return nil, fmt.Errorf("cannot address txi; (%s).%s", txiValue.Type().PkgPath(), txiValue.Type().Name())
	}

	// 5. reflect.Struct (pgxstdlib.wrapTx; not addressable)
	if !txiValue.CanInterface() {
		return nil, fmt.Errorf("cannot resolve txi interface; (%s).%s", txiValue.Type().PkgPath(), txiValue.Type().Name())
	}
	wrapTxValue := reflect.ValueOf(txiValue.Interface())

	// 6. reflect.Struct (pgxstdlib.wrapTx; addressable)
	wrapTxValue = CopyReflectValue(wrapTxValue)

	// 7. Verify `wrapTx`
	wrapTxType := wrapTxValue.Type()
	if wrapTxType.PkgPath() != "github.com/jackc/pgx/v5/stdlib" || wrapTxType.Name() != "wrapTx" {
		return nil, fmt.Errorf("unexpected type; (%s).%s", wrapTxType.PkgPath(), wrapTxType.Name())
	}

	// 8. Unsafely convert memory to `wrapTx`
	if !wrapTxValue.CanAddr() {
		return nil, errors.New("cannot address wrapTx")
	}
	p := unsafe.Pointer(wrapTxValue.UnsafeAddr())
	wt := (*wrapTx)(p)

	return wt.tx, nil
}
