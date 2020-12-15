package jsonpath

import (
	"regexp"
	"strconv"
)

type jsonPathParser struct {
	root          syntaxNode
	srcJSON       *interface{}
	resultPtr     *[]interface{}
	params        []interface{}
	thisError     error
	unescapeRegex *regexp.Regexp
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

func (p *jsonPathParser) updateResultPtr(checkNode syntaxNode, result **[]interface{}) {
	for checkNode != nil {
		checkNode.setResultPtr(result)
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
	child := p.pop().(syntaxNode)
	parent := p.pop().(syntaxNode)
	parent.setNext(child)
	p.push(parent)
}

func (p *jsonPathParser) setNodeText(text string) {
	node := p.pop().(syntaxNode)
	node.setText(text)
	p.push(node)
}

func (p *jsonPathParser) createRootIdentifier() syntaxNode {
	return &syntaxRootIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:   `$`,
			result: &p.resultPtr,
		},
		srcJSON: &p.srcJSON,
	}
}

func (p *jsonPathParser) createCurrentRootIdentifier() syntaxNode {
	return &syntaxCurrentRootIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:   `@`,
			result: &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createChildSingleIdentifier(text string) syntaxNode {
	return &syntaxChildSingleIdentifier{
		identifier: text,
		syntaxBasicNode: &syntaxBasicNode{
			text:       text,
			multiValue: false,
			result:     &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createChildMultiIdentifier(identifiers []string) syntaxNode {
	return &syntaxChildMultiIdentifier{
		identifiers: identifiers,
		syntaxBasicNode: &syntaxBasicNode{
			multiValue: true,
			result:     &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createChildAsteriskIdentifier(text string) syntaxNode {
	return &syntaxChildAsteriskIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:       text,
			multiValue: true,
			result:     &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createRecursiveChildIdentifier(node syntaxNode) syntaxNode {
	return &syntaxRecursiveChildIdentifier{
		syntaxBasicNode: &syntaxBasicNode{
			text:       `..`,
			multiValue: true,
			next:       node,
			result:     &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createUnionQualifier(subscript syntaxSubscript) *syntaxUnionQualifier {
	return &syntaxUnionQualifier{
		syntaxBasicNode: &syntaxBasicNode{
			multiValue: subscript.isMultiValue(),
			result:     &p.resultPtr,
		},
		subscripts: []syntaxSubscript{subscript},
	}
}

func (p *jsonPathParser) createFilterQualifier(query syntaxQuery) syntaxNode {
	return &syntaxFilterQualifier{
		query: query,
		syntaxBasicNode: &syntaxBasicNode{
			multiValue: true,
			result:     &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createScriptQualifier(text string) syntaxNode {
	return &syntaxScriptQualifier{
		command: text,
		syntaxBasicNode: &syntaxBasicNode{
			multiValue: true,
			result:     &p.resultPtr,
		},
	}
}

func (p *jsonPathParser) createSliceSubscript(isPositiveStep bool, start, end, step *syntaxIndexSubscript) syntaxSubscript {
	if isPositiveStep {
		return &syntaxSlicePositiveStepSubscript{
			syntaxBasicSubscript: &syntaxBasicSubscript{
				multiValue: true,
			},
			start: start,
			end:   end,
			step:  step,
		}
	}

	return &syntaxSliceNegativeStepSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: true,
		},
		start: start,
		end:   end,
		step:  step,
	}
}

func (p *jsonPathParser) createIndexSubscript(text string, isOmitted bool) syntaxSubscript {
	return &syntaxIndexSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: false,
		},
		number:    p.toInt(text),
		isOmitted: isOmitted,
	}
}

func (p *jsonPathParser) createAsteriskSubscript() syntaxSubscript {
	return &syntaxAsteriskSubscript{
		syntaxBasicSubscript: &syntaxBasicSubscript{
			multiValue: true,
		},
	}
}

func (p *jsonPathParser) createLogicalOr(query1, query2 syntaxQuery) syntaxQuery {
	return &syntaxLogicalOr{query1, query2}
}

func (p *jsonPathParser) createLogicalAnd(query1, query2 syntaxQuery) syntaxQuery {
	return &syntaxLogicalAnd{query1, query2}
}

func (p *jsonPathParser) createLogicalNot(jsonpathFilter syntaxQuery) syntaxQuery {
	return &syntaxLogicalNot{param: jsonpathFilter}
}

func (p *jsonPathParser) createBasicCompareQuery(
	leftParam, rightParam *syntaxBasicCompareParameter,
	comparator syntaxComparator) syntaxQuery {

	return &syntaxBasicCompareQuery{
		leftParam:  leftParam,
		rightParam: rightParam,
		comparator: comparator,
	}
}

func (p *jsonPathParser) createCompareEQ(
	leftParam, rightParam *syntaxBasicCompareParameter) syntaxQuery {
	return p.createBasicCompareQuery(leftParam, rightParam, &syntaxCompareEQ{})
}

func (p *jsonPathParser) createCompareNE(
	leftParam, rightParam *syntaxBasicCompareParameter) syntaxQuery {
	return &syntaxLogicalNot{
		param: p.createBasicCompareQuery(leftParam, rightParam, &syntaxCompareEQ{}),
	}
}

func (p *jsonPathParser) createCompareGE(
	leftParam, rightParam *syntaxBasicCompareParameter) syntaxQuery {
	return p.createBasicCompareQuery(leftParam, rightParam, &syntaxCompareGE{})
}

func (p *jsonPathParser) createCompareGT(
	leftParam, rightParam *syntaxBasicCompareParameter) syntaxQuery {
	return p.createBasicCompareQuery(leftParam, rightParam, &syntaxCompareGT{})
}

func (p *jsonPathParser) createCompareLE(
	leftParam, rightParam *syntaxBasicCompareParameter) syntaxQuery {
	return p.createBasicCompareQuery(leftParam, rightParam, &syntaxCompareLE{})
}

func (p *jsonPathParser) createCompareLT(
	leftParam, rightParam *syntaxBasicCompareParameter) syntaxQuery {
	return p.createBasicCompareQuery(leftParam, rightParam, &syntaxCompareLT{})
}

func (p *jsonPathParser) createCompareRegex(
	leftParam *syntaxBasicCompareParameter, regex string) syntaxQuery {
	return p.createBasicCompareQuery(
		leftParam, &syntaxBasicCompareParameter{
			param:     &syntaxQueryParamLiteral{literal: `regex`},
			isLiteral: true,
		},
		&syntaxCompareRegex{
			regex: regexp.MustCompile(regex),
		})
}

func (p *jsonPathParser) createBasicCompareParameter(
	parameter syntaxQueryParameter, isLiteral bool) *syntaxBasicCompareParameter {
	return &syntaxBasicCompareParameter{
		param:     parameter,
		isLiteral: isLiteral,
	}
}
func (p *jsonPathParser) createCompareParameterLiteral(text interface{}) *syntaxBasicCompareParameter {
	return p.createBasicCompareParameter(
		&syntaxQueryParamLiteral{text}, true)
}

func (p *jsonPathParser) createCompareParameterRoot(node syntaxNode) syntaxQueryParameter {
	param := &syntaxQueryParamRoot{
		param:     node,
		srcJSON:   &p.srcJSON,
		resultPtr: &[]interface{}{},
	}
	p.updateResultPtr(param.param, &param.resultPtr)
	return param
}

func (p *jsonPathParser) createCompareParameterCurrentRoot(node syntaxNode) syntaxQuery {
	param := &syntaxQueryParamCurrentRoot{
		param:     node,
		resultPtr: &[]interface{}{},
	}
	p.updateResultPtr(param.param, &param.resultPtr)
	return param
}
