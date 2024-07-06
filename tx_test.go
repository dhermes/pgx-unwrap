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
	"testing"

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
	assert.Nil(err)
	assert.NotNil(unwrapped)

	row := unwrapped.QueryRow(ctx, "SELECT current_user")
	searchPath := ""
	err = row.Scan(&searchPath)
	assert.Nil(err)
	assert.Equal("...", searchPath)
}
