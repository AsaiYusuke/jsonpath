package jsonpath

import (
	"strings"
	"regexp"
	"fmt"
	"math"
	"sort"
	"strconv"
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
	rulenodeFilter
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
	"nodeFilter",
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

type parser struct {
	jsonPathParser

	Buffer string
	buffer []rune
	rules  [100]func() bool
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
	p   *parser
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

func (p *parser) PrintSyntaxTree() {
	p.tokens32.PrintSyntaxTree(p.Buffer)
}

func (p *parser) Highlighter() {
	p.PrintSyntax()
}

func (p *parser) Execute() {
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
			root.setNext(&child)
			p.push(root)

		case ruleAction3:

			rootNode := p.pop().(syntaxNode)
			checkNode := rootNode
			for {
				if checkNode.isMultiValue() {
					rootNode.setMultiValue()
					break
				}
				next := checkNode.getNext()
				if next == nil {
					break
				}
				checkNode = *next
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
			p.push(syntaxRecursiveChildIdentifier{
				syntaxBasicNode: &syntaxBasicNode{
					text:       `..`,
					multiValue: true,
					next:       &node,
				},
			})

		case ruleAction7:

			identifier := p.pop().(syntaxNode)
			identifier.setText(text)
			p.push(identifier)

		case ruleAction8:

			child := p.pop().(syntaxNode)
			parent := p.pop().(syntaxNode)
			parent.setNext(&child)
			p.push(parent)

		case ruleAction9:

			node := p.pop().(syntaxNode)
			node.setText(text)
			p.push(node)

		case ruleAction10:

			p.push(syntaxRootIdentifier{
				syntaxBasicNode: &syntaxBasicNode{text: `$`},
			})

		case ruleAction11:

			p.push(syntaxCurrentRootIdentifier{
				syntaxBasicNode: &syntaxBasicNode{text: `@`},
			})

		case ruleAction12:

			unescapedText := p.unescape(text)
			if unescapedText == `*` {
				p.push(syntaxChildAsteriskIdentifier{
					syntaxBasicNode: &syntaxBasicNode{
						text:       unescapedText,
						multiValue: true,
					},
				})
			} else {
				p.push(syntaxChildSingleIdentifier{
					identifier: unescapedText,
					syntaxBasicNode: &syntaxBasicNode{
						text:       unescapedText,
						multiValue: false,
					},
				})
			}

		case ruleAction13:

			identifier := p.pop().([]string)
			if len(identifier) > 1 {
				p.push(syntaxChildMultiIdentifier{
					identifiers: identifier,
					syntaxBasicNode: &syntaxBasicNode{
						multiValue: true,
					},
				})
			} else {
				p.push(syntaxChildSingleIdentifier{
					identifier: identifier[0],
					syntaxBasicNode: &syntaxBasicNode{
						multiValue: false,
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
			union := syntaxUnion{
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: subscript.isMultiValue(),
				}}
			union.add(subscript)
			p.push(union)

		case ruleAction19:

			childIndexUnion := p.pop().(syntaxUnion)
			parentIndexUnion := p.pop().(syntaxUnion)
			parentIndexUnion.merge(childIndexUnion)
			parentIndexUnion.setMultiValue()
			p.push(parentIndexUnion)

		case ruleAction20:

			step := p.pop().(syntaxIndex)
			end := p.pop().(syntaxIndex)
			start := p.pop().(syntaxIndex)
			p.push(syntaxSlice{
				syntaxBasicSubscript: &syntaxBasicSubscript{
					multiValue: true,
				},
				start: start,
				end:   end,
				step:  step,
			})

		case ruleAction21:

			p.push(syntaxIndex{
				syntaxBasicSubscript: &syntaxBasicSubscript{
					multiValue: false,
				},
				number: p.toInt(text),
			})

		case ruleAction22:

			p.push(syntaxAsterisk{
				syntaxBasicSubscript: &syntaxBasicSubscript{
					multiValue: true,
				},
			})

		case ruleAction23:

			p.push(syntaxIndex{number: 1})

		case ruleAction24:

			if len(text) > 0 {
				p.push(syntaxIndex{number: p.toInt(text)})
			} else {
				p.push(syntaxIndex{number: 0, isOmitted: true})
			}

		case ruleAction25:

			p.push(syntaxScript{
				command: text,
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: true,
				},
			})

		case ruleAction26:

			p.push(syntaxFilter{
				query: p.pop().(syntaxQuery),
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: true,
				},
			})

		case ruleAction27:

			childQuery := p.pop().(syntaxQuery)
			parentQuery := p.pop().(syntaxQuery)
			p.push(syntaxLogicalOr{parentQuery, childQuery})

		case ruleAction28:

			childQuery := p.pop().(syntaxQuery)
			parentQuery := p.pop().(syntaxQuery)
			p.push(syntaxLogicalAnd{parentQuery, childQuery})

		case ruleAction29:

			if !p.hasErr() {
				query := p.pop().(syntaxQuery)

				var checkQuery syntaxBasicCompareQuery
				switch query.(type) {
				case syntaxBasicCompareQuery:
					checkQuery = query.(syntaxBasicCompareQuery)
				case syntaxLogicalNot:
					checkQuery = (query.(syntaxLogicalNot)).param.(syntaxBasicCompareQuery)
				}

				leftFilter, leftIsCurrentRoot := checkQuery.leftParam.(syntaxNodeFilter)
				rightFilter, rigthIsCurrentRoot := checkQuery.rightParam.(syntaxNodeFilter)
				if leftIsCurrentRoot && rigthIsCurrentRoot && leftFilter.isCurrentRoot() && rightFilter.isCurrentRoot() {
					p.syntaxErr(begin, msgErrorInvalidSyntaxTwoCurrentNode, buffer)
				}

				p.push(query)
			}

		case ruleAction30:

			p.push(strings.HasPrefix(text, `!`))

		case ruleAction31:

			nodeFilter := syntaxNodeFilter{p.pop().(syntaxNode)}
			isLogicalNot := p.pop().(bool)
			if isLogicalNot {
				p.push(syntaxLogicalNot{nodeFilter})
			} else {
				p.push(nodeFilter)
			}

		case ruleAction32:

			rightParam := p.pop().(syntaxQuery)
			leftParam := p.pop().(syntaxQuery)
			p.push(syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: syntaxCompareEQ{},
			})

		case ruleAction33:

			rightParam := p.pop().(syntaxQuery)
			leftParam := p.pop().(syntaxQuery)
			p.push(syntaxLogicalNot{syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: syntaxCompareEQ{},
			}})

		case ruleAction34:

			rightParam := p.pop().(syntaxQuery)
			leftParam := p.pop().(syntaxQuery)
			p.push(syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: syntaxCompareGE{},
			})

		case ruleAction35:

			rightParam := p.pop().(syntaxQuery)
			leftParam := p.pop().(syntaxQuery)
			p.push(syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: syntaxCompareGT{},
			})

		case ruleAction36:

			rightParam := p.pop().(syntaxQuery)
			leftParam := p.pop().(syntaxQuery)
			p.push(syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: syntaxCompareLE{},
			})

		case ruleAction37:

			rightParam := p.pop().(syntaxQuery)
			leftParam := p.pop().(syntaxQuery)
			p.push(syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: syntaxCompareLT{},
			})

		case ruleAction38:

			nodeFilter := syntaxNodeFilter{p.pop().(syntaxNode)}
			regex := regexp.MustCompile(text)
			p.push(syntaxBasicCompareQuery{
				leftParam:  nodeFilter,
				rightParam: syntaxCompareLiteral{literal: `regex`},
				comparator: syntaxCompareRegex{
					regex: regex,
				},
			})

		case ruleAction39:

			p.push(syntaxCompareLiteral{p.pop()})

		case ruleAction40:

			p.push(syntaxCompareLiteral{p.pop()})

		case ruleAction41:

			node := p.pop().(syntaxNode)
			p.push(syntaxNodeFilter{node})

			if !p.hasErr() && node.isMultiValue() {
				p.syntaxErr(begin, msgErrorInvalidSyntaxFilterMultiValuedNode, buffer)
			}

		case ruleAction42:

			p.push(p.toFloat(text))

		case ruleAction43:

			p.push(true)

		case ruleAction44:

			p.push(false)

		case ruleAction45:

			p.push(p.unescape(text))

		case ruleAction46:
			// '
			p.push(p.unescape(text))

		case ruleAction47:

			p.push(nil)

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *parser) Init() {
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
		/* 0 expression <- <((jsonpath END Action0) / (jsonpath? <.+> END Action1))> */
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
						if !matchDot() {
							goto l0
						}
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
		/* 8 dotChildIdentifier <- <(<(('\\' '\\') / ('\\' ('.' / '[' / ')' / '=' / '!' / '>' / '<' / ' ' / '\t' / '\r' / '\n')) / (!('.' / '[' / ')' / '=' / '!' / '>' / '<' / ' ' / '\t' / '\r' / '\n') .))+> Action12)> */
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
								if buffer[position] != rune('.') {
									goto l65
								}
								position++
								goto l64
							l65:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('[') {
									goto l66
								}
								position++
								goto l64
							l66:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune(')') {
									goto l67
								}
								position++
								goto l64
							l67:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('=') {
									goto l68
								}
								position++
								goto l64
							l68:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('!') {
									goto l69
								}
								position++
								goto l64
							l69:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('>') {
									goto l70
								}
								position++
								goto l64
							l70:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('<') {
									goto l71
								}
								position++
								goto l64
							l71:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune(' ') {
									goto l72
								}
								position++
								goto l64
							l72:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('\t') {
									goto l73
								}
								position++
								goto l64
							l73:
								position, tokenIndex, depth = position64, tokenIndex64, depth64
								if buffer[position] != rune('\r') {
									goto l74
								}
								position++
								goto l64
							l74:
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
							position75, tokenIndex75, depth75 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l76
							}
							position++
							if buffer[position] != rune('\\') {
								goto l76
							}
							position++
							goto l75
						l76:
							position, tokenIndex, depth = position75, tokenIndex75, depth75
							if buffer[position] != rune('\\') {
								goto l77
							}
							position++
							{
								position78, tokenIndex78, depth78 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l79
								}
								position++
								goto l78
							l79:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('[') {
									goto l80
								}
								position++
								goto l78
							l80:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune(')') {
									goto l81
								}
								position++
								goto l78
							l81:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('=') {
									goto l82
								}
								position++
								goto l78
							l82:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('!') {
									goto l83
								}
								position++
								goto l78
							l83:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('>') {
									goto l84
								}
								position++
								goto l78
							l84:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('<') {
									goto l85
								}
								position++
								goto l78
							l85:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune(' ') {
									goto l86
								}
								position++
								goto l78
							l86:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('\t') {
									goto l87
								}
								position++
								goto l78
							l87:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('\r') {
									goto l88
								}
								position++
								goto l78
							l88:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('\n') {
									goto l77
								}
								position++
							}
						l78:
							goto l75
						l77:
							position, tokenIndex, depth = position75, tokenIndex75, depth75
							{
								position89, tokenIndex89, depth89 := position, tokenIndex, depth
								{
									position90, tokenIndex90, depth90 := position, tokenIndex, depth
									if buffer[position] != rune('.') {
										goto l91
									}
									position++
									goto l90
								l91:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('[') {
										goto l92
									}
									position++
									goto l90
								l92:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune(')') {
										goto l93
									}
									position++
									goto l90
								l93:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('=') {
										goto l94
									}
									position++
									goto l90
								l94:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('!') {
										goto l95
									}
									position++
									goto l90
								l95:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('>') {
										goto l96
									}
									position++
									goto l90
								l96:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('<') {
										goto l97
									}
									position++
									goto l90
								l97:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune(' ') {
										goto l98
									}
									position++
									goto l90
								l98:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('\t') {
										goto l99
									}
									position++
									goto l90
								l99:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('\r') {
										goto l100
									}
									position++
									goto l90
								l100:
									position, tokenIndex, depth = position90, tokenIndex90, depth90
									if buffer[position] != rune('\n') {
										goto l89
									}
									position++
								}
							l90:
								goto l48
							l89:
								position, tokenIndex, depth = position89, tokenIndex89, depth89
							}
							if !matchDot() {
								goto l48
							}
						}
					l75:
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
			position101, tokenIndex101, depth101 := position, tokenIndex, depth
			{
				position102 := position
				depth++
				if !_rules[rulebracketNodeIdentifiers]() {
					goto l101
				}
				if !_rules[ruleAction13]() {
					goto l101
				}
				depth--
				add(rulebracketChildIdentifier, position102)
			}
			return true
		l101:
			position, tokenIndex, depth = position101, tokenIndex101, depth101
			return false
		},
		/* 10 bracketNodeIdentifiers <- <((singleQuotedNodeIdentifier / doubleQuotedNodeIdentifier) Action14 (sepBracketIdentifier bracketNodeIdentifiers Action15)?)> */
		func() bool {
			position103, tokenIndex103, depth103 := position, tokenIndex, depth
			{
				position104 := position
				depth++
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if !_rules[rulesingleQuotedNodeIdentifier]() {
						goto l106
					}
					goto l105
				l106:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
					if !_rules[ruledoubleQuotedNodeIdentifier]() {
						goto l103
					}
				}
			l105:
				if !_rules[ruleAction14]() {
					goto l103
				}
				{
					position107, tokenIndex107, depth107 := position, tokenIndex, depth
					if !_rules[rulesepBracketIdentifier]() {
						goto l107
					}
					if !_rules[rulebracketNodeIdentifiers]() {
						goto l107
					}
					if !_rules[ruleAction15]() {
						goto l107
					}
					goto l108
				l107:
					position, tokenIndex, depth = position107, tokenIndex107, depth107
				}
			l108:
				depth--
				add(rulebracketNodeIdentifiers, position104)
			}
			return true
		l103:
			position, tokenIndex, depth = position103, tokenIndex103, depth103
			return false
		},
		/* 11 singleQuotedNodeIdentifier <- <('\'' <(('\\' '\\') / ('\\' '\'') / (!'\'' .))*> '\'' Action16)> */
		func() bool {
			position109, tokenIndex109, depth109 := position, tokenIndex, depth
			{
				position110 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l109
				}
				position++
				{
					position111 := position
					depth++
				l112:
					{
						position113, tokenIndex113, depth113 := position, tokenIndex, depth
						{
							position114, tokenIndex114, depth114 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l115
							}
							position++
							if buffer[position] != rune('\\') {
								goto l115
							}
							position++
							goto l114
						l115:
							position, tokenIndex, depth = position114, tokenIndex114, depth114
							if buffer[position] != rune('\\') {
								goto l116
							}
							position++
							if buffer[position] != rune('\'') {
								goto l116
							}
							position++
							goto l114
						l116:
							position, tokenIndex, depth = position114, tokenIndex114, depth114
							{
								position117, tokenIndex117, depth117 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l117
								}
								position++
								goto l113
							l117:
								position, tokenIndex, depth = position117, tokenIndex117, depth117
							}
							if !matchDot() {
								goto l113
							}
						}
					l114:
						goto l112
					l113:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
					}
					depth--
					add(rulePegText, position111)
				}
				if buffer[position] != rune('\'') {
					goto l109
				}
				position++
				if !_rules[ruleAction16]() {
					goto l109
				}
				depth--
				add(rulesingleQuotedNodeIdentifier, position110)
			}
			return true
		l109:
			position, tokenIndex, depth = position109, tokenIndex109, depth109
			return false
		},
		/* 12 doubleQuotedNodeIdentifier <- <('"' <(('\\' '\\') / ('\\' '"') / (!'"' .))*> '"' Action17)> */
		func() bool {
			position118, tokenIndex118, depth118 := position, tokenIndex, depth
			{
				position119 := position
				depth++
				if buffer[position] != rune('"') {
					goto l118
				}
				position++
				{
					position120 := position
					depth++
				l121:
					{
						position122, tokenIndex122, depth122 := position, tokenIndex, depth
						{
							position123, tokenIndex123, depth123 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l124
							}
							position++
							if buffer[position] != rune('\\') {
								goto l124
							}
							position++
							goto l123
						l124:
							position, tokenIndex, depth = position123, tokenIndex123, depth123
							if buffer[position] != rune('\\') {
								goto l125
							}
							position++
							if buffer[position] != rune('"') {
								goto l125
							}
							position++
							goto l123
						l125:
							position, tokenIndex, depth = position123, tokenIndex123, depth123
							{
								position126, tokenIndex126, depth126 := position, tokenIndex, depth
								if buffer[position] != rune('"') {
									goto l126
								}
								position++
								goto l122
							l126:
								position, tokenIndex, depth = position126, tokenIndex126, depth126
							}
							if !matchDot() {
								goto l122
							}
						}
					l123:
						goto l121
					l122:
						position, tokenIndex, depth = position122, tokenIndex122, depth122
					}
					depth--
					add(rulePegText, position120)
				}
				if buffer[position] != rune('"') {
					goto l118
				}
				position++
				if !_rules[ruleAction17]() {
					goto l118
				}
				depth--
				add(ruledoubleQuotedNodeIdentifier, position119)
			}
			return true
		l118:
			position, tokenIndex, depth = position118, tokenIndex118, depth118
			return false
		},
		/* 13 sepBracketIdentifier <- <(space ',' space)> */
		func() bool {
			position127, tokenIndex127, depth127 := position, tokenIndex, depth
			{
				position128 := position
				depth++
				if !_rules[rulespace]() {
					goto l127
				}
				if buffer[position] != rune(',') {
					goto l127
				}
				position++
				if !_rules[rulespace]() {
					goto l127
				}
				depth--
				add(rulesepBracketIdentifier, position128)
			}
			return true
		l127:
			position, tokenIndex, depth = position127, tokenIndex127, depth127
			return false
		},
		/* 14 qualifier <- <(union / script / filter)> */
		func() bool {
			position129, tokenIndex129, depth129 := position, tokenIndex, depth
			{
				position130 := position
				depth++
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					if !_rules[ruleunion]() {
						goto l132
					}
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if !_rules[rulescript]() {
						goto l133
					}
					goto l131
				l133:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					if !_rules[rulefilter]() {
						goto l129
					}
				}
			l131:
				depth--
				add(rulequalifier, position130)
			}
			return true
		l129:
			position, tokenIndex, depth = position129, tokenIndex129, depth129
			return false
		},
		/* 15 union <- <(index Action18 (sepUnion union Action19)?)> */
		func() bool {
			position134, tokenIndex134, depth134 := position, tokenIndex, depth
			{
				position135 := position
				depth++
				if !_rules[ruleindex]() {
					goto l134
				}
				if !_rules[ruleAction18]() {
					goto l134
				}
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if !_rules[rulesepUnion]() {
						goto l136
					}
					if !_rules[ruleunion]() {
						goto l136
					}
					if !_rules[ruleAction19]() {
						goto l136
					}
					goto l137
				l136:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
				}
			l137:
				depth--
				add(ruleunion, position135)
			}
			return true
		l134:
			position, tokenIndex, depth = position134, tokenIndex134, depth134
			return false
		},
		/* 16 index <- <((slice Action20) / (<indexNumber> Action21) / ('*' Action22))> */
		func() bool {
			position138, tokenIndex138, depth138 := position, tokenIndex, depth
			{
				position139 := position
				depth++
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if !_rules[ruleslice]() {
						goto l141
					}
					if !_rules[ruleAction20]() {
						goto l141
					}
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					{
						position143 := position
						depth++
						if !_rules[ruleindexNumber]() {
							goto l142
						}
						depth--
						add(rulePegText, position143)
					}
					if !_rules[ruleAction21]() {
						goto l142
					}
					goto l140
				l142:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					if buffer[position] != rune('*') {
						goto l138
					}
					position++
					if !_rules[ruleAction22]() {
						goto l138
					}
				}
			l140:
				depth--
				add(ruleindex, position139)
			}
			return true
		l138:
			position, tokenIndex, depth = position138, tokenIndex138, depth138
			return false
		},
		/* 17 slice <- <(anyIndex sepSlice anyIndex ((sepSlice anyIndex) / (space Action23)))> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				if !_rules[ruleanyIndex]() {
					goto l144
				}
				if !_rules[rulesepSlice]() {
					goto l144
				}
				if !_rules[ruleanyIndex]() {
					goto l144
				}
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					if !_rules[rulesepSlice]() {
						goto l147
					}
					if !_rules[ruleanyIndex]() {
						goto l147
					}
					goto l146
				l147:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
					if !_rules[rulespace]() {
						goto l144
					}
					if !_rules[ruleAction23]() {
						goto l144
					}
				}
			l146:
				depth--
				add(ruleslice, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 18 anyIndex <- <(<indexNumber?> Action24)> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				{
					position150 := position
					depth++
					{
						position151, tokenIndex151, depth151 := position, tokenIndex, depth
						if !_rules[ruleindexNumber]() {
							goto l151
						}
						goto l152
					l151:
						position, tokenIndex, depth = position151, tokenIndex151, depth151
					}
				l152:
					depth--
					add(rulePegText, position150)
				}
				if !_rules[ruleAction24]() {
					goto l148
				}
				depth--
				add(ruleanyIndex, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 19 indexNumber <- <(('-' / '+')? [0-9]+)> */
		func() bool {
			position153, tokenIndex153, depth153 := position, tokenIndex, depth
			{
				position154 := position
				depth++
				{
					position155, tokenIndex155, depth155 := position, tokenIndex, depth
					{
						position157, tokenIndex157, depth157 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l158
						}
						position++
						goto l157
					l158:
						position, tokenIndex, depth = position157, tokenIndex157, depth157
						if buffer[position] != rune('+') {
							goto l155
						}
						position++
					}
				l157:
					goto l156
				l155:
					position, tokenIndex, depth = position155, tokenIndex155, depth155
				}
			l156:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l153
				}
				position++
			l159:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l160
					}
					position++
					goto l159
				l160:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
				}
				depth--
				add(ruleindexNumber, position154)
			}
			return true
		l153:
			position, tokenIndex, depth = position153, tokenIndex153, depth153
			return false
		},
		/* 20 sepUnion <- <(space ',' space)> */
		func() bool {
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				if !_rules[rulespace]() {
					goto l161
				}
				if buffer[position] != rune(',') {
					goto l161
				}
				position++
				if !_rules[rulespace]() {
					goto l161
				}
				depth--
				add(rulesepUnion, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 21 sepSlice <- <(space ':' space)> */
		func() bool {
			position163, tokenIndex163, depth163 := position, tokenIndex, depth
			{
				position164 := position
				depth++
				if !_rules[rulespace]() {
					goto l163
				}
				if buffer[position] != rune(':') {
					goto l163
				}
				position++
				if !_rules[rulespace]() {
					goto l163
				}
				depth--
				add(rulesepSlice, position164)
			}
			return true
		l163:
			position, tokenIndex, depth = position163, tokenIndex163, depth163
			return false
		},
		/* 22 script <- <(scriptStart <command> scriptEnd Action25)> */
		func() bool {
			position165, tokenIndex165, depth165 := position, tokenIndex, depth
			{
				position166 := position
				depth++
				if !_rules[rulescriptStart]() {
					goto l165
				}
				{
					position167 := position
					depth++
					if !_rules[rulecommand]() {
						goto l165
					}
					depth--
					add(rulePegText, position167)
				}
				if !_rules[rulescriptEnd]() {
					goto l165
				}
				if !_rules[ruleAction25]() {
					goto l165
				}
				depth--
				add(rulescript, position166)
			}
			return true
		l165:
			position, tokenIndex, depth = position165, tokenIndex165, depth165
			return false
		},
		/* 23 command <- <(!')' .)*> */
		func() bool {
			{
				position169 := position
				depth++
			l170:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					{
						position172, tokenIndex172, depth172 := position, tokenIndex, depth
						if buffer[position] != rune(')') {
							goto l172
						}
						position++
						goto l171
					l172:
						position, tokenIndex, depth = position172, tokenIndex172, depth172
					}
					if !matchDot() {
						goto l171
					}
					goto l170
				l171:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
				}
				depth--
				add(rulecommand, position169)
			}
			return true
		},
		/* 24 filter <- <(filterStart query filterEnd Action26)> */
		func() bool {
			position173, tokenIndex173, depth173 := position, tokenIndex, depth
			{
				position174 := position
				depth++
				if !_rules[rulefilterStart]() {
					goto l173
				}
				if !_rules[rulequery]() {
					goto l173
				}
				if !_rules[rulefilterEnd]() {
					goto l173
				}
				if !_rules[ruleAction26]() {
					goto l173
				}
				depth--
				add(rulefilter, position174)
			}
			return true
		l173:
			position, tokenIndex, depth = position173, tokenIndex173, depth173
			return false
		},
		/* 25 query <- <(andQuery (logicOr query Action27)?)> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				if !_rules[ruleandQuery]() {
					goto l175
				}
				{
					position177, tokenIndex177, depth177 := position, tokenIndex, depth
					if !_rules[rulelogicOr]() {
						goto l177
					}
					if !_rules[rulequery]() {
						goto l177
					}
					if !_rules[ruleAction27]() {
						goto l177
					}
					goto l178
				l177:
					position, tokenIndex, depth = position177, tokenIndex177, depth177
				}
			l178:
				depth--
				add(rulequery, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
		/* 26 andQuery <- <((subQueryStart query subQueryEnd) / (basicQuery (logicAnd andQuery Action28)?))> */
		func() bool {
			position179, tokenIndex179, depth179 := position, tokenIndex, depth
			{
				position180 := position
				depth++
				{
					position181, tokenIndex181, depth181 := position, tokenIndex, depth
					if !_rules[rulesubQueryStart]() {
						goto l182
					}
					if !_rules[rulequery]() {
						goto l182
					}
					if !_rules[rulesubQueryEnd]() {
						goto l182
					}
					goto l181
				l182:
					position, tokenIndex, depth = position181, tokenIndex181, depth181
					if !_rules[rulebasicQuery]() {
						goto l179
					}
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if !_rules[rulelogicAnd]() {
							goto l183
						}
						if !_rules[ruleandQuery]() {
							goto l183
						}
						if !_rules[ruleAction28]() {
							goto l183
						}
						goto l184
					l183:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
					}
				l184:
				}
			l181:
				depth--
				add(ruleandQuery, position180)
			}
			return true
		l179:
			position, tokenIndex, depth = position179, tokenIndex179, depth179
			return false
		},
		/* 27 basicQuery <- <((<comparator> Action29) / (<logicNot?> Action30 jsonpath Action31))> */
		func() bool {
			position185, tokenIndex185, depth185 := position, tokenIndex, depth
			{
				position186 := position
				depth++
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					{
						position189 := position
						depth++
						if !_rules[rulecomparator]() {
							goto l188
						}
						depth--
						add(rulePegText, position189)
					}
					if !_rules[ruleAction29]() {
						goto l188
					}
					goto l187
				l188:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
					{
						position190 := position
						depth++
						{
							position191, tokenIndex191, depth191 := position, tokenIndex, depth
							if !_rules[rulelogicNot]() {
								goto l191
							}
							goto l192
						l191:
							position, tokenIndex, depth = position191, tokenIndex191, depth191
						}
					l192:
						depth--
						add(rulePegText, position190)
					}
					if !_rules[ruleAction30]() {
						goto l185
					}
					if !_rules[rulejsonpath]() {
						goto l185
					}
					if !_rules[ruleAction31]() {
						goto l185
					}
				}
			l187:
				depth--
				add(rulebasicQuery, position186)
			}
			return true
		l185:
			position, tokenIndex, depth = position185, tokenIndex185, depth185
			return false
		},
		/* 28 logicOr <- <(space ('|' '|') space)> */
		func() bool {
			position193, tokenIndex193, depth193 := position, tokenIndex, depth
			{
				position194 := position
				depth++
				if !_rules[rulespace]() {
					goto l193
				}
				if buffer[position] != rune('|') {
					goto l193
				}
				position++
				if buffer[position] != rune('|') {
					goto l193
				}
				position++
				if !_rules[rulespace]() {
					goto l193
				}
				depth--
				add(rulelogicOr, position194)
			}
			return true
		l193:
			position, tokenIndex, depth = position193, tokenIndex193, depth193
			return false
		},
		/* 29 logicAnd <- <(space ('&' '&') space)> */
		func() bool {
			position195, tokenIndex195, depth195 := position, tokenIndex, depth
			{
				position196 := position
				depth++
				if !_rules[rulespace]() {
					goto l195
				}
				if buffer[position] != rune('&') {
					goto l195
				}
				position++
				if buffer[position] != rune('&') {
					goto l195
				}
				position++
				if !_rules[rulespace]() {
					goto l195
				}
				depth--
				add(rulelogicAnd, position196)
			}
			return true
		l195:
			position, tokenIndex, depth = position195, tokenIndex195, depth195
			return false
		},
		/* 30 logicNot <- <('!' space)> */
		func() bool {
			position197, tokenIndex197, depth197 := position, tokenIndex, depth
			{
				position198 := position
				depth++
				if buffer[position] != rune('!') {
					goto l197
				}
				position++
				if !_rules[rulespace]() {
					goto l197
				}
				depth--
				add(rulelogicNot, position198)
			}
			return true
		l197:
			position, tokenIndex, depth = position197, tokenIndex197, depth197
			return false
		},
		/* 31 comparator <- <((qParam space (('=' '=' space qParam Action32) / ('!' '=' space qParam Action33))) / (qNumericParam space (('<' '=' space qNumericParam Action34) / ('<' space qNumericParam Action35) / ('>' '=' space qNumericParam Action36) / ('>' space qNumericParam Action37))) / (jsonpath space ('=' '~') space '/' <regex> '/' Action38))> */
		func() bool {
			position199, tokenIndex199, depth199 := position, tokenIndex, depth
			{
				position200 := position
				depth++
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					if !_rules[ruleqParam]() {
						goto l202
					}
					if !_rules[rulespace]() {
						goto l202
					}
					{
						position203, tokenIndex203, depth203 := position, tokenIndex, depth
						if buffer[position] != rune('=') {
							goto l204
						}
						position++
						if buffer[position] != rune('=') {
							goto l204
						}
						position++
						if !_rules[rulespace]() {
							goto l204
						}
						if !_rules[ruleqParam]() {
							goto l204
						}
						if !_rules[ruleAction32]() {
							goto l204
						}
						goto l203
					l204:
						position, tokenIndex, depth = position203, tokenIndex203, depth203
						if buffer[position] != rune('!') {
							goto l202
						}
						position++
						if buffer[position] != rune('=') {
							goto l202
						}
						position++
						if !_rules[rulespace]() {
							goto l202
						}
						if !_rules[ruleqParam]() {
							goto l202
						}
						if !_rules[ruleAction33]() {
							goto l202
						}
					}
				l203:
					goto l201
				l202:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if !_rules[ruleqNumericParam]() {
						goto l205
					}
					if !_rules[rulespace]() {
						goto l205
					}
					{
						position206, tokenIndex206, depth206 := position, tokenIndex, depth
						if buffer[position] != rune('<') {
							goto l207
						}
						position++
						if buffer[position] != rune('=') {
							goto l207
						}
						position++
						if !_rules[rulespace]() {
							goto l207
						}
						if !_rules[ruleqNumericParam]() {
							goto l207
						}
						if !_rules[ruleAction34]() {
							goto l207
						}
						goto l206
					l207:
						position, tokenIndex, depth = position206, tokenIndex206, depth206
						if buffer[position] != rune('<') {
							goto l208
						}
						position++
						if !_rules[rulespace]() {
							goto l208
						}
						if !_rules[ruleqNumericParam]() {
							goto l208
						}
						if !_rules[ruleAction35]() {
							goto l208
						}
						goto l206
					l208:
						position, tokenIndex, depth = position206, tokenIndex206, depth206
						if buffer[position] != rune('>') {
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
						if !_rules[ruleqNumericParam]() {
							goto l209
						}
						if !_rules[ruleAction36]() {
							goto l209
						}
						goto l206
					l209:
						position, tokenIndex, depth = position206, tokenIndex206, depth206
						if buffer[position] != rune('>') {
							goto l205
						}
						position++
						if !_rules[rulespace]() {
							goto l205
						}
						if !_rules[ruleqNumericParam]() {
							goto l205
						}
						if !_rules[ruleAction37]() {
							goto l205
						}
					}
				l206:
					goto l201
				l205:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
					if !_rules[rulejsonpath]() {
						goto l199
					}
					if !_rules[rulespace]() {
						goto l199
					}
					if buffer[position] != rune('=') {
						goto l199
					}
					position++
					if buffer[position] != rune('~') {
						goto l199
					}
					position++
					if !_rules[rulespace]() {
						goto l199
					}
					if buffer[position] != rune('/') {
						goto l199
					}
					position++
					{
						position210 := position
						depth++
						if !_rules[ruleregex]() {
							goto l199
						}
						depth--
						add(rulePegText, position210)
					}
					if buffer[position] != rune('/') {
						goto l199
					}
					position++
					if !_rules[ruleAction38]() {
						goto l199
					}
				}
			l201:
				depth--
				add(rulecomparator, position200)
			}
			return true
		l199:
			position, tokenIndex, depth = position199, tokenIndex199, depth199
			return false
		},
		/* 32 qParam <- <((qLiteral Action39) / nodeFilter)> */
		func() bool {
			position211, tokenIndex211, depth211 := position, tokenIndex, depth
			{
				position212 := position
				depth++
				{
					position213, tokenIndex213, depth213 := position, tokenIndex, depth
					if !_rules[ruleqLiteral]() {
						goto l214
					}
					if !_rules[ruleAction39]() {
						goto l214
					}
					goto l213
				l214:
					position, tokenIndex, depth = position213, tokenIndex213, depth213
					if !_rules[rulenodeFilter]() {
						goto l211
					}
				}
			l213:
				depth--
				add(ruleqParam, position212)
			}
			return true
		l211:
			position, tokenIndex, depth = position211, tokenIndex211, depth211
			return false
		},
		/* 33 qNumericParam <- <((lNumber Action40) / nodeFilter)> */
		func() bool {
			position215, tokenIndex215, depth215 := position, tokenIndex, depth
			{
				position216 := position
				depth++
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l218
					}
					if !_rules[ruleAction40]() {
						goto l218
					}
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if !_rules[rulenodeFilter]() {
						goto l215
					}
				}
			l217:
				depth--
				add(ruleqNumericParam, position216)
			}
			return true
		l215:
			position, tokenIndex, depth = position215, tokenIndex215, depth215
			return false
		},
		/* 34 qLiteral <- <(lNumber / lBool / lString / lNull)> */
		func() bool {
			position219, tokenIndex219, depth219 := position, tokenIndex, depth
			{
				position220 := position
				depth++
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l222
					}
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if !_rules[rulelBool]() {
						goto l223
					}
					goto l221
				l223:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if !_rules[rulelString]() {
						goto l224
					}
					goto l221
				l224:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if !_rules[rulelNull]() {
						goto l219
					}
				}
			l221:
				depth--
				add(ruleqLiteral, position220)
			}
			return true
		l219:
			position, tokenIndex, depth = position219, tokenIndex219, depth219
			return false
		},
		/* 35 nodeFilter <- <(<jsonpath> Action41)> */
		func() bool {
			position225, tokenIndex225, depth225 := position, tokenIndex, depth
			{
				position226 := position
				depth++
				{
					position227 := position
					depth++
					if !_rules[rulejsonpath]() {
						goto l225
					}
					depth--
					add(rulePegText, position227)
				}
				if !_rules[ruleAction41]() {
					goto l225
				}
				depth--
				add(rulenodeFilter, position226)
			}
			return true
		l225:
			position, tokenIndex, depth = position225, tokenIndex225, depth225
			return false
		},
		/* 36 lNumber <- <(<(('-' / '+')? [0-9] ('-' / '+' / '.' / [0-9] / [a-z] / [A-Z])*)> Action42)> */
		func() bool {
			position228, tokenIndex228, depth228 := position, tokenIndex, depth
			{
				position229 := position
				depth++
				{
					position230 := position
					depth++
					{
						position231, tokenIndex231, depth231 := position, tokenIndex, depth
						{
							position233, tokenIndex233, depth233 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l234
							}
							position++
							goto l233
						l234:
							position, tokenIndex, depth = position233, tokenIndex233, depth233
							if buffer[position] != rune('+') {
								goto l231
							}
							position++
						}
					l233:
						goto l232
					l231:
						position, tokenIndex, depth = position231, tokenIndex231, depth231
					}
				l232:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l228
					}
					position++
				l235:
					{
						position236, tokenIndex236, depth236 := position, tokenIndex, depth
						{
							position237, tokenIndex237, depth237 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l238
							}
							position++
							goto l237
						l238:
							position, tokenIndex, depth = position237, tokenIndex237, depth237
							if buffer[position] != rune('+') {
								goto l239
							}
							position++
							goto l237
						l239:
							position, tokenIndex, depth = position237, tokenIndex237, depth237
							if buffer[position] != rune('.') {
								goto l240
							}
							position++
							goto l237
						l240:
							position, tokenIndex, depth = position237, tokenIndex237, depth237
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l241
							}
							position++
							goto l237
						l241:
							position, tokenIndex, depth = position237, tokenIndex237, depth237
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l242
							}
							position++
							goto l237
						l242:
							position, tokenIndex, depth = position237, tokenIndex237, depth237
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l236
							}
							position++
						}
					l237:
						goto l235
					l236:
						position, tokenIndex, depth = position236, tokenIndex236, depth236
					}
					depth--
					add(rulePegText, position230)
				}
				if !_rules[ruleAction42]() {
					goto l228
				}
				depth--
				add(rulelNumber, position229)
			}
			return true
		l228:
			position, tokenIndex, depth = position228, tokenIndex228, depth228
			return false
		},
		/* 37 lBool <- <(((('t' 'r' 'u' 'e') / ('T' 'r' 'u' 'e') / ('T' 'R' 'U' 'E')) Action43) / ((('f' 'a' 'l' 's' 'e') / ('F' 'a' 'l' 's' 'e') / ('F' 'A' 'L' 'S' 'E')) Action44))> */
		func() bool {
			position243, tokenIndex243, depth243 := position, tokenIndex, depth
			{
				position244 := position
				depth++
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					{
						position247, tokenIndex247, depth247 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l248
						}
						position++
						if buffer[position] != rune('r') {
							goto l248
						}
						position++
						if buffer[position] != rune('u') {
							goto l248
						}
						position++
						if buffer[position] != rune('e') {
							goto l248
						}
						position++
						goto l247
					l248:
						position, tokenIndex, depth = position247, tokenIndex247, depth247
						if buffer[position] != rune('T') {
							goto l249
						}
						position++
						if buffer[position] != rune('r') {
							goto l249
						}
						position++
						if buffer[position] != rune('u') {
							goto l249
						}
						position++
						if buffer[position] != rune('e') {
							goto l249
						}
						position++
						goto l247
					l249:
						position, tokenIndex, depth = position247, tokenIndex247, depth247
						if buffer[position] != rune('T') {
							goto l246
						}
						position++
						if buffer[position] != rune('R') {
							goto l246
						}
						position++
						if buffer[position] != rune('U') {
							goto l246
						}
						position++
						if buffer[position] != rune('E') {
							goto l246
						}
						position++
					}
				l247:
					if !_rules[ruleAction43]() {
						goto l246
					}
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					{
						position250, tokenIndex250, depth250 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l251
						}
						position++
						if buffer[position] != rune('a') {
							goto l251
						}
						position++
						if buffer[position] != rune('l') {
							goto l251
						}
						position++
						if buffer[position] != rune('s') {
							goto l251
						}
						position++
						if buffer[position] != rune('e') {
							goto l251
						}
						position++
						goto l250
					l251:
						position, tokenIndex, depth = position250, tokenIndex250, depth250
						if buffer[position] != rune('F') {
							goto l252
						}
						position++
						if buffer[position] != rune('a') {
							goto l252
						}
						position++
						if buffer[position] != rune('l') {
							goto l252
						}
						position++
						if buffer[position] != rune('s') {
							goto l252
						}
						position++
						if buffer[position] != rune('e') {
							goto l252
						}
						position++
						goto l250
					l252:
						position, tokenIndex, depth = position250, tokenIndex250, depth250
						if buffer[position] != rune('F') {
							goto l243
						}
						position++
						if buffer[position] != rune('A') {
							goto l243
						}
						position++
						if buffer[position] != rune('L') {
							goto l243
						}
						position++
						if buffer[position] != rune('S') {
							goto l243
						}
						position++
						if buffer[position] != rune('E') {
							goto l243
						}
						position++
					}
				l250:
					if !_rules[ruleAction44]() {
						goto l243
					}
				}
			l245:
				depth--
				add(rulelBool, position244)
			}
			return true
		l243:
			position, tokenIndex, depth = position243, tokenIndex243, depth243
			return false
		},
		/* 38 lString <- <(('\'' <(('\\' '\\') / ('\\' '\'') / (!'\'' .))*> '\'' Action45) / ('"' <(('\\' '\\') / ('\\' '"') / (!'"' .))*> '"' Action46))> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l256
					}
					position++
					{
						position257 := position
						depth++
					l258:
						{
							position259, tokenIndex259, depth259 := position, tokenIndex, depth
							{
								position260, tokenIndex260, depth260 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l261
								}
								position++
								if buffer[position] != rune('\\') {
									goto l261
								}
								position++
								goto l260
							l261:
								position, tokenIndex, depth = position260, tokenIndex260, depth260
								if buffer[position] != rune('\\') {
									goto l262
								}
								position++
								if buffer[position] != rune('\'') {
									goto l262
								}
								position++
								goto l260
							l262:
								position, tokenIndex, depth = position260, tokenIndex260, depth260
								{
									position263, tokenIndex263, depth263 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l263
									}
									position++
									goto l259
								l263:
									position, tokenIndex, depth = position263, tokenIndex263, depth263
								}
								if !matchDot() {
									goto l259
								}
							}
						l260:
							goto l258
						l259:
							position, tokenIndex, depth = position259, tokenIndex259, depth259
						}
						depth--
						add(rulePegText, position257)
					}
					if buffer[position] != rune('\'') {
						goto l256
					}
					position++
					if !_rules[ruleAction45]() {
						goto l256
					}
					goto l255
				l256:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
					if buffer[position] != rune('"') {
						goto l253
					}
					position++
					{
						position264 := position
						depth++
					l265:
						{
							position266, tokenIndex266, depth266 := position, tokenIndex, depth
							{
								position267, tokenIndex267, depth267 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l268
								}
								position++
								if buffer[position] != rune('\\') {
									goto l268
								}
								position++
								goto l267
							l268:
								position, tokenIndex, depth = position267, tokenIndex267, depth267
								if buffer[position] != rune('\\') {
									goto l269
								}
								position++
								if buffer[position] != rune('"') {
									goto l269
								}
								position++
								goto l267
							l269:
								position, tokenIndex, depth = position267, tokenIndex267, depth267
								{
									position270, tokenIndex270, depth270 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l270
									}
									position++
									goto l266
								l270:
									position, tokenIndex, depth = position270, tokenIndex270, depth270
								}
								if !matchDot() {
									goto l266
								}
							}
						l267:
							goto l265
						l266:
							position, tokenIndex, depth = position266, tokenIndex266, depth266
						}
						depth--
						add(rulePegText, position264)
					}
					if buffer[position] != rune('"') {
						goto l253
					}
					position++
					if !_rules[ruleAction46]() {
						goto l253
					}
				}
			l255:
				depth--
				add(rulelString, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 39 lNull <- <((('n' 'u' 'l' 'l') / ('N' 'u' 'l' 'l') / ('N' 'U' 'L' 'L')) Action47)> */
		func() bool {
			position271, tokenIndex271, depth271 := position, tokenIndex, depth
			{
				position272 := position
				depth++
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l274
					}
					position++
					if buffer[position] != rune('u') {
						goto l274
					}
					position++
					if buffer[position] != rune('l') {
						goto l274
					}
					position++
					if buffer[position] != rune('l') {
						goto l274
					}
					position++
					goto l273
				l274:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('N') {
						goto l275
					}
					position++
					if buffer[position] != rune('u') {
						goto l275
					}
					position++
					if buffer[position] != rune('l') {
						goto l275
					}
					position++
					if buffer[position] != rune('l') {
						goto l275
					}
					position++
					goto l273
				l275:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
					if buffer[position] != rune('N') {
						goto l271
					}
					position++
					if buffer[position] != rune('U') {
						goto l271
					}
					position++
					if buffer[position] != rune('L') {
						goto l271
					}
					position++
					if buffer[position] != rune('L') {
						goto l271
					}
					position++
				}
			l273:
				if !_rules[ruleAction47]() {
					goto l271
				}
				depth--
				add(rulelNull, position272)
			}
			return true
		l271:
			position, tokenIndex, depth = position271, tokenIndex271, depth271
			return false
		},
		/* 40 regex <- <(('\\' '\\') / ('\\' '/') / (!'/' .))*> */
		func() bool {
			{
				position277 := position
				depth++
			l278:
				{
					position279, tokenIndex279, depth279 := position, tokenIndex, depth
					{
						position280, tokenIndex280, depth280 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l281
						}
						position++
						if buffer[position] != rune('\\') {
							goto l281
						}
						position++
						goto l280
					l281:
						position, tokenIndex, depth = position280, tokenIndex280, depth280
						if buffer[position] != rune('\\') {
							goto l282
						}
						position++
						if buffer[position] != rune('/') {
							goto l282
						}
						position++
						goto l280
					l282:
						position, tokenIndex, depth = position280, tokenIndex280, depth280
						{
							position283, tokenIndex283, depth283 := position, tokenIndex, depth
							if buffer[position] != rune('/') {
								goto l283
							}
							position++
							goto l279
						l283:
							position, tokenIndex, depth = position283, tokenIndex283, depth283
						}
						if !matchDot() {
							goto l279
						}
					}
				l280:
					goto l278
				l279:
					position, tokenIndex, depth = position279, tokenIndex279, depth279
				}
				depth--
				add(ruleregex, position277)
			}
			return true
		},
		/* 41 squareBracketStart <- <('[' space)> */
		func() bool {
			position284, tokenIndex284, depth284 := position, tokenIndex, depth
			{
				position285 := position
				depth++
				if buffer[position] != rune('[') {
					goto l284
				}
				position++
				if !_rules[rulespace]() {
					goto l284
				}
				depth--
				add(rulesquareBracketStart, position285)
			}
			return true
		l284:
			position, tokenIndex, depth = position284, tokenIndex284, depth284
			return false
		},
		/* 42 squareBracketEnd <- <(space ']')> */
		func() bool {
			position286, tokenIndex286, depth286 := position, tokenIndex, depth
			{
				position287 := position
				depth++
				if !_rules[rulespace]() {
					goto l286
				}
				if buffer[position] != rune(']') {
					goto l286
				}
				position++
				depth--
				add(rulesquareBracketEnd, position287)
			}
			return true
		l286:
			position, tokenIndex, depth = position286, tokenIndex286, depth286
			return false
		},
		/* 43 scriptStart <- <('(' space)> */
		func() bool {
			position288, tokenIndex288, depth288 := position, tokenIndex, depth
			{
				position289 := position
				depth++
				if buffer[position] != rune('(') {
					goto l288
				}
				position++
				if !_rules[rulespace]() {
					goto l288
				}
				depth--
				add(rulescriptStart, position289)
			}
			return true
		l288:
			position, tokenIndex, depth = position288, tokenIndex288, depth288
			return false
		},
		/* 44 scriptEnd <- <(space ')')> */
		func() bool {
			position290, tokenIndex290, depth290 := position, tokenIndex, depth
			{
				position291 := position
				depth++
				if !_rules[rulespace]() {
					goto l290
				}
				if buffer[position] != rune(')') {
					goto l290
				}
				position++
				depth--
				add(rulescriptEnd, position291)
			}
			return true
		l290:
			position, tokenIndex, depth = position290, tokenIndex290, depth290
			return false
		},
		/* 45 filterStart <- <('?' '(' space)> */
		func() bool {
			position292, tokenIndex292, depth292 := position, tokenIndex, depth
			{
				position293 := position
				depth++
				if buffer[position] != rune('?') {
					goto l292
				}
				position++
				if buffer[position] != rune('(') {
					goto l292
				}
				position++
				if !_rules[rulespace]() {
					goto l292
				}
				depth--
				add(rulefilterStart, position293)
			}
			return true
		l292:
			position, tokenIndex, depth = position292, tokenIndex292, depth292
			return false
		},
		/* 46 filterEnd <- <(space ')')> */
		func() bool {
			position294, tokenIndex294, depth294 := position, tokenIndex, depth
			{
				position295 := position
				depth++
				if !_rules[rulespace]() {
					goto l294
				}
				if buffer[position] != rune(')') {
					goto l294
				}
				position++
				depth--
				add(rulefilterEnd, position295)
			}
			return true
		l294:
			position, tokenIndex, depth = position294, tokenIndex294, depth294
			return false
		},
		/* 47 subQueryStart <- <('(' space)> */
		func() bool {
			position296, tokenIndex296, depth296 := position, tokenIndex, depth
			{
				position297 := position
				depth++
				if buffer[position] != rune('(') {
					goto l296
				}
				position++
				if !_rules[rulespace]() {
					goto l296
				}
				depth--
				add(rulesubQueryStart, position297)
			}
			return true
		l296:
			position, tokenIndex, depth = position296, tokenIndex296, depth296
			return false
		},
		/* 48 subQueryEnd <- <(space ')')> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				if !_rules[rulespace]() {
					goto l298
				}
				if buffer[position] != rune(')') {
					goto l298
				}
				position++
				depth--
				add(rulesubQueryEnd, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 49 space <- <(' ' / '\t')*> */
		func() bool {
			{
				position301 := position
				depth++
			l302:
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					{
						position304, tokenIndex304, depth304 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l305
						}
						position++
						goto l304
					l305:
						position, tokenIndex, depth = position304, tokenIndex304, depth304
						if buffer[position] != rune('\t') {
							goto l303
						}
						position++
					}
				l304:
					goto l302
				l303:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
				}
				depth--
				add(rulespace, position301)
			}
			return true
		},
		/* 51 Action0 <- <{
		    p.root = p.pop().(syntaxNode)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 53 Action1 <- <{
		    p.syntaxErr(begin, msgErrorInvalidSyntaxUnrecognizedInput, buffer)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 54 Action2 <- <{
		    child := p.pop().(syntaxNode)
		    root := p.pop().(syntaxNode)
		    root.setNext(&child)
		    p.push(root)
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 55 Action3 <- <{
		    rootNode := p.pop().(syntaxNode)
		    checkNode := rootNode
		    for {
		        if checkNode.isMultiValue() {
		            rootNode.setMultiValue()
		            break
		        }
		        next := checkNode.getNext()
		        if next == nil {
		            break
		        }
		        checkNode = *next
		    }
		    p.push(rootNode)
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 56 Action4 <- <{
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
		/* 57 Action5 <- <{
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
		/* 58 Action6 <- <{
		    node := p.pop().(syntaxNode)
		    p.push(syntaxRecursiveChildIdentifier{
		        syntaxBasicNode: &syntaxBasicNode{
		            text: `..`,
		            multiValue: true,
		            next: &node,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 59 Action7 <- <{
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
		/* 60 Action8 <- <{
		    child := p.pop().(syntaxNode)
		    parent := p.pop().(syntaxNode)
		    parent.setNext(&child)
		    p.push(parent)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 61 Action9 <- <{
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
		/* 62 Action10 <- <{
		    p.push(syntaxRootIdentifier{
		        syntaxBasicNode: &syntaxBasicNode{text: `$`},
		    })
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 63 Action11 <- <{
		    p.push(syntaxCurrentRootIdentifier{
		        syntaxBasicNode: &syntaxBasicNode{text: `@`},
		    })
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 64 Action12 <- <{
		    unescapedText := p.unescape(text)
		    if unescapedText == `*` {
		        p.push(syntaxChildAsteriskIdentifier{
		            syntaxBasicNode: &syntaxBasicNode{
		                text: unescapedText,
		                multiValue: true,
		            },
		        })
		    } else {
		        p.push(syntaxChildSingleIdentifier{
		            identifier: unescapedText,
		            syntaxBasicNode: &syntaxBasicNode{
		                text: unescapedText,
		                multiValue: false,
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
		/* 65 Action13 <- <{
		    identifier := p.pop().([]string)
		    if len(identifier) > 1 {
		        p.push(syntaxChildMultiIdentifier{
		            identifiers: identifier,
		            syntaxBasicNode: &syntaxBasicNode{
		                multiValue: true,
		            },
		        })
		    } else {
		        p.push(syntaxChildSingleIdentifier{
		            identifier: identifier[0],
		            syntaxBasicNode: &syntaxBasicNode{
		                multiValue: false,
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
		/* 66 Action14 <- <{
		    p.push([]string{p.pop().(string)})
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 67 Action15 <- <{
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
		/* 68 Action16 <- <{
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 69 Action17 <- <{ // '
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 70 Action18 <- <{
		    subscript := p.pop().(syntaxSubscript)
		    union := syntaxUnion{
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: subscript.isMultiValue(),
		        }}
		    union.add(subscript)
		    p.push(union)
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 71 Action19 <- <{
		    childIndexUnion := p.pop().(syntaxUnion)
		    parentIndexUnion := p.pop().(syntaxUnion)
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
		/* 72 Action20 <- <{
		    step  := p.pop().(syntaxIndex)
		    end   := p.pop().(syntaxIndex)
		    start := p.pop().(syntaxIndex)
		    p.push(syntaxSlice{
		        syntaxBasicSubscript: &syntaxBasicSubscript{
		            multiValue: true,
		        },
		        start: start,
		        end: end,
		        step: step,
		    })
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 73 Action21 <- <{
		    p.push(syntaxIndex{
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
		/* 74 Action22 <- <{
		    p.push(syntaxAsterisk{
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
		/* 75 Action23 <- <{
		    p.push(syntaxIndex{number: 1})
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 76 Action24 <- <{
		    if len(text) > 0 {
		        p.push(syntaxIndex{number: p.toInt(text)})
		    } else {
		        p.push(syntaxIndex{number: 0, isOmitted: true})
		    }
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 77 Action25 <- <{
		    p.push(syntaxScript{
		        command: text,
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: true,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 78 Action26 <- <{
		    p.push(syntaxFilter{
		        query: p.pop().(syntaxQuery),
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: true,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 79 Action27 <- <{
		    childQuery := p.pop().(syntaxQuery)
		    parentQuery := p.pop().(syntaxQuery)
		    p.push(syntaxLogicalOr{parentQuery, childQuery})
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 80 Action28 <- <{
		    childQuery := p.pop().(syntaxQuery)
		    parentQuery := p.pop().(syntaxQuery)
		    p.push(syntaxLogicalAnd{parentQuery, childQuery})
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 81 Action29 <- <{
		    if !p.hasErr() {
		        query := p.pop().(syntaxQuery)

		        var checkQuery syntaxBasicCompareQuery
		        switch query.(type) {
		        case syntaxBasicCompareQuery:
		            checkQuery = query.(syntaxBasicCompareQuery)
		        case syntaxLogicalNot:
		            checkQuery = (query.(syntaxLogicalNot)).param.(syntaxBasicCompareQuery)
		        }

		        leftFilter, leftIsCurrentRoot := checkQuery.leftParam.(syntaxNodeFilter)
		        rightFilter, rigthIsCurrentRoot := checkQuery.rightParam.(syntaxNodeFilter)
		        if leftIsCurrentRoot && rigthIsCurrentRoot && leftFilter.isCurrentRoot() && rightFilter.isCurrentRoot() {
		            p.syntaxErr(begin, msgErrorInvalidSyntaxTwoCurrentNode, buffer)
		        }

		        p.push(query)
		    }
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 82 Action30 <- <{
		    p.push(strings.HasPrefix(text, `!`))
		}> */
		func() bool {
			{
				add(ruleAction30, position)
			}
			return true
		},
		/* 83 Action31 <- <{
		    nodeFilter := syntaxNodeFilter{p.pop().(syntaxNode)}
		    isLogicalNot := p.pop().(bool)
		    if isLogicalNot {
		        p.push(syntaxLogicalNot{nodeFilter})
		    } else {
		        p.push(nodeFilter)
		    }
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 84 Action32 <- <{
		    rightParam := p.pop().(syntaxQuery)
		    leftParam := p.pop().(syntaxQuery)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: syntaxCompareEQ{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 85 Action33 <- <{
		    rightParam := p.pop().(syntaxQuery)
		    leftParam := p.pop().(syntaxQuery)
		    p.push(syntaxLogicalNot{syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: syntaxCompareEQ{},
		    }})
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 86 Action34 <- <{
		    rightParam := p.pop().(syntaxQuery)
		    leftParam := p.pop().(syntaxQuery)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: syntaxCompareGE{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 87 Action35 <- <{
		    rightParam := p.pop().(syntaxQuery)
		    leftParam := p.pop().(syntaxQuery)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: syntaxCompareGT{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 88 Action36 <- <{
		    rightParam := p.pop().(syntaxQuery)
		    leftParam := p.pop().(syntaxQuery)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: syntaxCompareLE{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction36, position)
			}
			return true
		},
		/* 89 Action37 <- <{
		    rightParam := p.pop().(syntaxQuery)
		    leftParam := p.pop().(syntaxQuery)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: leftParam,
		        rightParam: rightParam,
		        comparator: syntaxCompareLT{},
		    })
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 90 Action38 <- <{
		    nodeFilter := syntaxNodeFilter{p.pop().(syntaxNode)}
		    regex := regexp.MustCompile(text)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: nodeFilter,
		        rightParam: syntaxCompareLiteral{literal: `regex`},
		        comparator: syntaxCompareRegex{
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
		/* 91 Action39 <- <{
		    p.push(syntaxCompareLiteral{p.pop()})
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 92 Action40 <- <{
		    p.push(syntaxCompareLiteral{p.pop()})
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 93 Action41 <- <{
		    node := p.pop().(syntaxNode)
		    p.push(syntaxNodeFilter{node})

		    if !p.hasErr() && node.isMultiValue() {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxFilterMultiValuedNode, buffer)
		    }
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 94 Action42 <- <{
		    p.push(p.toFloat(text))
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 95 Action43 <- <{
		    p.push(true)
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 96 Action44 <- <{
		    p.push(false)
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 97 Action45 <- <{
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 98 Action46 <- <{ // '
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 99 Action47 <- <{
		    p.push(nil)
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
	}
	p.rules = _rules
}
