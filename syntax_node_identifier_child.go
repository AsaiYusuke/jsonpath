package jsonpath

import (
	"reflect"
	"sort"
)

type syntaxChildIdentifier struct {
	*syntaxBasicNode

	identifiers []string
	isAsterisk  bool
}

func (i syntaxChildIdentifier) retrieve(root, current interface{}, result *resultContainer) error {
	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		if i.isAsteriskNode() {
			return i.processAsteriskValue(root, srcMap, result)
		}

		if len(i.identifiers) > 1 {
			return i.processMultiIdentifier(root, srcMap, result)
		}

		identifier := i.identifiers[0]
		child, ok := srcMap[identifier]
		if !ok {
			return ErrorMemberNotExist{i.text}
		}
		return i.retrieveNext(root, child, result)

	case []interface{}:
		if i.isAsteriskNode() {
			return i.processAsteriskValue(root, current, result)
		}

		if len(i.identifiers) > 1 || len(i.identifiers[0]) > 0 {
			return ErrorTypeUnmatched{`map`, reflect.TypeOf(current).String(), i.text}
		}

		return i.retrieveNext(root, current, result)
	}

	foundType := `null`
	if current != nil {
		foundType = reflect.TypeOf(current).String()
	}
	return ErrorTypeUnmatched{`map/array`, foundType, i.text}
}

func (i *syntaxChildIdentifier) processAsteriskValue(
	root, current interface{}, result *resultContainer) error {

	switch current.(type) {
	case map[string]interface{}:
		srcMap := current.(map[string]interface{})
		keys := make([]string, 0, len(srcMap))
		for key := range srcMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			i.retrieveNext(root, srcMap[key], result)
		}

	case []interface{}:
		for _, value := range current.([]interface{}) {
			i.retrieveNext(root, value, result)
		}
	}

	if !result.hasResult() {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}

func (i *syntaxChildIdentifier) processMultiIdentifier(
	root interface{}, srcMap map[string]interface{}, result *resultContainer) error {

	for _, key := range i.identifiers {
		if _, ok := srcMap[key]; ok {
			i.retrieveNext(root, srcMap[key], result)
		}
	}

	if !result.hasResult() {
		return ErrorNoneMatched{i.getConnectedText()}
	}

	return nil
}

func (i *syntaxChildIdentifier) isAsteriskNode() bool {
	return i.isAsterisk
}
