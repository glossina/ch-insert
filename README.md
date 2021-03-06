# ch-insert
[![Build Status](https://travis-ci.org/sirkon/ch-insert.svg?branch=master)](https://travis-ci.org/sirkon/ch-insert)

Clickhouse HTTP interface data inserter.

Clickhouse HTTP RowBinary insertion objects. They are meant to be used with [ch-encode](https://github.com/sirkon/ch-encode)-produced RowBinary data encoder.

Usage example:
### First create table test and generate encoder using [ch-encode](https://github.com/sirkon/ch-encode)
```bash
clickhouse-client --query "CREATE TABLE test (
    date Date,
    uid String,
    hidden UInt8
) ENGINE = MergeTree(date, (date, uid, hidden), 8192);" # Create table test

go get -u github.com/sirkon/ch-encode
go get -u github.com/sirkon/ch-insert
echo 'uid: UID' > dict.yaml   # We want uid to be represented as UID in Go code

bin/ch-encode --yaml-dict dict.yaml --date-field date test  # Generate encoder package in current directory
mv test src/                                                # and move it to src/ in order for go <cmd> to be able to use it
```

### Usage
```go
package main

import (
	"test"
	"time"

	chinsert "github.com/sirkon/ch-insert"
)

func main() {
	ins, err := chinsert.Open("localhost:8123/default", "test", 10*1024*1024, 1024*1024*1024)
	if err != nil {
		panic(err)
	}
	defer inserter.Close()
	inserter := ins.WithThreadSafe()
	encoder := test.NewTestRawEncoder(inserter)
	if err := encoder.Encode(test.Date.FromTime(time.Now()), test.UID("123"), test.Hidden(1)); err != nil {
		panic(err)
	}
	if err := encoder.Encode(test.Date.FromTime(time.Now()), test.UID("123"), test.Hidden(0)); err != nil {
		panic(err)
	}
}
```

Run it:
```bash
go run main.go
```

And see data in clickhouse test table:
```
SELECT *
FROM test

┌──────date─┬─uid─┬──hidden─┐
│ 2017-07-15 │ 123 │       0 │
│ 2017-07-15 │ 123 │       1 │
└───────────┴─────┴────────┘

2 rows in set. Elapsed: 0.004 sec.
```

### Lower level usage
```go
// file main.go
package main

import (
	"net/http"
	"test"
	"time"

	chinsert "github.com/sirkon/ch-insert"
)

func main() {
	rawInserter := chinsert.New(
		&http.Client{},
		chinsert.ConnParams{
			Host: "localhost",
			Port: 8123,
		},
		"test",  // Table name to insert data in
	)

	inserter := chinsert.NewBuf(rawInserter, 10*1024*1024) // 10Mb buffer
	defer inserter.Close()
	encoder := test.NewTestRawEncoder(inserter)
	if err := encoder.Encode(test.Date.FromTime(time.Now()), test.UID("123"), test.Hidden(1)); err != nil {
		panic(err)
	}
	if err := encoder.Encode(test.Date.FromTime(time.Now()), test.UID("123"), test.Hidden(0)); err != nil {
		panic(err)
	}
}
```


### Lower level with smart inserter
It is [not recommended](https://clickhouse.yandex/docs/en/introduction/performance.html#performance-on-data-insertion)
to insert more than 1 times per second (per writer I guess). There's an object that tries not to insert more than 1 time
per second, usage:
```go
package main

import (
	"net/http"
	"test"

	chinsert "github.com/sirkon/ch-insert"
)

func main() {
	rawInserter := chinsert.New(
		&http.Client{},
		chinsert.ConnParams{
			Host: "localhost",
			Port: 8123,
		},
		"test", // Table name to insert data in
	)

	epoch := chinsert.NewEpochDirect()
	inserter := chinsert.NewBuf(rawInserter, 1024*1024*1024) // 1Gb buffer is hard limit for insertion
	defer inserter.Close()

	si := chinsert.NewSmartInsert(inserter, 10*1024*1024, epoch)
	encoder := test.NewTestRawEncoder(si)
	for i := 0; i < 100000000; i++ {
		if err := encoder.Encode(test.Date.FromTimestamp(epoch.Seconds()), test.UID("123"), test.Hidden(0)); err != nil {
			panic(err)
		}
	}
}
```
This thing will insert 100 million records. Don't worry, these are similar rows and columns are not wide, thus it should
not take more than 10 seconds to complete.  
