package syntax

import (
	"encoding/json"
	"regexp"
	"strconv"

	"github.com/AsaiYusuke/jsonpath/errors"
)

type jsonPathParser struct {
	root               syntaxNode
	paramsList         [][]any
	params             []any
	unescapeRegex      *regexp.Regexp
	filterFunctions    map[string]func(any) (any, error)
	aggregateFunctions map[string]func([]any) (any, error)
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

func (p *jsonPathParser) push(param any) {
	p.params = append(p.params, param)
}

func (p *jsonPathParser) pop() any {
	var param any
	param, p.params = p.params[len(p.params)-1], p.params[:len(p.params)-1]
	return param
}

func (p *jsonPathParser) toInt(text string) int {
	value, err := strconv.Atoi(text)
	if err != nil {
		panic(errors.NewErrorInvalidArgument(text, err))
	}
	return value
}

func (p *jsonPathParser) toFloat(text string) float64 {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		panic(errors.NewErrorInvalidArgument(text, err))
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
		panic(errors.NewErrorInvalidArgument(text, err))
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
		panic(errors.NewErrorInvalidArgument(text, err))
	}

	return unescapedText
}

func (p *jsonPathParser) _unescapeJSONString(input []byte) (string, error) {
	var unescapedText string
	err := json.Unmarshal(input, &unescapedText)
	return unescapedText, err
}

func (p *jsonPathParser) syntaxErr(pos int, reason string, buffer string) error {
	return errors.NewErrorInvalidSyntax(pos, reason, buffer[pos:])
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
		p.params = []any{root}
	}
}

func (p *jsonPathParser) setConnectedPath(targetNode syntaxNode, postfix ...string) {
	appendPath := ``
	if targetNode.getNext() != nil {
		p.setConnectedPath(targetNode.getNext(), postfix...)
		appendPath = targetNode.getNext().getRemainingPath()
	} else {
		if len(postfix) > 0 {
			appendPath = postfix[0]
		}
	}

	targetNode.setRemainingPath(targetNode.getPath() + appendPath)

	if multiIdentifier, ok := targetNode.(*syntaxChildMultiIdentifier); ok {
		if multiIdentifier.isAllWildcard {
			multiIdentifier.unionQualifier.setRemainingPath(targetNode.getRemainingPath())
		}
	}

	if aggregate, ok := targetNode.(*syntaxAggregateFunction); ok {
		p.setConnectedPath(aggregate.param, aggregate.getRemainingPath())
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

func (p *jsonPathParser) deleteRootNodeIdentifier(targetNode syntaxNode) syntaxNode {
	if aggregateFunction, ok := targetNode.(*syntaxAggregateFunction); ok {
		aggregateFunction.param = p.deleteRootNodeIdentifier(aggregateFunction.param)
		return targetNode
	}

	if targetNode.getNext() != nil {
		if targetNode.isValueGroup() {
			targetNode.getNext().setValueGroup()
		}
		targetNode.setNext(nil)
		targetNode = targetNode.getNext()
	}

	return targetNode
}

func (p *jsonPathParser) setLastNodePath(path string) {
	node := p.params[len(p.params)-1].(syntaxNode)
	node.setPath(path)

	if multiIdentifier, ok := node.(*syntaxChildMultiIdentifier); ok {
		if multiIdentifier.isAllWildcard {
			multiIdentifier.unionQualifier.setPath(path)
		}
	}
}

func (p *jsonPathParser) updateAccessorMode(checkNode syntaxNode, mode bool) {
	for checkNode != nil {
		checkNode.setAccessorMode(mode)
		checkNode = checkNode.getNext()
	}
}

func (p *jsonPathParser) pushFunction(path string, funcName string) {
	if function, ok := p.filterFunctions[funcName]; ok {
		functionNode := syntaxFilterFunction{
			syntaxBasicNode: &syntaxBasicNode{
				path:         path,
				accessorMode: p.accessorMode,
			},
			function: function,
		}

		p.push(&functionNode)
		return
	}
	if function, ok := p.aggregateFunctions[funcName]; ok {
		functionNode := syntaxAggregateFunction{
			syntaxBasicNode: &syntaxBasicNode{
				path:         path,
				accessorMode: p.accessorMode,
			},
			function: function,
		}

		p.push(&functionNode)
		return
	}

	panic(errors.NewErrorFunctionNotFound(path))
}

func (p *jsonPathParser) pushRootNodeIdentifier() {
	p.push(&syntaxRootNodeIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			path:         `$`,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushCurrentNodeIdentifier() {
	p.push(&syntaxCurrentNodeIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			path:         `@`,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushChildSingleIdentifier(path string) {
	identifier := syntaxChildSingleIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			path:         path,
			valueGroup:   false,
			accessorMode: p.accessorMode,
		},
		identifier: path,
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

	}

	p.push(&identifier)
}

func (p *jsonPathParser) pushChildWildcardIdentifier() {
	identifier := syntaxChildWildcardIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			path:         `*`,
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
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
			path:         `..`,
			valueGroup:   true,
			next:         node,
			accessorMode: p.accessorMode,
		},
		nextMapRequired:  nextMapRequired,
		nextListRequired: nextListRequired,
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

	p.push(&qualifier)
}

func (p *jsonPathParser) pushScriptQualifier(text string) {
	panic(errors.NewErrorNotSupported("script", "[("+text+")]"))
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

func (p *jsonPathParser) pushOmittedIndexSubscript() {
	p._pushIndexSubscript(`0`, true)
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

func (p *jsonPathParser) _createCompareQuery(
	leftParam, rightParam syntaxCompareParameter,
	comparator syntaxComparator) syntaxQuery {

	return &syntaxCompareQuery{
		leftParam:  leftParam,
		rightParam: rightParam,
		comparator: comparator,
	}
}

func (p *jsonPathParser) pushCompareEQ(
	leftParam, rightParam syntaxCompareParameter) {
	if isLiteralParam(leftParam) {
		rightParam, leftParam = leftParam, rightParam
	}

	if rightLiteralParam, ok := rightParam.(*syntaxQueryParamLiteral); ok {
		switch rightLiteralParam.literal[0].(type) {
		case float64:
			p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{}))
		case bool:
			p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{}))
		case string:
			p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{}))
		case nil:
			p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareDirectEQ{}))
		}

		return
	}

	p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareDeepEQ{}))
}

