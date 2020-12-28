package jsonpath

import (
	"regexp"
	"strconv"
)

type jsonPathParser struct {
	root               syntaxNode
	paramsList         [][]interface{}
	params             []interface{}
	thisError          error
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
		p.thisError = ErrorInvalidArgument{text, err}
		return 0
	}
	return value
}

func (p *jsonPathParser) toFloat(text string) float64 {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		p.thisError = ErrorInvalidArgument{text, err}
		return 0
	}
	return value
}

func (p *jsonPathParser) unescape(text string) string {
	return p.unescapeRegex.ReplaceAllStringFunc(text, func(block string) string {
		varBlockSet := p.unescapeRegex.FindStringSubmatch(block)
		return varBlockSet[1]
	})
}

func (p *jsonPathParser) updateAccessorMode(checkNode syntaxNode, mode bool) {
	for checkNode != nil {
		checkNode.setAccessorMode(mode)
		checkNode = checkNode.getNext()
	}
}

func (p *jsonPathParser) syntaxErr(pos int, reason string, buffer string) {
	p.thisError = ErrorInvalidSyntax{pos, reason, buffer[pos:]}
}

func (p *jsonPathParser) hasErr() bool {
	return p.thisError != nil
}

func (p *jsonPathParser) setNodeChain() {
	if len(p.params) > 1 {
		root := p.params[0].(syntaxNode)
		last := root
		for _, next := range p.params[1:] {
			switch next.(type) {
			case *syntaxAggregateFunction:
				funcNode := next.(*syntaxAggregateFunction)
				funcNode.param = root
				root = funcNode
				last = root
			default:
				nextNode := next.(syntaxNode)
				last.setNext(nextNode)
				last = nextNode
			}
		}
		p.params = []interface{}{root}
	}
}

func (p *jsonPathParser) setLastNodeText(text string) {
	node := p.params[len(p.params)-1].(syntaxNode)
	node.setText(text)
}

func (p *jsonPathParser) setRecursiveMultiValue() {
	node := p.params[0].(syntaxNode)
	checkNode := node
	for checkNode != nil {
		if checkNode.isMultiValue() {
			node.setMultiValue()
			break
		}
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

	p.thisError = ErrorFunctionNotFound{
		function: text,
	}
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
		identifier: text,
		syntaxBasicNode: &syntaxBasicNode{
			text:         text,
			multiValue:   false,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushChildMultiIdentifier(identifiers []string) {
	p.push(&syntaxChildMultiIdentifier{
		identifiers: identifiers,
		syntaxBasicNode: &syntaxBasicNode{
			multiValue:   true,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushChildAsteriskIdentifier(text string) {
	p.push(&syntaxChildAsteriskIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         text,
			multiValue:   true,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushRecursiveChildIdentifier(node syntaxNode) {
	p.push(&syntaxRecursiveChildIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:         `..`,
			multiValue:   true,
			next:         node,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushUnionQualifier(subscript syntaxSubscript) {
	p.push(&syntaxUnionQualifier{
		syntaxBasicNode: &syntaxBasicNode{
			multiValue:   subscript.isMultiValue(),
			accessorMode: p.accessorMode,
		},
		subscripts: []syntaxSubscript{subscript},
	})
}

func (p *jsonPathParser) pushFilterQualifier(query syntaxQuery) {
	p.push(&syntaxFilterQualifier{
		query: query,
		syntaxBasicNode: &syntaxBasicNode{
			multiValue:   true,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushScriptQualifier(text string) {
	p.push(&syntaxScriptQualifier{
		command: text,
		syntaxBasicNode: &syntaxBasicNode{
			multiValue:   true,
			accessorMode: p.accessorMode,
		},
	})
}

func (p *jsonPathParser) pushSlicePositiveStepSubscript(start, end, step *syntaxIndexSubscript) {
	p.push(&syntaxSlicePositiveStepSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: true,
		},
		start: start,
		end:   end,
		step:  step,
	})
}

func (p *jsonPathParser) pushSliceNegativeStepSubscript(start, end, step *syntaxIndexSubscript) {
	p.push(&syntaxSliceNegativeStepSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: true,
		},
		start: start,
		end:   end,
		step:  step,
	})
}

func (p *jsonPathParser) _pushIndexSubscript(text string, isOmitted bool) {
	p.push(&syntaxIndexSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: false,
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

func (p *jsonPathParser) pushAsteriskSubscript() {
	p.push(&syntaxAsteriskSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: true,
		},
	})
}

func (p *jsonPathParser) pushLogicalOr(query1, query2 syntaxQuery) {
	p.push(&syntaxLogicalOr{query1, query2})
}

func (p *jsonPathParser) pushLogicalAnd(query1, query2 syntaxQuery) {
	p.push(&syntaxLogicalAnd{query1, query2})
}

func (p *jsonPathParser) pushLogicalNot(jsonpathFilter syntaxQuery) {
	p.push(&syntaxLogicalNot{param: jsonpathFilter})
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
		param: p._createBasicCompareQuery(leftParam, rightParam, &syntaxCompareEQ{}),
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
		p.thisError = ErrorInvalidArgument{regex, err}
	}

	p.push(p._createBasicCompareQuery(
		leftParam, &syntaxBasicCompareParameter{
			param:     &syntaxQueryParamLiteral{literal: `regex`},
			isLiteral: true,
		},
		&syntaxCompareRegex{
			regex: regexParam,
		}))
}

func (p *jsonPathParser) pushBasicCompareParameter(
	parameter syntaxQueryParameter, isLiteral bool) {
	p.push(&syntaxBasicCompareParameter{
		param:     parameter,
		isLiteral: isLiteral,
	})
}

func (p *jsonPathParser) pushCompareParameterLiteral(text interface{}) {
	p.pushBasicCompareParameter(
		&syntaxQueryParamLiteral{text}, true)
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
