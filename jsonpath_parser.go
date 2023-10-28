package jsonpath

import (
	"encoding/json"
	"regexp"
	"strconv"
)

type jsonPathParser struct {
	root               syntaxNode
	paramsList         [][]interface{}
	params             []interface{}
	unescapeRegex      *regexp.Regexp
	filterFunctions    map[string]func(interface{}) (interface{}, error)
	aggregateFunctions map[string]func([]interface{}) (interface{}, error)
	accessorMode       bool
}

func (p *jsonPathParser) saveParams() {
	if len(p.params) > 0 {
		p.paramsList = append(p.paramsList, p.params)
		p.params = nil
	}
}

func (p *jsonPathParser) loadParams() {
	if len(p.paramsList) > 0 {
		p.params = append(p.paramsList[len(p.paramsList)-1], p.params...)
		p.paramsList = p.paramsList[:len(p.paramsList)-1]
	}
}

func (p *jsonPathParser) push(param interface{}) {
	p.params = append(p.params, param)
}

func (p *jsonPathParser) pop() interface{} {
	var param interface{}
	param, p.params = p.params[len(p.params)-1], p.params[:len(p.params)-1]
	return param
}

func (p *jsonPathParser) toInt(text string) int {
	value, err := strconv.Atoi(text)
	if err != nil {
		panic(ErrorInvalidArgument{
			argument: text,
			err:      err,
		})
	}
	return value
}

func (p *jsonPathParser) toFloat(text string) float64 {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		panic(ErrorInvalidArgument{
			argument: text,
			err:      err,
		})
	}
	return value
}

func (p *jsonPathParser) unescape(text string) string {
	return p.unescapeRegex.ReplaceAllStringFunc(text, func(block string) string {
		varBlockSet := p.unescapeRegex.FindStringSubmatch(block)
		return varBlockSet[1]
	})
}

func (p *jsonPathParser) unescapeSingleQuotedString(text string) string {
	srcBytes := []byte(text)
	inputBytes := make([]byte, 0, 2+len(text))

	inputBytes = append(inputBytes, 0x22) // "

	var foundEscape bool
	for index := range srcBytes {
		switch srcBytes[index] {
		case 0x22: // "
			// " -> /"
			inputBytes = append(inputBytes, 0x5c, srcBytes[index])
		case 0x27: // '
			// \' -> '
			inputBytes = append(inputBytes, srcBytes[index])
			foundEscape = false
		case 0x5c: // \
			if foundEscape {
				inputBytes = append(inputBytes, 0x5c, 0x5c)
				foundEscape = false
			} else {
				foundEscape = true
			}
		default:
			if foundEscape {
				inputBytes = append(inputBytes, 0x5c, srcBytes[index])
				foundEscape = false
			} else {
				inputBytes = append(inputBytes, srcBytes[index])
			}
		}
	}

	inputBytes = append(inputBytes, 0x22) // "

	unescapedText, err := p._unescapeJSONString(inputBytes)

	if err != nil {
		panic(ErrorInvalidArgument{
			argument: text,
			err:      err,
		})
	}

	return unescapedText
}

func (p *jsonPathParser) unescapeDoubleQuotedString(text string) string {
	inputBytes := make([]byte, 0, 2+len(text))
	inputBytes = append(inputBytes, 0x22)
	inputBytes = append(inputBytes, []byte(text)...)
	inputBytes = append(inputBytes, 0x22)

	unescapedText, err := p._unescapeJSONString(inputBytes)

	if err != nil {
		panic(ErrorInvalidArgument{
			argument: text,
			err:      err,
		})
	}

	return unescapedText
}

func (p *jsonPathParser) _unescapeJSONString(input []byte) (string, error) {
	var unescapedText string
	err := json.Unmarshal(input, &unescapedText)
	return unescapedText, err
}

func (p *jsonPathParser) syntaxErr(pos int, reason string, buffer string) error {
	return ErrorInvalidSyntax{
		position: pos,
		reason:   reason,
		near:     buffer[pos:],
	}
}

func (p *jsonPathParser) setNodeChain() {
	if len(p.params) > 1 {
		root := p.params[0].(syntaxNode)
		last := root
		for _, next := range p.params[1:] {
			if funcNode, ok := next.(*syntaxAggregateFunction); ok {
				funcNode.param = root
				p.updateAccessorMode(funcNode.param, false)
				root = funcNode
				last = root
				continue
			}

			nextNode := next.(syntaxNode)

			if multiIdentifier, ok := last.(*syntaxChildMultiIdentifier); ok {
				for _, singleIdentifier := range multiIdentifier.identifiers {
					singleIdentifier.setNext(nextNode)
				}
				if multiIdentifier.isAllWildcard {
					multiIdentifier.unionQualifier.setNext(nextNode)
				}
			}

			last.setNext(nextNode)

			last = nextNode
		}
		p.params = []interface{}{root}
	}
}

