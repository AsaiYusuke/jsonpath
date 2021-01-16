# JSONPath

[![Build Status](https://travis-ci.com/AsaiYusuke/jsonpath.svg?branch=main)](https://travis-ci.com/AsaiYusuke/jsonpath)
[![Go Report Card](https://goreportcard.com/badge/github.com/AsaiYusuke/jsonpath)](https://goreportcard.com/report/github.com/AsaiYusuke/jsonpath)
[![Coverage Status](https://coveralls.io/repos/github/AsaiYusuke/jsonpath/badge.svg?branch=main)](https://coveralls.io/github/AsaiYusuke/jsonpath?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/AsaiYusuke/jsonpath.svg)](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is [Go](https://golang.org/) library providing the features that retrieves a part of the JSON object according to the query written in the JSONPath syntax.

The core JSONPath syntax on which this library based:

- [Stefan Gössner's JSONPath - XPath for JSON](https://goessner.net/articles/JsonPath/)
- [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison)
- [JSONPath Internet Draft Development](https://github.com/ietf-wg-jsonpath/draft-ietf-jsonpath-jsonpath)

#### Note:
Please check [my compare result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html) to know which responses are adapted.

## Getting started

```bash
go get github.com/AsaiYusuke/jsonpath
```

### Simple example:

```go
package main

import (
  "encoding/json"
  "fmt"

  "github.com/AsaiYusuke/jsonpath"
)

func main() {
  jsonPath, srcJSON := `$.key`, `{"key":"value"}`
  var src interface{}
  json.Unmarshal([]byte(srcJSON), &src)
  output, _ := jsonpath.Retrieve(jsonPath, src)
  outputJSON, _ := json.Marshal(output)
  fmt.Println(string(outputJSON))
  // Output:
  // ["value"]
}
```

## Basic design

- [PEG](https://github.com/pointlander/peg) separated the JSONPath syntax analyzer from functionality itself to simplify the source.
- The error specification enables the package user to handle errors correctly.
- Adopted more of the consensus behavior from the [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison).
  Adapted my own behavior to the other part of the such consensus behavior that found difficult to use.
- Equipped with a large number of unit tests to avoid bugs that lead to unexpected results.

## How to use

### * Retrieve one-time or repeated

The `Retrieve` function returns retrieved result using JSONPath and JSON object:

```go
output, err := jsonpath.Retrieve(jsonPath, src)
```

The `Parse` function returns a *parser-function* that completed to check JSONPath syntax.
By using *parser-function*, it can repeat to retrieve with the same JSONPath :

```go
jsonPath, err := jsonpath.Parse(jsonPath)
output1, err1 := jsonPath(src1)
output2, err2 := jsonPath(src2)
:
```

### * Error handling

If there is a problem with the execution of *APIs*, an error type returned.
These error types define the corresponding symptom, as listed below:

#### Syntax check errors from `Retrieve`, `Parse`

| Error type              | Message format                                     | Symptom                                                                                                       |
|-------------------------|----------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `ErrorInvalidSyntax`    | `invalid syntax (position=%d, reason=%s, near=%s)` | The invalid syntax found in the JSONPath.<br>The *reason* including in this message will tell you more about it. |
| `ErrorInvalidArgument`  | `invalid argument (argument=%s, error=%s)`         | The argument specified in the JSONPath treated as the invalid error in Go syntax.                         |
| `ErrorFunctionNotFound` | `function not found (function=%s)`                 | The function specified in the JSONPath is not found.                                                          |
| `ErrorNotSupported`     | `not supported (feature=%s, path=%s)`              | The unsupported syntaxes specified in the JSONPath.                                                           |

#### Runtime errors from `Retrieve`, *`parser-functions`*

| Error type             | Message format                                    | Symptom                                                                             |
|------------------------|---------------------------------------------------|-------------------------------------------------------------------------------------|
| `ErrorMemberNotExist`  | `member did not exist (path=%s)`                  | The object/array member specified in the JSONPath did not exist in the JSON object. |
| `ErrorIndexOutOfRange` | `index out of range (path=%s)`                    | The array indexes specified in the JSONPath were out of range.                      |
| `ErrorTypeUnmatched`   | `type unmatched (expected=%s, found=%s, path=%s)` | The node type specified in the JSONPath did not exist in the JSON object.           |
| `ErrorNoneMatched`     | `none matched (path=%s)`                          | The retrieving child paths specified in the JSONPath resulted in empty output.      |
| `ErrorFunctionFailed`  | `function failed (function=%s, error=%s)`         | The function specified in the JSONPath failed.                                      |

The type checking is convenient to recognize which error happened.

```go
  :
  _,err := jsonpath.Retrieve(jsonPath, srcJSON)
  if err != nil {
    switch err.(type) {
    case jsonpath.ErrorIndexOutOfRange:
      fmt.printf(`retry with other srcJSON: %v`, err)
      continue
    case jsonpath.ErrorInvalidArgumentFormat:
      return nil, fmt.errorf(`specified invalid argument: %v`, err)
    }
    :
  }
```

### * Function syntax

Function enables to format results by using user defined functions.
The function syntax comes after the JSONPath.

There are two ways to use function:

#### Filter function

The filter function applies a user function to each values in the result to get converted.

[Check the example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetFilterFunction)


#### Aggregate function

The aggregate function converts all values in the result into a single value.

[Check the example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetAggregateFunction)

### * Accessing JSON

You can get the accessors ( *Getters / Setters* ) of the input JSON instead of the retrieved values.
These accessors can use to update for the input JSON.

This feature can get enabled by giving `Config.SetAccessorMode()`.

[Check the example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetAccessorMode)

#### Note:
It is not possible to use *Setter* for some results, such as for JSONPath including function syntax.

Also, operations using accessors follow the map/slice manner of Go language.
If you use accessors after changing the structure of JSON, you need to pay attention to the behavior.
If you don't want to worry about it, get the accessor again every time you change the structure.

## Differences

Some behaviors that differ from the consensus exists in this package.
For the entire comparisons, please check [this result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html).

These behaviors will change in the future if appropriate ones found.

### Character types

The following character types can be available for identifiers in dot-child notation.

| Character type                                                                                                                           | Availabe | Escape |
|------------------------------------------------------------------------------------------------------------------------------------------|----------|--------|
| * Numbers and alphabets (`0-9` `A-Z` `a-z`)<br> * Hyphen and underscore (`-` `_`)<br> * Non-ASCII Unicode characters (`0x80 - 0x10FFFF`) | Yes      | No     |
| * Other printable symbols (`Space` `!` `"` `#` `$` `%` `&` `'` `(` `)` `*` `+` `,` `.` `/` `:` `;` `<` `=` `>` `?` `@` `[` `\` `]` `^` `` ` `` `{` <code>&#124;</code> `}` `~`) | ^ | Yes |
| * ~~Control code characters~~ (`0x00 - 0x1F`, `0x7F`) | No       | -      |

The printable symbols except hyphen and underscore can use by escaping them.

```text
JSONPath : $.abc\.def
srcJSON  : {"abc.def":1}
Output   : 1
```

### Wildcard in qualifier

The wildcards in qualifier can specify as a union of subscripts.

```text
JSONPath : $[0,1:3,*]
srcJSON  : [0,1,2,3,4,5]
Output   : [0,1,2,0,1,2,3,4,5]
```

### Regular expression

The regular expression syntax works as a regular expression in Go lang.
In particular, you can use "(?i)" to specify the regular expression as the ignore case option.

```text
JSONPath : $[?(@.a=~/(?i)CASE/)]
srcJSON  : ["Case","Hello"]
Output   : ["Case"]
```

### JSONPaths in the filter-qualifier

JSONPaths that returns value group cannot specify with `comparator` or `regular expression`.
But, `existence check` can use these.


| JSONPaths that return a value group | example       |
|-------------------------------------|---------------|
| Recursive descent                   | `@..a`        |
| Multiple identifier                 | `@['a','b']`  |
| Wildcard identifier                 | `@.*`         |
| Slice qualifier                     | `@[0:1]`      |
| Wildcard qualifier                  | `@[*]`        |
| Union in the qualifier              | `@[0,1]`      |
| Filter qualifier                    | `@.a[?(@.b)]` |

- comparator example (error)

```text
JSONPath : $[?(@..x == "hello world")]
srcJSON  : [{"a":1},{"b":{"x":"hello world"}}]
Error    : ErrorInvalidSyntax
```

- regular expression example (error)

```text
JSONPath : $[?(@..x=~/hello/)]
srcJSON  : [{"a":1},{"b":{"x":"hello world"}}]
Error    : ErrorInvalidSyntax
```

- existence check example

```text
JSONPath : $[?(@..x)]
srcJSON  : [{"a":1},{"b":{"x":"hello world"}}]
Output   : [{"b":{"x":"hello world"}}]
```

## Benchmarks

I benchmarked two JSONPaths using several libraries for the Go language.
What is being measured is the cost per job for a job that loops a lot after all the prep work done.

There was a performance differences.
But if the number of queries is little, there will not be a big difference between any of them.

Also, the results will vary depending on the data entered.
So this benchmark is for information only and should be re-measured at every time.

- [BenchmarkAsaiYusukeJSONPath](https://github.com/AsaiYusuke/jsonpath)
- [BenchmarkOhler55Ojg](https://github.com/ohler55/ojg/jp)
- [BenchmarkBhmjJSONSlice](https://github.com/bhmj/jsonslice)
- [BenchmarkPaesslerAGJSONPath](https://github.com/PaesslerAG/jsonpath)
- [BenchmarkOliveagleJsonpath](https://github.com/oliveagle/jsonpath)

### JSONPath for comparison with more libraries

This is the result of a JSONPath that all libraries were able to process.
Oliveagle/jsonpath is fastest.

```text
JSONPath : $.store.book[0].price

BenchmarkAsaiYusukeJSONPath_threeLevelsWithIndex-4         	 5674974	       207 ns/op	      56 B/op	       3 allocs/op
BenchmarkOhler55Ojg_threeLevelsWithIndex-4                 	 1728166	       683 ns/op	    1040 B/op	       2 allocs/op
BenchmarkBhmjJSONSlice_threeLevelsWithIndex-4              	  546034	      2213 ns/op	      32 B/op	       1 allocs/op
BenchmarkPaesslerAGJSONPath_threeLevelsWithIndex-4         	 1875916	       648 ns/op	     208 B/op	       7 allocs/op
BenchmarkOliveagleJsonpath_threeLevelsWithIndex-4          	13331599	        90.3 ns/op	       0 B/op	       0 allocs/op
```

### A slightly complex JSONPath

Libraries that can handle complex syntax limited to a few.
Among these libraries, my library is the fastest at the moment.

```text
JSONPath : $..book[?(@.price > $.store.bicycle.price)]

BenchmarkAsaiYusukeJSONPath_recursiveDescentWithFilter-4   	  387717	      3325 ns/op	     656 B/op	      26 allocs/op
BenchmarkOhler55Ojg_recursiveDescentWithFilter-4           	  231213	      5237 ns/op	    5240 B/op	      20 allocs/op
BenchmarkBhmjJSONSlice_recursiveDescentWithFilter-4        	   52902	     22661 ns/op	    3032 B/op	      57 allocs/op
BenchmarkPaesslerAGJSONPath                                	  not supported
BenchmarkOliveagleJsonpath                                 	  not supported
```

##### JSON used for the benchmark measurement:

<details>
<summary>Show results</summary>

```text
{ "store": {
  "book": [
    { "category": "reference",
    "author": "Nigel Rees",
    "title": "Sayings of the Century",
    "price": 8.95
    },
    { "category": "fiction",
    "author": "Evelyn Waugh",
    "title": "Sword of Honour",
    "price": 12.99
    },
    { "category": "fiction",
    "author": "Herman Melville",
    "title": "Moby Dick",
    "isbn": "0-553-21311-3",
    "price": 8.99
    },
    { "category": "fiction",
    "author": "J. R. R. Tolkien",
    "title": "The Lord of the Rings",
    "isbn": "0-395-19395-8",
    "price": 22.99
    }
  ],
  "bicycle": {
    "color": "red",
    "price": 19.95
  }
  }
}
```

</details>

##### Benchmark environment:

<details>
<summary>Show results</summary>

```text
Processor  : Intel Core i5-6267U 2.90GHz
Memory     : 16.0 GB
OS         : Windows 10
Go version : go1.15.6 windows/amd64
```
</details>

## Project progress

- Syntax
  - Identifier
    - [x] identifier in dot notations
    - [x] identifier in bracket notations
    - [x] wildcard
    - [x] multiple-identifier in bracket
    - [x] recursive retrieve
  - Qualifier
    - [x] index
    - [x] slice
    - [x] wildcard
    - Filter
      - [x] logical operation
      - [x] comparator
      - [x] JSONPath retrieve in filter
    - [ ] script
  - Function
    - [x] filter
    - [x] aggregate
  - [x] Refer to the consensus behaviors
- Archtecture
  - [x] PEG syntax analyzing
  - [x] Error handling
  - [x] Function
  - [x] Accessing JSON
- Go language manner
  - [x] retrieve with the object in interface unmarshal
  - [x] retrieve with the json.Number type
- Source code
  - [x] Release version
  - Unit tests
    - [x] syntax tests
    - [x] benchmark
    - [x] coverage >80%
  - [x] Examples
  - [x] CI automation
  - Documentation
    - [x] README
    - [ ] API doc
  - [x] comparison result (local)
- Development status
  - [x] determine requirements / functional design
  - [x] design-based coding
  - [ ] testing
  - [ ] documentation
- Future ToDo
  - [ ] Refer to the something standard
  - Go language affinity
    - [ ] retrieve with the object in struct unmarshal
    - [ ] retrieve with the struct tags
    - [ ] retrieve with the user defined objects
