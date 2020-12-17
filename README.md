# JSONPath

[![Build Status](https://travis-ci.com/AsaiYusuke/jsonpath.svg?branch=main)](https://travis-ci.com/AsaiYusuke/jsonpath)
[![Go Report Card](https://goreportcard.com/badge/github.com/AsaiYusuke/jsonpath)](https://goreportcard.com/report/github.com/AsaiYusuke/jsonpath)
[![Coverage Status](https://coveralls.io/repos/github/AsaiYusuke/jsonpath/badge.svg?branch=main)](https://coveralls.io/github/AsaiYusuke/jsonpath?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/AsaiYusuke/jsonpath.svg)](https://pkg.go.dev/github.com/AsaiYusuke/jsonpath)

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

| Error type             | Message format                                     | Symptom                                                                                                       |
|------------------------|----------------------------------------------------|---------------------------------------------------------------------------------------------------------------|
| `ErrorInvalidSyntax`   | `invalid syntax (position=%d, reason=%s, near=%s)` | The invalid syntax found in the JSONPath. The *reason* including in this message will tell you more about it. |
| `ErrorInvalidArgument` | `invalid argument (argument=%s, error=%s)`         | The argument specified in the JSONPath was treated as the invalid error in Go syntax.                         |
| `ErrorNotSupported`    | `not supported (feature=%s, path=%s)`              | The unsupported syntaxes specified in the JSONPath.                                                           |

#### Runtime errors from `Retrieve`, *`parser-functions`*

| Error type             | Message format                                    | Symptom                                                                        |
|------------------------|---------------------------------------------------|--------------------------------------------------------------------------------|
| `ErrorMemberNotExist`  | `member did not exist (path=%s)`                  | The object member specified in the JSONPath did not exist in the JSON object.  |
| `ErrorIndexOutOfRange` | `index out of range (path=%s)`                    | The array indexes specified in the JSONPath were out of range.                 |
| `ErrorTypeUnmatched`   | `type unmatched (expected=%s, found=%s, path=%s)` | The node type specified in the JSONPath did not exist in the JSON object.      |
| `ErrorNoneMatched`     | `none matched (path=%s)`                          | The retrieving child paths specified in the JSONPath resulted in empty output. |

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

### Asterisk in qualifier

The asterisk in qualifier can be specified mixed with other subscript syntaxes.

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
| Asterisk identifier | `@.*` |
| Slice qualifier | `@[0:1]` |
| Asterisk qualifier | `@[*]` |
| Union in the qualifier | `@[0,1]` |
| Filter qualifier | `@.a[?(@.b)]` |
## Benchmarks

<details>
<summary>Show results</summary>

```text
goos: windows
goarch: amd64
pkg: github.com/AsaiYusuke/jsonpath
BenchmarkRetrieve_dotNotation-4                             	   11193	     97395 ns/op	  533524 B/op	     144 allocs/op
BenchmarkRetrieve_bracketNotation-4                         	    9456	    109682 ns/op	  533610 B/op	     147 allocs/op
BenchmarkRetrieve_asterisk_identifier_dotNotation-4         	   10000	    108195 ns/op	  533538 B/op	     145 allocs/op
BenchmarkRetrieve_asterisk_identifier_bracketNotation-4     	   10000	    105758 ns/op	  533594 B/op	     146 allocs/op
BenchmarkRetrieve_multi_identifier-4                        	   10000	    109516 ns/op	  533889 B/op	     157 allocs/op
BenchmarkRetrieve_qualifier_index-4                         	   10000	    108211 ns/op	  533601 B/op	     148 allocs/op
BenchmarkRetrieve_qualifier_slice-4                         	   10000	    120385 ns/op	  533808 B/op	     152 allocs/op
BenchmarkRetrieve_qualifier_asterisk-4                      	   10000	    110426 ns/op	  533561 B/op	     146 allocs/op
BenchmarkRetrieve_qualifier_union-4                         	    9457	    118281 ns/op	  533892 B/op	     159 allocs/op
BenchmarkRetrieve_filter_logicalOR-4                        	    9456	    133341 ns/op	  534582 B/op	     159 allocs/op
BenchmarkRetrieve_filter_logicalAND-4                       	    8197	    132956 ns/op	  534583 B/op	     159 allocs/op
BenchmarkRetrieve_filter_nodeFilter-4                       	   10000	    120922 ns/op	  534319 B/op	     158 allocs/op
BenchmarkRetrieve_filter_logicalNOT-4                       	    9463	    130181 ns/op	  534389 B/op	     161 allocs/op
BenchmarkRetrieve_filter_compareEQ-4                        	    8197	    143824 ns/op	  534771 B/op	     166 allocs/op
BenchmarkRetrieve_filter_compareNE-4                        	    7681	    151789 ns/op	  534787 B/op	     167 allocs/op
BenchmarkRetrieve_filter_compareGE-4                        	    8194	    143069 ns/op	  534768 B/op	     166 allocs/op
BenchmarkRetrieve_filter_compareGT-4                        	    8782	    150355 ns/op	  534661 B/op	     165 allocs/op
BenchmarkRetrieve_filter_compareLE-4                        	    8778	    142800 ns/op	  534754 B/op	     165 allocs/op
BenchmarkRetrieve_filter_compareLT-4                        	    9458	    142593 ns/op	  534656 B/op	     164 allocs/op
BenchmarkRetrieve_filter_regex-4                            	    8781	    147454 ns/op	  543052 B/op	     178 allocs/op
BenchmarkParserFunc_dotNotation-4                           	 7589292	       163 ns/op	      48 B/op	       2 allocs/op
BenchmarkParserFunc_bracketNotation-4                       	 8721170	       136 ns/op	      48 B/op	       2 allocs/op
BenchmarkParserFunc_asterisk_identifier_dotNotation-4       	 4210114	       281 ns/op	      96 B/op	       4 allocs/op
BenchmarkParserFunc_asterisk_identifier_bracketNotation-4   	 3898969	       297 ns/op	      96 B/op	       4 allocs/op
BenchmarkParserFunc_multi_identifier-4                      	 4917340	       248 ns/op	      80 B/op	       3 allocs/op
BenchmarkParserFunc_qualifier_index-4                       	 6505556	       189 ns/op	      64 B/op	       4 allocs/op
BenchmarkParserFunc_qualifier_slice-4                       	 5469730	       225 ns/op	      64 B/op	       4 allocs/op
BenchmarkParserFunc_qualifier_asterisk-4                    	 6151656	       194 ns/op	      64 B/op	       4 allocs/op
BenchmarkParserFunc_qualifier_union-4                       	 3419431	       358 ns/op	     120 B/op	       7 allocs/op
BenchmarkParserFunc_filter_logicalOR-4                      	 1000000	      1181 ns/op	     912 B/op	      12 allocs/op
BenchmarkParserFunc_filter_logicalAND-4                     	  945879	      1193 ns/op	     912 B/op	      12 allocs/op
BenchmarkParserFunc_filter_nodeFilter-4                     	 1575550	       751 ns/op	     608 B/op	       8 allocs/op
BenchmarkParserFunc_filter_logicalNOT-4                     	 1215165	      1000 ns/op	     656 B/op	      10 allocs/op
BenchmarkParserFunc_filter_compareEQ-4                      	  876064	      1363 ns/op	     864 B/op	      10 allocs/op
BenchmarkParserFunc_filter_compareNE-4                      	  875866	      1418 ns/op	     864 B/op	      10 allocs/op
BenchmarkParserFunc_filter_compareGE-4                      	  942987	      1279 ns/op	     864 B/op	      10 allocs/op
BenchmarkParserFunc_filter_compareGT-4                      	  942847	      1267 ns/op	     864 B/op	      10 allocs/op
BenchmarkParserFunc_filter_compareLE-4                      	  945812	      1267 ns/op	     864 B/op	      10 allocs/op
BenchmarkParserFunc_filter_compareLT-4                      	  942994	      1274 ns/op	     864 B/op	      10 allocs/op
BenchmarkParserFunc_filter_regex-4                          	  862663	      1411 ns/op	     871 B/op	      10 allocs/op
```
</details>

## Project progress

- Syntax
  - Identifier
    - [x] identifier in dot notations
    - [x] identifier in bracket notations
    - [x] asterisk identifier
    - [x] multiple-identifier in bracket
    - [x] recursive retrieve
  - Qualifier
    - [x] index
    - [x] slice
    - [x] asterisk
    - Filter
      - [x] logical operation
      - [x] comparator
      - [x] JSONPath retrieve in filter
    - [ ] script
  - [x] Refer to the consensus behaviors
- Archtecture
  - [x] PEG syntax analyzing
  - [x] Error handling
- Go language manner
  - [x] retrieve with the object in interface unmarshal
  - [x] retrieve with the json.Number type
- Source code
  - [x] Release version
  - Unit tests
    - [x] syntax tests
    - [x] benchmark
    - [x] coverage >80%
  - [ ] Examples
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
  - [ ] Syntax expansion
  - [ ] Refer to the something standard
  - Go language affinity
    - [ ] retrieve with the object in struct unmarshal
    - [ ] retrieve with the struct tags
    - [ ] retrieve with the user defined objects