func (p *jsonPathParser) setConnectedText(targetNode syntaxNode, postfix ...string) {
	appendText := ``
	if targetNode.getNext() != nil {
		p.setConnectedText(targetNode.getNext(), postfix...)
		appendText = targetNode.getNext().getConnectedText()
	} else {
		if len(postfix) > 0 {
			appendText = postfix[0]
		}
	}

	targetNode.setConnectedText(targetNode.getText() + appendText)

	if multiIdentifier, ok := targetNode.(*syntaxChildMultiIdentifier); ok {
		if multiIdentifier.isAllWildcard {
			multiIdentifier.unionQualifier.setConnectedText(targetNode.getConnectedText())
		}
	}

	if aggregate, ok := targetNode.(*syntaxAggregateFunction); ok {
		p.setConnectedText(aggregate.param, aggregate.getConnectedText())
	}
}

func (p *jsonPathParser) updateRootValueGroup() {
	rootNode := p.params[0].(syntaxNode)
	checkNode := rootNode
	for checkNode != nil {
		if checkNode.isValueGroup() {
			rootNode.setValueGroup()
			break
		}
		checkNode = checkNode.getNext()
	}
}

func (p *jsonPathParser) deleteRootIdentifier(targetNode syntaxNode) syntaxNode {
	switch targetNode.(type) {
	case *syntaxRootIdentifier, *syntaxCurrentRootIdentifier:
		if targetNode.getNext() != nil {
			if targetNode.isValueGroup() {
				targetNode.getNext().setValueGroup()
			}
			targetNode.setNext(nil)
			targetNode = targetNode.getNext()
		}
		return targetNode
	}

	if aggregateFunction, ok := targetNode.(*syntaxAggregateFunction); ok {
		aggregateFunction.param = p.deleteRootIdentifier(aggregateFunction.param)
	}

	return targetNode
}

func (p *jsonPathParser) setLastNodeText(text string) {
	node := p.params[len(p.params)-1].(syntaxNode)
	node.setText(text)

	if multiIdentifier, ok := node.(*syntaxChildMultiIdentifier); ok {
		if multiIdentifier.isAllWildcard {
			multiIdentifier.unionQualifier.setText(text)
		}
	}
}

func (p *jsonPathParser) updateAccessorMode(checkNode syntaxNode, mode bool) {
	for checkNode != nil {
		checkNode.setAccessorMode(mode)
		checkNode = checkNode.getNext()
	}
}

func (p *jsonPathParser) pushFunction(text string, funcName string) {
	if function, ok := p.filterFunctions[funcName]; ok {
		functionNode := syntaxFilterFunction{
			syntaxBasicNode: &syntaxBasicNode{
				text:         text,
				accessorMode: p.accessorMode,
			},
			function: function,
		}

		functionNode.errorRuntime = &errorBasicRuntime{
			node: functionNode.syntaxBasicNode,
		}

		p.push(&functionNode)
		return
	}
	if function, ok := p.aggregateFunctions[funcName]; ok {
		functionNode := syntaxAggregateFunction{
			syntaxBasicNode: &syntaxBasicNode{
				text:         text,
				accessorMode: p.accessorMode,
			},
			function: function,
		}

		functionNode.errorRuntime = &errorBasicRuntime{
			node: functionNode.syntaxBasicNode,
		}

		p.push(&functionNode)
		return
	}

	panic(ErrorFunctionNotFound{
		function: text,
	})
}

