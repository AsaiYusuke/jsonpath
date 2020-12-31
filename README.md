# JSONPath

[![Build Status](https://travis-ci.com/AsaiYusuke/jsonpath.svg?branch=main)](https://travis-ci.com/AsaiYusuke/jsonpath)
[![Go Report Card](https://goreportcard.com/badge/github.com/AsaiYusuke/jsonpath)](https://goreportcard.com/report/github.com/AsaiYusuke/jsonpath)
[![Coverage Status](https://coveralls.io/repos/github/AsaiYusuke/jsonpath/badge.svg?branch=main)](https://coveralls.io/github/AsaiYusuke/jsonpath?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/AsaiYusuke/jsonpath.svg)](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This is [Go](https://golang.org/) package providing the features that retrieves a part of the JSON objects according to the statement written in the JSONPath syntax.

The core syntaxes of the JSONPath on which this package is based:  [JSONPath - XPath for JSON](https://goessner.net/articles/JsonPath/).

#### Note:
The unstated syntaxes found in "JSONPath - XPath for JSON" are implemented with reference to the test cases written in [cburmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison).
Please check [my compare result](https://asaiyusuke.github.io/jsonpath/cburgmer-json-path-comparison/docs/index.html) to know which responses are adapted.
Unfortunately, the proposals that is also discussing in "json-path-comparison" were not finalized at the start of development and were not adopted outright.

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
- Adopted more of the consensus behavior from the [cburmer's json-path-comparison](https://github.com/cburgmer/json-path-comparison).
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

If there is a problem with the execution of the `Retrieve`, `Prepare` or prepared *parser-functions*, an error type is returned.
These error types define the corresponding symptom, as listed below:

#### Syntax analyze errors from `Retrieve`, `Parse`

| Error type              | Message format                                     | Symptom                                                                                                       |
|-------------------------|----------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `ErrorInvalidSyntax`    | `invalid syntax (position=%d, reason=%s, near=%s)` | The invalid syntax found in the JSONPath. The *reason* including in this message will tell you more about it. |
| `ErrorInvalidArgument`  | `invalid argument (argument=%s, error=%s)`         | The argument specified in the JSONPath was treated as the invalid error in Go syntax.                         |
| `ErrorFunctionNotFound` | `function not found (function=%s)`                 | The function specified in the JSONPath is not found.                                                          |
| `ErrorNotSupported`     | `not supported (feature=%s, path=%s)`              | The unsupported syntaxes specified in the JSONPath.                                                           |

#### Runtime errors from `Retrieve`, *`parser-functions`*

| Error type             | Message format                                    | Symptom                                                                        |
|------------------------|---------------------------------------------------|--------------------------------------------------------------------------------|
| `ErrorMemberNotExist`  | `member did not exist (path=%s)`                  | The object member specified in the JSONPath did not exist in the JSON object.  |
| `ErrorIndexOutOfRange` | `index out of range (path=%s)`                    | The array indexes specified in the JSONPath were out of range.                 |
| `ErrorTypeUnmatched`   | `type unmatched (expected=%s, found=%s, path=%s)` | The node type specified in the JSONPath did not exist in the JSON object.      |
| `ErrorNoneMatched`     | `none matched (path=%s)`                          | The retrieving child paths specified in the JSONPath resulted in empty output. |
| `ErrorFunctionFailed`  | `function failed (function=%s, error=%s)`         | The function specified in the JSONPath failed.                                 |

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

The character types that can not be used for the identifiers in the dot notation are as follows :

```text
. [ ( ) = ! > < \t \r \n *SPACE*
```

You have to encode these characters when you enter them :

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

<details>
<summary>Show results</summary>

```text
goos: windows
goarch: amd64
pkg: github.com/AsaiYusuke/jsonpath
BenchmarkRetrieve_dotNotation-4                             	  802390	      1282 ns/op	     432 B/op	      18 allocs/op
BenchmarkRetrieve_bracketNotation-4                         	  801190	      1578 ns/op	     512 B/op	      21 allocs/op
BenchmarkRetrieve_wildcard_identifier_dotNotation-4         	  856964	      1420 ns/op	     456 B/op	      20 allocs/op
BenchmarkRetrieve_wildcard_identifier_bracketNotation-4     	  707596	      1614 ns/op	     504 B/op	      21 allocs/op
BenchmarkRetrieve_multi_identifier-4                        	  480096	      2443 ns/op	     856 B/op	      33 allocs/op
BenchmarkRetrieve_qualifier_index-4                         	  750308	      1621 ns/op	     568 B/op	      22 allocs/op
BenchmarkRetrieve_qualifier_slice-4                         	  600488	      2106 ns/op	     688 B/op	      29 allocs/op
BenchmarkRetrieve_qualifier_wildcard-4                      	  800319	      1550 ns/op	     504 B/op	      21 allocs/op
BenchmarkRetrieve_qualifier_union-4                         	  462747	      2539 ns/op	     904 B/op	      34 allocs/op
BenchmarkRetrieve_filter_logicalOR-4                        	  218623	      5463 ns/op	    1853 B/op	      45 allocs/op
BenchmarkRetrieve_filter_logicalAND-4                       	  226866	      5465 ns/op	    1853 B/op	      45 allocs/op
BenchmarkRetrieve_filter_nodeFilter-4                       	  260784	      4746 ns/op	    1488 B/op	      41 allocs/op
BenchmarkRetrieve_filter_logicalNOT-4                       	  261564	      4606 ns/op	    1840 B/op	      47 allocs/op
BenchmarkRetrieve_filter_compareEQ-4                        	  227083	      5236 ns/op	    1992 B/op	      51 allocs/op
BenchmarkRetrieve_filter_compareNE-4                        	  221324	      5512 ns/op	    2264 B/op	      54 allocs/op
BenchmarkRetrieve_filter_compareGE-4                        	  231394	      5343 ns/op	    1992 B/op	      51 allocs/op
BenchmarkRetrieve_filter_compareGT-4                        	  237956	      5232 ns/op	    1880 B/op	      50 allocs/op
BenchmarkRetrieve_filter_compareLE-4                        	  231217	      5232 ns/op	    1976 B/op	      50 allocs/op
BenchmarkRetrieve_filter_compareLT-4                        	  227025	      5219 ns/op	    1880 B/op	      49 allocs/op
BenchmarkRetrieve_filter_regex-4                            	  168806	      7237 ns/op	    2833 B/op	      62 allocs/op
BenchmarkParserFunc_dotNotation-4                           	 4852603	       248 ns/op	     144 B/op	       5 allocs/op
BenchmarkParserFunc_bracketNotation-4                       	 4877554	       247 ns/op	     144 B/op	       5 allocs/op
BenchmarkParserFunc_wildcard_identifier_dotNotation-4       	 3003252	       399 ns/op	     192 B/op	       7 allocs/op
BenchmarkParserFunc_wildcard_identifier_bracketNotation-4   	 2943961	       408 ns/op	     192 B/op	       7 allocs/op
BenchmarkParserFunc_multi_identifier-4                      	 2593132	       456 ns/op	     240 B/op	       8 allocs/op
BenchmarkParserFunc_qualifier_index-4                       	 3682868	       326 ns/op	     224 B/op	       7 allocs/op
BenchmarkParserFunc_qualifier_slice-4                       	 3730786	       330 ns/op	     192 B/op	       7 allocs/op
BenchmarkParserFunc_qualifier_wildcard-4                    	 3702061	       321 ns/op	     192 B/op	       7 allocs/op
BenchmarkParserFunc_qualifier_union-4                       	 2053358	       587 ns/op	     344 B/op	      12 allocs/op
BenchmarkParserFunc_filter_logicalOR-4                      	  816764	      1380 ns/op	    1104 B/op	      17 allocs/op
BenchmarkParserFunc_filter_logicalAND-4                     	  925590	      1375 ns/op	    1104 B/op	      17 allocs/op
BenchmarkParserFunc_filter_nodeFilter-4                     	 1000000	      1009 ns/op	     832 B/op	      14 allocs/op
BenchmarkParserFunc_filter_logicalNOT-4                     	  857142	      1533 ns/op	    1168 B/op	      19 allocs/op
BenchmarkParserFunc_filter_compareEQ-4                      	  750130	      1575 ns/op	    1088 B/op	      16 allocs/op
BenchmarkParserFunc_filter_compareNE-4                      	  667159	      1823 ns/op	    1344 B/op	      18 allocs/op
BenchmarkParserFunc_filter_compareGE-4                      	  858307	      1471 ns/op	    1088 B/op	      16 allocs/op
BenchmarkParserFunc_filter_compareGT-4                      	  802525	      1477 ns/op	    1088 B/op	      16 allocs/op
BenchmarkParserFunc_filter_compareLE-4                      	  833859	      1482 ns/op	    1088 B/op	      16 allocs/op
BenchmarkParserFunc_filter_compareLT-4                      	  707726	      1497 ns/op	    1088 B/op	      16 allocs/op
BenchmarkParserFunc_filter_regex-4                          	  799903	      1633 ns/op	    1095 B/op	      16 allocs/op
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
