libsql
======

# Project Overview

Package `libsql` provides wrappers for the Go standard library's "sql" package
that enables writing unit-testable and less error-prone SQL database code.

# Getting Started

The `libsql` package can be installed by running:

```
$ go get oss.indeed.com/go/libsql
```


The `libsql.Wrap` method can then be used to wrap a normal `*sql.DB`:

```go
import (
	"context"
	"database/sql"

	"oss.indeed.com/go/libsql"
)

func usageExample(sqldb *sql.DB) error {
	db := libsql.Wrap(sqldb)

	var column1, column2, column3 int64
	ctx := context.TODO()

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

`oss.indeed.com/go/libsql` is governed by the[Contributer Covenant v1.4.1](version/1/4/code-of-conduct.html)

For more information please contact opensource@indeed.com.

# License

The `oss.indeed.com/go/libsql` module is open source under the BSD-3-Clause
license.

