package dal

import (
	"database/sql"
	"testing"

	"github.com/magicalbanana/dal/mocks"
	"github.com/magicalbanana/vfs"

	"github.com/stretchr/testify/assert"
)

func TestQueryRow(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc      string
		assertion func(*testing.T, string)
	}{
		{
			desc: "no file store",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				d := &dal{db: &mocks.Db{}, fs: nil}

				// act
				_, err := d.QueryRow("", nil)

				// assertion
				assert.Error(t, err)
			},
		},
		{
			desc: "SQL File does not exist",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				fs, fsErr := vfs.LoadFiles("tests/fixtures/sqls")
				assert.NoError(t, fsErr)
				d := &dal{db: &mocks.Db{}, fs: fs}

				// act
				_, qryErr := d.QueryRow("manbearpig.sql", nil)

				// assertion
				assert.Error(t, qryErr)
			},
		},
		{
			// db.Prepare errors out because the SQL file is empty
			desc: "db.Prepare returned an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				fs, fsErr := vfs.LoadFiles("tests/fixtures/sqls")
				assert.NoError(t, fsErr)
				d := &dal{db: &mocks.Db{PrepareOk: false}, fs: fs}

				params := make(map[string]interface{})
				// act
				_, err := d.QueryRow("test.sql", params)

				// assertion
				assert.Error(t, err)
			},
		},
		{
			desc: "db.Prepare did not return an error",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				fs, fsErr := vfs.LoadFiles("tests/fixtures/sqls")
				assert.NoError(t, fsErr)
				db, openErr := sql.Open("postgres", "dbname=dbdal_test sslmode=disable")
				assert.NoError(t, openErr)
				params := map[string]interface{}{
					"first_name": "bearpig",
					"last_name":  "man",
					"address":    []byte(`{"test": "foo"}`),
				}
				d := &dal{db: db, fs: fs}

				// act
				_, qryErr := d.QueryRow("insert_customer.sql", params)

				// assertion
				assert.NoError(t, qryErr)

				// clean up
				_, qryErr = d.Query("delete_customer.sql", params)
				assert.NoError(t, qryErr)
			},
		},
		{
			desc: "no params passed",
			assertion: func(t *testing.T, desc string) {
				// arrangement
				fs, fsErr := vfs.LoadFiles("tests/fixtures/sqls")
				assert.NoError(t, fsErr)
				db, openErr := sql.Open("postgres", "dbname=dbdal_test sslmode=disable")
				assert.NoError(t, openErr)
				d := &dal{db: db, fs: fs}

				// act
				_, qryErr := d.QueryRow("select_all_customer.sql", nil)

				// assertion
				assert.NoError(t, qryErr)
			},
		},
	}

	for _, test := range tests {
		test.assertion(t, test.desc)
	}
}
