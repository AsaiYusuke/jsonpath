# AsaiYusuke/JSONPath

[![Test](https://github.com/AsaiYusuke/jsonpath/actions/workflows/test.yml/badge.svg)](https://github.com/AsaiYusuke/jsonpath/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/AsaiYusuke/jsonpath/v2)](https://goreportcard.com/report/github.com/AsaiYusuke/jsonpath/v2)
[![Coverage Status](https://coveralls.io/repos/github/AsaiYusuke/jsonpath/badge.svg?branch=main)](https://coveralls.io/github/AsaiYusuke/jsonpath?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/AsaiYusuke/jsonpath/v2.svg)](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath/v2)
[![Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#query-language)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

![AsaiYusuke/JSONPath](assets/logo.svg)

This [Go](https://golang.org/) library allows you to extract parts of a JSON object using the _JSONPath_ query syntax.

The core JSONPath syntax supported by this library is based on:

- [Stefan Gössner's JSONPath - XPath for JSON](https://goessner.net/articles/JsonPath/)
- [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison)
- [JSONPath Internet Draft Development](https://github.com/ietf-wg-jsonpath/draft-ietf-jsonpath-jsonpath)

## Note

For details on syntax compatibility with other libraries, see [:memo: my comparison results](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html).

## Table of Contents

- [Getting started](#getting-started)
- [Basic design](#basic-design)
- [How to use](#how-to-use)
  - [Retrieve one-time or repeated](#-retrieve-one-time-or-repeatedly)
  - [Error handling](#-error-handling)
  - [Function syntax](#-function-syntax)
  - [Accessing JSON](#-accessing-json)
- [Differences](#differences)
- [Benchmarks](#benchmarks)
- [Project progress](#project-progress)

## Getting started

### Install

```bash
go get github.com/AsaiYusuke/jsonpath/v2
```

### Simple example

```go
package main

import (
  "encoding/json"
  "fmt"

  "github.com/AsaiYusuke/jsonpath/v2"
)

func main() {
  jsonPath, srcJSON := `$.key`, `{"key":"value"}`
  var src any
  json.Unmarshal([]byte(srcJSON), &src)
  output, _ := jsonpath.Retrieve(jsonPath, src)
  outputJSON, _ := json.Marshal(output)
  fmt.Println(string(outputJSON))
  // Output:
  // ["value"]
}
```

## Basic design

### _Streamlined Development_

- The JSONPath syntax parser is implemented using [PEG](https://github.com/pointlander/peg), which helps keep the source code simple and maintainable.
- Robust unit tests are provided to prevent bugs and ensure consistent results.

### _User-Friendly Interface_

- The library provides comprehensive error types, making it easy for users to handle errors appropriately.

### _High Compatibility_

- The library incorporates consensus behaviors from [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison), ensuring high compatibility with other implementations.

## How to use

### \* Retrieve one-time or repeatedly

The `Retrieve` function extracts values from a JSON object using a JSONPath expression:

```go
output, err := jsonpath.Retrieve(jsonPath, src)
```

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Retrieve)

The `Parse` function returns a parser function that checks the JSONPath syntax in advance. You can use this parser function to repeatedly extract values with the same JSONPath:

```go
parsed, err := jsonpath.Parse(jsonPath)
output1, err1 := parsed(src1)
output2, err2 := parsed(src2)
```

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Parse)

### \* Error handling

If an error occurs during API execution, a specific error type is returned. The following error types help you identify the cause:

#### Syntax errors from `Retrieve` or `Parse`

| Error type              | Message format                                     | Symptom                                                                                                          | Ex                                                                                        |
| ----------------------- | -------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| `ErrorInvalidSyntax`    | `invalid syntax (position=%d, reason=%s, near=%s)` | The JSONPath contains invalid syntax. The _reason_ in the message provides more details. | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorInvalidSyntax)    |
| `ErrorInvalidArgument`  | `invalid argument (argument=%s, error=%s)`         | An argument in the JSONPath is invalid according to Go syntax.                                | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorInvalidArgument)  |
| `ErrorFunctionNotFound` | `function not found (path=%s)`                 | The specified function in the JSONPath was not found.                                                             | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorFunctionNotFound) |
| `ErrorNotSupported`     | `not supported (path=%s, feature=%s)`              | The JSONPath uses unsupported syntax.                                                              | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorNotSupported)     |

#### Runtime errors from `Retrieve` or parser functions

| Error type            | Message format                                    | Symptom                                                                             | Ex                                                                                      |
| --------------------- | ------------------------------------------------- | ----------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `ErrorMemberNotExist` | `member did not exist (path=%s)`                  | The specified object or array member does not exist in the JSON object. | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorMemberNotExist) |
| `ErrorTypeUnmatched`  | `type unmatched (path=%s, expected=%s, found=%s)` | The type of the node in the JSON object does not match what is expected by the JSONPath.           | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorTypeUnmatched)  |
| `ErrorFunctionFailed` | `function failed (path=%s, error=%s)`         | The function specified in the JSONPath failed to execute.                                      | [:memo:](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-ErrorFunctionFailed) |

Type checking makes it easy to determine which error occurred.

```go
import jsonpath "github.com/AsaiYusuke/jsonpath/v2"
import errors "github.com/AsaiYusuke/jsonpath/v2/errors"

_, err := jsonpath.Retrieve(jsonPath, srcJSON)
switch err.(type) {
case errors.ErrorMemberNotExist:
  fmt.Printf("retry with other srcJSON: %v", err)
  // handle or continue
case errors.ErrorInvalidArgument:
  return nil, fmt.Errorf("specified invalid argument: %v", err)
}
```

### \* Function syntax

You can use user-defined functions to format results. The function syntax is appended after the JSONPath expression.

There are two types of functions:

#### Filter function

A filter function applies a user-defined function to each value in the result, transforming them individually.

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetFilterFunction)

#### Aggregate function

An aggregate function combines all values in the result into a single value.

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetAggregateFunction)

### \* Accessing JSON

Instead of retrieving values directly, you can obtain accessors (_Getters_ / _Setters_) for the input JSON. These accessors allow you to update the original JSON object.

Enable this feature by calling `Config.SetAccessorMode()`.

[:memo: Example](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath#example-Config.SetAccessorMode)

#### Accessor limitations

Setters are not available for some results, such as when using function syntax in the JSONPath.

Accessor operations follow Go's map/slice semantics. If you modify the structure of the JSON, be aware that accessors may not behave as expected. To avoid issues, obtain a new accessor each time you change the structure.

## Differences

Some behaviors in this library differ from the consensus of other implementations.
For a full comparison, see [:memo: this result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html).

These behaviors may change in the future if more appropriate approaches are found.

### Character types

The following character types are allowed for identifiers in dot-child notation:

| Character type                                                                                                                                                                | Available | Escape required |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------- | --------------- |
| \* Numbers and alphabets (`0-9` `A-Z` `a-z`)| Yes       | No              |
| \* Hyphen and underscore (`-` `_`)| Yes       | No              |
| \* Non-ASCII Unicode characters (`0x80 - 0x10FFFF`)                                     | Yes       | No              |
| \* Other printable symbols (`Space` `!` `"` `#` `$` `%` `&` `'` `(` `)` `*` `+` `,` `.` `/` `:` `;` `<` `=` `>` `?` `@` `[` `\` `]` `^` `` ` `` `{` `\|` `}` `~`) | Yes       | Yes             |
| \* ~~Control code characters~~ (`0x00 - 0x1F`, `0x7F`)                                                                                                                        | No        | -               |

Printable symbols (except hyphen and underscore) can be used by escaping them.

```text
JSONPath : $.abc\.def
srcJSON  : {"abc.def":1}
Output   : [1]
```

### Wildcard in qualifier

Wildcards in qualifiers can be specified as a union with subscripts.

```text
JSONPath : $[0,1:3,*]
srcJSON  : [0,1,2,3,4,5]
Output   : [0,1,2,0,1,2,3,4,5]
```

### Regular expression

Regular expression syntax follows Go's regular expression rules.
In particular, you can use "(?i)" to make the regular expression case-insensitive.

```text
JSONPath : $[?(@=~/(?i)CASE/)]
srcJSON  : ["Case","Hello"]
Output   : ["Case"]
```

### JSONPaths in the filter-qualifier

JSONPaths that return a value group cannot be used with a `comparator` or `regular expression`. However, you can use these syntaxes for `existence checks`.

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

If a JSONPath filter begins with the Root, it performs a whole-match operation if any match is found.

```text
JSONPath : $[?($..x)]
srcJSON  : [{"a":1},{"b":{"x":"hello world"}}]
Output   : [{"a":1},{"b":{"x":"hello world"}}]
```

## Benchmarks

Benchmark results for various Go JSONPath libraries (measured by myself) are available in the following repository:

- [JSONPath Benchmark](https://github.com/AsaiYusuke/jsonpath-benchmark)

## Project progress

- Syntax
  - Identifier
    - [x] identifier in dot notation
    - [x] identifier in bracket notation
    - [x] wildcard
    - [x] multiple identifiers in brackets
    - [x] recursive retrieval
  - Qualifier
    - [x] index
    - [x] slice
    - [x] wildcard
    - Filter
      - [x] logical operations
      - [x] comparators
      - [x] JSONPath retrieval in filter
    - [ ] script
  - Function
    - [x] filter
    - [x] aggregate
  - [x] Refer to the consensus behaviors
- Architecture
  - [x] PEG syntax analysis
  - [x] Error handling
  - [x] Function support
  - [x] JSON accessors
- Go language manner
  - [x] retrieve with an object unmarshaled to interface
  - [x] retrieve with the json.Number type
- Source code
  - [x] Release version
  - Unit tests
    - [x] syntax tests
    - [x] benchmarks
    - [x] coverage >80%
  - [x] Examples
  - [x] CI automation
  - Documentation
    - [x] README
    - [ ] API documentation
  - [x] comparison results (local)
- Development status
  - [x] requirements and functional design
    - [x] Decided to follow a standard or reference implementation for JSONPath syntax
  - [x] design-based coding
  - [ ] testing
  - [ ] documentation
- Future ToDo
  - Go language affinity
    - [ ] retrieve with an object unmarshaled to struct
    - [ ] retrieve with struct tags
    - [ ] retrieve with user-defined objects