func (p *jsonPathParser) pushCompareNE(
	leftParam, rightParam syntaxCompareParameter) {
	p.pushCompareEQ(leftParam, rightParam)
	p.push(&syntaxLogicalNot{query: p.pop().(syntaxQuery)})
}

func (p *jsonPathParser) pushCompareGE(
	leftParam, rightParam syntaxCompareParameter) {
	if isLiteralParam(leftParam) {
		p.push(p._createCompareQuery(rightParam, leftParam, &syntaxCompareLE{}))
		return
	}
	p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareGE{}))
}

func (p *jsonPathParser) pushCompareGT(
	leftParam, rightParam syntaxCompareParameter) {
	if isLiteralParam(leftParam) {
		p.push(p._createCompareQuery(rightParam, leftParam, &syntaxCompareLT{}))
		return
	}
	p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareGT{}))
}

func (p *jsonPathParser) pushCompareLE(
	leftParam, rightParam syntaxCompareParameter) {
	if isLiteralParam(leftParam) {
		p.push(p._createCompareQuery(rightParam, leftParam, &syntaxCompareGE{}))
		return
	}
	p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareLE{}))
}

func (p *jsonPathParser) pushCompareLT(
	leftParam, rightParam syntaxCompareParameter) {
	if isLiteralParam(leftParam) {
		p.push(p._createCompareQuery(rightParam, leftParam, &syntaxCompareGT{}))
		return
	}
	p.push(p._createCompareQuery(leftParam, rightParam, &syntaxCompareLT{}))
}

func (p *jsonPathParser) pushCompareRegex(
	leftParam syntaxCompareParameter, regex string) {
	regexParam, err := regexp.Compile(regex)
	if err != nil {
		panic(errors.NewErrorInvalidArgument(regex, err))
	}

	p.push(p._createCompareQuery(
		leftParam, &syntaxQueryParamLiteral{
			literal: []any{`regex`},
		},
		&syntaxCompareRegex{
			regex: regexParam,
		}))
}

func (p *jsonPathParser) pushCompareParameterLiteral(text any) {
	p.push(
		&syntaxQueryParamLiteral{
			literal: []any{text},
		})
}

func (p *jsonPathParser) pushCompareParameterRoot(node syntaxNode) {
	p.updateAccessorMode(node, false)
	if _, ok := node.(*syntaxRootNodeIdentifier); ok {
		// Fast path: parameter is the root node '$' itself.
		p.push(&syntaxQueryParamRootNode{
			param: node,
		})
		return
	}
	p.push(&syntaxQueryParamRootNodePath{
		param: node,
	})
}

func (p *jsonPathParser) pushCompareParameterCurrentNode(node syntaxNode) {
	p.updateAccessorMode(node, false)
	if _, ok := node.(*syntaxCurrentNodeIdentifier); ok {
		// Fast path: parameter is the current node '@' itself.
		p.push(&syntaxQueryParamCurrentNode{
			param: node,
		})
		return
	}
	p.push(&syntaxQueryParamCurrentNodePath{
		param: node,
	})
}
