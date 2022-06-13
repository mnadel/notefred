package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkQueryTitlesAll(t *testing.B) {
	bearDB, err := NewDB()
	assert.Nil(t, err, "cannot create db")

	for i := 0; i < t.N; i++ {
		_, err = bearDB.QueryTitles("2022", "")
		assert.Nil(t, err, "error searching titles all")
	}
}

func BenchmarkQueryTitlesSome(t *testing.B) {
	bearDB, err := NewDB()
	assert.Nil(t, err, "cannot create db")

	for i := 0; i < t.N; i++ {
		_, err = bearDB.QueryTitles("2022", "qqq")
		assert.Nil(t, err, "error searching titles some")
	}
}
