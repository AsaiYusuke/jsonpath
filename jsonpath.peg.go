package jsonpath

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleexpression
	ruleEND
	rulejsonpath
	rulerootNode
	rulechildNodes
	rulebracketNode
	rulerootIdentifier
	rulecurrentRootIdentifier
	ruledotChildIdentifier
	rulebracketChildIdentifier
	rulebracketNodeIdentifiers
	rulesingleQuotedNodeIdentifier
	ruledoubleQuotedNodeIdentifier
	rulesepBracketIdentifier
	rulequalifier
	ruleunion
	ruleindex
	ruleslice
	ruleanyIndex
	ruleindexNumber
	rulesepUnion
	rulesepSlice
	rulescript
	rulecommand
	rulefilter
	rulequery
	ruleandQuery
	rulebasicQuery
	rulelogicOr
	rulelogicAnd
	rulelogicNot
	rulecomparator
	ruleqParam
	ruleqNumericParam
	ruleqLiteral
	rulesingleJsonpathFilter
	rulejsonpathFilter
	rulelNumber
	rulelBool
	rulelString
	rulelNull
	ruleregex
	rulesquareBracketStart
	rulesquareBracketEnd
	rulescriptStart
	rulescriptEnd
	rulefilterStart
	rulefilterEnd
	rulesubQueryStart
	rulesubQueryEnd
	rulespace
	ruleAction0
	rulePegText
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12
	ruleAction13
	ruleAction14
	ruleAction15
	ruleAction16
	ruleAction17
	ruleAction18
	ruleAction19
	ruleAction20
	ruleAction21
	ruleAction22
	ruleAction23
	ruleAction24
	ruleAction25
	ruleAction26
	ruleAction27
	ruleAction28
	ruleAction29
	ruleAction30
	ruleAction31
	ruleAction32
	ruleAction33
	ruleAction34
	ruleAction35
	ruleAction36
	ruleAction37
	ruleAction38
	ruleAction39
	ruleAction40
	ruleAction41
	ruleAction42
	ruleAction43
	ruleAction44
	ruleAction45
	ruleAction46
	ruleAction47
	ruleAction48

	rulePre
	ruleIn
	ruleSuf
)

