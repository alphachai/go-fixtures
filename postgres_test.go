package fixtures

import (
	"testing"

	"github.com/charlieparkes/go-fixtures/test"
	"github.com/stretchr/testify/assert"
)

func TestPostgres(t *testing.T) {
	timer := test.NewTimer()
	fixtures := Fixtures{}
	defer fixtures.TearDown()

	d := &Docker{NamePrefix: "gofixtures"}
	fixtures.AddByName("docker", d)
	timer.PrintSplit("Docker")

	p1 := &PostgresWithSchema{
		Docker:   d,
		PathGlob: "./testdata/migrations/*.sql",
	}
	fixtures.AddByName("postgres_with_schema", p1)
	timer.PrintSplit("PostgresWithSchema")

	db, err := p1.GetConnection()
	assert.NoError(t, err)

	tables := []string{}
	assert.NoError(t, db.Select(&tables, "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'information_schema' AND schemaname != 'pg_catalog';"))
	assert.Len(t, tables, 2) // migrations define two tables
	assert.NoError(t, db.Close())
	timer.PrintSplit("Migrations check")

	p1.Postgres.DumpDatabase("testdata/tmp", "test.pgdump")
	timer.PrintSplit("Postgres.DumpDatabase")

	c := &PostgresDatabaseCopy{
		Postgres:     p1.Postgres,
		SkipTearDown: true,
	}
	fixtures.AddByName("db_copy", c)
	timer.PrintSplit("PostgresDatabaseCopy")

	db, err = c.GetConnection()
	assert.NoError(t, err)

	tables = []string{}
	assert.NoError(t, db.Select(&tables, "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'information_schema' AND schemaname != 'pg_catalog';"))
	assert.Len(t, tables, 2)
	timer.PrintSplit("Migrations check")

	p2 := &Postgres{
		Docker: d,
	}
	fixtures.AddByName("postgres", p2)
	timer.PrintSplit("Postgres")

	err = p2.RestoreDatabase("testdata/tmp", "test.pgdump")
	assert.NoError(t, err)
	timer.PrintSplit("Postgres.RestoreDatabase")

	db, err = p2.GetConnection("")
	assert.NoError(t, err)

	tables = []string{}
	assert.NoError(t, db.Select(&tables, "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'information_schema' AND schemaname != 'pg_catalog';"))
	assert.Len(t, tables, 2)
	assert.NoError(t, db.Close())
	timer.PrintSplit("Migrations check")
}
