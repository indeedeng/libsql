libsql
======

[![Go Report Card](https://goreportcard.com/badge/oss.indeed.com/go/libsql)](https://goreportcard.com/report/oss.indeed.com/go/libsql)
[![Build Status](https://travis-ci.org/indeedeng/libsql.svg?branch=master)](https://travis-ci.org/indeedeng/libsql)
[![GoDoc](https://godoc.org/oss.indeed.com/go/libsql?status.svg)](https://godoc.org/oss.indeed.com/go/libsql)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/indeedeng/libsql.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/indeedeng/libsql.svg)](LICENSE)

# Project Overview

Package `libsql` provides wrappers for the Go standard library's "sql" package
that enables writing unit-testable and less error-prone SQL database code.

# Getting Started

The `libsql` package can be installed by running:

```
$ go get oss.indeed.com/go/libsql
```


To get started, use the `libsql.Wrap` method and pass a `*sql.DB`:
```go
import (
	"context"
	"database/sql"

	"oss.indeed.com/go/libsql"
)

func usageExample(ctx context.Context, sqldb *sql.DB) error {
	db := libsql.Wrap(sqldb)

	var column1, column2, column3 int64

	err := db.ScanOne(
		ctx,
		libsql.Into(
			&column1,
			&column2,
			&column3,
		),
		"SELECT 1, 2, 3 FROM dual WHERE 1 = ? AND 2 = ?",
		1,
		2,
	)
	if err != nil {
		return err
	}
	
	// use the values of column1, column2, column3
}
```

To scan multiple rows, provide a custom implementation of `libsql.RowScanner`: 
```go

import (
	"context"

	"oss.indeed.com/go/libsql"
)

type elephant struct {
	name string
}

type elephantRowScanner struct {
	rawElephantName string

	elephants []elephant
}

var _ libsql.RowScanner = (*elephantRowScanner)(nil)

func (e *elephantRowScanner) Into() []interface{} {
	return []interface{}{&e.rawElephantName}
}

func (e *elephantRowScanner) RowScanned() error {
	if e.rawElephantName == "" {
		return errors.New("empty elephant name in the database")
	}

	e.elephants = append(e.elephants, elephant{name: e.rawElephantName})

	return nil
}

func elephantsFrom(ctx context.Context, q libsql.Queryer) ([]elephant, error) {
	elephantRS := elephantRowScanner{}
	err := q.Scan(
		ctx,
		&elephantRS,
		`SELECT * FROM (VALUES ROW ('Dumbo'), ROW ('Horton')) AS elephants(name)`,
	)
	return elephantsRS.elephants, err
}
```

# Asking Questions

For technical questions about `libsql`, just file an issue in the GitHub tracker.

For questions about Open Source in Indeed Engineering, send us an email at
opensource@indeed.com

# Contributing

We welcome contributions! Feel free to help make `libsql` better.

### Process

- Open an issue and describe the desired feature / bug fix before making
changes. It's useful to get a second pair of eyes before investing development
effort.
- Make the change. If adding a new feature, remember to provide tests that
demonstrate the new feature works, including any error paths. If contributing
a bug fix, add tests that demonstrate the erroneous behavior is fixed.
- Open a pull request. Automated CI tests will run. If the tests fail, please
make changes to fix the behavior, and repeat until the tests pass.
- Once everything looks good, one of the indeedeng members will review the
PR and provide feedback.

# Maintainers

The `oss.indeed.com/go/libsql` module is maintained by Indeed Engineering.

While we are always busy helping people get jobs, we will try to respond to
GitHub issues, pull requests, and questions within a couple of business days.

# Code of Conduct

`oss.indeed.com/go/libsql` is governed by the[Contributer Covenant v1.4.1](CODE_OF_CONDUCT.md)

For more information please contact opensource@indeed.com.

# License

The `oss.indeed.com/go/libsql` module is open source under the [BSD-3-Clause](LICENSE)
license.

