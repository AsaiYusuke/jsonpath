package jsonpath

import (
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
			}

			last.setNext(nextNode)

			last = nextNode
		}
		p.params = []interface{}{root}
	}
}

func (p *jsonPathParser) setConnectedText(targetNode syntaxNode) {
	childText := ``
	if targetNode.getNext() != nil {
		p.setConnectedText(targetNode.getNext())
		childText = targetNode.getNext().getConnectedText()
	}
	targetNode.setConnectedText(targetNode.getText() + childText)
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
	_, isRootIdentifier := targetNode.(*syntaxRootIdentifier)
	_, isCurrentRootIdentifier := targetNode.(*syntaxCurrentRootIdentifier)
	if isRootIdentifier || isCurrentRootIdentifier {
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
}

func (p *jsonPathParser) updateAccessorMode(checkNode syntaxNode, mode bool) {
	for checkNode != nil {
		checkNode.setAccessorMode(mode)
		checkNode = checkNode.getNext()
	}
}

func (p *jsonPathParser) pushFunction(text string, funcName string) {
	if function, ok := p.filterFunctions[funcName]; ok {
		p.push(&syntaxFilterFunction{
			syntaxBasicNode: &syntaxBasicNode{
				text:         text,
				accessorMode: p.accessorMode,
			},
			function: function,
		})
		return
	}
	if function, ok := p.aggregateFunctions[funcName]; ok {
		p.push(&syntaxAggregateFunction{
			syntaxBasicNode: &syntaxBasicNode{
				text:         text,
				accessorMode: p.accessorMode,
			},
			function: function,
		})
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
	p.push(&syntaxChildSingleIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         text,
			valueGroup:   false,
			accessorMode: p.accessorMode,
		},
		identifier: text,
	})
}

func (p *jsonPathParser) pushChildMultiIdentifier(
	node syntaxNode, appendNode syntaxNode) {

	if multiIdentifier, ok := node.(*syntaxChildMultiIdentifier); ok {
		multiIdentifier.identifiers = append(multiIdentifier.identifiers, appendNode)
		p.push(multiIdentifier)
		return
	}

	p.push(&syntaxChildMultiIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
		identifiers: []syntaxNode{
			node,
			appendNode,
		},
	})
}

func (p *jsonPathParser) pushChildWildcardIdentifier() {
	p.push(&syntaxChildWildcardIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `*`,
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
	})
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
	p.push(&syntaxRecursiveChildIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `..`,
			valueGroup:   true,
			next:         node,
			accessorMode: p.accessorMode,
		},
		nextMapRequired:  nextMapRequired,
		nextListRequired: nextListRequired,
	})
}

func (p *jsonPathParser) pushUnionQualifier(subscript syntaxSubscript) {
	p.push(&syntaxUnionQualifier{
		syntaxBasicNode: &syntaxBasicNode{
			valueGroup:   subscript.isValueGroup(),
			accessorMode: p.accessorMode,
		},
		subscripts: []syntaxSubscript{subscript},
	})
}

func (p *jsonPathParser) pushFilterQualifier(query syntaxQuery) {
	p.push(&syntaxFilterQualifier{
		syntaxBasicNode: &syntaxBasicNode{
			valueGroup:   true,
			accessorMode: p.accessorMode,
		},
		query: query,
	})
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
	p.push(&syntaxLogicalNot{query: query})
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
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareEQ{}))
}

func (p *jsonPathParser) pushCompareNE(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	p.push(&syntaxLogicalNot{
		query: p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareEQ{}),
	})
}

func (p *jsonPathParser) pushCompareGE(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareGE{}))
}

func (p *jsonPathParser) pushCompareGT(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareGT{}))
}

func (p *jsonPathParser) pushCompareLE(
	leftParam, rightParam *syntaxBasicCompareParameter) {
	p.push(p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareLE{}))
}

func (p *jsonPathParser) pushCompareLT(
	leftParam, rightParam *syntaxBasicCompareParameter) {
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
	param := &syntaxQueryParamRoot{
		param: node,
	}
	p.updateAccessorMode(param.param, false)
	p.push(param)
}

func (p *jsonPathParser) pushCompareParameterCurrentRoot(node syntaxNode) {
	param := &syntaxQueryParamCurrentRoot{
		param: node,
	}
	p.updateAccessorMode(param.param, false)
	p.push(param)
}