func (p *jsonPathParser) pushRootIdentifier() {
	p.push(&syntaxRootIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `$`,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushCurrentRootIdentifier() {
	p.push(&syntaxCurrentRootIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `@`,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushChildSingleIdentifier(text string) {
	identifier := syntaxChildSingleIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         text,
			valueGroup:   false,
			accessorMode: p.accessorMode,
		},
		identifier: text,
	}

	identifier.errorRuntime = &errorBasicRuntime{
		node: identifier.syntaxBasicNode,
	}

	p.push(&identifier)
}

func (p *jsonPathParser) pushChildMultiIdentifier(
	node syntaxNode, appendNode syntaxNode) {

	if multiIdentifier, ok := node.(*syntaxChildMultiIdentifier); ok {
		multiIdentifier.identifiers = append(multiIdentifier.identifiers, appendNode)

		_, isWildcard := appendNode.(*syntaxChildWildcardIdentifier)
		multiIdentifier.isAllWildcard = multiIdentifier.isAllWildcard && isWildcard

		if multiIdentifier.isAllWildcard {
			multiIdentifier.unionQualifier.subscripts = append(
				multiIdentifier.unionQualifier.subscripts,
				&syntaxWildcardSubscript{},
			)
		} else {
			multiIdentifier.unionQualifier = syntaxUnionQualifier{}
		}

		p.push(multiIdentifier)
		return
	}

	_, isNodeWildcard := node.(*syntaxChildWildcardIdentifier)
	_, isAppendNodeWildcard := appendNode.(*syntaxChildWildcardIdentifier)

	identifier := syntaxChildMultiIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
		identifiers: []syntaxNode{
			node,
			appendNode,
		},
		isAllWildcard: isNodeWildcard && isAppendNodeWildcard,
	}

	identifier.errorRuntime = &errorBasicRuntime{
		node: identifier.syntaxBasicNode,
	}

	if identifier.isAllWildcard {
		identifier.unionQualifier = syntaxUnionQualifier{
			syntaxBasicNode: &syntaxBasicNode{
				valueGroup:   true,
				accessorMode: p.accessorMode,
			},
			subscripts: []syntaxSubscript{
				&syntaxWildcardSubscript{},
				&syntaxWildcardSubscript{},
			},
		}

		identifier.unionQualifier.errorRuntime = &errorBasicRuntime{
			node: identifier.unionQualifier.syntaxBasicNode,
		}
	}

	p.push(&identifier)
}

func (p *jsonPathParser) pushChildWildcardIdentifier() {
	identifier := syntaxChildWildcardIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `*`,
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
	}

	identifier.errorRuntime = &errorBasicRuntime{
		node: identifier.syntaxBasicNode,
	}

	p.push(&identifier)
}

func (p *jsonPathParser) pushRecursiveChildIdentifier(node syntaxNode) {
	var nextMapRequired, nextListRequired bool
	switch node.(type) {
	case *syntaxChildWildcardIdentifier, *syntaxChildMultiIdentifier, *syntaxFilterQualifier:
		nextMapRequired = true
		nextListRequired = true
	case *syntaxChildSingleIdentifier:
		nextMapRequired = true
	case *syntaxUnionQualifier:
		nextListRequired = true
	}

	identifier := syntaxRecursiveChildIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `..`,
			valueGroup:   true,
			next:         node,
			accessorMode: p.accessorMode,
		},
		nextMapRequired:  nextMapRequired,
		nextListRequired: nextListRequired,
	}

	identifier.errorRuntime = &errorBasicRuntime{
		node: identifier.syntaxBasicNode,
	}

	p.push(&identifier)
}

func (p *jsonPathParser) pushUnionQualifier(subscript syntaxSubscript) {
	qualifier := syntaxUnionQualifier{
		syntaxBasicNode: &syntaxBasicNode{
			valueGroup:   subscript.isValueGroup(),
			accessorMode: p.accessorMode,
		},
		subscripts: []syntaxSubscript{subscript},
	}

	qualifier.errorRuntime = &errorBasicRuntime{
		node: qualifier.syntaxBasicNode,
	}

	p.push(&qualifier)
}

func (p *jsonPathParser) pushFilterQualifier(query syntaxQuery) {
	qualifier := syntaxFilterQualifier{
		syntaxBasicNode: &syntaxBasicNode{
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
		query: query,
	}

	qualifier.errorRuntime = &errorBasicRuntime{
		node: qualifier.syntaxBasicNode,
	}

	p.push(&qualifier)
}

func (p *jsonPathParser) pushScriptQualifier(text string) {
	panic(ErrorNotSupported{
		feature: `script`,
		path:    `[(` + text + `)]`,
	})
}

func (p *jsonPathParser) pushSlicePositiveStepSubscript(start, end, step *syntaxIndexSubscript) {
	p.push(&syntaxSlicePositiveStepSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			valueGroup: true,
		},
		start: start,
		end:   end,
		step:  step,
	})
}

func (p *jsonPathParser) pushSliceNegativeStepSubscript(start, end, step *syntaxIndexSubscript) {
	p.push(&syntaxSliceNegativeStepSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			valueGroup: true,
		},
		start: start,
		end:   end,
		step:  step,
	})
}

func (p *jsonPathParser) _pushIndexSubscript(text string, isOmitted bool) {
	p.push(&syntaxIndexSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			valueGroup: false,
		},
		number:    p.toInt(text),
		isOmitted: isOmitted,
	})
}

