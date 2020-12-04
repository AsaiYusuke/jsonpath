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
			if identifier, ok := node.(syntaxChildIdentifier); ok {
				identifier.setText(text)
				identifier.setMultiValue()
				p.push(syntaxRecursiveChildIdentifier{
					syntaxChildIdentifier: identifier,
				})

			} else {
				p.push(syntaxRecursiveChildIdentifier{
					syntaxChildIdentifier: syntaxChildIdentifier{
						identifiers: []string{``},
						syntaxBasicNode: &syntaxBasicNode{
							text:       `..`,
							multiValue: true,
							next:       &node,
						},
					},
				})
			}

		case ruleAction7:

			identifier := p.pop().(syntaxChildIdentifier)
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
			p.push(syntaxChildIdentifier{
				identifiers: []string{unescapedText},
				isAsterisk:  unescapedText == `*`,
				syntaxBasicNode: &syntaxBasicNode{
					text:       unescapedText,
					multiValue: unescapedText == `*`,
				},
			})

		case ruleAction13:

			identifier := p.pop().([]string)
			p.push(syntaxChildIdentifier{
				identifiers: identifier,
				isAsterisk:  false,
				syntaxBasicNode: &syntaxBasicNode{
					multiValue: len(identifier) > 1,
				},
			})

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

			node := p.pop().(syntaxNode)
			regex := regexp.MustCompile(text)
			p.push(syntaxBasicCompareQuery{
				leftParam:  syntaxNodeFilter{node},
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

			p.push(p.toFloat(text, 0))

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
		/* 4 childNodes <- <(((<('.' '.' (bracketNode / dotChildIdentifier))> Action6) / (<('.' dotChildIdentifier)> Action7) / bracketNode) (childNodes Action8)?)> */
		func() bool {
			position25, tokenIndex25, depth25 := position, tokenIndex, depth
			{
				position26 := position
				depth++
				{
					position27, tokenIndex27, depth27 := position, tokenIndex, depth
					{
						position29 := position
						depth++
						if buffer[position] != rune('.') {
							goto l28
						}
						position++
						if buffer[position] != rune('.') {
							goto l28
						}
						position++
						{
							position30, tokenIndex30, depth30 := position, tokenIndex, depth
							if !_rules[rulebracketNode]() {
								goto l31
							}
							goto l30
						l31:
							position, tokenIndex, depth = position30, tokenIndex30, depth30
							if !_rules[ruledotChildIdentifier]() {
								goto l28
							}
						}
					l30:
						depth--
						add(rulePegText, position29)
					}
					if !_rules[ruleAction6]() {
						goto l28
					}
					goto l27
				l28:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					{
						position33 := position
						depth++
						if buffer[position] != rune('.') {
							goto l32
						}
						position++
						if !_rules[ruledotChildIdentifier]() {
							goto l32
						}
						depth--
						add(rulePegText, position33)
					}
					if !_rules[ruleAction7]() {
						goto l32
					}
					goto l27
				l32:
					position, tokenIndex, depth = position27, tokenIndex27, depth27
					if !_rules[rulebracketNode]() {
						goto l25
					}
				}
			l27:
				{
					position34, tokenIndex34, depth34 := position, tokenIndex, depth
					if !_rules[rulechildNodes]() {
						goto l34
					}
					if !_rules[ruleAction8]() {
						goto l34
					}
					goto l35
				l34:
					position, tokenIndex, depth = position34, tokenIndex34, depth34
				}
			l35:
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
			position36, tokenIndex36, depth36 := position, tokenIndex, depth
			{
				position37 := position
				depth++
				{
					position38 := position
					depth++
					if !_rules[rulesquareBracketStart]() {
						goto l36
					}
					{
						position39, tokenIndex39, depth39 := position, tokenIndex, depth
						if !_rules[rulebracketChildIdentifier]() {
							goto l40
						}
						goto l39
					l40:
						position, tokenIndex, depth = position39, tokenIndex39, depth39
						if !_rules[rulequalifier]() {
							goto l36
						}
					}
				l39:
					if !_rules[rulesquareBracketEnd]() {
						goto l36
					}
					depth--
					add(rulePegText, position38)
				}
				if !_rules[ruleAction9]() {
					goto l36
				}
				depth--
				add(rulebracketNode, position37)
			}
			return true
		l36:
			position, tokenIndex, depth = position36, tokenIndex36, depth36
			return false
		},
		/* 6 rootIdentifier <- <('$' Action10)> */
		func() bool {
			position41, tokenIndex41, depth41 := position, tokenIndex, depth
			{
				position42 := position
				depth++
				if buffer[position] != rune('$') {
					goto l41
				}
				position++
				if !_rules[ruleAction10]() {
					goto l41
				}
				depth--
				add(rulerootIdentifier, position42)
			}
			return true
		l41:
			position, tokenIndex, depth = position41, tokenIndex41, depth41
			return false
		},
		/* 7 currentRootIdentifier <- <('@' Action11)> */
		func() bool {
			position43, tokenIndex43, depth43 := position, tokenIndex, depth
			{
				position44 := position
				depth++
				if buffer[position] != rune('@') {
					goto l43
				}
				position++
				if !_rules[ruleAction11]() {
					goto l43
				}
				depth--
				add(rulecurrentRootIdentifier, position44)
			}
			return true
		l43:
			position, tokenIndex, depth = position43, tokenIndex43, depth43
			return false
		},
		/* 8 dotChildIdentifier <- <(<(('\\' '\\') / ('\\' ('.' / '[' / ')' / '=' / '!' / '>' / '<' / ' ' / '\t' / '\r' / '\n')) / (!('.' / '[' / ')' / '=' / '!' / '>' / '<' / ' ' / '\t' / '\r' / '\n') .))+> Action12)> */
		func() bool {
			position45, tokenIndex45, depth45 := position, tokenIndex, depth
			{
				position46 := position
				depth++
				{
					position47 := position
					depth++
					{
						position50, tokenIndex50, depth50 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l51
						}
						position++
						if buffer[position] != rune('\\') {
							goto l51
						}
						position++
						goto l50
					l51:
						position, tokenIndex, depth = position50, tokenIndex50, depth50
						if buffer[position] != rune('\\') {
							goto l52
						}
						position++
						{
							position53, tokenIndex53, depth53 := position, tokenIndex, depth
							if buffer[position] != rune('.') {
								goto l54
							}
							position++
							goto l53
						l54:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('[') {
								goto l55
							}
							position++
							goto l53
						l55:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune(')') {
								goto l56
							}
							position++
							goto l53
						l56:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('=') {
								goto l57
							}
							position++
							goto l53
						l57:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('!') {
								goto l58
							}
							position++
							goto l53
						l58:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('>') {
								goto l59
							}
							position++
							goto l53
						l59:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('<') {
								goto l60
							}
							position++
							goto l53
						l60:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune(' ') {
								goto l61
							}
							position++
							goto l53
						l61:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('\t') {
								goto l62
							}
							position++
							goto l53
						l62:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('\r') {
								goto l63
							}
							position++
							goto l53
						l63:
							position, tokenIndex, depth = position53, tokenIndex53, depth53
							if buffer[position] != rune('\n') {
								goto l52
							}
							position++
						}
					l53:
						goto l50
					l52:
						position, tokenIndex, depth = position50, tokenIndex50, depth50
						{
							position64, tokenIndex64, depth64 := position, tokenIndex, depth
							{
								position65, tokenIndex65, depth65 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l66
								}
								position++
								goto l65
							l66:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('[') {
									goto l67
								}
								position++
								goto l65
							l67:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune(')') {
									goto l68
								}
								position++
								goto l65
							l68:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('=') {
									goto l69
								}
								position++
								goto l65
							l69:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('!') {
									goto l70
								}
								position++
								goto l65
							l70:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('>') {
									goto l71
								}
								position++
								goto l65
							l71:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('<') {
									goto l72
								}
								position++
								goto l65
							l72:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune(' ') {
									goto l73
								}
								position++
								goto l65
							l73:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('\t') {
									goto l74
								}
								position++
								goto l65
							l74:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('\r') {
									goto l75
								}
								position++
								goto l65
							l75:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('\n') {
									goto l64
								}
								position++
							}
						l65:
							goto l45
						l64:
							position, tokenIndex, depth = position64, tokenIndex64, depth64
						}
						if !matchDot() {
							goto l45
						}
					}
				l50:
				l48:
					{
						position49, tokenIndex49, depth49 := position, tokenIndex, depth
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
									if buffer[position] != rune('.') {
										goto l92
									}
									position++
									goto l91
								l92:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('[') {
										goto l93
									}
									position++
									goto l91
								l93:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune(')') {
										goto l94
									}
									position++
									goto l91
								l94:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('=') {
										goto l95
									}
									position++
									goto l91
								l95:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('!') {
										goto l96
									}
									position++
									goto l91
								l96:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('>') {
										goto l97
									}
									position++
									goto l91
								l97:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('<') {
										goto l98
									}
									position++
									goto l91
								l98:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune(' ') {
										goto l99
									}
									position++
									goto l91
								l99:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('\t') {
										goto l100
									}
									position++
									goto l91
								l100:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('\r') {
										goto l101
									}
									position++
									goto l91
								l101:
									position, tokenIndex, depth = position91, tokenIndex91, depth91
									if buffer[position] != rune('\n') {
										goto l90
									}
									position++
								}
							l91:
								goto l49
							l90:
								position, tokenIndex, depth = position90, tokenIndex90, depth90
							}
							if !matchDot() {
								goto l49
							}
						}
					l76:
						goto l48
					l49:
						position, tokenIndex, depth = position49, tokenIndex49, depth49
					}
					depth--
					add(rulePegText, position47)
				}
				if !_rules[ruleAction12]() {
					goto l45
				}
				depth--
				add(ruledotChildIdentifier, position46)
			}
			return true
		l45:
			position, tokenIndex, depth = position45, tokenIndex45, depth45
			return false
		},
		/* 9 bracketChildIdentifier <- <(bracketNodeIdentifiers Action13)> */
		func() bool {
			position102, tokenIndex102, depth102 := position, tokenIndex, depth
			{
				position103 := position
				depth++
				if !_rules[rulebracketNodeIdentifiers]() {
					goto l102
				}
				if !_rules[ruleAction13]() {
					goto l102
				}
				depth--
				add(rulebracketChildIdentifier, position103)
			}
			return true
		l102:
			position, tokenIndex, depth = position102, tokenIndex102, depth102
			return false
		},
		/* 10 bracketNodeIdentifiers <- <((singleQuotedNodeIdentifier / doubleQuotedNodeIdentifier) Action14 (sepBracketIdentifier bracketNodeIdentifiers Action15)?)> */
		func() bool {
			position104, tokenIndex104, depth104 := position, tokenIndex, depth
			{
				position105 := position
				depth++
				{
					position106, tokenIndex106, depth106 := position, tokenIndex, depth
					if !_rules[rulesingleQuotedNodeIdentifier]() {
						goto l107
					}
					goto l106
				l107:
					position, tokenIndex, depth = position106, tokenIndex106, depth106
					if !_rules[ruledoubleQuotedNodeIdentifier]() {
						goto l104
					}
				}
			l106:
				if !_rules[ruleAction14]() {
					goto l104
				}
				{
					position108, tokenIndex108, depth108 := position, tokenIndex, depth
					if !_rules[rulesepBracketIdentifier]() {
						goto l108
					}
					if !_rules[rulebracketNodeIdentifiers]() {
						goto l108
					}
					if !_rules[ruleAction15]() {
						goto l108
					}
					goto l109
				l108:
					position, tokenIndex, depth = position108, tokenIndex108, depth108
				}
			l109:
				depth--
				add(rulebracketNodeIdentifiers, position105)
			}
			return true
		l104:
			position, tokenIndex, depth = position104, tokenIndex104, depth104
			return false
		},
		/* 11 singleQuotedNodeIdentifier <- <('\'' <(('\\' '\\') / ('\\' '\'') / (!'\'' .))*> '\'' Action16)> */
		func() bool {
			position110, tokenIndex110, depth110 := position, tokenIndex, depth
			{
				position111 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l110
				}
				position++
				{
					position112 := position
					depth++
				l113:
					{
						position114, tokenIndex114, depth114 := position, tokenIndex, depth
						{
							position115, tokenIndex115, depth115 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l116
							}
							position++
							if buffer[position] != rune('\\') {
								goto l116
							}
							position++
							goto l115
						l116:
							position, tokenIndex, depth = position115, tokenIndex115, depth115
							if buffer[position] != rune('\\') {
								goto l117
							}
							position++
							if buffer[position] != rune('\'') {
								goto l117
							}
							position++
							goto l115
						l117:
							position, tokenIndex, depth = position115, tokenIndex115, depth115
							{
								position118, tokenIndex118, depth118 := position, tokenIndex, depth
								if buffer[position] != rune('\'') {
									goto l118
								}
								position++
								goto l114
							l118:
								position, tokenIndex, depth = position118, tokenIndex118, depth118
							}
							if !matchDot() {
								goto l114
							}
						}
					l115:
						goto l113
					l114:
						position, tokenIndex, depth = position114, tokenIndex114, depth114
					}
					depth--
					add(rulePegText, position112)
				}
				if buffer[position] != rune('\'') {
					goto l110
				}
				position++
				if !_rules[ruleAction16]() {
					goto l110
				}
				depth--
				add(rulesingleQuotedNodeIdentifier, position111)
			}
			return true
		l110:
			position, tokenIndex, depth = position110, tokenIndex110, depth110
			return false
		},
		/* 12 doubleQuotedNodeIdentifier <- <('"' <(('\\' '\\') / ('\\' '"') / (!'"' .))*> '"' Action17)> */
		func() bool {
			position119, tokenIndex119, depth119 := position, tokenIndex, depth
			{
				position120 := position
				depth++
				if buffer[position] != rune('"') {
					goto l119
				}
				position++
				{
					position121 := position
					depth++
				l122:
					{
						position123, tokenIndex123, depth123 := position, tokenIndex, depth
						{
							position124, tokenIndex124, depth124 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l125
							}
							position++
							if buffer[position] != rune('\\') {
								goto l125
							}
							position++
							goto l124
						l125:
							position, tokenIndex, depth = position124, tokenIndex124, depth124
							if buffer[position] != rune('\\') {
								goto l126
							}
							position++
							if buffer[position] != rune('"') {
								goto l126
							}
							position++
							goto l124
						l126:
							position, tokenIndex, depth = position124, tokenIndex124, depth124
							{
								position127, tokenIndex127, depth127 := position, tokenIndex, depth
								if buffer[position] != rune('"') {
									goto l127
								}
								position++
								goto l123
							l127:
								position, tokenIndex, depth = position127, tokenIndex127, depth127
							}
							if !matchDot() {
								goto l123
							}
						}
					l124:
						goto l122
					l123:
						position, tokenIndex, depth = position123, tokenIndex123, depth123
					}
					depth--
					add(rulePegText, position121)
				}
				if buffer[position] != rune('"') {
					goto l119
				}
				position++
				if !_rules[ruleAction17]() {
					goto l119
				}
				depth--
				add(ruledoubleQuotedNodeIdentifier, position120)
			}
			return true
		l119:
			position, tokenIndex, depth = position119, tokenIndex119, depth119
			return false
		},
		/* 13 sepBracketIdentifier <- <(space ',' space)> */
		func() bool {
			position128, tokenIndex128, depth128 := position, tokenIndex, depth
			{
				position129 := position
				depth++
				if !_rules[rulespace]() {
					goto l128
				}
				if buffer[position] != rune(',') {
					goto l128
				}
				position++
				if !_rules[rulespace]() {
					goto l128
				}
				depth--
				add(rulesepBracketIdentifier, position129)
			}
			return true
		l128:
			position, tokenIndex, depth = position128, tokenIndex128, depth128
			return false
		},
		/* 14 qualifier <- <(union / script / filter)> */
		func() bool {
			position130, tokenIndex130, depth130 := position, tokenIndex, depth
			{
				position131 := position
				depth++
				{
					position132, tokenIndex132, depth132 := position, tokenIndex, depth
					if !_rules[ruleunion]() {
						goto l133
					}
					goto l132
				l133:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if !_rules[rulescript]() {
						goto l134
					}
					goto l132
				l134:
					position, tokenIndex, depth = position132, tokenIndex132, depth132
					if !_rules[rulefilter]() {
						goto l130
					}
				}
			l132:
				depth--
				add(rulequalifier, position131)
			}
			return true
		l130:
			position, tokenIndex, depth = position130, tokenIndex130, depth130
			return false
		},
		/* 15 union <- <(index Action18 (sepUnion union Action19)?)> */
		func() bool {
			position135, tokenIndex135, depth135 := position, tokenIndex, depth
			{
				position136 := position
				depth++
				if !_rules[ruleindex]() {
					goto l135
				}
				if !_rules[ruleAction18]() {
					goto l135
				}
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if !_rules[rulesepUnion]() {
						goto l137
					}
					if !_rules[ruleunion]() {
						goto l137
					}
					if !_rules[ruleAction19]() {
						goto l137
					}
					goto l138
				l137:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
				}
			l138:
				depth--
				add(ruleunion, position136)
			}
			return true
		l135:
			position, tokenIndex, depth = position135, tokenIndex135, depth135
			return false
		},
		/* 16 index <- <((slice Action20) / (<indexNumber> Action21) / ('*' Action22))> */
		func() bool {
			position139, tokenIndex139, depth139 := position, tokenIndex, depth
			{
				position140 := position
				depth++
				{
					position141, tokenIndex141, depth141 := position, tokenIndex, depth
					if !_rules[ruleslice]() {
						goto l142
					}
					if !_rules[ruleAction20]() {
						goto l142
					}
					goto l141
				l142:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					{
						position144 := position
						depth++
						if !_rules[ruleindexNumber]() {
							goto l143
						}
						depth--
						add(rulePegText, position144)
					}
					if !_rules[ruleAction21]() {
						goto l143
					}
					goto l141
				l143:
					position, tokenIndex, depth = position141, tokenIndex141, depth141
					if buffer[position] != rune('*') {
						goto l139
					}
					position++
					if !_rules[ruleAction22]() {
						goto l139
					}
				}
			l141:
				depth--
				add(ruleindex, position140)
			}
			return true
		l139:
			position, tokenIndex, depth = position139, tokenIndex139, depth139
			return false
		},
		/* 17 slice <- <(anyIndex sepSlice anyIndex ((sepSlice anyIndex) / (space Action23)))> */
		func() bool {
			position145, tokenIndex145, depth145 := position, tokenIndex, depth
			{
				position146 := position
				depth++
				if !_rules[ruleanyIndex]() {
					goto l145
				}
				if !_rules[rulesepSlice]() {
					goto l145
				}
				if !_rules[ruleanyIndex]() {
					goto l145
				}
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					if !_rules[rulesepSlice]() {
						goto l148
					}
					if !_rules[ruleanyIndex]() {
						goto l148
					}
					goto l147
				l148:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
					if !_rules[rulespace]() {
						goto l145
					}
					if !_rules[ruleAction23]() {
						goto l145
					}
				}
			l147:
				depth--
				add(ruleslice, position146)
			}
			return true
		l145:
			position, tokenIndex, depth = position145, tokenIndex145, depth145
			return false
		},
		/* 18 anyIndex <- <(<indexNumber?> Action24)> */
		func() bool {
			position149, tokenIndex149, depth149 := position, tokenIndex, depth
			{
				position150 := position
				depth++
				{
					position151 := position
					depth++
					{
						position152, tokenIndex152, depth152 := position, tokenIndex, depth
						if !_rules[ruleindexNumber]() {
							goto l152
						}
						goto l153
					l152:
						position, tokenIndex, depth = position152, tokenIndex152, depth152
					}
				l153:
					depth--
					add(rulePegText, position151)
				}
				if !_rules[ruleAction24]() {
					goto l149
				}
				depth--
				add(ruleanyIndex, position150)
			}
			return true
		l149:
			position, tokenIndex, depth = position149, tokenIndex149, depth149
			return false
		},
		/* 19 indexNumber <- <(('-' / '+')? [0-9]+)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
				{
					position156, tokenIndex156, depth156 := position, tokenIndex, depth
					{
						position158, tokenIndex158, depth158 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l159
						}
						position++
						goto l158
					l159:
						position, tokenIndex, depth = position158, tokenIndex158, depth158
						if buffer[position] != rune('+') {
							goto l156
						}
						position++
					}
				l158:
					goto l157
				l156:
					position, tokenIndex, depth = position156, tokenIndex156, depth156
				}
			l157:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l154
				}
				position++
			l160:
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l161
					}
					position++
					goto l160
				l161:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
				}
				depth--
				add(ruleindexNumber, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 20 sepUnion <- <(space ',' space)> */
		func() bool {
			position162, tokenIndex162, depth162 := position, tokenIndex, depth
			{
				position163 := position
				depth++
				if !_rules[rulespace]() {
					goto l162
				}
				if buffer[position] != rune(',') {
					goto l162
				}
				position++
				if !_rules[rulespace]() {
					goto l162
				}
				depth--
				add(rulesepUnion, position163)
			}
			return true
		l162:
			position, tokenIndex, depth = position162, tokenIndex162, depth162
			return false
		},
		/* 21 sepSlice <- <(space ':' space)> */
		func() bool {
			position164, tokenIndex164, depth164 := position, tokenIndex, depth
			{
				position165 := position
				depth++
				if !_rules[rulespace]() {
					goto l164
				}
				if buffer[position] != rune(':') {
					goto l164
				}
				position++
				if !_rules[rulespace]() {
					goto l164
				}
				depth--
				add(rulesepSlice, position165)
			}
			return true
		l164:
			position, tokenIndex, depth = position164, tokenIndex164, depth164
			return false
		},
		/* 22 script <- <(scriptStart <command> scriptEnd Action25)> */
		func() bool {
			position166, tokenIndex166, depth166 := position, tokenIndex, depth
			{
				position167 := position
				depth++
				if !_rules[rulescriptStart]() {
					goto l166
				}
				{
					position168 := position
					depth++
					if !_rules[rulecommand]() {
						goto l166
					}
					depth--
					add(rulePegText, position168)
				}
				if !_rules[rulescriptEnd]() {
					goto l166
				}
				if !_rules[ruleAction25]() {
					goto l166
				}
				depth--
				add(rulescript, position167)
			}
			return true
		l166:
			position, tokenIndex, depth = position166, tokenIndex166, depth166
			return false
		},
		/* 23 command <- <(!')' .)*> */
		func() bool {
			{
				position170 := position
				depth++
			l171:
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					{
						position173, tokenIndex173, depth173 := position, tokenIndex, depth
						if buffer[position] != rune(')') {
							goto l173
						}
						position++
						goto l172
					l173:
						position, tokenIndex, depth = position173, tokenIndex173, depth173
					}
					if !matchDot() {
						goto l172
					}
					goto l171
				l172:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
				}
				depth--
				add(rulecommand, position170)
			}
			return true
		},
		/* 24 filter <- <(filterStart query filterEnd Action26)> */
		func() bool {
			position174, tokenIndex174, depth174 := position, tokenIndex, depth
			{
				position175 := position
				depth++
				if !_rules[rulefilterStart]() {
					goto l174
				}
				if !_rules[rulequery]() {
					goto l174
				}
				if !_rules[rulefilterEnd]() {
					goto l174
				}
				if !_rules[ruleAction26]() {
					goto l174
				}
				depth--
				add(rulefilter, position175)
			}
			return true
		l174:
			position, tokenIndex, depth = position174, tokenIndex174, depth174
			return false
		},
		/* 25 query <- <(andQuery (logicOr query Action27)?)> */
		func() bool {
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				if !_rules[ruleandQuery]() {
					goto l176
				}
				{
					position178, tokenIndex178, depth178 := position, tokenIndex, depth
					if !_rules[rulelogicOr]() {
						goto l178
					}
					if !_rules[rulequery]() {
						goto l178
					}
					if !_rules[ruleAction27]() {
						goto l178
					}
					goto l179
				l178:
					position, tokenIndex, depth = position178, tokenIndex178, depth178
				}
			l179:
				depth--
				add(rulequery, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 26 andQuery <- <((subQueryStart query subQueryEnd) / (basicQuery (logicAnd andQuery Action28)?))> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					if !_rules[rulesubQueryStart]() {
						goto l183
					}
					if !_rules[rulequery]() {
						goto l183
					}
					if !_rules[rulesubQueryEnd]() {
						goto l183
					}
					goto l182
				l183:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
					if !_rules[rulebasicQuery]() {
						goto l180
					}
					{
						position184, tokenIndex184, depth184 := position, tokenIndex, depth
						if !_rules[rulelogicAnd]() {
							goto l184
						}
						if !_rules[ruleandQuery]() {
							goto l184
						}
						if !_rules[ruleAction28]() {
							goto l184
						}
						goto l185
					l184:
						position, tokenIndex, depth = position184, tokenIndex184, depth184
					}
				l185:
				}
			l182:
				depth--
				add(ruleandQuery, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 27 basicQuery <- <((<comparator> Action29) / (<logicNot?> Action30 jsonpath Action31))> */
		func() bool {
			position186, tokenIndex186, depth186 := position, tokenIndex, depth
			{
				position187 := position
				depth++
				{
					position188, tokenIndex188, depth188 := position, tokenIndex, depth
					{
						position190 := position
						depth++
						if !_rules[rulecomparator]() {
							goto l189
						}
						depth--
						add(rulePegText, position190)
					}
					if !_rules[ruleAction29]() {
						goto l189
					}
					goto l188
				l189:
					position, tokenIndex, depth = position188, tokenIndex188, depth188
					{
						position191 := position
						depth++
						{
							position192, tokenIndex192, depth192 := position, tokenIndex, depth
							if !_rules[rulelogicNot]() {
								goto l192
							}
							goto l193
						l192:
							position, tokenIndex, depth = position192, tokenIndex192, depth192
						}
					l193:
						depth--
						add(rulePegText, position191)
					}
					if !_rules[ruleAction30]() {
						goto l186
					}
					if !_rules[rulejsonpath]() {
						goto l186
					}
					if !_rules[ruleAction31]() {
						goto l186
					}
				}
			l188:
				depth--
				add(rulebasicQuery, position187)
			}
			return true
		l186:
			position, tokenIndex, depth = position186, tokenIndex186, depth186
			return false
		},
		/* 28 logicOr <- <(space ('|' '|') space)> */
		func() bool {
			position194, tokenIndex194, depth194 := position, tokenIndex, depth
			{
				position195 := position
				depth++
				if !_rules[rulespace]() {
					goto l194
				}
				if buffer[position] != rune('|') {
					goto l194
				}
				position++
				if buffer[position] != rune('|') {
					goto l194
				}
				position++
				if !_rules[rulespace]() {
					goto l194
				}
				depth--
				add(rulelogicOr, position195)
			}
			return true
		l194:
			position, tokenIndex, depth = position194, tokenIndex194, depth194
			return false
		},
		/* 29 logicAnd <- <(space ('&' '&') space)> */
		func() bool {
			position196, tokenIndex196, depth196 := position, tokenIndex, depth
			{
				position197 := position
				depth++
				if !_rules[rulespace]() {
					goto l196
				}
				if buffer[position] != rune('&') {
					goto l196
				}
				position++
				if buffer[position] != rune('&') {
					goto l196
				}
				position++
				if !_rules[rulespace]() {
					goto l196
				}
				depth--
				add(rulelogicAnd, position197)
			}
			return true
		l196:
			position, tokenIndex, depth = position196, tokenIndex196, depth196
			return false
		},
		/* 30 logicNot <- <('!' space)> */
		func() bool {
			position198, tokenIndex198, depth198 := position, tokenIndex, depth
			{
				position199 := position
				depth++
				if buffer[position] != rune('!') {
					goto l198
				}
				position++
				if !_rules[rulespace]() {
					goto l198
				}
				depth--
				add(rulelogicNot, position199)
			}
			return true
		l198:
			position, tokenIndex, depth = position198, tokenIndex198, depth198
			return false
		},
		/* 31 comparator <- <((qParam space (('=' '=' space qParam Action32) / ('!' '=' space qParam Action33))) / (qNumericParam space (('<' '=' space qNumericParam Action34) / ('<' space qNumericParam Action35) / ('>' '=' space qNumericParam Action36) / ('>' space qNumericParam Action37))) / (jsonpath space ('=' '~') space '/' <regex> '/' Action38))> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if !_rules[ruleqParam]() {
						goto l203
					}
					if !_rules[rulespace]() {
						goto l203
					}
					{
						position204, tokenIndex204, depth204 := position, tokenIndex, depth
						if buffer[position] != rune('=') {
							goto l205
						}
						position++
						if buffer[position] != rune('=') {
							goto l205
						}
						position++
						if !_rules[rulespace]() {
							goto l205
						}
						if !_rules[ruleqParam]() {
							goto l205
						}
						if !_rules[ruleAction32]() {
							goto l205
						}
						goto l204
					l205:
						position, tokenIndex, depth = position204, tokenIndex204, depth204
						if buffer[position] != rune('!') {
							goto l203
						}
						position++
						if buffer[position] != rune('=') {
							goto l203
						}
						position++
						if !_rules[rulespace]() {
							goto l203
						}
						if !_rules[ruleqParam]() {
							goto l203
						}
						if !_rules[ruleAction33]() {
							goto l203
						}
					}
				l204:
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if !_rules[ruleqNumericParam]() {
						goto l206
					}
					if !_rules[rulespace]() {
						goto l206
					}
					{
						position207, tokenIndex207, depth207 := position, tokenIndex, depth
						if buffer[position] != rune('<') {
							goto l208
						}
						position++
						if buffer[position] != rune('=') {
							goto l208
						}
						position++
						if !_rules[rulespace]() {
							goto l208
						}
						if !_rules[ruleqNumericParam]() {
							goto l208
						}
						if !_rules[ruleAction34]() {
							goto l208
						}
						goto l207
					l208:
						position, tokenIndex, depth = position207, tokenIndex207, depth207
						if buffer[position] != rune('<') {
							goto l209
						}
						position++
						if !_rules[rulespace]() {
							goto l209
						}
						if !_rules[ruleqNumericParam]() {
							goto l209
						}
						if !_rules[ruleAction35]() {
							goto l209
						}
						goto l207
					l209:
						position, tokenIndex, depth = position207, tokenIndex207, depth207
						if buffer[position] != rune('>') {
							goto l210
						}
						position++
						if buffer[position] != rune('=') {
							goto l210
						}
						position++
						if !_rules[rulespace]() {
							goto l210
						}
						if !_rules[ruleqNumericParam]() {
							goto l210
						}
						if !_rules[ruleAction36]() {
							goto l210
						}
						goto l207
					l210:
						position, tokenIndex, depth = position207, tokenIndex207, depth207
						if buffer[position] != rune('>') {
							goto l206
						}
						position++
						if !_rules[rulespace]() {
							goto l206
						}
						if !_rules[ruleqNumericParam]() {
							goto l206
						}
						if !_rules[ruleAction37]() {
							goto l206
						}
					}
				l207:
					goto l202
				l206:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if !_rules[rulejsonpath]() {
						goto l200
					}
					if !_rules[rulespace]() {
						goto l200
					}
					if buffer[position] != rune('=') {
						goto l200
					}
					position++
					if buffer[position] != rune('~') {
						goto l200
					}
					position++
					if !_rules[rulespace]() {
						goto l200
					}
					if buffer[position] != rune('/') {
						goto l200
					}
					position++
					{
						position211 := position
						depth++
						if !_rules[ruleregex]() {
							goto l200
						}
						depth--
						add(rulePegText, position211)
					}
					if buffer[position] != rune('/') {
						goto l200
					}
					position++
					if !_rules[ruleAction38]() {
						goto l200
					}
				}
			l202:
				depth--
				add(rulecomparator, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 32 qParam <- <((qLiteral Action39) / nodeFilter)> */
		func() bool {
			position212, tokenIndex212, depth212 := position, tokenIndex, depth
			{
				position213 := position
				depth++
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					if !_rules[ruleqLiteral]() {
						goto l215
					}
					if !_rules[ruleAction39]() {
						goto l215
					}
					goto l214
				l215:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
					if !_rules[rulenodeFilter]() {
						goto l212
					}
				}
			l214:
				depth--
				add(ruleqParam, position213)
			}
			return true
		l212:
			position, tokenIndex, depth = position212, tokenIndex212, depth212
			return false
		},
		/* 33 qNumericParam <- <((lNumber Action40) / nodeFilter)> */
		func() bool {
			position216, tokenIndex216, depth216 := position, tokenIndex, depth
			{
				position217 := position
				depth++
				{
					position218, tokenIndex218, depth218 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l219
					}
					if !_rules[ruleAction40]() {
						goto l219
					}
					goto l218
				l219:
					position, tokenIndex, depth = position218, tokenIndex218, depth218
					if !_rules[rulenodeFilter]() {
						goto l216
					}
				}
			l218:
				depth--
				add(ruleqNumericParam, position217)
			}
			return true
		l216:
			position, tokenIndex, depth = position216, tokenIndex216, depth216
			return false
		},
		/* 34 qLiteral <- <(lNumber / lBool / lString / lNull)> */
		func() bool {
			position220, tokenIndex220, depth220 := position, tokenIndex, depth
			{
				position221 := position
				depth++
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l223
					}
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if !_rules[rulelBool]() {
						goto l224
					}
					goto l222
				l224:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if !_rules[rulelString]() {
						goto l225
					}
					goto l222
				l225:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					if !_rules[rulelNull]() {
						goto l220
					}
				}
			l222:
				depth--
				add(ruleqLiteral, position221)
			}
			return true
		l220:
			position, tokenIndex, depth = position220, tokenIndex220, depth220
			return false
		},
		/* 35 nodeFilter <- <(<jsonpath> Action41)> */
		func() bool {
			position226, tokenIndex226, depth226 := position, tokenIndex, depth
			{
				position227 := position
				depth++
				{
					position228 := position
					depth++
					if !_rules[rulejsonpath]() {
						goto l226
					}
					depth--
					add(rulePegText, position228)
				}
				if !_rules[ruleAction41]() {
					goto l226
				}
				depth--
				add(rulenodeFilter, position227)
			}
			return true
		l226:
			position, tokenIndex, depth = position226, tokenIndex226, depth226
			return false
		},
		/* 36 lNumber <- <(<(('-' / '+')? [0-9] ('-' / '+' / '.' / [0-9] / [a-z] / [A-Z])*)> Action42)> */
		func() bool {
			position229, tokenIndex229, depth229 := position, tokenIndex, depth
			{
				position230 := position
				depth++
				{
					position231 := position
					depth++
					{
						position232, tokenIndex232, depth232 := position, tokenIndex, depth
						{
							position234, tokenIndex234, depth234 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l235
							}
							position++
							goto l234
						l235:
							position, tokenIndex, depth = position234, tokenIndex234, depth234
							if buffer[position] != rune('+') {
								goto l232
							}
							position++
						}
					l234:
						goto l233
					l232:
						position, tokenIndex, depth = position232, tokenIndex232, depth232
					}
				l233:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l229
					}
					position++
				l236:
					{
						position237, tokenIndex237, depth237 := position, tokenIndex, depth
						{
							position238, tokenIndex238, depth238 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l239
							}
							position++
							goto l238
						l239:
							position, tokenIndex, depth = position238, tokenIndex238, depth238
							if buffer[position] != rune('+') {
								goto l240
							}
							position++
							goto l238
						l240:
							position, tokenIndex, depth = position238, tokenIndex238, depth238
							if buffer[position] != rune('.') {
								goto l241
							}
							position++
							goto l238
						l241:
							position, tokenIndex, depth = position238, tokenIndex238, depth238
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l242
							}
							position++
							goto l238
						l242:
							position, tokenIndex, depth = position238, tokenIndex238, depth238
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l243
							}
							position++
							goto l238
						l243:
							position, tokenIndex, depth = position238, tokenIndex238, depth238
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l237
							}
							position++
						}
					l238:
						goto l236
					l237:
						position, tokenIndex, depth = position237, tokenIndex237, depth237
					}
					depth--
					add(rulePegText, position231)
				}
				if !_rules[ruleAction42]() {
					goto l229
				}
				depth--
				add(rulelNumber, position230)
			}
			return true
		l229:
			position, tokenIndex, depth = position229, tokenIndex229, depth229
			return false
		},
		/* 37 lBool <- <(((('t' 'r' 'u' 'e') / ('T' 'r' 'u' 'e') / ('T' 'R' 'U' 'E')) Action43) / ((('f' 'a' 'l' 's' 'e') / ('F' 'a' 'l' 's' 'e') / ('F' 'A' 'L' 'S' 'E')) Action44))> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					{
						position248, tokenIndex248, depth248 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
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
						goto l248
					l249:
						position, tokenIndex, depth = position248, tokenIndex248, depth248
						if buffer[position] != rune('T') {
							goto l250
						}
						position++
						if buffer[position] != rune('r') {
							goto l250
						}
						position++
						if buffer[position] != rune('u') {
							goto l250
						}
						position++
						if buffer[position] != rune('e') {
							goto l250
						}
						position++
						goto l248
					l250:
						position, tokenIndex, depth = position248, tokenIndex248, depth248
						if buffer[position] != rune('T') {
							goto l247
						}
						position++
						if buffer[position] != rune('R') {
							goto l247
						}
						position++
						if buffer[position] != rune('U') {
							goto l247
						}
						position++
						if buffer[position] != rune('E') {
							goto l247
						}
						position++
					}
				l248:
					if !_rules[ruleAction43]() {
						goto l247
					}
					goto l246
				l247:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
					{
						position251, tokenIndex251, depth251 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
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
						goto l251
					l252:
						position, tokenIndex, depth = position251, tokenIndex251, depth251
						if buffer[position] != rune('F') {
							goto l253
						}
						position++
						if buffer[position] != rune('a') {
							goto l253
						}
						position++
						if buffer[position] != rune('l') {
							goto l253
						}
						position++
						if buffer[position] != rune('s') {
							goto l253
						}
						position++
						if buffer[position] != rune('e') {
							goto l253
						}
						position++
						goto l251
					l253:
						position, tokenIndex, depth = position251, tokenIndex251, depth251
						if buffer[position] != rune('F') {
							goto l244
						}
						position++
						if buffer[position] != rune('A') {
							goto l244
						}
						position++
						if buffer[position] != rune('L') {
							goto l244
						}
						position++
						if buffer[position] != rune('S') {
							goto l244
						}
						position++
						if buffer[position] != rune('E') {
							goto l244
						}
						position++
					}
				l251:
					if !_rules[ruleAction44]() {
						goto l244
					}
				}
			l246:
				depth--
				add(rulelBool, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
			return false
		},
		/* 38 lString <- <(('\'' <(('\\' '\\') / ('\\' '\'') / (!'\'' .))*> '\'' Action45) / ('"' <(('\\' '\\') / ('\\' '"') / (!'"' .))*> '"' Action46))> */
		func() bool {
			position254, tokenIndex254, depth254 := position, tokenIndex, depth
			{
				position255 := position
				depth++
				{
					position256, tokenIndex256, depth256 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l257
					}
					position++
					{
						position258 := position
						depth++
					l259:
						{
							position260, tokenIndex260, depth260 := position, tokenIndex, depth
							{
								position261, tokenIndex261, depth261 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l262
								}
								position++
								if buffer[position] != rune('\\') {
									goto l262
								}
								position++
								goto l261
							l262:
								position, tokenIndex, depth = position261, tokenIndex261, depth261
								if buffer[position] != rune('\\') {
									goto l263
								}
								position++
								if buffer[position] != rune('\'') {
									goto l263
								}
								position++
								goto l261
							l263:
								position, tokenIndex, depth = position261, tokenIndex261, depth261
								{
									position264, tokenIndex264, depth264 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l264
									}
									position++
									goto l260
								l264:
									position, tokenIndex, depth = position264, tokenIndex264, depth264
								}
								if !matchDot() {
									goto l260
								}
							}
						l261:
							goto l259
						l260:
							position, tokenIndex, depth = position260, tokenIndex260, depth260
						}
						depth--
						add(rulePegText, position258)
					}
					if buffer[position] != rune('\'') {
						goto l257
					}
					position++
					if !_rules[ruleAction45]() {
						goto l257
					}
					goto l256
				l257:
					position, tokenIndex, depth = position256, tokenIndex256, depth256
					if buffer[position] != rune('"') {
						goto l254
					}
					position++
					{
						position265 := position
						depth++
					l266:
						{
							position267, tokenIndex267, depth267 := position, tokenIndex, depth
							{
								position268, tokenIndex268, depth268 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l269
								}
								position++
								if buffer[position] != rune('\\') {
									goto l269
								}
								position++
								goto l268
							l269:
								position, tokenIndex, depth = position268, tokenIndex268, depth268
								if buffer[position] != rune('\\') {
									goto l270
								}
								position++
								if buffer[position] != rune('"') {
									goto l270
								}
								position++
								goto l268
							l270:
								position, tokenIndex, depth = position268, tokenIndex268, depth268
								{
									position271, tokenIndex271, depth271 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l271
									}
									position++
									goto l267
								l271:
									position, tokenIndex, depth = position271, tokenIndex271, depth271
								}
								if !matchDot() {
									goto l267
								}
							}
						l268:
							goto l266
						l267:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
						}
						depth--
						add(rulePegText, position265)
					}
					if buffer[position] != rune('"') {
						goto l254
					}
					position++
					if !_rules[ruleAction46]() {
						goto l254
					}
				}
			l256:
				depth--
				add(rulelString, position255)
			}
			return true
		l254:
			position, tokenIndex, depth = position254, tokenIndex254, depth254
			return false
		},
		/* 39 lNull <- <((('n' 'u' 'l' 'l') / ('N' 'u' 'l' 'l') / ('N' 'U' 'L' 'L')) Action47)> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
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
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('N') {
						goto l276
					}
					position++
					if buffer[position] != rune('u') {
						goto l276
					}
					position++
					if buffer[position] != rune('l') {
						goto l276
					}
					position++
					if buffer[position] != rune('l') {
						goto l276
					}
					position++
					goto l274
				l276:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					if buffer[position] != rune('N') {
						goto l272
					}
					position++
					if buffer[position] != rune('U') {
						goto l272
					}
					position++
					if buffer[position] != rune('L') {
						goto l272
					}
					position++
					if buffer[position] != rune('L') {
						goto l272
					}
					position++
				}
			l274:
				if !_rules[ruleAction47]() {
					goto l272
				}
				depth--
				add(rulelNull, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 40 regex <- <(('\\' '\\') / ('\\' '/') / (!'/' .))*> */
		func() bool {
			{
				position278 := position
				depth++
			l279:
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					{
						position281, tokenIndex281, depth281 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l282
						}
						position++
						if buffer[position] != rune('\\') {
							goto l282
						}
						position++
						goto l281
					l282:
						position, tokenIndex, depth = position281, tokenIndex281, depth281
						if buffer[position] != rune('\\') {
							goto l283
						}
						position++
						if buffer[position] != rune('/') {
							goto l283
						}
						position++
						goto l281
					l283:
						position, tokenIndex, depth = position281, tokenIndex281, depth281
						{
							position284, tokenIndex284, depth284 := position, tokenIndex, depth
							if buffer[position] != rune('/') {
								goto l284
							}
							position++
							goto l280
						l284:
							position, tokenIndex, depth = position284, tokenIndex284, depth284
						}
						if !matchDot() {
							goto l280
						}
					}
				l281:
					goto l279
				l280:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
				}
				depth--
				add(ruleregex, position278)
			}
			return true
		},
		/* 41 squareBracketStart <- <('[' space)> */
		func() bool {
			position285, tokenIndex285, depth285 := position, tokenIndex, depth
			{
				position286 := position
				depth++
				if buffer[position] != rune('[') {
					goto l285
				}
				position++
				if !_rules[rulespace]() {
					goto l285
				}
				depth--
				add(rulesquareBracketStart, position286)
			}
			return true
		l285:
			position, tokenIndex, depth = position285, tokenIndex285, depth285
			return false
		},
		/* 42 squareBracketEnd <- <(space ']')> */
		func() bool {
			position287, tokenIndex287, depth287 := position, tokenIndex, depth
			{
				position288 := position
				depth++
				if !_rules[rulespace]() {
					goto l287
				}
				if buffer[position] != rune(']') {
					goto l287
				}
				position++
				depth--
				add(rulesquareBracketEnd, position288)
			}
			return true
		l287:
			position, tokenIndex, depth = position287, tokenIndex287, depth287
			return false
		},
		/* 43 scriptStart <- <('(' space)> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{
				position290 := position
				depth++
				if buffer[position] != rune('(') {
					goto l289
				}
				position++
				if !_rules[rulespace]() {
					goto l289
				}
				depth--
				add(rulescriptStart, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 44 scriptEnd <- <(space ')')> */
		func() bool {
			position291, tokenIndex291, depth291 := position, tokenIndex, depth
			{
				position292 := position
				depth++
				if !_rules[rulespace]() {
					goto l291
				}
				if buffer[position] != rune(')') {
					goto l291
				}
				position++
				depth--
				add(rulescriptEnd, position292)
			}
			return true
		l291:
			position, tokenIndex, depth = position291, tokenIndex291, depth291
			return false
		},
		/* 45 filterStart <- <('?' '(' space)> */
		func() bool {
			position293, tokenIndex293, depth293 := position, tokenIndex, depth
			{
				position294 := position
				depth++
				if buffer[position] != rune('?') {
					goto l293
				}
				position++
				if buffer[position] != rune('(') {
					goto l293
				}
				position++
				if !_rules[rulespace]() {
					goto l293
				}
				depth--
				add(rulefilterStart, position294)
			}
			return true
		l293:
			position, tokenIndex, depth = position293, tokenIndex293, depth293
			return false
		},
		/* 46 filterEnd <- <(space ')')> */
		func() bool {
			position295, tokenIndex295, depth295 := position, tokenIndex, depth
			{
				position296 := position
				depth++
				if !_rules[rulespace]() {
					goto l295
				}
				if buffer[position] != rune(')') {
					goto l295
				}
				position++
				depth--
				add(rulefilterEnd, position296)
			}
			return true
		l295:
			position, tokenIndex, depth = position295, tokenIndex295, depth295
			return false
		},
		/* 47 subQueryStart <- <('(' space)> */
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
				add(rulesubQueryStart, position298)
			}
			return true
		l297:
			position, tokenIndex, depth = position297, tokenIndex297, depth297
			return false
		},
		/* 48 subQueryEnd <- <(space ')')> */
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
				add(rulesubQueryEnd, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 49 space <- <(' ' / '\t')*> */
		func() bool {
			{
				position302 := position
				depth++
			l303:
				{
					position304, tokenIndex304, depth304 := position, tokenIndex, depth
					{
						position305, tokenIndex305, depth305 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l306
						}
						position++
						goto l305
					l306:
						position, tokenIndex, depth = position305, tokenIndex305, depth305
						if buffer[position] != rune('\t') {
							goto l304
						}
						position++
					}
				l305:
					goto l303
				l304:
					position, tokenIndex, depth = position304, tokenIndex304, depth304
				}
				depth--
				add(rulespace, position302)
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
		    if identifier, ok := node.(syntaxChildIdentifier); ok {
		        identifier.setText(text)
		        identifier.setMultiValue()
		        p.push(syntaxRecursiveChildIdentifier{
		            syntaxChildIdentifier: identifier,
		        })

		    } else {
		        p.push(syntaxRecursiveChildIdentifier{
		            syntaxChildIdentifier: syntaxChildIdentifier{
		                identifiers: []string{``},
		                syntaxBasicNode: &syntaxBasicNode{
		                    text: `..`,
		                    multiValue: true,
		                    next: &node,
		                },
		            },
		        })
		    }
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 59 Action7 <- <{
		    identifier := p.pop().(syntaxChildIdentifier)
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
		    p.push(syntaxChildIdentifier{
		        identifiers: []string{unescapedText},
		        isAsterisk: unescapedText == `*`,
		        syntaxBasicNode: &syntaxBasicNode{
		            text: unescapedText,
		            multiValue: unescapedText == `*`,
		        },
		    })
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 65 Action13 <- <{
		    identifier := p.pop().([]string)
		    p.push(syntaxChildIdentifier{
		        identifiers: identifier,
		        isAsterisk: false,
		        syntaxBasicNode: &syntaxBasicNode{
		            multiValue: len(identifier) > 1,
		        },
		    })
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
		    node := p.pop().(syntaxNode)
		    regex := regexp.MustCompile(text)
		    p.push(syntaxBasicCompareQuery{
		        leftParam: syntaxNodeFilter{node},
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
		    p.push(p.toFloat(text, 0))
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