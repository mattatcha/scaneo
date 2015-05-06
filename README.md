# scaneo

[![Build Status](https://drone.io/github.com/variadico/scaneo/status.png)](https://drone.io/github.com/variadico/scaneo/latest)

Generate code to convert `*sql.Row` and `*sql.Rows` into arbitrary structs.

## Installation

```
go get github.com/variadico/scaneo
```

## Usage

```
usage: scaneo [options] paths...

    -c    Clobber (overwrite) file if exists. Default is append to file.

    -o    Name of output file. Default is scans.go.

    -p    Package name. Default is current directory name.

    -u    Unexport functions. Default is export all.

    -h    Display usage information and exit.
```

## Examples

Let's say you have a file called `tables.go` that looks like this.

```
package models

import "time"

type Post struct {
	ID          int64
	UUID        []byte
	Title       []byte
	Slug        string
	Markdown    []byte
	HTML        []byte
	IsFeatured  bool
	IsPage      bool
	IsPublished bool
	Date        *time.Time
	Tags        []byte
	Author      string
	Image       []byte
}
```

Run `scaneo tables.go` and this will generate a new file called `scans.go`.
`scans.go` will look like this.

```
package models

import "database/sql"

func ScanPost(r *sql.Row) (Post, error) {
	var s Post

	if err := r.Scan(
		&s.ID,
		&s.UUID,
		&s.Title,
		&s.Slug,
		&s.Markdown,
		&s.HTML,
		&s.IsFeatured,
		&s.IsPage,
		&s.IsPublished,
		&s.Date,
		&s.Tags,
		&s.Author,
		&s.Image,
	); err != nil {
		return Post{}, err
	}

	return s, nil
}

func ScanPosts(rs *sql.Rows) ([]Post, error) {
	structs := make([]Post, 0, 16)

	var err error
	for rs.Next() {
		var s Post

		if err = rs.Scan(
			&s.ID,
			&s.UUID,
			&s.Title,
			&s.Slug,
			&s.Markdown,
			&s.HTML,
			&s.IsFeatured,
			&s.IsPage,
			&s.IsPublished,
			&s.Date,
			&s.Tags,
			&s.Author,
			&s.Image,
		); err != nil {
			return nil, err
		}

		structs = append(structs, s)
	}

	return structs, nil
}
```

Then you can call those functions from other code, like this.

```
func serveHome(resp http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("select * from post")
	if err != nil {
		log.Println(err)
		return
	}

	posts, err := models.ScanPosts(rows)
	if err != nil {
		log.Println(err)
	}

	// ... send posts to template or whatever...
}
```

### Go Generate

You can integrate `scaneo` with `go generate` by adding the generate comment to
the beginning of `tables.go`. `$GOFILE` is the name of the current
file, in this case, `tables.go`.

```
//go:generate scaneo -c $GOFILE

package models
... rest of code...
```

Then just call `go generate` from within the package and `scans.go` will be
created.
