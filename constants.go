package jsonpath

const (
	msgErrorInvalidSyntaxUnrecognizedInput     string = `unrecognized input`
	msgErrorInvalidSyntaxUseBeginAtsign        string = `the use of '@' at the beginning is prohibited`
	msgErrorInvalidSyntaxOmitDollar            string = `the omission of '$' allowed only at the beginning`
	msgErrorInvalidSyntaxTwoCurrentNode        string = `comparison between two current nodes is not allowed`
	msgErrorInvalidSyntaxFilterMultiValuedNode string = `multi-valued node retrieve into the filter is prohibited`
)