func (p *jsonPathParser) pushIndexSubscript(text string) {
	p._pushIndexSubscript(text, false)
}

func (p *jsonPathParser) pushOmittedIndexSubscript(text string) {
	p._pushIndexSubscript(text, true)
}

func (p *jsonPathParser) pushWildcardSubscript() {
	p.push(&syntaxWildcardSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			valueGroup: true,
		},
	})
}

func (p *jsonPathParser) pushLogicalOr(leftQuery, rightQuery syntaxQuery) {
	p.push(&syntaxLogicalOr{
		leftQuery:  leftQuery,
		rightQuery: rightQuery,
	})
}

func (p *jsonPathParser) pushLogicalAnd(leftQuery, rightQuery syntaxQuery) {
	p.push(&syntaxLogicalAnd{
		leftQuery:  leftQuery,
		rightQuery: rightQuery,
	})
}

func (p *jsonPathParser) pushLogicalNot(query syntaxQuery) {
	p.push(&syntaxLogicalNot{
		query: query,
	})
}

func (p *jsonPathParser) _createBasicCompareQuery(
	leftParam, rightParam *syntaxBasicCompareParameter,
	comparator syntaxComparator) syntaxQuery {

	return &syntaxBasicCompareQuery{
		leftParam:  leftParam,
		rightParam: rightParam,
		comparator: comparator,
	}
}

func (p *jsonPathParser) pushCompareEQ(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	if leftParam.isLiteral {
		rightParam, leftParam = leftParam, rightParam
	}

	if rightLiteralParam, ok := rightParam.param.(*syntaxQueryParamLiteral); ok {
		switch rightLiteralParam.literal[0].(type) {
		case float64:
			p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{
				syntaxTypeValidator: &syntaxBasicNumericTypeValidator{},
			}))
		case bool:
			p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{
				syntaxTypeValidator: &syntaxBasicBoolTypeValidator{},
			}))
		case string:
			p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{
				syntaxTypeValidator: &syntaxBasicStringTypeValidator{},
			}))
		case nil:
			p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{
				syntaxTypeValidator: &syntaxBasicNilTypeValidator{},
			}))
		}

		return
	}

	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareDeepEQ{}))
}

func (p *jsonPathParser) pushCompareNE(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	p.pushCompareEQ(leftParam, rightParam)
	p.push(&syntaxLogicalNot{query: p.pop().(syntaxQuery)})
}

func (p *jsonPathParser) pushCompareGE(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	if leftParam.isLiteral {
		p.pushCompareLE(rightParam, leftParam)
		return
	}
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareGE{}))
}

func (p *jsonPathParser) pushCompareGT(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	if leftParam.isLiteral {
		p.pushCompareLT(rightParam, leftParam)
		return
	}
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareGT{}))
}

func (p *jsonPathParser) pushCompareLE(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	if leftParam.isLiteral {
		p.pushCompareGE(rightParam, leftParam)
		return
	}
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareLE{}))
}

func (p *jsonPathParser) pushCompareLT(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	if leftParam.isLiteral {
		p.pushCompareGT(rightParam, leftParam)
		return
	}
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareLT{}))
}

func (p *jsonPathParser) pushCompareRegex(
	leftParam *syntaxBasicCompareParameter, regex string) {
	regexParam, err := regexp.Compile(regex)
	if err != nil {
		panic(ErrorInvalidArgument{
			argument: regex,
			err:      err,
		})
	}

	p.push(p._createBasicCompareQuery(
		leftParam, &syntaxBasicCompareParameter{
			param: &syntaxQueryParamLiteral{
				literal: []interface{}{`regex`},
			},
			isLiteral: true,
		},
		&syntaxCompareRegex{
			regex: regexParam,
		}))
}

func (p *jsonPathParser) pushBasicCompareParameter(
	parameter syntaxQuery, isLiteral bool) {
	p.push(&syntaxBasicCompareParameter{
		param:     parameter,
		isLiteral: isLiteral,
	})
}

func (p *jsonPathParser) pushCompareParameterLiteral(text interface{}) {
	p.pushBasicCompareParameter(
		&syntaxQueryParamLiteral{
			literal: []interface{}{text},
		}, true)
}

func (p *jsonPathParser) pushCompareParameterRoot(node syntaxNode) {
	p.updateAccessorMode(node, false)
	p.push(&syntaxQueryParamRoot{
		param: node,
	})
}

func (p *jsonPathParser) pushCompareParameterCurrentRoot(node syntaxNode) {
	p.updateAccessorMode(node, false)
	p.push(&syntaxQueryParamCurrentRoot{
		param: node,
	})
}
