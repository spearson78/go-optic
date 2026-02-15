package olex

import "strings"

type AstNode[I any] interface {
	Format(yield func(index Span[I], focus Token) bool) bool
}

func AstNodeToString[I any](n AstNode[I]) string {
	var sb strings.Builder
	n.Format(func(index Span[I], focus Token) bool {
		sb.WriteString(focus.value)
		return true
	})
	return sb.String()
}
