package syntax

import (
	"github.com/AsaiYusuke/jsonpath/v2/errors"
)

type syntaxChildSingleIdentifier struct {
	*syntaxBasicNode

	identifier string
}

func (i *syntaxChildSingleIdentifier) retrieve(
	root, current any, results *[]any) errors.ErrorRuntime {

	if srcMap, ok := current.(map[string]any); ok {
		return i.retrieveMapNext(root, srcMap, i.identifier, results)
	}

	return i.newErrTypeUnmatched(msgTypeObject, current)
}
