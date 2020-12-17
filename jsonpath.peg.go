package jsonpath

import (
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
	rulechildNode
	rulefunction
	rulefunctionName
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
	ruleAction49

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
	"childNode",
	"function",
	"functionName",
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
	"Action49",

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
	rules  [105]func() bool
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

			p.saveParams()

		case ruleAction3:

			p.setNodeChain()
			p.setRecursiveMultiValue()
			p.loadParams()

		case ruleAction4:

			if len(p.paramsList) == 0 {
				p.syntaxErr(begin, msgErrorInvalidSyntaxUseBeginAtsign, buffer)
			}

		case ruleAction5:

			if len(p.paramsList) != 0 {
				p.syntaxErr(begin, msgErrorInvalidSyntaxOmitDollar, buffer)
			}

		case ruleAction6:

			p.pushRecursiveChildIdentifier(p.pop().(syntaxNode))

		case ruleAction7:

			p.setLastNodeText(text)

		case ruleAction8:

			funcName := p.pop().(string)
			p.pushFunction(text, funcName)

		case ruleAction9:

			p.push(text)

		case ruleAction10:

			p.setLastNodeText(text)

		case ruleAction11:

			p.pushRootIdentifier()

		case ruleAction12:

			p.pushCurrentRootIdentifier()

		case ruleAction13:

			unescapedText := p.unescape(text)
			if unescapedText == `*` {
				p.pushChildAsteriskIdentifier(unescapedText)
			} else {
				p.pushChildSingleIdentifier(unescapedText)
			}

		case ruleAction14:

			identifier := p.pop().([]string)
			if len(identifier) > 1 {
				p.pushChildMultiIdentifier(identifier)
			} else {
				p.pushChildSingleIdentifier(identifier[0])
			}

		case ruleAction15:

			p.push([]string{p.pop().(string)})

		case ruleAction16:

			identifier2 := p.pop().([]string)
			identifier1 := p.pop().([]string)
			identifier1 = append(identifier1, identifier2...)
			p.push(identifier1)

		case ruleAction17:

			p.push(p.unescape(text))

		case ruleAction18:
			// '
			p.push(p.unescape(text))

		case ruleAction19:

			subscript := p.pop().(syntaxSubscript)
			p.pushUnionQualifier(subscript)

		case ruleAction20:

			childIndexUnion := p.pop().(*syntaxUnionQualifier)
			parentIndexUnion := p.pop().(*syntaxUnionQualifier)
			parentIndexUnion.merge(childIndexUnion)
			parentIndexUnion.setMultiValue()
			p.push(parentIndexUnion)

		case ruleAction21:

			step := p.pop().(*syntaxIndexSubscript)
			end := p.pop().(*syntaxIndexSubscript)
			start := p.pop().(*syntaxIndexSubscript)

			if step.isOmitted || step.number == 0 {
				step.number = 1
			}

			if step.number > 0 {
				p.pushSliceSubscript(true, start, end, step)
			} else {
				p.pushSliceSubscript(false, start, end, step)
			}

		case ruleAction22:

			p.pushIndexSubscript(text, false)

		case ruleAction23:

			p.pushAsteriskSubscript()

		case ruleAction24:

			p.pushIndexSubscript(`1`, false)

		case ruleAction25:

			if len(text) > 0 {
				p.pushIndexSubscript(text, false)
			} else {
				p.pushIndexSubscript(`0`, true)
			}

		case ruleAction26:

			p.pushScriptQualifier(text)

		case ruleAction27:

			p.pushFilterQualifier(p.pop().(syntaxQuery))

		case ruleAction28:

			childQuery := p.pop().(syntaxQuery)
			parentQuery := p.pop().(syntaxQuery)
			p.pushLogicalOr(parentQuery, childQuery)

		case ruleAction29:

			childQuery := p.pop().(syntaxQuery)
			parentQuery := p.pop().(syntaxQuery)
			p.pushLogicalAnd(parentQuery, childQuery)

		case ruleAction30:

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

		case ruleAction31:

			p.push(len(text) > 0 && text[0:1] == `!`)

		case ruleAction32:

			_ = p.pop().(bool)
			jsonpathFilter := p.pop().(syntaxQuery)
			isLogicalNot := p.pop().(bool)
			if isLogicalNot {
				p.pushLogicalNot(jsonpathFilter)
			} else {
				p.push(jsonpathFilter)
			}

		case ruleAction33:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.pushCompareEQ(leftParam, rightParam)

		case ruleAction34:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.pushCompareNE(leftParam, rightParam)

		case ruleAction35:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.pushCompareGE(leftParam, rightParam)

		case ruleAction36:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.push(&syntaxBasicCompareQuery{
				leftParam:  leftParam,
				rightParam: rightParam,
				comparator: &syntaxCompareGT{},
			})

		case ruleAction37:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.pushCompareLE(leftParam, rightParam)

		case ruleAction38:

			rightParam := p.pop().(*syntaxBasicCompareParameter)
			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.pushCompareLT(leftParam, rightParam)

		case ruleAction39:

			leftParam := p.pop().(*syntaxBasicCompareParameter)
			p.pushCompareRegex(leftParam, text)

		case ruleAction40:

			p.pushCompareParameterLiteral(p.pop())

		case ruleAction41:

			p.pushCompareParameterLiteral(p.pop())

		case ruleAction42:

			isLiteral := p.pop().(bool)
			param := p.pop().(syntaxQueryParameter)
			if !p.hasErr() && param.isMultiValueParameter() {
				p.syntaxErr(begin, msgErrorInvalidSyntaxFilterValueGroup, buffer)
			}
			p.pushBasicCompareParameter(param, isLiteral)

		case ruleAction43:

			node := p.pop().(syntaxNode)

			switch node.(type) {
			case *syntaxRootIdentifier:
				p.pushCompareParameterRoot(node)
				p.push(true)
			case *syntaxCurrentRootIdentifier:
				p.pushCompareParameterCurrentRoot(node)
				p.push(false)
			default:
				p.push(&syntaxQueryParamRoot{})
				p.push(true)
			}

		case ruleAction44:

			p.push(p.toFloat(text))

		case ruleAction45:

			p.push(true)

		case ruleAction46:

			p.push(false)

		case ruleAction47:

			p.push(p.unescape(text))

		case ruleAction48:
			// '
			p.push(p.unescape(text))

		case ruleAction49:

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
		/* 2 jsonpath <- <(space Action2 rootNode childNode* function* space Action3)> */
		func() bool {
			position12, tokenIndex12, depth12 := position, tokenIndex, depth
			{
				position13 := position
				depth++
				if !_rules[rulespace]() {
					goto l12
				}
				if !_rules[ruleAction2]() {
					goto l12
				}
				if !_rules[rulerootNode]() {
					goto l12
				}
			l14:
				{
					position15, tokenIndex15, depth15 := position, tokenIndex, depth
					if !_rules[rulechildNode]() {
						goto l15
					}
					goto l14
				l15:
					position, tokenIndex, depth = position15, tokenIndex15, depth15
				}
			l16:
				{
					position17, tokenIndex17, depth17 := position, tokenIndex, depth
					if !_rules[rulefunction]() {
						goto l17
					}
					goto l16
				l17:
					position, tokenIndex, depth = position17, tokenIndex17, depth17
				}
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
			position18, tokenIndex18, depth18 := position, tokenIndex, depth
			{
				position19 := position
				depth++
				{
					position20, tokenIndex20, depth20 := position, tokenIndex, depth
					if !_rules[rulerootIdentifier]() {
						goto l21
					}
					goto l20
				l21:
					position, tokenIndex, depth = position20, tokenIndex20, depth20
					{
						position23 := position
						depth++
						if !_rules[rulecurrentRootIdentifier]() {
							goto l22
						}
						depth--
						add(rulePegText, position23)
					}
					if !_rules[ruleAction4]() {
						goto l22
					}
					goto l20
				l22:
					position, tokenIndex, depth = position20, tokenIndex20, depth20
					{
						position24 := position
						depth++
						{
							position25, tokenIndex25, depth25 := position, tokenIndex, depth
							if !_rules[rulebracketNode]() {
								goto l26
							}
							goto l25
						l26:
							position, tokenIndex, depth = position25, tokenIndex25, depth25
							if !_rules[ruledotChildIdentifier]() {
								goto l18
							}
						}
					l25:
						depth--
						add(rulePegText, position24)
					}
					if !_rules[ruleAction5]() {
						goto l18
					}
				}
			l20:
				depth--
				add(rulerootNode, position19)
			}
			return true
		l18:
			position, tokenIndex, depth = position18, tokenIndex18, depth18
			return false
		},
		/* 4 childNode <- <(('.' '.' (bracketNode / dotChildIdentifier) Action6) / (<('.' dotChildIdentifier)> Action7) / bracketNode)> */
		func() bool {
			position27, tokenIndex27, depth27 := position, tokenIndex, depth
			{
				position28 := position
				depth++
				{
					position29, tokenIndex29, depth29 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l30
					}
					position++
					if buffer[position] != rune('.') {
						goto l30
					}
					position++
					{
						position31, tokenIndex31, depth31 := position, tokenIndex, depth
						if !_rules[rulebracketNode]() {
							goto l32
						}
						goto l31
					l32:
						position, tokenIndex, depth = position31, tokenIndex31, depth31
						if !_rules[ruledotChildIdentifier]() {
							goto l30
						}
					}
				l31:
					if !_rules[ruleAction6]() {
						goto l30
					}
					goto l29
				l30:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					{
						position34 := position
						depth++
						if buffer[position] != rune('.') {
							goto l33
						}
						position++
						if !_rules[ruledotChildIdentifier]() {
							goto l33
						}
						depth--
						add(rulePegText, position34)
					}
					if !_rules[ruleAction7]() {
						goto l33
					}
					goto l29
				l33:
					position, tokenIndex, depth = position29, tokenIndex29, depth29
					if !_rules[rulebracketNode]() {
						goto l27
					}
				}
			l29:
				depth--
				add(rulechildNode, position28)
			}
			return true
		l27:
			position, tokenIndex, depth = position27, tokenIndex27, depth27
			return false
		},
		/* 5 function <- <(<('.' functionName ('(' ')'))> Action8)> */
		func() bool {
			position35, tokenIndex35, depth35 := position, tokenIndex, depth
			{
				position36 := position
				depth++
				{
					position37 := position
					depth++
					if buffer[position] != rune('.') {
						goto l35
					}
					position++
					if !_rules[rulefunctionName]() {
						goto l35
					}
					if buffer[position] != rune('(') {
						goto l35
					}
					position++
					if buffer[position] != rune(')') {
						goto l35
					}
					position++
					depth--
					add(rulePegText, position37)
				}
				if !_rules[ruleAction8]() {
					goto l35
				}
				depth--
				add(rulefunction, position36)
			}
			return true
		l35:
			position, tokenIndex, depth = position35, tokenIndex35, depth35
			return false
		},
		/* 6 functionName <- <(<('-' / '_' / [a-z] / [A-Z])+> Action9)> */
		func() bool {
			position38, tokenIndex38, depth38 := position, tokenIndex, depth
			{
				position39 := position
				depth++
				{
					position40 := position
					depth++
					{
						position43, tokenIndex43, depth43 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l44
						}
						position++
						goto l43
					l44:
						position, tokenIndex, depth = position43, tokenIndex43, depth43
						if buffer[position] != rune('_') {
							goto l45
						}
						position++
						goto l43
					l45:
						position, tokenIndex, depth = position43, tokenIndex43, depth43
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l46
						}
						position++
						goto l43
					l46:
						position, tokenIndex, depth = position43, tokenIndex43, depth43
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l38
						}
						position++
					}
				l43:
				l41:
					{
						position42, tokenIndex42, depth42 := position, tokenIndex, depth
						{
							position47, tokenIndex47, depth47 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l48
							}
							position++
							goto l47
						l48:
							position, tokenIndex, depth = position47, tokenIndex47, depth47
							if buffer[position] != rune('_') {
								goto l49
							}
							position++
							goto l47
						l49:
							position, tokenIndex, depth = position47, tokenIndex47, depth47
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l50
							}
							position++
							goto l47
						l50:
							position, tokenIndex, depth = position47, tokenIndex47, depth47
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l42
							}
							position++
						}
					l47:
						goto l41
					l42:
						position, tokenIndex, depth = position42, tokenIndex42, depth42
					}
					depth--
					add(rulePegText, position40)
				}
				if !_rules[ruleAction9]() {
					goto l38
				}
				depth--
				add(rulefunctionName, position39)
			}
			return true
		l38:
			position, tokenIndex, depth = position38, tokenIndex38, depth38
			return false
		},
		/* 7 bracketNode <- <(<(squareBracketStart (bracketChildIdentifier / qualifier) squareBracketEnd)> Action10)> */
		func() bool {
			position51, tokenIndex51, depth51 := position, tokenIndex, depth
			{
				position52 := position
				depth++
				{
					position53 := position
					depth++
					if !_rules[rulesquareBracketStart]() {
						goto l51
					}
					{
						position54, tokenIndex54, depth54 := position, tokenIndex, depth
						if !_rules[rulebracketChildIdentifier]() {
							goto l55
						}
						goto l54
					l55:
						position, tokenIndex, depth = position54, tokenIndex54, depth54
						if !_rules[rulequalifier]() {
							goto l51
						}
					}
				l54:
					if !_rules[rulesquareBracketEnd]() {
						goto l51
					}
					depth--
					add(rulePegText, position53)
				}
				if !_rules[ruleAction10]() {
					goto l51
				}
				depth--
				add(rulebracketNode, position52)
			}
			return true
		l51:
			position, tokenIndex, depth = position51, tokenIndex51, depth51
			return false
		},
		/* 8 rootIdentifier <- <('$' Action11)> */
		func() bool {
			position56, tokenIndex56, depth56 := position, tokenIndex, depth
			{
				position57 := position
				depth++
				if buffer[position] != rune('$') {
					goto l56
				}
				position++
				if !_rules[ruleAction11]() {
					goto l56
				}
				depth--
				add(rulerootIdentifier, position57)
			}
			return true
		l56:
			position, tokenIndex, depth = position56, tokenIndex56, depth56
			return false
		},
		/* 9 currentRootIdentifier <- <('@' Action12)> */
		func() bool {
			position58, tokenIndex58, depth58 := position, tokenIndex, depth
			{
				position59 := position
				depth++
				if buffer[position] != rune('@') {
					goto l58
				}
				position++
				if !_rules[ruleAction12]() {
					goto l58
				}
				depth--
				add(rulecurrentRootIdentifier, position59)
			}
			return true
		l58:
			position, tokenIndex, depth = position58, tokenIndex58, depth58
			return false
		},
		/* 10 dotChildIdentifier <- <(<(('\\' '\\') / ('\\' ('.' / '[' / '(' / ')' / '=' / '!' / '>' / '<' / '\t' / '\r' / '\n' / ' ')) / (!('\\' / '.' / '[' / '(' / ')' / '=' / '!' / '>' / '<' / '\t' / '\r' / '\n' / ' ') .))+> !('(' ')') Action13)> */
		func() bool {
			position60, tokenIndex60, depth60 := position, tokenIndex, depth
			{
				position61 := position
				depth++
				{
					position62 := position
					depth++
					{
						position65, tokenIndex65, depth65 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l66
						}
						position++
						if buffer[position] != rune('\\') {
							goto l66
						}
						position++
						goto l65
					l66:
						position, tokenIndex, depth = position65, tokenIndex65, depth65
						if buffer[position] != rune('\\') {
							goto l67
						}
						position++
						{
							position68, tokenIndex68, depth68 := position, tokenIndex, depth
							if buffer[position] != rune('.') {
								goto l69
							}
							position++
							goto l68
						l69:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('[') {
								goto l70
							}
							position++
							goto l68
						l70:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('(') {
								goto l71
							}
							position++
							goto l68
						l71:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune(')') {
								goto l72
							}
							position++
							goto l68
						l72:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('=') {
								goto l73
							}
							position++
							goto l68
						l73:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('!') {
								goto l74
							}
							position++
							goto l68
						l74:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('>') {
								goto l75
							}
							position++
							goto l68
						l75:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('<') {
								goto l76
							}
							position++
							goto l68
						l76:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('\t') {
								goto l77
							}
							position++
							goto l68
						l77:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('\r') {
								goto l78
							}
							position++
							goto l68
						l78:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune('\n') {
								goto l79
							}
							position++
							goto l68
						l79:
							position, tokenIndex, depth = position68, tokenIndex68, depth68
							if buffer[position] != rune(' ') {
								goto l67
							}
							position++
						}
					l68:
						goto l65
					l67:
						position, tokenIndex, depth = position65, tokenIndex65, depth65
						{
							position80, tokenIndex80, depth80 := position, tokenIndex, depth
							{
								position81, tokenIndex81, depth81 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l82
								}
								position++
								goto l81
							l82:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('.') {
									goto l83
								}
								position++
								goto l81
							l83:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('[') {
									goto l84
								}
								position++
								goto l81
							l84:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('(') {
									goto l85
								}
								position++
								goto l81
							l85:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune(')') {
									goto l86
								}
								position++
								goto l81
							l86:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('=') {
									goto l87
								}
								position++
								goto l81
							l87:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('!') {
									goto l88
								}
								position++
								goto l81
							l88:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('>') {
									goto l89
								}
								position++
								goto l81
							l89:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('<') {
									goto l90
								}
								position++
								goto l81
							l90:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('\t') {
									goto l91
								}
								position++
								goto l81
							l91:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('\r') {
									goto l92
								}
								position++
								goto l81
							l92:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune('\n') {
									goto l93
								}
								position++
								goto l81
							l93:
								position, tokenIndex, depth = position81, tokenIndex81, depth81
								if buffer[position] != rune(' ') {
									goto l80
								}
								position++
							}
						l81:
							goto l60
						l80:
							position, tokenIndex, depth = position80, tokenIndex80, depth80
						}
						if !matchDot() {
							goto l60
						}
					}
				l65:
				l63:
					{
						position64, tokenIndex64, depth64 := position, tokenIndex, depth
						{
							position94, tokenIndex94, depth94 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l95
							}
							position++
							if buffer[position] != rune('\\') {
								goto l95
							}
							position++
							goto l94
						l95:
							position, tokenIndex, depth = position94, tokenIndex94, depth94
							if buffer[position] != rune('\\') {
								goto l96
							}
							position++
							{
								position97, tokenIndex97, depth97 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l98
								}
								position++
								goto l97
							l98:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('[') {
									goto l99
								}
								position++
								goto l97
							l99:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('(') {
									goto l100
								}
								position++
								goto l97
							l100:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune(')') {
									goto l101
								}
								position++
								goto l97
							l101:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('=') {
									goto l102
								}
								position++
								goto l97
							l102:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('!') {
									goto l103
								}
								position++
								goto l97
							l103:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('>') {
									goto l104
								}
								position++
								goto l97
							l104:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('<') {
									goto l105
								}
								position++
								goto l97
							l105:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('\t') {
									goto l106
								}
								position++
								goto l97
							l106:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('\r') {
									goto l107
								}
								position++
								goto l97
							l107:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('\n') {
									goto l108
								}
								position++
								goto l97
							l108:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune(' ') {
									goto l96
								}
								position++
							}
						l97:
							goto l94
						l96:
							position, tokenIndex, depth = position94, tokenIndex94, depth94
							{
								position109, tokenIndex109, depth109 := position, tokenIndex, depth
								{
									position110, tokenIndex110, depth110 := position, tokenIndex, depth
									if buffer[position] != rune('\\') {
										goto l111
									}
									position++
									goto l110
								l111:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('.') {
										goto l112
									}
									position++
									goto l110
								l112:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('[') {
										goto l113
									}
									position++
									goto l110
								l113:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('(') {
										goto l114
									}
									position++
									goto l110
								l114:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune(')') {
										goto l115
									}
									position++
									goto l110
								l115:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('=') {
										goto l116
									}
									position++
									goto l110
								l116:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('!') {
										goto l117
									}
									position++
									goto l110
								l117:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('>') {
										goto l118
									}
									position++
									goto l110
								l118:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('<') {
										goto l119
									}
									position++
									goto l110
								l119:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('\t') {
										goto l120
									}
									position++
									goto l110
								l120:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('\r') {
										goto l121
									}
									position++
									goto l110
								l121:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune('\n') {
										goto l122
									}
									position++
									goto l110
								l122:
									position, tokenIndex, depth = position110, tokenIndex110, depth110
									if buffer[position] != rune(' ') {
										goto l109
									}
									position++
								}
							l110:
								goto l64
							l109:
								position, tokenIndex, depth = position109, tokenIndex109, depth109
							}
							if !matchDot() {
								goto l64
							}
						}
					l94:
						goto l63
					l64:
						position, tokenIndex, depth = position64, tokenIndex64, depth64
					}
					depth--
					add(rulePegText, position62)
				}
				{
					position123, tokenIndex123, depth123 := position, tokenIndex, depth
					if buffer[position] != rune('(') {
						goto l123
					}
					position++
					if buffer[position] != rune(')') {
						goto l123
					}
					position++
					goto l60
				l123:
					position, tokenIndex, depth = position123, tokenIndex123, depth123
				}
				if !_rules[ruleAction13]() {
					goto l60
				}
				depth--
				add(ruledotChildIdentifier, position61)
			}
			return true
		l60:
			position, tokenIndex, depth = position60, tokenIndex60, depth60
			return false
		},
		/* 11 bracketChildIdentifier <- <(bracketNodeIdentifiers Action14)> */
		func() bool {
			position124, tokenIndex124, depth124 := position, tokenIndex, depth
			{
				position125 := position
				depth++
				if !_rules[rulebracketNodeIdentifiers]() {
					goto l124
				}
				if !_rules[ruleAction14]() {
					goto l124
				}
				depth--
				add(rulebracketChildIdentifier, position125)
			}
			return true
		l124:
			position, tokenIndex, depth = position124, tokenIndex124, depth124
			return false
		},
		/* 12 bracketNodeIdentifiers <- <((singleQuotedNodeIdentifier / doubleQuotedNodeIdentifier) Action15 (sepBracketIdentifier bracketNodeIdentifiers Action16)?)> */
		func() bool {
			position126, tokenIndex126, depth126 := position, tokenIndex, depth
			{
				position127 := position
				depth++
				{
					position128, tokenIndex128, depth128 := position, tokenIndex, depth
					if !_rules[rulesingleQuotedNodeIdentifier]() {
						goto l129
					}
					goto l128
				l129:
					position, tokenIndex, depth = position128, tokenIndex128, depth128
					if !_rules[ruledoubleQuotedNodeIdentifier]() {
						goto l126
					}
				}
			l128:
				if !_rules[ruleAction15]() {
					goto l126
				}
				{
					position130, tokenIndex130, depth130 := position, tokenIndex, depth
					if !_rules[rulesepBracketIdentifier]() {
						goto l130
					}
					if !_rules[rulebracketNodeIdentifiers]() {
						goto l130
					}
					if !_rules[ruleAction16]() {
						goto l130
					}
					goto l131
				l130:
					position, tokenIndex, depth = position130, tokenIndex130, depth130
				}
			l131:
				depth--
				add(rulebracketNodeIdentifiers, position127)
			}
			return true
		l126:
			position, tokenIndex, depth = position126, tokenIndex126, depth126
			return false
		},
		/* 13 singleQuotedNodeIdentifier <- <('\'' <(('\\' '\\') / ('\\' '\'') / (!('\\' / '\'') .))*> '\'' Action17)> */
		func() bool {
			position132, tokenIndex132, depth132 := position, tokenIndex, depth
			{
				position133 := position
				depth++
				if buffer[position] != rune('\'') {
					goto l132
				}
				position++
				{
					position134 := position
					depth++
				l135:
					{
						position136, tokenIndex136, depth136 := position, tokenIndex, depth
						{
							position137, tokenIndex137, depth137 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l138
							}
							position++
							if buffer[position] != rune('\\') {
								goto l138
							}
							position++
							goto l137
						l138:
							position, tokenIndex, depth = position137, tokenIndex137, depth137
							if buffer[position] != rune('\\') {
								goto l139
							}
							position++
							if buffer[position] != rune('\'') {
								goto l139
							}
							position++
							goto l137
						l139:
							position, tokenIndex, depth = position137, tokenIndex137, depth137
							{
								position140, tokenIndex140, depth140 := position, tokenIndex, depth
								{
									position141, tokenIndex141, depth141 := position, tokenIndex, depth
									if buffer[position] != rune('\\') {
										goto l142
									}
									position++
									goto l141
								l142:
									position, tokenIndex, depth = position141, tokenIndex141, depth141
									if buffer[position] != rune('\'') {
										goto l140
									}
									position++
								}
							l141:
								goto l136
							l140:
								position, tokenIndex, depth = position140, tokenIndex140, depth140
							}
							if !matchDot() {
								goto l136
							}
						}
					l137:
						goto l135
					l136:
						position, tokenIndex, depth = position136, tokenIndex136, depth136
					}
					depth--
					add(rulePegText, position134)
				}
				if buffer[position] != rune('\'') {
					goto l132
				}
				position++
				if !_rules[ruleAction17]() {
					goto l132
				}
				depth--
				add(rulesingleQuotedNodeIdentifier, position133)
			}
			return true
		l132:
			position, tokenIndex, depth = position132, tokenIndex132, depth132
			return false
		},
		/* 14 doubleQuotedNodeIdentifier <- <('"' <(('\\' '\\') / ('\\' '"') / (!('\\' / '"') .))*> '"' Action18)> */
		func() bool {
			position143, tokenIndex143, depth143 := position, tokenIndex, depth
			{
				position144 := position
				depth++
				if buffer[position] != rune('"') {
					goto l143
				}
				position++
				{
					position145 := position
					depth++
				l146:
					{
						position147, tokenIndex147, depth147 := position, tokenIndex, depth
						{
							position148, tokenIndex148, depth148 := position, tokenIndex, depth
							if buffer[position] != rune('\\') {
								goto l149
							}
							position++
							if buffer[position] != rune('\\') {
								goto l149
							}
							position++
							goto l148
						l149:
							position, tokenIndex, depth = position148, tokenIndex148, depth148
							if buffer[position] != rune('\\') {
								goto l150
							}
							position++
							if buffer[position] != rune('"') {
								goto l150
							}
							position++
							goto l148
						l150:
							position, tokenIndex, depth = position148, tokenIndex148, depth148
							{
								position151, tokenIndex151, depth151 := position, tokenIndex, depth
								{
									position152, tokenIndex152, depth152 := position, tokenIndex, depth
									if buffer[position] != rune('\\') {
										goto l153
									}
									position++
									goto l152
								l153:
									position, tokenIndex, depth = position152, tokenIndex152, depth152
									if buffer[position] != rune('"') {
										goto l151
									}
									position++
								}
							l152:
								goto l147
							l151:
								position, tokenIndex, depth = position151, tokenIndex151, depth151
							}
							if !matchDot() {
								goto l147
							}
						}
					l148:
						goto l146
					l147:
						position, tokenIndex, depth = position147, tokenIndex147, depth147
					}
					depth--
					add(rulePegText, position145)
				}
				if buffer[position] != rune('"') {
					goto l143
				}
				position++
				if !_rules[ruleAction18]() {
					goto l143
				}
				depth--
				add(ruledoubleQuotedNodeIdentifier, position144)
			}
			return true
		l143:
			position, tokenIndex, depth = position143, tokenIndex143, depth143
			return false
		},
		/* 15 sepBracketIdentifier <- <(space ',' space)> */
		func() bool {
			position154, tokenIndex154, depth154 := position, tokenIndex, depth
			{
				position155 := position
				depth++
				if !_rules[rulespace]() {
					goto l154
				}
				if buffer[position] != rune(',') {
					goto l154
				}
				position++
				if !_rules[rulespace]() {
					goto l154
				}
				depth--
				add(rulesepBracketIdentifier, position155)
			}
			return true
		l154:
			position, tokenIndex, depth = position154, tokenIndex154, depth154
			return false
		},
		/* 16 qualifier <- <(union / script / filter)> */
		func() bool {
			position156, tokenIndex156, depth156 := position, tokenIndex, depth
			{
				position157 := position
				depth++
				{
					position158, tokenIndex158, depth158 := position, tokenIndex, depth
					if !_rules[ruleunion]() {
						goto l159
					}
					goto l158
				l159:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if !_rules[rulescript]() {
						goto l160
					}
					goto l158
				l160:
					position, tokenIndex, depth = position158, tokenIndex158, depth158
					if !_rules[rulefilter]() {
						goto l156
					}
				}
			l158:
				depth--
				add(rulequalifier, position157)
			}
			return true
		l156:
			position, tokenIndex, depth = position156, tokenIndex156, depth156
			return false
		},
		/* 17 union <- <(index Action19 (sepUnion union Action20)?)> */
		func() bool {
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				if !_rules[ruleindex]() {
					goto l161
				}
				if !_rules[ruleAction19]() {
					goto l161
				}
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if !_rules[rulesepUnion]() {
						goto l163
					}
					if !_rules[ruleunion]() {
						goto l163
					}
					if !_rules[ruleAction20]() {
						goto l163
					}
					goto l164
				l163:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
				}
			l164:
				depth--
				add(ruleunion, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 18 index <- <((slice Action21) / (<indexNumber> Action22) / ('*' Action23))> */
		func() bool {
			position165, tokenIndex165, depth165 := position, tokenIndex, depth
			{
				position166 := position
				depth++
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					if !_rules[ruleslice]() {
						goto l168
					}
					if !_rules[ruleAction21]() {
						goto l168
					}
					goto l167
				l168:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					{
						position170 := position
						depth++
						if !_rules[ruleindexNumber]() {
							goto l169
						}
						depth--
						add(rulePegText, position170)
					}
					if !_rules[ruleAction22]() {
						goto l169
					}
					goto l167
				l169:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
					if buffer[position] != rune('*') {
						goto l165
					}
					position++
					if !_rules[ruleAction23]() {
						goto l165
					}
				}
			l167:
				depth--
				add(ruleindex, position166)
			}
			return true
		l165:
			position, tokenIndex, depth = position165, tokenIndex165, depth165
			return false
		},
		/* 19 slice <- <(anyIndex sepSlice anyIndex ((sepSlice anyIndex) / (space Action24)))> */
		func() bool {
			position171, tokenIndex171, depth171 := position, tokenIndex, depth
			{
				position172 := position
				depth++
				if !_rules[ruleanyIndex]() {
					goto l171
				}
				if !_rules[rulesepSlice]() {
					goto l171
				}
				if !_rules[ruleanyIndex]() {
					goto l171
				}
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					if !_rules[rulesepSlice]() {
						goto l174
					}
					if !_rules[ruleanyIndex]() {
						goto l174
					}
					goto l173
				l174:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
					if !_rules[rulespace]() {
						goto l171
					}
					if !_rules[ruleAction24]() {
						goto l171
					}
				}
			l173:
				depth--
				add(ruleslice, position172)
			}
			return true
		l171:
			position, tokenIndex, depth = position171, tokenIndex171, depth171
			return false
		},
		/* 20 anyIndex <- <(<indexNumber?> Action25)> */
		func() bool {
			position175, tokenIndex175, depth175 := position, tokenIndex, depth
			{
				position176 := position
				depth++
				{
					position177 := position
					depth++
					{
						position178, tokenIndex178, depth178 := position, tokenIndex, depth
						if !_rules[ruleindexNumber]() {
							goto l178
						}
						goto l179
					l178:
						position, tokenIndex, depth = position178, tokenIndex178, depth178
					}
				l179:
					depth--
					add(rulePegText, position177)
				}
				if !_rules[ruleAction25]() {
					goto l175
				}
				depth--
				add(ruleanyIndex, position176)
			}
			return true
		l175:
			position, tokenIndex, depth = position175, tokenIndex175, depth175
			return false
		},
		/* 21 indexNumber <- <(('-' / '+')? [0-9]+)> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					{
						position184, tokenIndex184, depth184 := position, tokenIndex, depth
						if buffer[position] != rune('-') {
							goto l185
						}
						position++
						goto l184
					l185:
						position, tokenIndex, depth = position184, tokenIndex184, depth184
						if buffer[position] != rune('+') {
							goto l182
						}
						position++
					}
				l184:
					goto l183
				l182:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
				}
			l183:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l180
				}
				position++
			l186:
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l187
					}
					position++
					goto l186
				l187:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
				}
				depth--
				add(ruleindexNumber, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 22 sepUnion <- <(space ',' space)> */
		func() bool {
			position188, tokenIndex188, depth188 := position, tokenIndex, depth
			{
				position189 := position
				depth++
				if !_rules[rulespace]() {
					goto l188
				}
				if buffer[position] != rune(',') {
					goto l188
				}
				position++
				if !_rules[rulespace]() {
					goto l188
				}
				depth--
				add(rulesepUnion, position189)
			}
			return true
		l188:
			position, tokenIndex, depth = position188, tokenIndex188, depth188
			return false
		},
		/* 23 sepSlice <- <(space ':' space)> */
		func() bool {
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				if !_rules[rulespace]() {
					goto l190
				}
				if buffer[position] != rune(':') {
					goto l190
				}
				position++
				if !_rules[rulespace]() {
					goto l190
				}
				depth--
				add(rulesepSlice, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 24 script <- <(scriptStart <command> scriptEnd Action26)> */
		func() bool {
			position192, tokenIndex192, depth192 := position, tokenIndex, depth
			{
				position193 := position
				depth++
				if !_rules[rulescriptStart]() {
					goto l192
				}
				{
					position194 := position
					depth++
					if !_rules[rulecommand]() {
						goto l192
					}
					depth--
					add(rulePegText, position194)
				}
				if !_rules[rulescriptEnd]() {
					goto l192
				}
				if !_rules[ruleAction26]() {
					goto l192
				}
				depth--
				add(rulescript, position193)
			}
			return true
		l192:
			position, tokenIndex, depth = position192, tokenIndex192, depth192
			return false
		},
		/* 25 command <- <(!')' .)+> */
		func() bool {
			position195, tokenIndex195, depth195 := position, tokenIndex, depth
			{
				position196 := position
				depth++
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if buffer[position] != rune(')') {
						goto l199
					}
					position++
					goto l195
				l199:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
				}
				if !matchDot() {
					goto l195
				}
			l197:
				{
					position198, tokenIndex198, depth198 := position, tokenIndex, depth
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if buffer[position] != rune(')') {
							goto l200
						}
						position++
						goto l198
					l200:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
					}
					if !matchDot() {
						goto l198
					}
					goto l197
				l198:
					position, tokenIndex, depth = position198, tokenIndex198, depth198
				}
				depth--
				add(rulecommand, position196)
			}
			return true
		l195:
			position, tokenIndex, depth = position195, tokenIndex195, depth195
			return false
		},
		/* 26 filter <- <(filterStart query filterEnd Action27)> */
		func() bool {
			position201, tokenIndex201, depth201 := position, tokenIndex, depth
			{
				position202 := position
				depth++
				if !_rules[rulefilterStart]() {
					goto l201
				}
				if !_rules[rulequery]() {
					goto l201
				}
				if !_rules[rulefilterEnd]() {
					goto l201
				}
				if !_rules[ruleAction27]() {
					goto l201
				}
				depth--
				add(rulefilter, position202)
			}
			return true
		l201:
			position, tokenIndex, depth = position201, tokenIndex201, depth201
			return false
		},
		/* 27 query <- <(andQuery (logicOr query Action28)?)> */
		func() bool {
			position203, tokenIndex203, depth203 := position, tokenIndex, depth
			{
				position204 := position
				depth++
				if !_rules[ruleandQuery]() {
					goto l203
				}
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					if !_rules[rulelogicOr]() {
						goto l205
					}
					if !_rules[rulequery]() {
						goto l205
					}
					if !_rules[ruleAction28]() {
						goto l205
					}
					goto l206
				l205:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
				}
			l206:
				depth--
				add(rulequery, position204)
			}
			return true
		l203:
			position, tokenIndex, depth = position203, tokenIndex203, depth203
			return false
		},
		/* 28 andQuery <- <((subQueryStart query subQueryEnd) / (basicQuery (logicAnd andQuery Action29)?))> */
		func() bool {
			position207, tokenIndex207, depth207 := position, tokenIndex, depth
			{
				position208 := position
				depth++
				{
					position209, tokenIndex209, depth209 := position, tokenIndex, depth
					if !_rules[rulesubQueryStart]() {
						goto l210
					}
					if !_rules[rulequery]() {
						goto l210
					}
					if !_rules[rulesubQueryEnd]() {
						goto l210
					}
					goto l209
				l210:
					position, tokenIndex, depth = position209, tokenIndex209, depth209
					if !_rules[rulebasicQuery]() {
						goto l207
					}
					{
						position211, tokenIndex211, depth211 := position, tokenIndex, depth
						if !_rules[rulelogicAnd]() {
							goto l211
						}
						if !_rules[ruleandQuery]() {
							goto l211
						}
						if !_rules[ruleAction29]() {
							goto l211
						}
						goto l212
					l211:
						position, tokenIndex, depth = position211, tokenIndex211, depth211
					}
				l212:
				}
			l209:
				depth--
				add(ruleandQuery, position208)
			}
			return true
		l207:
			position, tokenIndex, depth = position207, tokenIndex207, depth207
			return false
		},
		/* 29 basicQuery <- <((<comparator> Action30) / (<logicNot?> Action31 jsonpathFilter Action32))> */
		func() bool {
			position213, tokenIndex213, depth213 := position, tokenIndex, depth
			{
				position214 := position
				depth++
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					{
						position217 := position
						depth++
						if !_rules[rulecomparator]() {
							goto l216
						}
						depth--
						add(rulePegText, position217)
					}
					if !_rules[ruleAction30]() {
						goto l216
					}
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					{
						position218 := position
						depth++
						{
							position219, tokenIndex219, depth219 := position, tokenIndex, depth
							if !_rules[rulelogicNot]() {
								goto l219
							}
							goto l220
						l219:
							position, tokenIndex, depth = position219, tokenIndex219, depth219
						}
					l220:
						depth--
						add(rulePegText, position218)
					}
					if !_rules[ruleAction31]() {
						goto l213
					}
					if !_rules[rulejsonpathFilter]() {
						goto l213
					}
					if !_rules[ruleAction32]() {
						goto l213
					}
				}
			l215:
				depth--
				add(rulebasicQuery, position214)
			}
			return true
		l213:
			position, tokenIndex, depth = position213, tokenIndex213, depth213
			return false
		},
		/* 30 logicOr <- <(space ('|' '|') space)> */
		func() bool {
			position221, tokenIndex221, depth221 := position, tokenIndex, depth
			{
				position222 := position
				depth++
				if !_rules[rulespace]() {
					goto l221
				}
				if buffer[position] != rune('|') {
					goto l221
				}
				position++
				if buffer[position] != rune('|') {
					goto l221
				}
				position++
				if !_rules[rulespace]() {
					goto l221
				}
				depth--
				add(rulelogicOr, position222)
			}
			return true
		l221:
			position, tokenIndex, depth = position221, tokenIndex221, depth221
			return false
		},
		/* 31 logicAnd <- <(space ('&' '&') space)> */
		func() bool {
			position223, tokenIndex223, depth223 := position, tokenIndex, depth
			{
				position224 := position
				depth++
				if !_rules[rulespace]() {
					goto l223
				}
				if buffer[position] != rune('&') {
					goto l223
				}
				position++
				if buffer[position] != rune('&') {
					goto l223
				}
				position++
				if !_rules[rulespace]() {
					goto l223
				}
				depth--
				add(rulelogicAnd, position224)
			}
			return true
		l223:
			position, tokenIndex, depth = position223, tokenIndex223, depth223
			return false
		},
		/* 32 logicNot <- <('!' space)> */
		func() bool {
			position225, tokenIndex225, depth225 := position, tokenIndex, depth
			{
				position226 := position
				depth++
				if buffer[position] != rune('!') {
					goto l225
				}
				position++
				if !_rules[rulespace]() {
					goto l225
				}
				depth--
				add(rulelogicNot, position226)
			}
			return true
		l225:
			position, tokenIndex, depth = position225, tokenIndex225, depth225
			return false
		},
		/* 33 comparator <- <((qParam space (('=' '=' space qParam Action33) / ('!' '=' space qParam Action34))) / (qNumericParam space (('<' '=' space qNumericParam Action35) / ('<' space qNumericParam Action36) / ('>' '=' space qNumericParam Action37) / ('>' space qNumericParam Action38))) / (singleJsonpathFilter space ('=' '~') space '/' <regex> '/' Action39))> */
		func() bool {
			position227, tokenIndex227, depth227 := position, tokenIndex, depth
			{
				position228 := position
				depth++
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if !_rules[ruleqParam]() {
						goto l230
					}
					if !_rules[rulespace]() {
						goto l230
					}
					{
						position231, tokenIndex231, depth231 := position, tokenIndex, depth
						if buffer[position] != rune('=') {
							goto l232
						}
						position++
						if buffer[position] != rune('=') {
							goto l232
						}
						position++
						if !_rules[rulespace]() {
							goto l232
						}
						if !_rules[ruleqParam]() {
							goto l232
						}
						if !_rules[ruleAction33]() {
							goto l232
						}
						goto l231
					l232:
						position, tokenIndex, depth = position231, tokenIndex231, depth231
						if buffer[position] != rune('!') {
							goto l230
						}
						position++
						if buffer[position] != rune('=') {
							goto l230
						}
						position++
						if !_rules[rulespace]() {
							goto l230
						}
						if !_rules[ruleqParam]() {
							goto l230
						}
						if !_rules[ruleAction34]() {
							goto l230
						}
					}
				l231:
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if !_rules[ruleqNumericParam]() {
						goto l233
					}
					if !_rules[rulespace]() {
						goto l233
					}
					{
						position234, tokenIndex234, depth234 := position, tokenIndex, depth
						if buffer[position] != rune('<') {
							goto l235
						}
						position++
						if buffer[position] != rune('=') {
							goto l235
						}
						position++
						if !_rules[rulespace]() {
							goto l235
						}
						if !_rules[ruleqNumericParam]() {
							goto l235
						}
						if !_rules[ruleAction35]() {
							goto l235
						}
						goto l234
					l235:
						position, tokenIndex, depth = position234, tokenIndex234, depth234
						if buffer[position] != rune('<') {
							goto l236
						}
						position++
						if !_rules[rulespace]() {
							goto l236
						}
						if !_rules[ruleqNumericParam]() {
							goto l236
						}
						if !_rules[ruleAction36]() {
							goto l236
						}
						goto l234
					l236:
						position, tokenIndex, depth = position234, tokenIndex234, depth234
						if buffer[position] != rune('>') {
							goto l237
						}
						position++
						if buffer[position] != rune('=') {
							goto l237
						}
						position++
						if !_rules[rulespace]() {
							goto l237
						}
						if !_rules[ruleqNumericParam]() {
							goto l237
						}
						if !_rules[ruleAction37]() {
							goto l237
						}
						goto l234
					l237:
						position, tokenIndex, depth = position234, tokenIndex234, depth234
						if buffer[position] != rune('>') {
							goto l233
						}
						position++
						if !_rules[rulespace]() {
							goto l233
						}
						if !_rules[ruleqNumericParam]() {
							goto l233
						}
						if !_rules[ruleAction38]() {
							goto l233
						}
					}
				l234:
					goto l229
				l233:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if !_rules[rulesingleJsonpathFilter]() {
						goto l227
					}
					if !_rules[rulespace]() {
						goto l227
					}
					if buffer[position] != rune('=') {
						goto l227
					}
					position++
					if buffer[position] != rune('~') {
						goto l227
					}
					position++
					if !_rules[rulespace]() {
						goto l227
					}
					if buffer[position] != rune('/') {
						goto l227
					}
					position++
					{
						position238 := position
						depth++
						if !_rules[ruleregex]() {
							goto l227
						}
						depth--
						add(rulePegText, position238)
					}
					if buffer[position] != rune('/') {
						goto l227
					}
					position++
					if !_rules[ruleAction39]() {
						goto l227
					}
				}
			l229:
				depth--
				add(rulecomparator, position228)
			}
			return true
		l227:
			position, tokenIndex, depth = position227, tokenIndex227, depth227
			return false
		},
		/* 34 qParam <- <((qLiteral Action40) / singleJsonpathFilter)> */
		func() bool {
			position239, tokenIndex239, depth239 := position, tokenIndex, depth
			{
				position240 := position
				depth++
				{
					position241, tokenIndex241, depth241 := position, tokenIndex, depth
					if !_rules[ruleqLiteral]() {
						goto l242
					}
					if !_rules[ruleAction40]() {
						goto l242
					}
					goto l241
				l242:
					position, tokenIndex, depth = position241, tokenIndex241, depth241
					if !_rules[rulesingleJsonpathFilter]() {
						goto l239
					}
				}
			l241:
				depth--
				add(ruleqParam, position240)
			}
			return true
		l239:
			position, tokenIndex, depth = position239, tokenIndex239, depth239
			return false
		},
		/* 35 qNumericParam <- <((lNumber Action41) / singleJsonpathFilter)> */
		func() bool {
			position243, tokenIndex243, depth243 := position, tokenIndex, depth
			{
				position244 := position
				depth++
				{
					position245, tokenIndex245, depth245 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l246
					}
					if !_rules[ruleAction41]() {
						goto l246
					}
					goto l245
				l246:
					position, tokenIndex, depth = position245, tokenIndex245, depth245
					if !_rules[rulesingleJsonpathFilter]() {
						goto l243
					}
				}
			l245:
				depth--
				add(ruleqNumericParam, position244)
			}
			return true
		l243:
			position, tokenIndex, depth = position243, tokenIndex243, depth243
			return false
		},
		/* 36 qLiteral <- <(lNumber / lBool / lString / lNull)> */
		func() bool {
			position247, tokenIndex247, depth247 := position, tokenIndex, depth
			{
				position248 := position
				depth++
				{
					position249, tokenIndex249, depth249 := position, tokenIndex, depth
					if !_rules[rulelNumber]() {
						goto l250
					}
					goto l249
				l250:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if !_rules[rulelBool]() {
						goto l251
					}
					goto l249
				l251:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if !_rules[rulelString]() {
						goto l252
					}
					goto l249
				l252:
					position, tokenIndex, depth = position249, tokenIndex249, depth249
					if !_rules[rulelNull]() {
						goto l247
					}
				}
			l249:
				depth--
				add(ruleqLiteral, position248)
			}
			return true
		l247:
			position, tokenIndex, depth = position247, tokenIndex247, depth247
			return false
		},
		/* 37 singleJsonpathFilter <- <(jsonpathFilter Action42)> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				if !_rules[rulejsonpathFilter]() {
					goto l253
				}
				if !_rules[ruleAction42]() {
					goto l253
				}
				depth--
				add(rulesingleJsonpathFilter, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 38 jsonpathFilter <- <(<jsonpath> Action43)> */
		func() bool {
			position255, tokenIndex255, depth255 := position, tokenIndex, depth
			{
				position256 := position
				depth++
				{
					position257 := position
					depth++
					if !_rules[rulejsonpath]() {
						goto l255
					}
					depth--
					add(rulePegText, position257)
				}
				if !_rules[ruleAction43]() {
					goto l255
				}
				depth--
				add(rulejsonpathFilter, position256)
			}
			return true
		l255:
			position, tokenIndex, depth = position255, tokenIndex255, depth255
			return false
		},
		/* 39 lNumber <- <(<(('-' / '+')? [0-9] ('-' / '+' / '.' / [0-9] / [a-z] / [A-Z])*)> Action44)> */
		func() bool {
			position258, tokenIndex258, depth258 := position, tokenIndex, depth
			{
				position259 := position
				depth++
				{
					position260 := position
					depth++
					{
						position261, tokenIndex261, depth261 := position, tokenIndex, depth
						{
							position263, tokenIndex263, depth263 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l264
							}
							position++
							goto l263
						l264:
							position, tokenIndex, depth = position263, tokenIndex263, depth263
							if buffer[position] != rune('+') {
								goto l261
							}
							position++
						}
					l263:
						goto l262
					l261:
						position, tokenIndex, depth = position261, tokenIndex261, depth261
					}
				l262:
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l258
					}
					position++
				l265:
					{
						position266, tokenIndex266, depth266 := position, tokenIndex, depth
						{
							position267, tokenIndex267, depth267 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l268
							}
							position++
							goto l267
						l268:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if buffer[position] != rune('+') {
								goto l269
							}
							position++
							goto l267
						l269:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if buffer[position] != rune('.') {
								goto l270
							}
							position++
							goto l267
						l270:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l271
							}
							position++
							goto l267
						l271:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l272
							}
							position++
							goto l267
						l272:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l266
							}
							position++
						}
					l267:
						goto l265
					l266:
						position, tokenIndex, depth = position266, tokenIndex266, depth266
					}
					depth--
					add(rulePegText, position260)
				}
				if !_rules[ruleAction44]() {
					goto l258
				}
				depth--
				add(rulelNumber, position259)
			}
			return true
		l258:
			position, tokenIndex, depth = position258, tokenIndex258, depth258
			return false
		},
		/* 40 lBool <- <(((('t' 'r' 'u' 'e') / ('T' 'r' 'u' 'e') / ('T' 'R' 'U' 'E')) Action45) / ((('f' 'a' 'l' 's' 'e') / ('F' 'a' 'l' 's' 'e') / ('F' 'A' 'L' 'S' 'E')) Action46))> */
		func() bool {
			position273, tokenIndex273, depth273 := position, tokenIndex, depth
			{
				position274 := position
				depth++
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					{
						position277, tokenIndex277, depth277 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l278
						}
						position++
						if buffer[position] != rune('r') {
							goto l278
						}
						position++
						if buffer[position] != rune('u') {
							goto l278
						}
						position++
						if buffer[position] != rune('e') {
							goto l278
						}
						position++
						goto l277
					l278:
						position, tokenIndex, depth = position277, tokenIndex277, depth277
						if buffer[position] != rune('T') {
							goto l279
						}
						position++
						if buffer[position] != rune('r') {
							goto l279
						}
						position++
						if buffer[position] != rune('u') {
							goto l279
						}
						position++
						if buffer[position] != rune('e') {
							goto l279
						}
						position++
						goto l277
					l279:
						position, tokenIndex, depth = position277, tokenIndex277, depth277
						if buffer[position] != rune('T') {
							goto l276
						}
						position++
						if buffer[position] != rune('R') {
							goto l276
						}
						position++
						if buffer[position] != rune('U') {
							goto l276
						}
						position++
						if buffer[position] != rune('E') {
							goto l276
						}
						position++
					}
				l277:
					if !_rules[ruleAction45]() {
						goto l276
					}
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					{
						position280, tokenIndex280, depth280 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l281
						}
						position++
						if buffer[position] != rune('a') {
							goto l281
						}
						position++
						if buffer[position] != rune('l') {
							goto l281
						}
						position++
						if buffer[position] != rune('s') {
							goto l281
						}
						position++
						if buffer[position] != rune('e') {
							goto l281
						}
						position++
						goto l280
					l281:
						position, tokenIndex, depth = position280, tokenIndex280, depth280
						if buffer[position] != rune('F') {
							goto l282
						}
						position++
						if buffer[position] != rune('a') {
							goto l282
						}
						position++
						if buffer[position] != rune('l') {
							goto l282
						}
						position++
						if buffer[position] != rune('s') {
							goto l282
						}
						position++
						if buffer[position] != rune('e') {
							goto l282
						}
						position++
						goto l280
					l282:
						position, tokenIndex, depth = position280, tokenIndex280, depth280
						if buffer[position] != rune('F') {
							goto l273
						}
						position++
						if buffer[position] != rune('A') {
							goto l273
						}
						position++
						if buffer[position] != rune('L') {
							goto l273
						}
						position++
						if buffer[position] != rune('S') {
							goto l273
						}
						position++
						if buffer[position] != rune('E') {
							goto l273
						}
						position++
					}
				l280:
					if !_rules[ruleAction46]() {
						goto l273
					}
				}
			l275:
				depth--
				add(rulelBool, position274)
			}
			return true
		l273:
			position, tokenIndex, depth = position273, tokenIndex273, depth273
			return false
		},
		/* 41 lString <- <(('\'' <(('\\' '\\') / ('\\' '\'') / (!'\'' .))*> '\'' Action47) / ('"' <(('\\' '\\') / ('\\' '"') / (!'"' .))*> '"' Action48))> */
		func() bool {
			position283, tokenIndex283, depth283 := position, tokenIndex, depth
			{
				position284 := position
				depth++
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if buffer[position] != rune('\'') {
						goto l286
					}
					position++
					{
						position287 := position
						depth++
					l288:
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							{
								position290, tokenIndex290, depth290 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l291
								}
								position++
								if buffer[position] != rune('\\') {
									goto l291
								}
								position++
								goto l290
							l291:
								position, tokenIndex, depth = position290, tokenIndex290, depth290
								if buffer[position] != rune('\\') {
									goto l292
								}
								position++
								if buffer[position] != rune('\'') {
									goto l292
								}
								position++
								goto l290
							l292:
								position, tokenIndex, depth = position290, tokenIndex290, depth290
								{
									position293, tokenIndex293, depth293 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l293
									}
									position++
									goto l289
								l293:
									position, tokenIndex, depth = position293, tokenIndex293, depth293
								}
								if !matchDot() {
									goto l289
								}
							}
						l290:
							goto l288
						l289:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
						}
						depth--
						add(rulePegText, position287)
					}
					if buffer[position] != rune('\'') {
						goto l286
					}
					position++
					if !_rules[ruleAction47]() {
						goto l286
					}
					goto l285
				l286:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
					if buffer[position] != rune('"') {
						goto l283
					}
					position++
					{
						position294 := position
						depth++
					l295:
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							{
								position297, tokenIndex297, depth297 := position, tokenIndex, depth
								if buffer[position] != rune('\\') {
									goto l298
								}
								position++
								if buffer[position] != rune('\\') {
									goto l298
								}
								position++
								goto l297
							l298:
								position, tokenIndex, depth = position297, tokenIndex297, depth297
								if buffer[position] != rune('\\') {
									goto l299
								}
								position++
								if buffer[position] != rune('"') {
									goto l299
								}
								position++
								goto l297
							l299:
								position, tokenIndex, depth = position297, tokenIndex297, depth297
								{
									position300, tokenIndex300, depth300 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l300
									}
									position++
									goto l296
								l300:
									position, tokenIndex, depth = position300, tokenIndex300, depth300
								}
								if !matchDot() {
									goto l296
								}
							}
						l297:
							goto l295
						l296:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
						}
						depth--
						add(rulePegText, position294)
					}
					if buffer[position] != rune('"') {
						goto l283
					}
					position++
					if !_rules[ruleAction48]() {
						goto l283
					}
				}
			l285:
				depth--
				add(rulelString, position284)
			}
			return true
		l283:
			position, tokenIndex, depth = position283, tokenIndex283, depth283
			return false
		},
		/* 42 lNull <- <((('n' 'u' 'l' 'l') / ('N' 'u' 'l' 'l') / ('N' 'U' 'L' 'L')) Action49)> */
		func() bool {
			position301, tokenIndex301, depth301 := position, tokenIndex, depth
			{
				position302 := position
				depth++
				{
					position303, tokenIndex303, depth303 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l304
					}
					position++
					if buffer[position] != rune('u') {
						goto l304
					}
					position++
					if buffer[position] != rune('l') {
						goto l304
					}
					position++
					if buffer[position] != rune('l') {
						goto l304
					}
					position++
					goto l303
				l304:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('N') {
						goto l305
					}
					position++
					if buffer[position] != rune('u') {
						goto l305
					}
					position++
					if buffer[position] != rune('l') {
						goto l305
					}
					position++
					if buffer[position] != rune('l') {
						goto l305
					}
					position++
					goto l303
				l305:
					position, tokenIndex, depth = position303, tokenIndex303, depth303
					if buffer[position] != rune('N') {
						goto l301
					}
					position++
					if buffer[position] != rune('U') {
						goto l301
					}
					position++
					if buffer[position] != rune('L') {
						goto l301
					}
					position++
					if buffer[position] != rune('L') {
						goto l301
					}
					position++
				}
			l303:
				if !_rules[ruleAction49]() {
					goto l301
				}
				depth--
				add(rulelNull, position302)
			}
			return true
		l301:
			position, tokenIndex, depth = position301, tokenIndex301, depth301
			return false
		},
		/* 43 regex <- <(('\\' '\\') / ('\\' '/') / (!'/' .))*> */
		func() bool {
			{
				position307 := position
				depth++
			l308:
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					{
						position310, tokenIndex310, depth310 := position, tokenIndex, depth
						if buffer[position] != rune('\\') {
							goto l311
						}
						position++
						if buffer[position] != rune('\\') {
							goto l311
						}
						position++
						goto l310
					l311:
						position, tokenIndex, depth = position310, tokenIndex310, depth310
						if buffer[position] != rune('\\') {
							goto l312
						}
						position++
						if buffer[position] != rune('/') {
							goto l312
						}
						position++
						goto l310
					l312:
						position, tokenIndex, depth = position310, tokenIndex310, depth310
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('/') {
								goto l313
							}
							position++
							goto l309
						l313:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
						}
						if !matchDot() {
							goto l309
						}
					}
				l310:
					goto l308
				l309:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
				}
				depth--
				add(ruleregex, position307)
			}
			return true
		},
		/* 44 squareBracketStart <- <('[' space)> */
		func() bool {
			position314, tokenIndex314, depth314 := position, tokenIndex, depth
			{
				position315 := position
				depth++
				if buffer[position] != rune('[') {
					goto l314
				}
				position++
				if !_rules[rulespace]() {
					goto l314
				}
				depth--
				add(rulesquareBracketStart, position315)
			}
			return true
		l314:
			position, tokenIndex, depth = position314, tokenIndex314, depth314
			return false
		},
		/* 45 squareBracketEnd <- <(space ']')> */
		func() bool {
			position316, tokenIndex316, depth316 := position, tokenIndex, depth
			{
				position317 := position
				depth++
				if !_rules[rulespace]() {
					goto l316
				}
				if buffer[position] != rune(']') {
					goto l316
				}
				position++
				depth--
				add(rulesquareBracketEnd, position317)
			}
			return true
		l316:
			position, tokenIndex, depth = position316, tokenIndex316, depth316
			return false
		},
		/* 46 scriptStart <- <('(' space)> */
		func() bool {
			position318, tokenIndex318, depth318 := position, tokenIndex, depth
			{
				position319 := position
				depth++
				if buffer[position] != rune('(') {
					goto l318
				}
				position++
				if !_rules[rulespace]() {
					goto l318
				}
				depth--
				add(rulescriptStart, position319)
			}
			return true
		l318:
			position, tokenIndex, depth = position318, tokenIndex318, depth318
			return false
		},
		/* 47 scriptEnd <- <(space ')')> */
		func() bool {
			position320, tokenIndex320, depth320 := position, tokenIndex, depth
			{
				position321 := position
				depth++
				if !_rules[rulespace]() {
					goto l320
				}
				if buffer[position] != rune(')') {
					goto l320
				}
				position++
				depth--
				add(rulescriptEnd, position321)
			}
			return true
		l320:
			position, tokenIndex, depth = position320, tokenIndex320, depth320
			return false
		},
		/* 48 filterStart <- <('?' '(' space)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				if buffer[position] != rune('?') {
					goto l322
				}
				position++
				if buffer[position] != rune('(') {
					goto l322
				}
				position++
				if !_rules[rulespace]() {
					goto l322
				}
				depth--
				add(rulefilterStart, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 49 filterEnd <- <(space ')')> */
		func() bool {
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				if !_rules[rulespace]() {
					goto l324
				}
				if buffer[position] != rune(')') {
					goto l324
				}
				position++
				depth--
				add(rulefilterEnd, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 50 subQueryStart <- <('(' space)> */
		func() bool {
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				if buffer[position] != rune('(') {
					goto l326
				}
				position++
				if !_rules[rulespace]() {
					goto l326
				}
				depth--
				add(rulesubQueryStart, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 51 subQueryEnd <- <(space ')')> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				if !_rules[rulespace]() {
					goto l328
				}
				if buffer[position] != rune(')') {
					goto l328
				}
				position++
				depth--
				add(rulesubQueryEnd, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 52 space <- <(' ' / '\t')*> */
		func() bool {
			{
				position331 := position
				depth++
			l332:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					{
						position334, tokenIndex334, depth334 := position, tokenIndex, depth
						if buffer[position] != rune(' ') {
							goto l335
						}
						position++
						goto l334
					l335:
						position, tokenIndex, depth = position334, tokenIndex334, depth334
						if buffer[position] != rune('\t') {
							goto l333
						}
						position++
					}
				l334:
					goto l332
				l333:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
				}
				depth--
				add(rulespace, position331)
			}
			return true
		},
		/* 54 Action0 <- <{
		    p.root = p.pop().(syntaxNode)
		}> */
		func() bool {
			{
				add(ruleAction0, position)
			}
			return true
		},
		nil,
		/* 56 Action1 <- <{
		    p.syntaxErr(begin, msgErrorInvalidSyntaxUnrecognizedInput, buffer)
		}> */
		func() bool {
			{
				add(ruleAction1, position)
			}
			return true
		},
		/* 57 Action2 <- <{
		    p.saveParams()
		}> */
		func() bool {
			{
				add(ruleAction2, position)
			}
			return true
		},
		/* 58 Action3 <- <{
		    p.setNodeChain()
		    p.setRecursiveMultiValue()
		    p.loadParams()
		}> */
		func() bool {
			{
				add(ruleAction3, position)
			}
			return true
		},
		/* 59 Action4 <- <{
		    if len(p.paramsList) == 0 {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxUseBeginAtsign, buffer)
		    }
		}> */
		func() bool {
			{
				add(ruleAction4, position)
			}
			return true
		},
		/* 60 Action5 <- <{
		    if len(p.paramsList) != 0 {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxOmitDollar, buffer)
		    }
		}> */
		func() bool {
			{
				add(ruleAction5, position)
			}
			return true
		},
		/* 61 Action6 <- <{
		    p.pushRecursiveChildIdentifier(p.pop().(syntaxNode))
		}> */
		func() bool {
			{
				add(ruleAction6, position)
			}
			return true
		},
		/* 62 Action7 <- <{
		    p.setLastNodeText(text)
		}> */
		func() bool {
			{
				add(ruleAction7, position)
			}
			return true
		},
		/* 63 Action8 <- <{
		    funcName := p.pop().(string)
		    p.pushFunction(text, funcName)
		}> */
		func() bool {
			{
				add(ruleAction8, position)
			}
			return true
		},
		/* 64 Action9 <- <{
		    p.push(text)
		}> */
		func() bool {
			{
				add(ruleAction9, position)
			}
			return true
		},
		/* 65 Action10 <- <{
		    p.setLastNodeText(text)
		}> */
		func() bool {
			{
				add(ruleAction10, position)
			}
			return true
		},
		/* 66 Action11 <- <{
		    p.pushRootIdentifier()
		}> */
		func() bool {
			{
				add(ruleAction11, position)
			}
			return true
		},
		/* 67 Action12 <- <{
		    p.pushCurrentRootIdentifier()
		}> */
		func() bool {
			{
				add(ruleAction12, position)
			}
			return true
		},
		/* 68 Action13 <- <{
		    unescapedText := p.unescape(text)
		    if unescapedText == `*` {
		        p.pushChildAsteriskIdentifier(unescapedText)
		    } else {
		        p.pushChildSingleIdentifier(unescapedText)
		    }
		}> */
		func() bool {
			{
				add(ruleAction13, position)
			}
			return true
		},
		/* 69 Action14 <- <{
		    identifier := p.pop().([]string)
		    if len(identifier) > 1 {
		        p.pushChildMultiIdentifier(identifier)
		    } else {
		        p.pushChildSingleIdentifier(identifier[0])
		    }
		}> */
		func() bool {
			{
				add(ruleAction14, position)
			}
			return true
		},
		/* 70 Action15 <- <{
		    p.push([]string{p.pop().(string)})
		}> */
		func() bool {
			{
				add(ruleAction15, position)
			}
			return true
		},
		/* 71 Action16 <- <{
		    identifier2 := p.pop().([]string)
		    identifier1 := p.pop().([]string)
		    identifier1 = append(identifier1, identifier2...)
		    p.push(identifier1)
		}> */
		func() bool {
			{
				add(ruleAction16, position)
			}
			return true
		},
		/* 72 Action17 <- <{
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction17, position)
			}
			return true
		},
		/* 73 Action18 <- <{ // '
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction18, position)
			}
			return true
		},
		/* 74 Action19 <- <{
		    subscript := p.pop().(syntaxSubscript)
		    p.pushUnionQualifier(subscript)
		}> */
		func() bool {
			{
				add(ruleAction19, position)
			}
			return true
		},
		/* 75 Action20 <- <{
		    childIndexUnion := p.pop().(*syntaxUnionQualifier)
		    parentIndexUnion := p.pop().(*syntaxUnionQualifier)
		    parentIndexUnion.merge(childIndexUnion)
		    parentIndexUnion.setMultiValue()
		    p.push(parentIndexUnion)
		}> */
		func() bool {
			{
				add(ruleAction20, position)
			}
			return true
		},
		/* 76 Action21 <- <{
		    step  := p.pop().(*syntaxIndexSubscript)
		    end   := p.pop().(*syntaxIndexSubscript)
		    start := p.pop().(*syntaxIndexSubscript)

		    if step.isOmitted || step.number == 0 {
		        step.number = 1
		    }

		    if step.number > 0 {
		        p.pushSliceSubscript(true, start, end, step)
		    } else {
		        p.pushSliceSubscript(false, start, end, step)
		    }
		}> */
		func() bool {
			{
				add(ruleAction21, position)
			}
			return true
		},
		/* 77 Action22 <- <{
		    p.pushIndexSubscript(text, false)
		}> */
		func() bool {
			{
				add(ruleAction22, position)
			}
			return true
		},
		/* 78 Action23 <- <{
		    p.pushAsteriskSubscript()
		}> */
		func() bool {
			{
				add(ruleAction23, position)
			}
			return true
		},
		/* 79 Action24 <- <{
		    p.pushIndexSubscript(`1`, false)
		}> */
		func() bool {
			{
				add(ruleAction24, position)
			}
			return true
		},
		/* 80 Action25 <- <{
		    if len(text) > 0 {
		        p.pushIndexSubscript(text, false)
		    } else {
		        p.pushIndexSubscript(`0`, true)
		    }
		}> */
		func() bool {
			{
				add(ruleAction25, position)
			}
			return true
		},
		/* 81 Action26 <- <{
		    p.pushScriptQualifier(text)
		}> */
		func() bool {
			{
				add(ruleAction26, position)
			}
			return true
		},
		/* 82 Action27 <- <{
		    p.pushFilterQualifier(p.pop().(syntaxQuery))
		}> */
		func() bool {
			{
				add(ruleAction27, position)
			}
			return true
		},
		/* 83 Action28 <- <{
		    childQuery := p.pop().(syntaxQuery)
		    parentQuery := p.pop().(syntaxQuery)
		    p.pushLogicalOr(parentQuery, childQuery)
		}> */
		func() bool {
			{
				add(ruleAction28, position)
			}
			return true
		},
		/* 84 Action29 <- <{
		    childQuery := p.pop().(syntaxQuery)
		    parentQuery := p.pop().(syntaxQuery)
		    p.pushLogicalAnd(parentQuery, childQuery)
		}> */
		func() bool {
			{
				add(ruleAction29, position)
			}
			return true
		},
		/* 85 Action30 <- <{
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
				add(ruleAction30, position)
			}
			return true
		},
		/* 86 Action31 <- <{
		    p.push(len(text) > 0 && text[0:1] == `!`)
		}> */
		func() bool {
			{
				add(ruleAction31, position)
			}
			return true
		},
		/* 87 Action32 <- <{
		    _ = p.pop().(bool)
		    jsonpathFilter := p.pop().(syntaxQuery)
		    isLogicalNot := p.pop().(bool)
		    if isLogicalNot {
		        p.pushLogicalNot(jsonpathFilter)
		    } else {
		        p.push(jsonpathFilter)
		    }
		}> */
		func() bool {
			{
				add(ruleAction32, position)
			}
			return true
		},
		/* 88 Action33 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.pushCompareEQ(leftParam, rightParam)
		}> */
		func() bool {
			{
				add(ruleAction33, position)
			}
			return true
		},
		/* 89 Action34 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.pushCompareNE(leftParam, rightParam)
		}> */
		func() bool {
			{
				add(ruleAction34, position)
			}
			return true
		},
		/* 90 Action35 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.pushCompareGE(leftParam, rightParam)
		}> */
		func() bool {
			{
				add(ruleAction35, position)
			}
			return true
		},
		/* 91 Action36 <- <{
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
				add(ruleAction36, position)
			}
			return true
		},
		/* 92 Action37 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.pushCompareLE(leftParam, rightParam)
		}> */
		func() bool {
			{
				add(ruleAction37, position)
			}
			return true
		},
		/* 93 Action38 <- <{
		    rightParam := p.pop().(*syntaxBasicCompareParameter)
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.pushCompareLT(leftParam, rightParam)
		}> */
		func() bool {
			{
				add(ruleAction38, position)
			}
			return true
		},
		/* 94 Action39 <- <{
		    leftParam := p.pop().(*syntaxBasicCompareParameter)
		    p.pushCompareRegex(leftParam, text)
		}> */
		func() bool {
			{
				add(ruleAction39, position)
			}
			return true
		},
		/* 95 Action40 <- <{
		    p.pushCompareParameterLiteral(p.pop())
		}> */
		func() bool {
			{
				add(ruleAction40, position)
			}
			return true
		},
		/* 96 Action41 <- <{
		    p.pushCompareParameterLiteral(p.pop())
		}> */
		func() bool {
			{
				add(ruleAction41, position)
			}
			return true
		},
		/* 97 Action42 <- <{
		    isLiteral := p.pop().(bool)
		    param := p.pop().(syntaxQueryParameter)
		    if !p.hasErr() && param.isMultiValueParameter() {
		        p.syntaxErr(begin, msgErrorInvalidSyntaxFilterValueGroup, buffer)
		    }
		    p.pushBasicCompareParameter(param, isLiteral)
		}> */
		func() bool {
			{
				add(ruleAction42, position)
			}
			return true
		},
		/* 98 Action43 <- <{
		    node := p.pop().(syntaxNode)

		    switch node.(type) {
		    case *syntaxRootIdentifier:
		        p.pushCompareParameterRoot(node)
		        p.push(true)
		    case *syntaxCurrentRootIdentifier:
		        p.pushCompareParameterCurrentRoot(node)
		        p.push(false)
		    default:
		        p.push(&syntaxQueryParamRoot{})
		        p.push(true)
		    }
		}> */
		func() bool {
			{
				add(ruleAction43, position)
			}
			return true
		},
		/* 99 Action44 <- <{
		    p.push(p.toFloat(text))
		}> */
		func() bool {
			{
				add(ruleAction44, position)
			}
			return true
		},
		/* 100 Action45 <- <{
		    p.push(true)
		}> */
		func() bool {
			{
				add(ruleAction45, position)
			}
			return true
		},
		/* 101 Action46 <- <{
		    p.push(false)
		}> */
		func() bool {
			{
				add(ruleAction46, position)
			}
			return true
		},
		/* 102 Action47 <- <{
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction47, position)
			}
			return true
		},
		/* 103 Action48 <- <{ // '
		    p.push(p.unescape(text))
		}> */
		func() bool {
			{
				add(ruleAction48, position)
			}
			return true
		},
		/* 104 Action49 <- <{
		    p.push(nil)
		}> */
		func() bool {
			{
				add(ruleAction49, position)
			}
			return true
		},
	}
	p.rules = _rules
}
