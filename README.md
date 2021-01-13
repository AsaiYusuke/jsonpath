# JSONPath

[![Build Status](https://travis-ci.com/AsaiYusuke/jsonpath.svg?branch=main)](https://travis-ci.com/AsaiYusuke/jsonpath)
[![Go Report Card](https://goreportcard.com/badge/github.com/AsaiYusuke/jsonpath)](https://goreportcard.com/report/github.com/AsaiYusuke/jsonpath)
[![Coverage Status](https://coveralls.io/repos/github/AsaiYusuke/jsonpath/badge.svg?branch=main)](https://coveralls.io/github/AsaiYusuke/jsonpath?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/AsaiYusuke/jsonpath.svg)](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is [Go](https://golang.org/) package providing the features that retrieves a part of the JSON objects according to the query written in the JSONPath syntax.

The core syntaxes of the JSONPath on which this package is based:

- [Stefan Gössner's JSONPath - XPath for JSON](https://goessner.net/articles/JsonPath/)
- [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison)
- [JSONPath Internet Draft Development](https://github.com/ietf-wg-jsonpath/draft-ietf-jsonpath-jsonpath)

#### Note:
Please check [my compare result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html) to know which responses are adapted.
Unfortunately, the proposals that is also discussing in "json-path-comparison"  and the draft of the Internet Draft were not finalized at the start of development and are not adopted outright.

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
- The error specification allows package users to handle errors appropriately.
- Adopted more of the consensus behavior from the [Christoph Burgmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison).
  Adapted my own behavior to the other part of the such consensus behavior that found difficult to use.
- Equipped with numerous unit tests and tried to eliminate the bugs that return strange result.

## How to use

### * Retrieve one-time, or successively

The `Retrieve` function returns a retrieved JSON object by a one-time sequential operation (analyzing syntax and retrieving objects) using the given JSONPath and the source JSON object :

```go
output, err := jsonpath.Retrieve(jsonPath, src)
```

The `Parse` function returns a *parser-function* that completed to analyze the JSONPath syntax.
By using this returned *parser-function* it can be performed successively a retrieve with the same JSONPath syntax :

```go
jsonPath, err := jsonpath.Parse(jsonPath)
output1, err1 := jsonPath(src1)
output2, err2 := jsonPath(src2)
:
```

### * Error handling

If there is a problem with the execution of the `Retrieve`, `Parse` or prepared *parser-functions*, an error type is returned.
These error types define the corresponding symptom, as listed below:

#### Syntax analyze errors from `Retrieve`, `Parse`

| Error type              | Message format                                     | Symptom                                                                                                       |
|-------------------------|----------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `ErrorInvalidSyntax`    | `invalid syntax (position=%d, reason=%s, near=%s)` | The invalid syntax found in the JSONPath. The *reason* including in this message will tell you more about it. |
| `ErrorInvalidArgument`  | `invalid argument (argument=%s, error=%s)`         | The argument specified in the JSONPath was treated as the invalid error in Go syntax.                         |
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

Function is a feature that allows you to format JSONPath results by using pre-registered user functions and the instruction syntaxes at the end of the JSONPath statement.

There are two ways to use function:

#### Filter function

The filter function applies a user function to each values in the JSONPath result to get converted.

```Go
  config := jsonpath.Config{}
  config.SetFilterFunction(`twice`, func(param interface{}) (interface{}, error) {
    if floatParam, ok := param.(float64); ok {
      return floatParam * 2, nil
    }
    return nil, fmt.Errorf(`type error`)
  })
  jsonPath, srcJSON := `$[*].twice()`, `[1,3]`
  var src interface{}
  json.Unmarshal([]byte(srcJSON), &src)
  output, _ := jsonpath.Retrieve(jsonPath, src, config)
  outputJSON, _ := json.Marshal(output)
  fmt.Println(string(outputJSON))
  // Output:
  // [2,6]
```

#### Aggregate function

Aggregate function converts all values in the JSONPath result into a single value by applying them to a user function.

```Go
  config := jsonpath.Config{}
  config.SetAggregateFunction(`max`, func(params []interface{}) (interface{}, error) {
    var result float64
    for _, param := range params {
      if floatParam, ok := param.(float64); ok {
        if result < floatParam {
          result = floatParam
        }
        continue
      }
      return nil, fmt.Errorf(`type error`)
    }
    return result, nil
  })
  jsonPath, srcJSON := `$[*].max()`, `[1,3]`
  var src interface{}
  json.Unmarshal([]byte(srcJSON), &src)
  output, _ := jsonpath.Retrieve(jsonPath, src, config)
  outputJSON, _ := json.Marshal(output)
  fmt.Println(string(outputJSON))
  // Output:
  // [3]
```

### * Accessing JSON

You can get a collection of accessors ( *Getters* / *Setters* ) to the input JSON instead of the retrieved values by giving `Config.SetAccessorMode()`.
These accessors can be used to update the original nodes retrieved by JSONPath in the input JSON.
See the Example for usage.

#### Note:
It is not possible to use *Setter* for some execution results, such as including function syntax.

Also, operations using accessors follow the map/slice manner of Go language, so if you use accessors after changing the structure of JSON, you need to pay attention to the behavior caused by the operation.
If you want to handle it casually, you may want to retrieve the accessor again each time you change the structure of JSON.

## Differences

Some behaviors that differ from the consensus exists in this package.
For the entire comparisons, please check [this result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html) to see which responses are different.
These behaviors will be changed in the future if appropriate ones are found.

### Character types

The character types that can be used for identifiers in dot child notation are as follows:

| Character type                   | Availabe | Escape |
|----------------------------------|----------|--------|
| Control code (0x00 - 0x1F, 0x7F) | No       | -      |
| `- _`                            | Yes      | No     |
| ``Space ! " # $ % & ' ( ) * + , . / : ; < = > ? @ [ \ ] ^ ` { \| } ~`` | Yes | Yes |

The character types that can be used for the proposal (hyphen, underscore) can be used without escape, and other printable symbols can be used by escaping them.

```text
JSONPath : $.abc\.def
srcJSON  : {"abc.def":1}
Output   : 1
```

### Wildcard in qualifier

The wildcard in qualifier can be specified mixed with other subscript syntaxes.

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

In the case of the `comparators` and `regular expressions` in the filter qualifier, the following JSONPaths that return a value group cannot be specified.
On the other hand, in the case of the `existence check` in the filter qualifier, it can be specified.

| JSONPaths that return a value group | example |
| :------- | :------ |
| Recursive descent | `@..a` |
| Multiple identifier  | `@['a','b']` |
| Wildcard identifier | `@.*` |
| Slice qualifier | `@[0:1]` |
| Wildcard qualifier | `@[*]` |
| Union in the qualifier | `@[0,1]` |
| Filter qualifier | `@.a[?(@.b)]` |

## Benchmarks

I benchmarked three JSONPaths using several libraries for the Go language.
There was a difference in execution performance between the libraries, but if the number of queries using JSONPaths is little, I don't think there will be a big difference between any of them.

<details>
<summary>Show results</summary>

- [BenchmarkAsaiYusukeJSONPath](https://github.com/AsaiYusuke/jsonpath)
- [BenchmarkOhler55Ojg](https://github.com/ohler55/ojg/jp)
- [BenchmarkBhmjJSONSlice](https://github.com/bhmj/jsonslice)
- [BenchmarkPaesslerAGJSONPath](https://github.com/PaesslerAG/jsonpath)
- [BenchmarkOliveagleJsonpath](https://github.com/oliveagle/jsonpath)
- [BenchmarkSpyzhovAjson](https://github.com/spyzhov/ajson)

```text
JSONPath : $.store.book[0].price

BenchmarkAsaiYusukeJSONPath_threeLevelsWithIndex-4          4557238       265 ns/op      64 B/op       4 allocs/op
BenchmarkOhler55Ojg_threeLevelsWithIndex-4                  1711050       716 ns/op    1040 B/op       2 allocs/op
BenchmarkBhmjJSONSlice_threeLevelsWithIndex-4                601647      2087 ns/op      32 B/op       1 allocs/op
BenchmarkPaesslerAGJSONPath_threeLevelsWithIndex-4          1814469       651 ns/op     208 B/op       7 allocs/op
BenchmarkOliveagleJsonpath_threeLevelsWithIndex-4          12950515        90.3 ns/op     0 B/op       0 allocs/op
BenchmarkSpyzhovAjson                                      not supported
```

```text
JSONPath : $.store..price

BenchmarkAsaiYusukeJSONPath_recursiveDescent-4               429724      2716 ns/op     640 B/op      13 allocs/op
BenchmarkOhler55Ojg_recursiveDescent-4                       481300      2513 ns/op    1368 B/op      18 allocs/op
BenchmarkBhmjJSONSlice_recursiveDescent-4                    129428      9289 ns/op     688 B/op      27 allocs/op
BenchmarkPaesslerAGJSONPath_recursiveDescent-4                42361     28190 ns/op   12167 B/op     397 allocs/op
BenchmarkOliveagleJsonpath                                 not supported
BenchmarkSpyzhovAjson_recursiveDescent-4                     158319      7473 ns/op    2280 B/op      72 allocs/op
```

```text
JSONPath : $..book[?(@.price > $.store.bicycle.price)]

BenchmarkAsaiYusukeJSONPath_recursiveDescentWithFilter-4     331218      3572 ns/op     816 B/op      27 allocs/op
BenchmarkOhler55Ojg_recursiveDescentWithFilter-4             222663      5361 ns/op    5240 B/op      20 allocs/op
BenchmarkBhmjJSONSlice_recursiveDescentWithFilter-4           55168     26743 ns/op    3032 B/op      57 allocs/op
BenchmarkPaesslerAGJSONPath                                not supported
BenchmarkOliveagleJsonpath                                 not supported
BenchmarkSpyzhovAjson                                      not supported
```

JSON used for the benchmark measurement

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

Environment

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