var rul3s = [...]string{
	"Unknown",
	"expression",
	"END",
	"jsonpath",
	"rootNode",
	"childNodes",
	"bracketNode",
	"rootIdentifier",
	"currentRootIdentifier",
	"dotChildIdentifier",
	"bracketChildIdentifier",
	"bracketNodeIdentifiers",
	"singleQuotedNodeIdentifier",
	"doubleQuotedNodeIdentifier",
	"sepBracketIdentifier",
	"qualifier",
	"union",
	"index",
	"slice",
	"anyIndex",
	"indexNumber",
	"sepUnion",
	"sepSlice",
	"script",
	"command",
	"filter",
	"query",
	"andQuery",
	"basicQuery",
	"logicOr",
	"logicAnd",
	"logicNot",
	"comparator",
	"qParam",
	"qNumericParam",
	"qLiteral",
	"singleJsonpathFilter",
	"jsonpathFilter",
	"lNumber",
	"lBool",
	"lString",
	"lNull",
	"regex",
	"squareBracketStart",
	"squareBracketEnd",
	"scriptStart",
	"scriptEnd",
	"filterStart",
	"filterEnd",
	"subQueryStart",
	"subQueryEnd",
	"space",
	"Action0",
	"PegText",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
	"Action10",
	"Action11",
	"Action12",
	"Action13",
	"Action14",
	"Action15",
	"Action16",
	"Action17",
	"Action18",
	"Action19",
	"Action20",
	"Action21",
	"Action22",
	"Action23",
	"Action24",
	"Action25",
	"Action26",
	"Action27",
	"Action28",
	"Action29",
	"Action30",
	"Action31",
	"Action32",
	"Action33",
	"Action34",
	"Action35",
	"Action36",
	"Action37",
	"Action38",
	"Action39",
	"Action40",
	"Action41",
	"Action42",
	"Action43",
	"Action44",
	"Action45",
	"Action46",
	"Action47",
	"Action48",

	"Pre_",
	"_In_",
	"_Suf",
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (node *node32) Print(buffer string) {
	node.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: ruleIn, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: ruleSuf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

func (t *tokens32) Expand(index int) {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
}

type pegJSONPathParser struct {
	jsonPathParser

	Buffer string
	buffer []rune
	rules  [102]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	Pretty bool
	tokens32
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *pegJSONPathParser
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *pegJSONPathParser) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *pegJSONPathParser) Highlighter() {
	p.PrintSyntax()
}

func (p *pegJSONPathParser) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:

			p.root = p.pop().(syntaxNode)

		case ruleAction1:

			p.syntaxErr(begin, msgErrorInvalidSyntaxUnrecognizedInput, buffer)

		case ruleAction2:

			child := p.pop().(syntaxNode)
			root := p.pop().(syntaxNode)
			root.setNext(child)
			p.push(root)

		case ruleAction3:

			rootNode := p.pop().(syntaxNode)
			checkNode := rootNode
			for checkNode != nil {
				if checkNode.isMultiValue() {
					rootNode.setMultiValue()
					break
				}
				checkNode = checkNode.getNext()
			}
			p.push(rootNode)

		case ruleAction4:

			if len(p.params) == 1 {
				p.syntaxErr(begin, msgErrorInvalidSyntaxUseBeginAtsign, buffer)
			}

		case ruleAction5:

			if len(p.params) != 1 {
				p.syntaxErr(begin, msgErrorInvalidSyntaxOmitDollar, buffer)
			}

		case ruleAction6:

			node := p.pop().(syntaxNode)
			p.push(&syntaxRecursiveChildIdentifier{
				syntaxBasicNode: &syntaxBasicNode{
					text:       `..`,
					multiValue: true,
					next:       node,
					result:     &p.resultPtr,
				},
			})

		case ruleAction7:

			identifier := p.pop().(syntaxNode)
			identifier.setText(text)
			p.push(identifier)

		case ruleAction8:

			child := p.pop().(syntaxNode)
			parent := p.pop().(syntaxNode)
			parent.setNext(child)
			p.push(parent)

		case ruleAction9:

			node := p.pop().(syntaxNode)
			node.setText(text)
			p.push(node)

		case ruleAction10:

			p.push(&syntaxRootIdentifier{
				syntaxBasicNode: &syntaxBasicNode{
					text:   `$`,
					result: &p.resultPtr,
				},
			})

		case ruleAction11:

			p.push(&syntaxCurrentRootIdentifier{
				syntaxBasicNode: &syntaxBasicNode{
					text:   `@`,
					result: &p.resultPtr,
				},
			})

		case ruleAction12:

			unescapedText := p.unescape(text)
			if unescapedText == `*` {
				p.push(&syntaxChildAsteriskIdentifier{
					syntaxBasicNode: &syntaxBasicNode{
						text:       unescapedText,
						multiValue: true,
						result:     &p.resultPtr,
					},
				})
			} else {
				p.push(&syntaxChildSingleIdentifier{
					identifier: unescapedText,
					syntaxBasicNode: &syntaxBasicNode{
						text:       unescapedText,
						multiValue: false,
						result:     &p.resultPtr,
					},
				})
			}

		case ruleAction13:

			identifier := p.pop().([]string)
			if len(identifier) > 1 {
				p.push(&syntaxChildMultiIdentifier{
					identifiers: identifier,
					syntaxBasicNode: &syntaxBasicNode{
						multiValue: true,
						result:     &p.resultPtr,
					},
				})
			} else {
				p.push(&syntaxChildSingleIdentifier{
					identifier: identifier[0],
					syntaxBasicNode: &syntaxBasicNode{
						multiValue: false,
						result:     &p.resultPtr,
					},
				})
			}

		case ruleAction14:

			p.push([]string{p.pop().(string)})

		case ruleAction15:

			identifier2 := p.pop().([]string)
			identifier1 := p.pop().([]string)
			identifier1 = append(identifier1, identifier2...)
			p.push(identifier1)

		case ruleAction16:

			p.push(p.unescape(text))

		case ruleAction17:
			// '
			p.push(p.unescape(text))

		case ruleAction18:

			subscript := p.pop().(syntaxSubscript)
			p.push(&syntaxUnionQualifier{
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: subscript.isMultiValue(),
					result:     &p.resultPtr,
				},
				subscripts: []syntaxSubscript{subscript},
			})

		case ruleAction19:

			childIndexUnion := p.pop().(*syntaxUnionQualifier)
			parentIndexUnion := p.pop().(*syntaxUnionQualifier)
			parentIndexUnion.merge(childIndexUnion)
			parentIndexUnion.setMultiValue()
			p.push(parentIndexUnion)

		case ruleAction20:

			step := p.pop().(*syntaxIndex)
			end := p.pop().(*syntaxIndex)
			start := p.pop().(*syntaxIndex)

			if step.isOmitted || step.number == 0 {
				step.number = 1
			}

			if step.number > 0 {
				p.push(&syntaxSlicePositiveStep{
					syntaxBasicSubscript: &syntaxBasicSubscript{
						multiValue: true,
					},
					start: start,
					end:   end,
					step:  step,
				})
			} else {
				p.push(&syntaxSliceNegativeStep{
					syntaxBasicSubscript: &syntaxBasicSubscript{
						multiValue: true,
					},
					start: start,
					end:   end,
					step:  step,
				})
			}

		case ruleAction21:

			p.push(&syntaxIndex{
				syntaxBasicSubscript: &syntaxBasicSubscript{
					multiValue: false,
				},
				number: p.toInt(text),
			})

		case ruleAction22:

			p.push(&syntaxAsterisk{
				syntaxBasicSubscript: &syntaxBasicSubscript{
					multiValue: true,
				},
			})

		case ruleAction23:

			p.push(&syntaxIndex{number: 1})

		case ruleAction24:

			if len(text) > 0 {
				p.push(&syntaxIndex{number: p.toInt(text)})
			} else {
				p.push(&syntaxIndex{number: 0, isOmitted: true})
			}

		case ruleAction25:

			p.push(&syntaxScriptQualifier{
				command: text,
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: true,
					result:     &p.resultPtr,
				},
			})

		case ruleAction26:

			query := p.pop().(syntaxQuery)
			p.push(&syntaxFilterQualifier{
				query: query,
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: true,
					result:     &p.resultPtr,
				},
			})

		case ruleAction27:

			childQuery := p.pop().(syntaxQuery)
			parentQuery := p.pop().(syntaxQuery)
			p.push(&syntaxLogicalOr{parentQuery, childQuery})

		case ruleAction28:

			childQuery := p.pop().(syntaxQuery)
			parentQuery := p.pop().(syntaxQuery)
			p.push(&syntaxLogicalAnd{parentQuery, childQuery})

		case ruleAction29:

			if !p.hasErr() {
				query := p.pop().(syntaxQuery)
				p.push(query)

				if logicalNot, ok := query.(*syntaxLogicalNot); ok {
					query = (*logicalNot).param
				}
				if checkQuery, ok := query.(*syntaxBasicCompareQuery); ok {
					_, leftIsCurrentRoot := checkQuery.leftParam.param.(*syntaxQueryParamCurrentRoot)
					_, rigthIsCurrentRoot := checkQuery.rightParam.param.(*syntaxQueryParamCurrentRoot)
					if leftIsCurrentRoot && rigthIsCurrentRoot {
						p.syntaxErr(begin, msgErrorInvalidSyntaxTwoCurrentNode, buffer)
					}
				}
			}

		case ruleAction30:

			p.push(strings.HasPrefix(text, `!`))

		case ruleAction31:

			_ = p.pop().(bool)
			jsonpathFilter := p.pop().(syntaxQuery)
			isLogicalNot := p.pop().(bool)
			if isLogicalNot {
				p.push(&syntaxLogicalNot{
					param: jsonpathFilter,
				})
			} else {
				p.push(jsonpathFilter)
			}

		case ruleAction32:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: &syntaxCompareEQ{},
			})

		case ruleAction33:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxLogicalNot{
				param: &syntaxBasicCompareQuery{
					leftParam:  leftParam,
					rightParam: rightParam,
					comparator: &syntaxCompareEQ{},
				},
			})

		case ruleAction34:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: &syntaxCompareGE{},
			})

		case ruleAction35:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: &syntaxCompareGT{},
			})

		case ruleAction36:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: &syntaxCompareLE{},
			})

		case ruleAction37:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: &syntaxCompareLT{},
			})

		case ruleAction38:

			leftParam := p.pop().(*syntaxBasicCompareParameter)
			regex := regexp.MustCompile(text)
			p.push(&syntaxBasicCompareQuery{
				leftParam: leftParam,
				rightParam: &syntaxBasicCompareParameter{
					param:     &syntaxQueryParamLiteral{literal: `regex`},
					isLiteral: true,
				},
				comparator: &syntaxCompareRegex{
					regex: regex,
				},
			})

		case ruleAction39:

			p.push(&syntaxBasicCompareParameter{
				param:     &syntaxQueryParamLiteral{p.pop()},
				isLiteral: true,
			})

		case ruleAction40:

			p.push(&syntaxBasicCompareParameter{
				param:     &syntaxQueryParamLiteral{p.pop()},
				isLiteral: true,
			})

		case ruleAction41:

			isLiteral := p.pop().(bool)
			param := p.pop().(syntaxQueryParameter)
			if !p.hasErr() && param.isMultiValueParameter() {
				p.syntaxErr(begin, msgErrorInvalidSyntaxFilterValueGroup, buffer)
			}
			p.push(&syntaxBasicCompareParameter{
				param:     param,
				isLiteral: isLiteral,
			})

		case ruleAction42:

			node := p.pop().(syntaxNode)
			switch node.(type) {
			case *syntaxRootIdentifier:
				param := &syntaxQueryParamRoot{
					param:     node,
					resultPtr: &[]interface{}{},
				}
				p.updateResultPtr(param.param, &param.resultPtr)
				p.push(param)
				p.push(true)
			case *syntaxCurrentRootIdentifier:
				param := &syntaxQueryParamCurrentRoot{
					param:     node,
					resultPtr: &[]interface{}{},
				}
				p.updateResultPtr(param.param, &param.resultPtr)
				p.push(param)
				p.push(false)
			default:
				p.push(&syntaxQueryParamRoot{})
				p.push(true)
			}

		case ruleAction43:

			p.push(p.toFloat(text))

		case ruleAction44:

			p.push(true)

		case ruleAction45:

			p.push(false)

		case ruleAction46:

			p.push(p.unescape(text))

		case ruleAction47:
			// '
			p.push(p.unescape(text))

		case ruleAction48:

			p.push(nil)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *pegJSONPathParser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
		p.buffer = append(p.buffer, endSymbol)
	}

	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	var max token32
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		tree.Expand(tokenIndex)
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position, depth}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 expression <- <((jsonpath END Action0) / (jsonpath? <.*> END Action1))> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				{
					position2, tokenIndex2, depth2 := position, tokenIndex, depth
					if !_rules[rulejsonpath]() {
						goto l3
					}
					if !_rules[ruleEND]() {
						goto l3
					}
					if !_rules[ruleAction0]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex, depth = position2, tokenIndex2, depth2
					{
						position4, tokenIndex4, depth4 := position, tokenIndex, depth
						if !_rules[rulejsonpath]() {
							goto l4
						}
						goto l5
					l4:
						position, tokenIndex, depth = position4, tokenIndex4, depth4
					}
				l5:
					{
						position6 := position
						depth++
					l7:
						{
							position8, tokenIndex8, depth8 := position, tokenIndex, depth
							if !matchDot() {
								goto l8
							}
							goto l7
						l8:
							position, tokenIndex, depth = position8, tokenIndex8, depth8
						}
						depth--
						add(rulePegText, position6)
					}
					if !_rules[ruleEND]() {
						goto l0
					}
					if !_rules[ruleAction1]() {
						goto l0
					}
				}
			l2:
				depth--
				add(ruleexpression, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 END <- <!.> */
		func() bool {
			position9, tokenIndex9, depth9 := position, tokenIndex, depth
			{
				position10 := position
				depth++
				{
					position11, tokenIndex11, depth11 := position, tokenIndex, depth
					if !matchDot() {
						goto l11
					}
					goto l9
				l11:
					position, tokenIndex, depth = position11, tokenIndex11, depth11
				}
				depth--
				add(ruleEND, position10)
			}
			return true
		l9:
			position, tokenIndex, depth = position9, tokenIndex9, depth9
			return false
		},
		/* 2 jsonpath <- <(space rootNode (childNodes Action2)? space Action3)> */
		func() bool {
			position12, tokenIndex12, depth12 := position, tokenIndex, depth
			{
				position13 := position
				depth++
				if !_rules[rulespace]() {
					goto l12
				}
				if !_rules[rulerootNode]() {
					goto l12
				}
				{
					position14, tokenIndex14, depth14 := position, tokenIndex, depth
					if !_rules[rulechildNodes]() {
						goto l14
					}
					if !_rules[ruleAction2]() {
						goto l14
					}
					goto l15
				l14:
					position, tokenIndex, depth = position14, tokenIndex14, depth14
				}
			l15:
				if !_rules[rulespace]() {
					goto l12
				}
				if !_rules[ruleAction3]() {
					goto l12
				}
				depth--
				add(rulejsonpath, position13)
			}
			return true
		l12:
			position, tokenIndex, depth = position12, tokenIndex12, depth12
			return false
		},
		/* 3 rootNode <- <(rootIdentifier / (<currentRootIdentifier> Action4) / (<(bracketNode / dotChildIdentifier)> Action5))> */
		func() bool {
			position16, tokenIndex16, depth16 := position, tokenIndex, depth
			{
				position17 := position
				depth++
				{
					position18, tokenIndex18, depth18 := position, tokenIndex, depth
					if !_rules[rulerootIdentifier]() {
						goto l19
					}
					goto l18
				l19:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					{
						position21 := position
						depth++
						if !_rules[rulecurrentRootIdentifier]() {
							goto l20
						}
						depth--
						add(rulePegText, position21)
					}
					if !_rules[ruleAction4]() {
						goto l20
					}
					goto l18
				l20:
					position, tokenIndex, depth = position18, tokenIndex18, depth18
					{
						position22 := position
						depth++
						{
							position23, tokenIndex23, depth23 := position, tokenIndex, depth
							if !_rules[rulebracketNode]() {
								goto l24
							}
							goto l23
						l24:
							position, tokenIndex, depth = position23, tokenIndex23, depth23
							if !_rules[ruledotChildIdentifier]() {
								goto l16
							}
						}
					l23:
						depth--
						add(rulePegText, position22)
					}
					if !_rules[ruleAction5]() {
						goto l16
					}
				}
			l18:
				depth--
				add(rulerootNode, position17)
			}
			return true
		l16:
			position, tokenIndex, depth = position16, tokenIndex16, depth16
			return false
		},
		/* 4 childNodes <- <((('.' '.' (bracketNode / dotChildIdentifier) Action6) / (<('.' dotChildIdentifier)> Action7) / bracketNode) (childNodes Action8)?)> */
		func() bool {
			position25, tokenIndex25, depth25 := position, tokenIndex, depth
			{
				position26 := position
				depth++
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l28
					}
					position++
					if buffer[position] != rune('.') {
						goto l28
					}
					position++
					{
						position29, tokenIndex29, depth29 := position, tokenIndex, depth
						if !_rules[rulebracketNode]() {
							goto l30
						}
						goto l29
					l30:
						position, tokenIndex, depth = position29, tokenIndex29, depth29
						if !_rules[ruledotChildIdentifier]() {
							goto l28
						}
					}
				l29:
					if !_rules[ruleAction6]() {
						goto l28
					}
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					{
						position32 := position
						depth++
						if buffer[position] != rune('.') {
							goto l31
						}
						position++
						if !_rules[ruledotChildIdentifier]() {
							goto l31
						}
						depth--
						add(rulePegText, position32)
					}
					if !_rules[ruleAction7]() {
						goto l31
					}
					goto l27
				l31:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if !_rules[rulebracketNode]() {
						goto l25
					}
				}
			l27:
				{
					position33, tokenIndex33, depth33 := position, tokenIndex, depth
					if !_rules[rulechildNodes]() {
						goto l33
					}
					if !_rules[ruleAction8]() {
						goto l33
					}
					goto l34
				l33:
					position, tokenIndex, depth = position33, tokenIndex33, depth33
				}
			l34:
				depth--
				add(rulechildNodes, position26)
			}
			return true
		l25:
			position, tokenIndex, depth = position25, tokenIndex25, depth25
			return false
		},
		/* 5 bracketNode <- <(<(squareBracketStart (bracketChildIdentifier / qualifier) squareBracketEnd)> Action9)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				{
					position37 := position
					depth++
					if !_rules[rulesquareBracketStart]() {
						goto l35
					}
					{
						position38, tokenIndex38, depth38 := position, tokenIndex, depth
						if !_rules[rulebracketChildIdentifier]() {
							goto l39
						}
						goto l38
					l39:
						position, tokenIndex, depth = position38, tokenIndex38, depth38
						if !_rules[rulequalifier]() {
							goto l35
						}
					}
				l38:
					if !_rules[rulesquareBracketEnd]() {
						goto l35
					}
					depth--
					add(rulePegText, position37)
				}
				if !_rules[ruleAction9]() {
					goto l35
				}
				depth--
				add(rulebracketNode, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 6 rootIdentifier <- <('$' Action10)> */
		func() bool {
			position40, tokenIndex40, depth40 := position, tokenIndex, depth
			{
				position41 := position
				depth++
				if buffer[position] != rune('$') {
					goto l40
				}
				position++
				if !_rules[ruleAction10]() {
					goto l40
				}
				depth--
				add(rulerootIdentifier, position41)
			}
			return true
		l40:
			position, tokenIndex, depth = position40, tokenIndex40, depth40
			return false
		},
		/* 7 currentRootIdentifier <- <('@' Action11)> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				if buffer[position] != rune('@') {
					goto l42
				}
				position++
				if !_rules[ruleAction11]() {
					goto l42
				}
				depth--
				add(rulecurrentRootIdentifier, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 8 dotChildIdentifier <- <(<(('\\' '\\') / ('\\' ('.' / '[' / ')' / '=' / '!' / '>' / '<' / ' ' / '\t' / '\r' / '\n')) / (!('\\' / '.' / '[' / ')' / '=' / '!' / '>' / '<' / ' ' / '\t' / '\r' / '\n') .))+> Action12)> */
		func() bool {
			position44, tokenIndex44, depth44 := position, tokenIndex, depth
			{
				position45 := position
				depth++
				{
					position46 := position
					depth++
					{
						position49, tokenIndex49, depth49 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l50
						}
						position++
						if buffer[position] != rune('\\') {
							goto l50
						}
						position++
						goto l49
					l50:
						position, tokenIndex, depth = position49, tokenIndex49, depth49
						if buffer[position] != rune('\\') {
							goto l51
						}
						position++
						{
							position52, tokenIndex52, depth52 := position, tokenIndex, depth
							if buffer[position] != rune('.') {
								goto l53
							}
							position++
							goto l52
						l53:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('[') {
								goto l54
							}
							position++
							goto l52
						l54:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune(')') {
								goto l55
							}
							position++
							goto l52
						l55:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('=') {
								goto l56
							}
							position++
							goto l52
						l56:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('!') {
								goto l57
							}
							position++
							goto l52
						l57:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('>') {
								goto l58
							}
							position++
							goto l52
						l58:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('<') {
								goto l59
							}
							position++
							goto l52
						l59:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune(' ') {
								goto l60
							}
							position++
							goto l52
						l60:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('\t') {
								goto l61
							}
							position++
							goto l52
						l61:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('\r') {
								goto l62
							}
							position++
							goto l52
						l62:
							position, tokenIndex, depth = position52, tokenIndex52, depth52
							if buffer[position] != rune('\n') {
								goto l51
							}
							position++
						}
					l52:
						goto l49
					l51:
						position, tokenIndex, depth = position49, tokenIndex49, depth49
						{
							position63, tokenIndex63, depth63 := position, tokenIndex, depth
							{
								position64, tokenIndex64, depth64 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l65
								}
								position++
								goto l64
							l65:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('.') {
									goto l66
								}
								position++
								goto l64
							l66:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('[') {
									goto l67
								}
								position++
								goto l64
							l67:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune(')') {
									goto l68
								}
								position++
								goto l64
							l68:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('=') {
									goto l69
								}
								position++
								goto l64
							l69:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('!') {
									goto l70
								}
								position++
								goto l64
							l70:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('>') {
									goto l71
								}
								position++
								goto l64
							l71:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('<') {
									goto l72
								}
								position++
								goto l64
							l72:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune(' ') {
									goto l73
								}
								position++
								goto l64
							l73:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('\t') {
									goto l74
								}
								position++
								goto l64
							l74:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('\r') {
									goto l75
								}
								position++
								goto l64
							l75:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('\n') {
									goto l63
								}
								position++
							}
						l64:
							goto l44
						l63:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
						}
						if !matchDot() {
							goto l44
						}
					}
				l49:
				l47:
					{
						position48, tokenIndex48, depth48 := position, tokenIndex, depth
						{
							position76, tokenIndex76, depth76 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l77
							}
							position++
							if buffer[position] != rune('\\') {
								goto l77
							}
							position++
							goto l76
						l77:
							position, tokenIndex, depth = position76, tokenIndex76, depth76
							if buffer[position] != rune('\\') {
								goto l78
							}
							position++
							{
								position79, tokenIndex79, depth79 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l80
								}
								position++
								goto l79
							l80:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('[') {
									goto l81
								}
								position++
								goto l79
							l81:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune(')') {
									goto l82
								}
								position++
								goto l79
							l82:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('=') {
									goto l83
								}
								position++
								goto l79
							l83:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('!') {
									goto l84
								}
								position++
								goto l79
							l84:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('>') {
									goto l85
								}
								position++
								goto l79
							l85:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('<') {
									goto l86
								}
								position++
								goto l79
							l86:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune(' ') {
									goto l87
								}
								position++
								goto l79
							l87:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('\t') {
									goto l88
								}
								position++
								goto l79
							l88:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('\r') {
									goto l89
								}
								position++
								goto l79
							l89:
								position, tokenIndex, depth = position79, tokenIndex79, depth79
								if buffer[position] != rune('\n') {
									goto l78
								}
								position++
							}
						l79:
							goto l76
						l78:
							position, tokenIndex, depth = position76, tokenIndex76, depth76
							{
								position90, tokenIndex90, depth90 := position, tokenIndex, depth
								{
									position91, tokenIndex91, depth91 := position, tokenIndex, depth
									if buffer[position] != rune('\\') {
										goto l92
									}
									position++
									goto l91
								l92:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('.') {
										goto l93
									}
									position++
									goto l91
								l93:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('[') {
										goto l94
									}
									position++
									goto l91
								l94:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune(')') {
										goto l95
									}
									position++
									goto l91
								l95:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('=') {
										goto l96
									}
									position++
									goto l91
								l96:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('!') {
										goto l97
									}
									position++
									goto l91
								l97:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('>') {
										goto l98
									}
									position++
									goto l91
								l98:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('<') {
										goto l99
									}
									position++
									goto l91
								l99:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune(' ') {
										goto l100
									}
									position++
									goto l91
								l100:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('\t') {
										goto l101
									}
									position++
									goto l91
								l101:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('\r') {
										goto l102
									}
									position++
									goto l91
								l102:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('\n') {
										goto l90
									}
									position++
								}
							l91:
								goto l48
							l90:
								position, tokenIndex, depth = position90, tokenIndex90, depth90
							}
							if !matchDot() {
								goto l48
							}
						}
					l76:
						goto l47
					l48:
						position, tokenIndex, depth = position48, tokenIndex48, depth48
					}
					depth--
					add(rulePegText, position46)
				}
				if !_rules[ruleAction12]() {
					goto l44
				}
				depth--
				add(ruledotChildIdentifier, position45)
			}
			return true
		l44:
			position, tokenIndex, depth = position44, tokenIndex44, depth44
			return false
		},
		/* 9 bracketChildIdentifier <- <(bracketNodeIdentifiers Action13)> */
		func() bool {
			position103, tokenIndex103, depth103 := position, tokenIndex, depth
			{
				position104 := position
				depth++
				if !_rules[rulebracketNodeIdentifiers]() {
					goto l103
				}
				if !_rules[ruleAction13]() {
					goto l103
				}
				depth--
				add(rulebracketChildIdentifier, position104)
			}
			return true
		l103:
			position, tokenIndex, depth = position103, tokenIndex103, depth103
			return false
		},
		/* 10 bracketNodeIdentifiers <- <((singleQuotedNodeIdentifier / doubleQuotedNodeIdentifier) Action14 (sepBracketIdentifier bracketNodeIdentifiers Action15)?)> */
		func() bool {
			position105, tokenIndex105, depth105 := position, tokenIndex, depth
			{
				position106 := position
				depth++
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if !_rules[rulesingleQuotedNodeIdentifier]() {
						goto l108
					}
					goto l107
				l108:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
					if !_rules[ruledoubleQuotedNodeIdentifier]() {
						goto l105
					}
				}
			l107:
				if !_rules[ruleAction14]() {
					goto l105
				}
				{
					position109, tokenIndex109, depth109 := position, tokenIndex, depth
					if !_rules[rulesepBracketIdentifier]() {
						goto l109
					}
					if !_rules[rulebracketNodeIdentifiers]() {
						goto l109
					}
					if !_rules[ruleAction15]() {
						goto l109
					}
					goto l110
				l109:
					position, tokenIndex, depth = position109, tokenIndex109, depth109
				}
			l110:
				depth--
				add(rulebracketNodeIdentifiers, position106)
			}
			return true
		l105:
			position, tokenIndex, depth = position105, tokenIndex105, depth105
			return false
		},
		/* 11 singleQuotedNodeIdentifier <- <('\'' <(('\\' '\\') / ('\\' '\'') / (!('\\' / '\'') .))*> '\'' Action16)> */
		func() bool {
			position111, tokenIndex111, depth111 := position, tokenIndex, depth
			{
				position112 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l111
				}
				position++
				{
					position113 := position
					depth++
				l114:
					{
						position115, tokenIndex115, depth115 := position, tokenIndex, depth
						{
							position116, tokenIndex116, depth116 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l117
							}
							position++
							if buffer[position] != rune('\\') {
								goto l117
							}
							position++
							goto l116
						l117:
							position, tokenIndex, depth = position116, tokenIndex116, depth116
							if buffer[position] != rune('\\') {
								goto l118
							}
							position++
							if buffer[position] != rune('\'') {
								goto l118
							}
							position++
							goto l116
						l118:
							position, tokenIndex, depth = position116, tokenIndex116, depth116
							{
								position119, tokenIndex119, depth119 := position, tokenIndex, depth
								{
									position120, tokenIndex120, depth120 := position, tokenIndex, depth
									if buffer[position] != rune('\\') {
										goto l121
									}
									position++
									goto l120
								l121:
									position, tokenIndex, depth = position120, tokenIndex120, depth120
									if buffer[position] != rune('\'') {
										goto l119
									}
									position++
								}
							l120:
								goto l115
							l119:
								position, tokenIndex, depth = position119, tokenIndex119, depth119
							}
							if !matchDot() {
								goto l115
							}
						}
					l116:
						goto l114
					l115:
						position, tokenIndex, depth = position115, tokenIndex115, depth115
					}
					depth--
					add(rulePegText, position113)
				}
				if buffer[position] != rune('\'') {
					goto l111
				}
				position++
				if !_rules[ruleAction16]() {
					goto l111
				}
				depth--
				add(rulesingleQuotedNodeIdentifier, position112)
			}
			return true
		l111:
			position, tokenIndex, depth = position111, tokenIndex111, depth111
			return false
		},
		/* 12 doubleQuotedNodeIdentifier <- <('"' <(('\\' '\\') / ('\\' '"') / (!('\\' / '"') .))*> '"' Action17)> */
		func() bool {
			position122, tokenIndex122, depth122 := position, tokenIndex, depth
			{
				position123 := position
				depth++
				if buffer[position] != rune('"') {
					goto l122
				}
				position++
				{
					position124 := position
					depth++
				l125:
					{
						position126, tokenIndex126, depth126 := position, tokenIndex, depth
						{
							position127, tokenIndex127, depth127 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l128
							}
							position++
							if buffer[position] != rune('\\') {
								goto l128
							}
							position++
							goto l127
						l128:
							position, tokenIndex, depth = position127, tokenIndex127, depth127
							if buffer[position] != rune('\\') {
								goto l129
							}
							position++
							if buffer[position] != rune('"') {
								goto l129
							}
							position++
							goto l127
						l129:
							position, tokenIndex, depth = position127, tokenIndex127, depth127
							{
								position130, tokenIndex130, depth130 := position, tokenIndex, depth
								{
									position131, tokenIndex131, depth131 := position, tokenIndex, depth
									if buffer[position] != rune('\\') {
										goto l132
									}
									position++
									goto l131
								l132:
									position, tokenIndex, depth = position131, tokenIndex131, depth131
									if buffer[position] != rune('"') {
										goto l130
									}
									position++
								}
							l131:
								goto l126
							l130:
								position, tokenIndex, depth = position130, tokenIndex130, depth130
							}
							if !matchDot() {
								goto l126
							}
						}
					l127:
						goto l125
					l126:
						position, tokenIndex, depth = position126, tokenIndex126, depth126
					}
					depth--
					add(rulePegText, position124)
				}
				if buffer[position] != rune('"') {
					goto l122
				}
				position++
				if !_rules[ruleAction17]() {
					goto l122
				}
				depth--
				add(ruledoubleQuotedNodeIdentifier, position123)
			}
			return true
		l122:
			position, tokenIndex, depth = position122, tokenIndex122, depth122
			return false
		},
		/* 13 sepBracketIdentifier <- <(space ',' space)> */
		func() bool {
			position133, tokenIndex133, depth133 := position, tokenIndex, depth
			{
				position134 := position
				depth++
				if !_rules[rulespace]() {
					goto l133
				}
				if buffer[position] != rune(',') {
					goto l133
				}
				position++
				if !_rules[rulespace]() {
					goto l133
				}
				depth--
				add(rulesepBracketIdentifier, position134)
			}
			return true
		l133:
			position, tokenIndex, depth = position133, tokenIndex133, depth133
			return false
		},
		/* 14 qualifier <- <(union / script / filter)> */
		func() bool {
			position135, tokenIndex135, depth135 := position, tokenIndex, depth
			{
				position136 := position
				depth++
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if !_rules[ruleunion]() {
						goto l138
					}
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if !_rules[rulescript]() {
						goto l139
					}
					goto l137
				l139:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					if !_rules[rulefilter]() {
						goto l135
					}
				}
			l137:
				depth--
				add(rulequalifier, position136)
			}
			return true
		l135:
			position, tokenIndex, depth = position135, tokenIndex135, depth135
			return false
		},
		/* 15 union <- <(index Action18 (sepUnion union Action19)?)> */
		func() bool {
			position140, tokenIndex140, depth140 := position, tokenIndex, depth
			{
				position141 := position
				depth++
				if !_rules[ruleindex]() {
					goto l140
				}
				if !_rules[ruleAction18]() {
					goto l140
				}
				{
					position142, tokenIndex142, depth142 := position, tokenIndex, depth
					if !_rules[rulesepUnion]() {
						goto l142
					}
					if !_rules[ruleunion]() {
						goto l142
					}
					if !_rules[ruleAction19]() {
						goto l142
					}
					goto l143
				l142:
					position, tokenIndex, depth = position142, tokenIndex142, depth142
				}
			l143:
				depth--
				add(ruleunion, position141)
			}
			return true
		l140:
			position, tokenIndex, depth = position140, tokenIndex140, depth140
			return false
		},
		/* 16 index <- <((slice Action20) / (<indexNumber> Action21) / ('*' Action22))> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if !_rules[ruleslice]() {
						goto l147
					}
					if !_rules[ruleAction20]() {
						goto l147
					}
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					{
						position149 := position
						depth++
						if !_rules[ruleindexNumber]() {
							goto l148
						}
						depth--
						add(rulePegText, position149)
					}
					if !_rules[ruleAction21]() {
						goto l148
					}
					goto l146
				l148:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if buffer[position] != rune('*') {
						goto l144
					}
					position++
					if !_rules[ruleAction22]() {
						goto l144
					}
				}
			l146:
				depth--
				add(ruleindex, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 17 slice <- <(anyIndex sepSlice anyIndex ((sepSlice anyIndex) / (space Action23)))> */
		func() bool {
			position150, tokenIndex150, depth150 := position, tokenIndex, depth
			{
				position151 := position
				depth++
				if !_rules[ruleanyIndex]() {
					goto l150
				}
				if !_rules[rulesepSlice]() {
					goto l150
				}
				if !_rules[ruleanyIndex]() {
					goto l150
				}
				{
					position152, tokenIndex152, depth152 := position, tokenIndex, depth
					if !_rules[rulesepSlice]() {
						goto l153
					}
					if !_rules[ruleanyIndex]() {
						goto l153
					}
					goto l152
				l153:
					position, tokenIndex, depth = position152, tokenIndex152, depth152
					if !_rules[rulespace]() {
						goto l150
					}
					if !_rules[ruleAction23]() {
						goto l150
					}
				}
			l152:
				depth--
				add(ruleslice, position151)
			}
			return true
		l150:
			position, tokenIndex, depth = position150, tokenIndex150, depth150
			return false
		},
		/* 18 anyIndex <- <(<indexNumber?> Action24)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
				{
					position156 := position
					depth++
					{
						position157, tokenIndex157, depth157 := position, tokenIndex, depth
						if !_rules[ruleindexNumber]() {
							goto l157
						}
						goto l158
					l157:
						position, tokenIndex, depth = position157, tokenIndex157, depth157
					}
				l158:
					depth--
					add(rulePegText, position156)
				}
				if !_rules[ruleAction24]() {
					goto l154
				}
				depth--
				add(ruleanyIndex, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 19 indexNumber <- <(('-' / '+')? [0-9]+)> */
		func() bool {
			position159, tokenIndex159, depth159 := position, tokenIndex, depth
			{
				position160 := position
				depth++
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l164
						}
						position++
						goto l163
					l164:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
						if buffer[position] != rune('+') {
							goto l161
						}
						position++
					}
				l163:
					goto l162
				l161:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
				}
			l162:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l159
				}
				position++
			l165:
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l166
					}
					position++
					goto l165
				l166:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
				}
				depth--
				add(ruleindexNumber, position160)
			}
			return true
		l159:
			position, tokenIndex, depth = position159, tokenIndex159, depth159
			return false
		},
		/* 20 sepUnion <- <(space ',' space)> */
		func() bool {
			position167, tokenIndex167, depth167 := position, tokenIndex, depth
			{
				position168 := position
				depth++
				if !_rules[rulespace]() {
					goto l167
				}
				if buffer[position] != rune(',') {
					goto l167
				}
				position++
				if !_rules[rulespace]() {
					goto l167
				}
				depth--
				add(rulesepUnion, position168)
			}
			return true
		l167:
			position, tokenIndex, depth = position167, tokenIndex167, depth167
			return false
		},
		/* 21 sepSlice <- <(space ':' space)> */
		func() bool {
			position169, tokenIndex169, depth169 := position, tokenIndex, depth
			{
				position170 := position
				depth++
				if !_rules[rulespace]() {
					goto l169
				}
				if buffer[position] != rune(':') {
					goto l169
				}
				position++
				if !_rules[rulespace]() {
					goto l169
				}
				depth--
				add(rulesepSlice, position170)
			}
			return true
		l169:
			position, tokenIndex, depth = position169, tokenIndex169, depth169
			return false
		},
		/* 22 script <- <(scriptStart <command> scriptEnd Action25)> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				if !_rules[rulescriptStart]() {
					goto l171
				}
				{
					position173 := position
					depth++
					if !_rules[rulecommand]() {
						goto l171
					}
					depth--
					add(rulePegText, position173)
				}
				if !_rules[rulescriptEnd]() {
					goto l171
				}
				if !_rules[ruleAction25]() {
					goto l171
				}
				depth--
				add(rulescript, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 23 command <- <(!')' .)+> */
		func() bool {
			position174, tokenIndex174, depth174 := position, tokenIndex, depth
			{
				position175 := position
				depth++
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if buffer[position] != rune(')') {
						goto l178
					}
					position++
					goto l174
				l178:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
				}
				if !matchDot() {
					goto l174
				}
			l176:
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					{
						position179, tokenIndex179, depth179 := position, tokenIndex, depth
						if buffer[position] != rune(')') {
							goto l179
						}
						position++
						goto l177
					l179:
						position, tokenIndex, depth = position179, tokenIndex179, depth179
					}
					if !matchDot() {
						goto l177
					}
					goto l176
				l177:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
				}
				depth--
				add(rulecommand, position175)
			}
			return true
		l174:
			position, tokenIndex, depth = position174, tokenIndex174, depth174
			return false
		},
		/* 24 filter <- <(filterStart query filterEnd Action26)> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				if !_rules[rulefilterStart]() {
					goto l180
				}
				if !_rules[rulequery]() {
					goto l180
				}
				if !_rules[rulefilterEnd]() {
					goto l180
				}
				if !_rules[ruleAction26]() {
					goto l180
				}
				depth--
				add(rulefilter, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 25 query <- <(andQuery (logicOr query Action27)?)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				if !_rules[ruleandQuery]() {
					goto l182
				}
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					if !_rules[rulelogicOr]() {
						goto l184
					}
					if !_rules[rulequery]() {
						goto l184
					}
					if !_rules[ruleAction27]() {
						goto l184
					}
					goto l185
				l184:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
				}
			l185:
				depth--
				add(rulequery, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
			return false
		},
		/* 26 andQuery <- <((subQueryStart query subQueryEnd) / (basicQuery (logicAnd andQuery Action28)?))> */
		func() bool {
			position186, tokenIndex186, depth186 := position, tokenIndex, depth
			{
				position187 := position
				depth++
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					if !_rules[rulesubQueryStart]() {
						goto l189
					}
					if !_rules[rulequery]() {
						goto l189
					}
					if !_rules[rulesubQueryEnd]() {
						goto l189
					}
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					if !_rules[rulebasicQuery]() {
						goto l186
					}
					{
						position190, tokenIndex190, depth190 := position, tokenIndex, depth
						if !_rules[rulelogicAnd]() {
							goto l190
						}
						if !_rules[ruleandQuery]() {
							goto l190
						}
						if !_rules[ruleAction28]() {
							goto l190
						}
						goto l191
					l190:
						position, tokenIndex, depth = position190, tokenIndex190, depth190
					}
				l191:
				}
			l188:
				depth--
				add(ruleandQuery, position187)
			}
			return true
		l186:
			position, tokenIndex, depth = position186, tokenIndex186, depth186
			return false
		},
		/* 27 basicQuery <- <((<comparator> Action29) / (<logicNot?> Action30 jsonpathFilter Action31))> */
		func() bool {
			position192, tokenIndex192, depth192 := position, tokenIndex, depth
			{
				position193 := position
				depth++
				{
					position194, tokenIndex194, depth194 := position, tokenIndex, depth
					{
						position196 := position
						depth++
						if !_rules[rulecomparator]() {
							goto l195
						}
						depth--
						add(rulePegText, position196)
					}
					if !_rules[ruleAction29]() {
						goto l195
					}
					goto l194
				l195:
					position, tokenIndex, depth = position194, tokenIndex194, depth194
					{
						position197 := position
						depth++
						{
							position198, tokenIndex198, depth198 := position, tokenIndex, depth
							if !_rules[rulelogicNot]() {
								goto l198
							}
							goto l199
						l198:
							position, tokenIndex, depth = position198, tokenIndex198, depth198
						}
					l199:
						depth--
						add(rulePegText, position197)
					}
					if !_rules[ruleAction30]() {
						goto l192
					}
					if !_rules[rulejsonpathFilter]() {
						goto l192
					}
					if !_rules[ruleAction31]() {
						goto l192
					}
				}
			l194:
				depth--
				add(rulebasicQuery, position193)
			}
			return true
		l192:
			position, tokenIndex, depth = position192, tokenIndex192, depth192
			return false
		},
		/* 28 logicOr <- <(space ('|' '|') space)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				if !_rules[rulespace]() {
					goto l200
				}
				if buffer[position] != rune('|') {
					goto l200
				}
				position++
				if buffer[position] != rune('|') {
					goto l200
				}
				position++
				if !_rules[rulespace]() {
					goto l200
				}
				depth--
				add(rulelogicOr, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 29 logicAnd <- <(space ('&' '&') space)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				if !_rules[rulespace]() {
					goto l202
				}
				if buffer[position] != rune('&') {
					goto l202
				}
				position++
				if buffer[position] != rune('&') {
					goto l202
				}
				position++
				if !_rules[rulespace]() {
					goto l202
				}
				depth--
				add(rulelogicAnd, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 30 logicNot <- <('!' space)> */
		func() bool {
			position204, tokenIndex204, depth204 := position, tokenIndex, depth
			{
				position205 := position
				depth++
				if buffer[position] != rune('!') {
					goto l204
				}
				position++
				if !_rules[rulespace]() {
					goto l204
				}
				depth--
				add(rulelogicNot, position205)
			}
			return true
		l204:
			position, tokenIndex, depth = position204, tokenIndex204, depth204
			return false
		},
		/* 31 comparator <- <((qParam space (('=' '=' space qParam Action32) / ('!' '=' space qParam Action33))) / (qNumericParam space (('<' '=' space qNumericParam Action34) / ('<' space qNumericParam Action35) / ('>' '=' space qNumericParam Action36) / ('>' space qNumericParam Action37))) / (singleJsonpathFilter space ('=' '~') space '/' <regex> '/' Action38))> */
		func() bool {
			position206, tokenIndex206, depth206 := position, tokenIndex, depth
			{
				position207 := position
				depth++
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					if !_rules[ruleqParam]() {
						goto l209
					}
					if !_rules[rulespace]() {
						goto l209
					}
					{
						position210, tokenIndex210, depth210 := position, tokenIndex, depth
						if buffer[position] != rune('=') {
							goto l211
						}
						position++
						if buffer[position] != rune('=') {
							goto l211
						}
						position++
						if !_rules[rulespace]() {
							goto l211
						}
						if !_rules[ruleqParam]() {
							goto l211
						}
						if !_rules[ruleAction32]() {
							goto l211
						}
						goto l210
					l211:
						position, tokenIndex, depth = position210, tokenIndex210, depth210
						if buffer[position] != rune('!') {
							goto l209
						}
						position++
						if buffer[position] != rune('=') {
							goto l209
						}
						position++
						if !_rules[rulespace]() {
							goto l209
						}
						if !_rules[ruleqParam]() {
							goto l209
						}
						if !_rules[ruleAction33]() {
							goto l209
						}
					}
				l210:
					goto l208
				l209:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if !_rules[ruleqNumericParam]() {
						goto l212
					}
					if !_rules[rulespace]() {
						goto l212
					}
					{
						position213, tokenIndex213, depth213 := position, tokenIndex, depth
						if buffer[position] != rune('<') {
							goto l214
						}
						position++
						if buffer[position] != rune('=') {
							goto l214
						}
						position++
						if !_rules[rulespace]() {
							goto l214
						}
						if !_rules[ruleqNumericParam]() {
							goto l214
						}
						if !_rules[ruleAction34]() {
							goto l214
						}
						goto l213
					l214:
						position, tokenIndex, depth = position213, tokenIndex213, depth213
						if buffer[position] != rune('<') {
							goto l215
						}
						position++
						if !_rules[rulespace]() {
							goto l215
						}
						if !_rules[ruleqNumericParam]() {
							goto l215
						}
						if !_rules[ruleAction35]() {
							goto l215
						}
						goto l213
					l215:
						position, tokenIndex, depth = position213, tokenIndex213, depth213
						if buffer[position] != rune('>') {
							goto l216
						}
						position++
						if buffer[position] != rune('=') {
							goto l216
						}
						position++
						if !_rules[rulespace]() {
							goto l216
						}
						if !_rules[ruleqNumericParam]() {
							goto l216
						}
						if !_rules[ruleAction36]() {
							goto l216
						}
						goto l213
					l216:
						position, tokenIndex, depth = position213, tokenIndex213, depth213
						if buffer[position] != rune('>') {
							goto l212
						}
						position++
						if !_rules[rulespace]() {
							goto l212
						}
						if !_rules[ruleqNumericParam]() {
							goto l212
						}
						if !_rules[ruleAction37]() {
							goto l212
						}
					}
				l213:
					goto l208
				l212:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
					if !_rules[rulesingleJsonpathFilter]() {
						goto l206
					}
					if !_rules[rulespace]() {
						goto l206
					}
					if buffer[position] != rune('=') {
						goto l206
					}
					position++
					if buffer[position] != rune('~') {
						goto l206
					}
					position++
					if !_rules[rulespace]() {
						goto l206
					}
					if buffer[position] != rune('/') {
						goto l206
					}
					position++
					{
						position217 := position
						depth++
						if !_rules[ruleregex]() {
							goto l206
						}
						depth--
						add(rulePegText, position217)
					}
					if buffer[position] != rune('/') {
						goto l206
					}
					position++
					if !_rules[ruleAction38]() {
						goto l206
					}
				}
			l208:
				depth--
				add(rulecomparator, position207)
			}
			return true
		l206:
			position, tokenIndex, depth = position206, tokenIndex206, depth206
			return false
		},
		/* 32 qParam <- <((qLiteral Action39) / singleJsonpathFilter)> */
		func() bool {
			position218, tokenIndex218, depth218 := position, tokenIndex, depth
			{
				position219 := position
				depth++
				{
					position220, tokenIndex220, depth220 := position, tokenIndex, depth
					if !_rules[ruleqLiteral]() {
						goto l221
					}
					if !_rules[ruleAction39]() {
						goto l221
					}
					goto l220
				l221:
					position, tokenIndex, depth = position220, tokenIndex220, depth220
					if !_rules[rulesingleJsonpathFilter]() {
						goto l218
					}
				}
			l220:
				depth--
				add(ruleqParam, position219)
			}
			return true
		l218:
			position, tokenIndex, depth = position218, tokenIndex218, depth218
			return false
		},
		/* 33 qNumericParam <- <((lNumber Action40) / singleJsonpathFilter)> */
		func() bool {
			position222, tokenIndex222, depth222 := position, tokenIndex, depth
			{
				position223 := position
				depth++
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l225
					}
					if !_rules[ruleAction40]() {
						goto l225
					}
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					if !_rules[rulesingleJsonpathFilter]() {
						goto l222
					}
				}
			l224:
				depth--
				add(ruleqNumericParam, position223)
			}
			return true
		l222:
			position, tokenIndex, depth = position222, tokenIndex222, depth222
			return false
		},
		/* 34 qLiteral <- <(lNumber / lBool / lString / lNull)> */
		func() bool {
			position226, tokenIndex226, depth226 := position, tokenIndex, depth
			{
				position227 := position
				depth++
				{
					position228, tokenIndex228, depth228 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l229
					}
					goto l228
				l229:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if !_rules[rulelBool]() {
						goto l230
					}
					goto l228
				l230:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if !_rules[rulelString]() {
						goto l231
					}
					goto l228
				l231:
					position, tokenIndex, depth = position228, tokenIndex228, depth228
					if !_rules[rulelNull]() {
						goto l226
					}
				}
			l228:
				depth--
				add(ruleqLiteral, position227)
			}
			return true
		l226:
			position, tokenIndex, depth = position226, tokenIndex226, depth226
			return false
		},
		/* 35 singleJsonpathFilter <- <(jsonpathFilter Action41)> */
		func() bool {
			position232, tokenIndex232, depth232 := position, tokenIndex, depth
			{
				position233 := position
				depth++
				if !_rules[rulejsonpathFilter]() {
					goto l232
				}
				if !_rules[ruleAction41]() {
					goto l232
				}
				depth--
				add(rulesingleJsonpathFilter, position233)
			}
			return true
		l232:
			position, tokenIndex, depth = position232, tokenIndex232, depth232
			return false
		},
		/* 36 jsonpathFilter <- <(<jsonpath> Action42)> */
		func() bool {
			position234, tokenIndex234, depth234 := position, tokenIndex, depth
			{
				position235 := position
				depth++
				{
					position236 := position
					depth++
					if !_rules[rulejsonpath]() {
						goto l234
					}
					depth--
					add(rulePegText, position236)
				}
				if !_rules[ruleAction42]() {
					goto l234
				}
				depth--
				add(rulejsonpathFilter, position235)
			}
			return true
		l234:
			position, tokenIndex, depth = position234, tokenIndex234, depth234
			return false
		},
		/* 37 lNumber <- <(<(('-' / '+')? [0-9] ('-' / '+' / '.' / [0-9] / [a-z] / [A-Z])*)> Action43)> */
		func() bool {
			position237, tokenIndex237, depth237 := position, tokenIndex, depth
			{
				position238 := position
				depth++
				{
					position239 := position
					depth++
					{
						position240, tokenIndex240, depth240 := position, tokenIndex, depth
						{
							position242, tokenIndex242, depth242 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l243
							}
							position++
							goto l242
						l243:
							position, tokenIndex, depth = position242, tokenIndex242, depth242
							if buffer[position] != rune('+') {
								goto l240
							}
							position++
						}
					l242:
						goto l241
					l240:
						position, tokenIndex, depth = position240, tokenIndex240, depth240
					}
				l241:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l237
					}
					position++
				l244:
					{
						position245, tokenIndex245, depth245 := position, tokenIndex, depth
						{
							position246, tokenIndex246, depth246 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l247
							}
							position++
							goto l246
						l247:
							position, tokenIndex, depth = position246, tokenIndex246, depth246
							if buffer[position] != rune('+') {
								goto l248
							}
							position++
							goto l246
						l248:
							position, tokenIndex, depth = position246, tokenIndex246, depth246
							if buffer[position] != rune('.') {
								goto l249
							}
							position++
							goto l246
						l249:
							position, tokenIndex, depth = position246, tokenIndex246, depth246
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l250
							}
							position++
							goto l246
						l250:
							position, tokenIndex, depth = position246, tokenIndex246, depth246
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l251
							}
							position++
							goto l246
						l251:
							position, tokenIndex, depth = position246, tokenIndex246, depth246
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l245
							}
							position++
						}
					l246:
						goto l244
					l245:
						position, tokenIndex, depth = position245, tokenIndex245, depth245
					}
					depth--
					add(rulePegText, position239)
				}
				if !_rules[ruleAction43]() {
					goto l237
				}
				depth--
				add(rulelNumber, position238)
			}
			return true
		l237:
			position, tokenIndex, depth = position237, tokenIndex237, depth237
			return false
		},
		/* 38 lBool <- <(((('t' 'r' 'u' 'e') / ('T' 'r' 'u' 'e') / ('T' 'R' 'U' 'E')) Action44) / ((('f' 'a' 'l' 's' 'e') / ('F' 'a' 'l' 's' 'e') / ('F' 'A' 'L' 'S' 'E')) Action45))> */
		func() bool {
			position252, tokenIndex252, depth252 := position, tokenIndex, depth
			{
				position253 := position
				depth++
				{
					position254, tokenIndex254, depth254 := position, tokenIndex, depth
					{
						position256, tokenIndex256, depth256 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l257
						}
						position++
						if buffer[position] != rune('r') {
							goto l257
						}
						position++
						if buffer[position] != rune('u') {
							goto l257
						}
						position++
						if buffer[position] != rune('e') {
							goto l257
						}
						position++
						goto l256
					l257:
						position, tokenIndex, depth = position256, tokenIndex256, depth256
						if buffer[position] != rune('T') {
							goto l258
						}
						position++
						if buffer[position] != rune('r') {
							goto l258
						}
						position++
						if buffer[position] != rune('u') {
							goto l258
						}
						position++
						if buffer[position] != rune('e') {
							goto l258
						}
						position++
						goto l256
					l258:
						position, tokenIndex, depth = position256, tokenIndex256, depth256
						if buffer[position] != rune('T') {
							goto l255
						}
						position++
						if buffer[position] != rune('R') {
							goto l255
						}
						position++
						if buffer[position] != rune('U') {
							goto l255
						}
						position++
						if buffer[position] != rune('E') {
							goto l255
						}
						position++
					}
				l256:
					if !_rules[ruleAction44]() {
						goto l255
					}
					goto l254
				l255:
					position, tokenIndex, depth = position254, tokenIndex254, depth254
					{
						position259, tokenIndex259, depth259 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l260
						}
						position++
						if buffer[position] != rune('a') {
							goto l260
						}
						position++
						if buffer[position] != rune('l') {
							goto l260
						}
						position++
						if buffer[position] != rune('s') {
							goto l260
						}
						position++
						if buffer[position] != rune('e') {
							goto l260
						}
						position++
						goto l259
					l260:
						position, tokenIndex, depth = position259, tokenIndex259, depth259
						if buffer[position] != rune('F') {
							goto l261
						}
						position++
						if buffer[position] != rune('a') {
							goto l261
						}
						position++
						if buffer[position] != rune('l') {
							goto l261
						}
						position++
						if buffer[position] != rune('s') {
							goto l261
						}
						position++
						if buffer[position] != rune('e') {
							goto l261
						}
						position++
						goto l259
					l261:
						position, tokenIndex, depth = position259, tokenIndex259, depth259
						if buffer[position] != rune('F') {
							goto l252
						}
						position++
						if buffer[position] != rune('A') {
							goto l252
						}
						position++
						if buffer[position] != rune('L') {
							goto l252
						}
						position++
						if buffer[position] != rune('S') {
							goto l252
						}
						position++
						if buffer[position] != rune('E') {
							goto l252
						}
						position++
					}
				l259:
					if !_rules[ruleAction45]() {
						goto l252
					}
				}
			l254:
				depth--
				add(rulelBool, position253)
			}
			return true
		l252:
			position, tokenIndex, depth = position252, tokenIndex252, depth252
			return false
		},
		/* 39 lString <- <(('\'' <(('\\' '\\') / ('\\' '\'') / (!'\'' .))*> '\'' Action46) / ('"' <(('\\' '\\') / ('\\' '"') / (!'"' .))*> '"' Action47))> */
		func() bool {
			position262, tokenIndex262, depth262 := position, tokenIndex, depth
			{
				position263 := position
				depth++
				{
					position264, tokenIndex264, depth264 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l265
					}
					position++
					{
						position266 := position
						depth++
					l267:
						{
							position268, tokenIndex268, depth268 := position, tokenIndex, depth
							{
								position269, tokenIndex269, depth269 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l270
								}
								position++
								if buffer[position] != rune('\\') {
									goto l270
								}
								position++
								goto l269
							l270:
								position, tokenIndex, depth = position269, tokenIndex269, depth269
								if buffer[position] != rune('\\') {
									goto l271
								}
								position++
								if buffer[position] != rune('\'') {
									goto l271
								}
								position++
								goto l269
							l271:
								position, tokenIndex, depth = position269, tokenIndex269, depth269
								{
									position272, tokenIndex272, depth272 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l272
									}
									position++
									goto l268
								l272:
									position, tokenIndex, depth = position272, tokenIndex272, depth272
								}
								if !matchDot() {
									goto l268
								}
							}
						l269:
							goto l267
						l268:
							position, tokenIndex, depth = position268, tokenIndex268, depth268
						}
						depth--
						add(rulePegText, position266)
					}
					if buffer[position] != rune('\'') {
						goto l265
					}
					position++
					if !_rules[ruleAction46]() {
						goto l265
					}
					goto l264
				l265:
					position, tokenIndex, depth = position264, tokenIndex264, depth264
					if buffer[position] != rune('"') {
						goto l262
					}
					position++
					{
						position273 := position
						depth++
					l274:
						{
							position275, tokenIndex275, depth275 := position, tokenIndex, depth
							{
								position276, tokenIndex276, depth276 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l277
								}
								position++
								if buffer[position] != rune('\\') {
									goto l277
								}
								position++
								goto l276
							l277:
								position, tokenIndex, depth = position276, tokenIndex276, depth276
								if buffer[position] != rune('\\') {
									goto l278
								}
								position++
								if buffer[position] != rune('"') {
									goto l278
								}
								position++
								goto l276
							l278:
								position, tokenIndex, depth = position276, tokenIndex276, depth276
								{
									position279, tokenIndex279, depth279 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l279
									}
									position++
									goto l275
								l279:
									position, tokenIndex, depth = position279, tokenIndex279, depth279
								}
								if !matchDot() {
									goto l275
								}
							}
						l276:
							goto l274
						l275:
							position, tokenIndex, depth = position275, tokenIndex275, depth275
						}
						depth--
						add(rulePegText, position273)
					}
					if buffer[position] != rune('"') {
						goto l262
					}
					position++
					if !_rules[ruleAction47]() {
						goto l262
					}
				}
			l264:
				depth--
				add(rulelString, position263)
			}
			return true
		l262:
			position, tokenIndex, depth = position262, tokenIndex262, depth262
			return false
		},
		/* 40 lNull <- <((('n' 'u' 'l' 'l') / ('N' 'u' 'l' 'l') / ('N' 'U' 'L' 'L')) Action48)> */
		func() bool {
			position280, tokenIndex280, depth280 := position, tokenIndex, depth
			{
				position281 := position
				depth++
				{
					position282, tokenIndex282, depth282 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l283
					}
					position++
					if buffer[position] != rune('u') {
						goto l283
					}
					position++
					if buffer[position] != rune('l') {
						goto l283
					}
					position++
					if buffer[position] != rune('l') {
						goto l283
					}
					position++
					goto l282
				l283:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('N') {
						goto l284
					}
					position++
					if buffer[position] != rune('u') {
						goto l284
					}
					position++
					if buffer[position] != rune('l') {
						goto l284
					}
					position++
					if buffer[position] != rune('l') {
						goto l284
					}
					position++
					goto l282
				l284:
					position, tokenIndex, depth = position282, tokenIndex282, depth282
					if buffer[position] != rune('N') {
						goto l280
					}
					position++
					if buffer[position] != rune('U') {
						goto l280
					}
					position++
					if buffer[position] != rune('L') {
						goto l280
					}
					position++
					if buffer[position] != rune('L') {
						goto l280
					}
					position++
				}
			l282:
				if !_rules[ruleAction48]() {
					goto l280
				}
				depth--
				add(rulelNull, position281)
			}
			return true
		l280:
			position, tokenIndex, depth = position280, tokenIndex280, depth280
			return false
		},
		/* 41 regex <- <(('\\' '\\') / ('\\' '/') / (!'/' .))*> */
		func() bool {
			{
				position286 := position
				depth++
			l287:
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					{
						position289, tokenIndex289, depth289 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l290
						}
						position++
						if buffer[position] != rune('\\') {
							goto l290
						}
						position++
						goto l289
					l290:
						position, tokenIndex, depth = position289, tokenIndex289, depth289
						if buffer[position] != rune('\\') {
							goto l291
						}
						position++
						if buffer[position] != rune('/') {
							goto l291
						}
						position++
						goto l289
					l291:
						position, tokenIndex, depth = position289, tokenIndex289, depth289
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							if buffer[position] != rune('/') {
								goto l292
							}
							position++
							goto l288
						l292:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
						}
						if !matchDot() {
							goto l288
						}
					}
				l289:
					goto l287
				l288:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
				}
				depth--
				add(ruleregex, position286)
			}
			return true
		},
		/* 42 squareBracketStart <- <('[' space)> */
		func() bool {
			position293, tokenIndex293, depth293 := position, tokenIndex, depth
			{
				position294 := position
				depth++
				if buffer[position] != rune('[') {
					goto l293
				}
				position++
				if !_rules[rulespace]() {
					goto l293
				}
				depth--
				add(rulesquareBracketStart, position294)
			}
			return true
		l293:
			position, tokenIndex, depth = position293, tokenIndex293, depth293
			return false
		},
		/* 43 squareBracketEnd <- <(space ']')> */
		func() bool {
			position295, tokenIndex295, depth295 := position, tokenIndex, depth
			{
				position296 := position
				depth++
				if !_rules[rulespace]() {
					goto l295
				}
				if buffer[position] != rune(']') {
					goto l295
				}
				position++
				depth--
				add(rulesquareBracketEnd, position296)
			}
			return true
		l295:
			position, tokenIndex, depth = position295, tokenIndex295, depth295
			return false
		},
		/* 44 scriptStart <- <('(' space)> */
		func() bool {
			position297, tokenIndex297, depth297 := position, tokenIndex, depth
			{
				position298 := position
				depth++
				if buffer[position] != rune('(') {
					goto l297
				}
				position++
				if !_rules[rulespace]() {
					goto l297
				}
				depth--
				add(rulescriptStart, position298)
			}
			return true
		l297:
			position, tokenIndex, depth = position297, tokenIndex297, depth297
			return false
		},
		/* 45 scriptEnd <- <(space ')')> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				if !_rules[rulespace]() {
					goto l299
				}
				if buffer[position] != rune(')') {
					goto l299
				}
				position++
				depth--
				add(rulescriptEnd, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 46 filterStart <- <('?' '(' space)> */
		func() bool {
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				if buffer[position] != rune('?') {
					goto l301
				}
				position++
				if buffer[position] != rune('(') {
					goto l301
				}
				position++
				if !_rules[rulespace]() {
					goto l301
				}
				depth--
				add(rulefilterStart, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
			return false
		},
		/* 47 filterEnd <- <(space ')')> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				if !_rules[rulespace]() {
					goto l303
				}
				if buffer[position] != rune(')') {
					goto l303
				}
				position++
				depth--
				add(rulefilterEnd, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 48 subQueryStart <- <('(' space)> */
		func() bool {
			position305, tokenIndex305, depth305 := position, tokenIndex, depth
			{
				position306 := position
				depth++
				if buffer[position] != rune('(') {
					goto l305
				}
				position++
				if !_rules[rulespace]() {
					goto l305
				}
				depth--
				add(rulesubQueryStart, position306)
			}
			return true
		l305:
			position, tokenIndex, depth = position305, tokenIndex305, depth305
			return false
		},
		/* 49 subQueryEnd <- <(space ')')> */
		func() bool {
			position307, tokenIndex307, depth307 := position, tokenIndex, depth
			{
				position308 := position
				depth++
				if !_rules[rulespace]() {
					goto l307
				}
				if buffer[position] != rune(')') {
					goto l307
				}
				position++
				depth--
				add(rulesubQueryEnd, position308)
			}
			return true
		l307:
			position, tokenIndex, depth = position307, tokenIndex307, depth307
			return false
		},
		/* 50 space <- <(' ' / '\t')*> */
		func() bool {
			{
				position310 := position
				depth++
			l311:
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					{
						position313, tokenIndex313, depth313 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l314
						}
						position++
						goto l313
					l314:
						position, tokenIndex, depth = position313, tokenIndex313, depth313
						if buffer[position] != rune('\t') {
							goto l312
						}
						position++
					}
				l313:
					goto l311
				l312:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
				}
				depth--
				add(rulespace, position310)
			}
			return true
		},
		/* 52 Action0 <- <{
		    p.root = p.pop().(syntaxNode)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 54 Action1 <- <{
		    p.syntaxErr(begin, msgErrorInvalidSyntaxUnrecognizedInput, buffer)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 55 Action2 <- <{
		    child := p.pop().(syntaxNode)
		    root := p.pop().(syntaxNode)
		    root.setNext(child)
		    p.push(root)
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 56 Action3 <- <{
		        rootNode := p.pop().(syntaxNode)
		        checkNode := rootNode
		        for checkNode != nil {
					if checkNode.isMultiValue() {
		                rootNode.setMultiValue()
		                break
		            }
		            checkNode =  checkNode.getNext()
		        }
		        p.push(rootNode)
		    }> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 57 Action4 <- <{
		    if len(p.params) == 1 {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxUseBeginAtsign, buffer)
		    }
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 58 Action5 <- <{
		    if len(p.params) != 1 {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxOmitDollar, buffer)
		    }
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 59 Action6 <- <{
		    node := p.pop().(syntaxNode)
		    p.push(&syntaxRecursiveChildIdentifier{
		        syntaxBasicNode: &syntaxBasicNode{
		            text: `..`,
		            multiValue: true,
		            next: node,
		            result: &p.resultPtr,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 60 Action7 <- <{
		    identifier := p.pop().(syntaxNode)
		    identifier.setText(text)
		    p.push(identifier)
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 61 Action8 <- <{
		    child := p.pop().(syntaxNode)
		    parent := p.pop().(syntaxNode)
		    parent.setNext(child)
		    p.push(parent)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 62 Action9 <- <{
		    node := p.pop().(syntaxNode)
		    node.setText(text)
		    p.push(node)
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 63 Action10 <- <{
		    p.push(&syntaxRootIdentifier{
		        syntaxBasicNode: &syntaxBasicNode{
		            text: `$`,
		            result: &p.resultPtr,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 64 Action11 <- <{
		    p.push(&syntaxCurrentRootIdentifier{
		        syntaxBasicNode: &syntaxBasicNode{
		            text: `@`,
		            result: &p.resultPtr,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 65 Action12 <- <{
		    unescapedText := p.unescape(text)
		    if unescapedText == `*` {
		        p.push(&syntaxChildAsteriskIdentifier{
		            syntaxBasicNode: &syntaxBasicNode{
		                text: unescapedText,
		                multiValue: true,
		                result: &p.resultPtr,
		            },
		        })
		    } else {
		        p.push(&syntaxChildSingleIdentifier{
		            identifier: unescapedText,
		            syntaxBasicNode: &syntaxBasicNode{
		                text: unescapedText,
		                multiValue: false,
		                result: &p.resultPtr,
		            },
		        })
		    }
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 66 Action13 <- <{
		    identifier := p.pop().([]string)
		    if len(identifier) > 1 {
		        p.push(&syntaxChildMultiIdentifier{
		            identifiers: identifier,
		            syntaxBasicNode: &syntaxBasicNode{
		                multiValue: true,
		                result: &p.resultPtr,
		            },
		        })
		    } else {
		        p.push(&syntaxChildSingleIdentifier{
		            identifier: identifier[0],
		            syntaxBasicNode: &syntaxBasicNode{
		                multiValue: false,
		                result: &p.resultPtr,
		            },
		        })
		    }
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 67 Action14 <- <{
		    p.push([]string{p.pop().(string)})
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 68 Action15 <- <{
		    identifier2 := p.pop().([]string)
		    identifier1 := p.pop().([]string)
		    identifier1 = append(identifier1, identifier2...)
		    p.push(identifier1)
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 69 Action16 <- <{
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 70 Action17 <- <{ // '
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 71 Action18 <- <{
		    subscript := p.pop().(syntaxSubscript)
		    p.push(&syntaxUnionQualifier{
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: subscript.isMultiValue(),
		            result: &p.resultPtr,
		        },
		        subscripts: []syntaxSubscript{subscript},
		    })
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 72 Action19 <- <{
		    childIndexUnion := p.pop().(*syntaxUnionQualifier)
		    parentIndexUnion := p.pop().(*syntaxUnionQualifier)
		    parentIndexUnion.merge(childIndexUnion)
		    parentIndexUnion.setMultiValue()
		    p.push(parentIndexUnion)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 73 Action20 <- <{
		    step  := p.pop().(*syntaxIndex)
		    end   := p.pop().(*syntaxIndex)
		    start := p.pop().(*syntaxIndex)

		    if step.isOmitted || step.number == 0 {
		        step.number = 1
		    }

		    if step.number > 0 {
		        p.push(&syntaxSlicePositiveStep{
		            syntaxBasicSubscript: &syntaxBasicSubscript{
		                multiValue: true,
		            },
		            start: start,
		            end: end,
		            step: step,
		        })
		    } else {
		        p.push(&syntaxSliceNegativeStep{
		            syntaxBasicSubscript: &syntaxBasicSubscript{
		                multiValue: true,
		            },
		            start: start,
		            end: end,
		            step: step,
		        })
		    }
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 74 Action21 <- <{
		    p.push(&syntaxIndex{
		        syntaxBasicSubscript: &syntaxBasicSubscript{
		            multiValue: false,
		        },
		        number: p.toInt(text),
		    })
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 75 Action22 <- <{
		    p.push(&syntaxAsterisk{
		        syntaxBasicSubscript: &syntaxBasicSubscript{
		            multiValue: true,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 76 Action23 <- <{
		    p.push(&syntaxIndex{number: 1})
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 77 Action24 <- <{
		    if len(text) > 0 {
		        p.push(&syntaxIndex{number: p.toInt(text)})
		    } else {
		        p.push(&syntaxIndex{number: 0, isOmitted: true})
		    }
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 78 Action25 <- <{
		    p.push(&syntaxScriptQualifier{
		        command: text,
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: true,
		            result: &p.resultPtr,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 79 Action26 <- <{
		    query := p.pop().(syntaxQuery)
		    p.push(&syntaxFilterQualifier{
		        query: query,
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: true,
		            result: &p.resultPtr,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 80 Action27 <- <{
		    childQuery := p.pop().(syntaxQuery)
		    parentQuery := p.pop().(syntaxQuery)
		    p.push(&syntaxLogicalOr{parentQuery, childQuery})
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 81 Action28 <- <{
		    childQuery := p.pop().(syntaxQuery)
		    parentQuery := p.pop().(syntaxQuery)
		    p.push(&syntaxLogicalAnd{parentQuery, childQuery})
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 82 Action29 <- <{
		        if !p.hasErr() {
		            query := p.pop().(syntaxQuery)
		            p.push(query)

					if logicalNot, ok := query.(*syntaxLogicalNot); ok {
						query = (*logicalNot).param
					}
		            if checkQuery, ok := query.(*syntaxBasicCompareQuery); ok {
		                _, leftIsCurrentRoot := checkQuery.leftParam.param.(*syntaxQueryParamCurrentRoot)
		                _, rigthIsCurrentRoot := checkQuery.rightParam.param.(*syntaxQueryParamCurrentRoot)
		                if leftIsCurrentRoot && rigthIsCurrentRoot {
		                    p.syntaxErr(begin, msgErrorInvalidSyntaxTwoCurrentNode, buffer)
		                }
					}
		        }
		    }> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 83 Action30 <- <{
		    p.push(strings.HasPrefix(text, `!`))
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 84 Action31 <- <{
		    _ = p.pop().(bool)
		    jsonpathFilter := p.pop().(syntaxQuery)
		    isLogicalNot := p.pop().(bool)
		    if isLogicalNot {
		        p.push(&syntaxLogicalNot{
		            param: jsonpathFilter,
		        })
		    } else {
		        p.push(jsonpathFilter)
		    }
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 85 Action32 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.push(&syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: &syntaxCompareEQ{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 86 Action33 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.push(&syntaxLogicalNot{
		        param: &syntaxBasicCompareQuery{
		            leftParam: leftParam,
		            rightParam: rightParam,
		            comparator: &syntaxCompareEQ{},
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 87 Action34 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.push(&syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: &syntaxCompareGE{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 88 Action35 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.push(&syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: &syntaxCompareGT{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 89 Action36 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.push(&syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: &syntaxCompareLE{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 90 Action37 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.push(&syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: &syntaxCompareLT{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 91 Action38 <- <{
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    regex := regexp.MustCompile(text)
		    p.push(&syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: &syntaxBasicCompareParameter{
		            param: &syntaxQueryParamLiteral{literal: `regex`},
		            isLiteral: true,
		        },
		        comparator: &syntaxCompareRegex{
		            regex: regex,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 92 Action39 <- <{
		    p.push(&syntaxBasicCompareParameter{
		        param: &syntaxQueryParamLiteral{p.pop()},
		        isLiteral: true,
		    })
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 93 Action40 <- <{
		    p.push(&syntaxBasicCompareParameter{
		        param: &syntaxQueryParamLiteral{p.pop()},
		        isLiteral: true,
		    })
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 94 Action41 <- <{
		    isLiteral := p.pop().(bool)
		    param := p.pop().(syntaxQueryParameter)
		    if !p.hasErr() && param.isMultiValueParameter() {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxFilterValueGroup, buffer)
		    }
		    p.push(&syntaxBasicCompareParameter{
		        param: param,
		        isLiteral: isLiteral,
		    })
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 95 Action42 <- <{
		    node := p.pop().(syntaxNode)
		    switch node.(type) {
		    case *syntaxRootIdentifier:
		        param := &syntaxQueryParamRoot{
		            param: node,
		            resultPtr: &[]interface{}{},
		        }
		        p.updateResultPtr(param.param, &param.resultPtr)
		        p.push(param)
		        p.push(true)
		    case *syntaxCurrentRootIdentifier:
		        param := &syntaxQueryParamCurrentRoot{
		            param: node,
		            resultPtr: &[]interface{}{},
		        }
		        p.updateResultPtr(param.param, &param.resultPtr)
		        p.push(param)
		        p.push(false)
		    default:
		        p.push(&syntaxQueryParamRoot{})
		        p.push(true)
		    }
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 96 Action43 <- <{
		    p.push(p.toFloat(text))
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 97 Action44 <- <{
		    p.push(true)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 98 Action45 <- <{
		    p.push(false)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 99 Action46 <- <{
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 100 Action47 <- <{ // '
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 101 Action48 <- <{
		    p.push(nil)
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
	}
	p.rules = _rules
}
