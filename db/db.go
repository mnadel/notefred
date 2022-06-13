package db

import (
	"os"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mnadel/notefred/util"
	"github.com/pkg/errors"

	"database/sql"
)

var (
	systemIgnoreFolders = strings.Split(systemIgnoreFolderCSV, ",")
)

const (
	systemIgnoreFolderCSV = "Recently Deleted"

	dbFile = `/Library/Group Containers/group.com.apple.notes/NoteStore.sqlite?mode=ro`

	sqlTitle = `
		SELECT DISTINCT
			c1.ZIDENTIFIER as uuid,
			c1.ZTITLE1 as title,
    		c2.ZTITLE2 as folder
  		FROM
    		ZICNOTEDATA as n
    		LEFT JOIN ZICCLOUDSYNCINGOBJECT as c1 ON c1.ZNOTEDATA = n.Z_PK
    		LEFT JOIN ZICCLOUDSYNCINGOBJECT as c2 ON c2.Z_PK = c1.ZFOLDER
  		WHERE
		  	c1.ZTITLE1 LIKE ?
    		AND c2.ZTITLE2 NOT IN ({ignore})
	`

	sqlPragma = `
		PRAGMA query_only = on;
		PRAGMA synchronous = normal;
		PRAGMA temp_store = memory;
		PRAGMA mmap_size = 30000000000;
		PRAGMA cache_size = -64000;
	`
)

// DB represents the Bear Notes database
type DB struct {
	db *sql.DB
}

// Result references a specific note: its identifier and title
type Result struct {
	UUID   string
	Title  string
	Folder string
}

// Results is a list of *Result, and represents a collection of notes in the database
type Results []*Result

// Create a new DB, referencing the user's Bear Notes database
func NewDB() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db, err := sql.Open("sqlite3", path.Join(home, dbFile))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if _, err := db.Exec(sqlPragma); err != nil {
		return nil, errors.WithStack(err)
	}

	return &DB{db}, nil
}

// Close cleans up our database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// QueryTitles searches for a term within the titles of notes within the database, set
// ignores to a CSV of Folder names to ignore notes in those folders.
func (d *DB) QueryTitles(term string, ignores string) (Results, error) {
	searchBind := substringSearch(term)

	var ignoreFolders []string

	if ignores == "" {
		ignoreFolders = systemIgnoreFolders
	} else {
		ignoreFolders = strings.Split(ignores, ",")
		ignoreFolders = append(ignoreFolders, systemIgnoreFolders...)
	}

	bindString := "?" + strings.Repeat(",?", len(ignoreFolders)-1)
	query := strings.Replace(sqlTitle, "{ignore}", bindString, 1)

	args := make([]interface{}, len(ignoreFolders)+1)
	args[0] = searchBind
	for i, f := range ignoreFolders {
		args[i+1] = f
	}

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	return rowsToResults(rows)
}

// TitleCase returns a Alfred-safe version of the proper title casing
func (r *Result) TitleCase() string {
	return util.ToSafeString(util.ToTitleCase(r.Title))
}

func rowsToResults(rows *sql.Rows) (Results, error) {
	var uuid string
	var title string
	var folder string

	results := make(Results, 0)

	for rows.Next() {
		err := rows.Scan(&uuid, &title, &folder)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		results = append(results, &Result{UUID: uuid, Title: title, Folder: folder})
	}

	return results, errors.WithStack(rows.Err())
}

func substringSearch(term string) string {
	bind := strings.Builder{}
	bind.WriteString(`%`)
	bind.WriteString(term)
	bind.WriteString(`%`)
	return bind.String()
}
