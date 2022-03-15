// Package str 字符串辅助方法
package str

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Plural 转为复数 user -> users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Singular 转为单数 users -> user
func Singular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// Snake 转为 snake_case，如 TopicComment -> topic_comment
func Snake(word string) string {
	return strcase.ToSnake(word)
}

// Camel 转为 CamelCase，如 topic_comment -> TopicComment
func Camel(word string) string {
	return strcase.ToCamel(word)
}

// LowerCamel 转为 lowerCamelCase，如 TopicComment -> topicComment
func LowerCamel(word string) string {
	return strcase.ToLowerCamel(word)
}
