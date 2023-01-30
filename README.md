# AsaiYusuke/JSONPath

[![Build Status](https://app.travis-ci.com/AsaiYusuke/jsonpath.svg?branch=main)](https://app.travis-ci.com/AsaiYusuke)
[![Go Report Card](https://goreportcard.com/badge/github.com/AsaiYusuke/jsonpath)](https://goreportcard.com/report/github.com/AsaiYusuke/jsonpath)
[![Coverage Status](https://coveralls.io/repos/github/AsaiYusuke/jsonpath/badge.svg?branch=main)](https://coveralls.io/github/AsaiYusuke/jsonpath?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/AsaiYusuke/jsonpath.svg)](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath)
[![Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#query-language)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![AsaiYusuke/JSONPath](assets/logo.svg)

This [Go](https://golang.org/) library is for retrieving a part of JSON according to the _JSONPath_ query syntax.

The core JSONPath syntax on which this library based:

- [Stefan Gössner's JSONPath - XPath for JSON](https://goessner.net/articles/JsonPath/)
- [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison)
- [JSONPath Internet Draft Development](https://github.com/ietf-wg-jsonpath/draft-ietf-jsonpath-jsonpath)

#### Note:

For syntax compatibility among other libraries, please check [:memo: my comparison results](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html).

## Table of Contents

- [Getting started](#getting-started)
- [Basic design](#basic-design)
- [How to use](#how-to-use)
  - [Retrieve one-time or repeated](#-retrieve-one-time-or-repeated)
  - [Error handling](#-error-handling)
  - [Function syntax](#-function-syntax)
  - [Accessing JSON](#-accessing-json)
- [Differences](#differences)
- [Benchmarks](#benchmarks)
- [Project progress](#project-progress)

## Getting started

### Install:

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

#### _Streamlined Development_

- The JSONPath syntax analysis functionality has been separated using [PEG](https://github.com/pointlander/peg), resulting in a more simplified source code.
- Robust unit testing has been implemented to prevent bugs and ensure consistent outcomes.

#### _User-Friendly Interface_

- The library is equipped with a comprehensive error specification, allowing users to effectively handle any errors that may arise.

#### _Unwavering Compatibility_

- The library has integrated a greater level of consensus behavior from [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison), ensuring seamless compatibility with related technologies.

## How to use

### \* Retrieve one-time or repeated

The `Retrieve` function returns retrieved result using JSONPath and JSON object:

```go
output, err := jsonpath.Retrieve(jsonPath, src)
```

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Retrieve)

The `Parse` function returns a _parser-function_ that completed to check JSONPath syntax.
By using _parser-function_, it can repeat to retrieve with the same JSONPath :

```go
jsonPath, err := jsonpath.Parse(jsonPath)
output1, err1 := jsonPath(src1)
output2, err2 := jsonPath(src2)
:
```

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Parse)

### \* Error handling

If there is a problem with the execution of _APIs_, an error type returned.
These error types define the corresponding symptom, as listed below:

#### Syntax check errors from `Retrieve`, `Parse`

| Error type              | Message format                                     | Symptom                                                                                                          | Ex                                                                                        |
| ----------------------- | -------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `ErrorInvalidSyntax`    | `invalid syntax (position=%d, reason=%s, near=%s)` | The invalid syntax found in the JSONPath.<br>The _reason_ including in this message will tell you more about it. | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorInvalidSyntax)    |
| `ErrorInvalidArgument`  | `invalid argument (argument=%s, error=%s)`         | The argument specified in the JSONPath treated as the invalid error in Go syntax.                                | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorInvalidArgument)  |
| `ErrorFunctionNotFound` | `function not found (function=%s)`                 | The function specified in the JSONPath is not found.                                                             | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorFunctionNotFound) |
| `ErrorNotSupported`     | `not supported (feature=%s, path=%s)`              | The unsupported syntaxes specified in the JSONPath.                                                              | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorNotSupported)     |

#### Runtime errors from `Retrieve`, _`parser-functions`_

| Error type            | Message format                                    | Symptom                                                                             | Ex                                                                                      |
| --------------------- | ------------------------------------------------- | ----------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `ErrorMemberNotExist` | `member did not exist (path=%s)`                  | The object/array member specified in the JSONPath did not exist in the JSON object. | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorMemberNotExist) |
| `ErrorTypeUnmatched`  | `type unmatched (expected=%s, found=%s, path=%s)` | The node type specified in the JSONPath did not exist in the JSON object.           | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorTypeUnmatched)  |
| `ErrorFunctionFailed` | `function failed (function=%s, error=%s)`         | The function specified in the JSONPath failed.                                      | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorFunctionFailed) |

The type checking is convenient to recognize which error happened.

```go
  :
  _,err := jsonpath.Retrieve(jsonPath, srcJSON)
  if err != nil {
    switch err.(type) {
    case jsonpath.ErrorMemberNotExist:
      fmt.printf(`retry with other srcJSON: %v`, err)
      continue
    case jsonpath.ErrorInvalidArgumentFormat:
      return nil, fmt.errorf(`specified invalid argument: %v`, err)
    }
    :
  }
```

### \* Function syntax

Function enables to format results by using user defined functions.
The function syntax comes after the JSONPath.

There are two ways to use function:

#### Filter function

The filter function applies a user function to each values in the result to get converted.

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetFilterFunction)

#### Aggregate function

The aggregate function converts all values in the result into a single value.

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetAggregateFunction)

### \* Accessing JSON

You can get the accessors ( _Getters / Setters_ ) of the input JSON instead of the retrieved values.
These accessors can use to update for the input JSON.

This feature can get enabled by giving `Config.SetAccessorMode()`.

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetAccessorMode)

#### Note:

It is not possible to use _Setter_ for some results, such as for JSONPath including function syntax.

Also, operations using accessors follow the map/slice manner of Go language.
If you use accessors after changing the structure of JSON, you need to pay attention to the behavior.
If you don't want to worry about it, get the accessor again every time you change the structure.

## Differences

Some behaviors that differ from the consensus exists in this library.
For the entire comparisons, please check [:memo: this result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html).

These behaviors will change in the future if appropriate ones found.

### Character types

The following character types can be available for identifiers in dot-child notation.

| Character type                                                                                                                                                                | Available | Escape required |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- | --------------- |
| \* Numbers and alphabets (`0-9` `A-Z` `a-z`)<br> \* Hyphen and underscore (`-` `_`)<br> \* Non-ASCII Unicode characters (`0x80 - 0x10FFFF`)                                     | Yes       | No              |
| \* Other printable symbols (`Space` `!` `"` `#` `$` `%` `&` `'` `(` `)` `*` `+` `,` `.` `/` `:` `;` `<` `=` `>` `?` `@` `[` `\` `]` `^` `` ` `` `{` <code>&#124;</code> `}` `~`) | Yes       | Yes             |
| \* ~~Control code characters~~ (`0x00 - 0x1F`, `0x7F`)                                                                                                                        | No        | -               |

The printable symbols except hyphen and underscore can use by escaping them.

```text
JSONPath : $.abc\.def
srcJSON  : {"abc.def":1}
Output   : [1]
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
JSONPath : $[?(@=~/(?i)CASE/)]
srcJSON  : ["Case","Hello"]
Output   : ["Case"]
```

### JSONPaths in the filter-qualifier

JSONPaths that return a value group cannot use with `comparator` or `regular expression`. However, `existence check` can use these syntaxes.

| JSONPaths that return a value group | example       |
| ----------------------------------- | ------------- |
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

JSONPath filter that begins with Root is a whole-match operation when any one is detected.

```text
JSONPath : $[?($..x)]
srcJSON  : [{"a":1},{"b":{"x":"hello world"}}]
Output   : [{"a":1},{"b":{"x":"hello world"}}]
```

## Benchmarks

The benchmarks for various JSONPath libraries in Go language can be compared in the following repository.

- [JSONPath Benchmark](https://github.com/AsaiYusuke/jsonpath-benchmark)

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
- Architecture
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
