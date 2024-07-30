/*

package xgorm

import (
	"strings"
	"unicode"

	"github.com/gnomego/sdk/errors"
	"gorm.io/gorm"
)

type Token struct {
	Type  string
	Value string
}

type filterExpression struct {
	Field string
	Op    string
	Value string
}

type groupExpression struct {
	Op    string
	Items []FilterExpression
}

type Predicate interface {
	Apply(db *gorm.DB) *gorm.DB
}

func handleTokens(db *gorm.DB, filter string) error {
	tokens := tokenize(filter)
	filter2 := "(name == 'john' and age > 20) or name startswith 'j'"
	group := []string{"and", "or"}

	logical := []string{"==", "!=", ">", ">", "<", "gte"}
	functions := []string{"startswith", "endswith", "contains"}

	if len(tokens) < 3 {
		return errors.NewWithCodef("invalid_filter", "filter must have at least 3 tokens %s", filter)
	}

	nextFilter := &FilterExpression{}

	for _, token := range tokens {

	}

	return nil
}

func tokenize(filter string) []string {

	quote := false
	sb := strings.Builder{}
	tokens := []string{}

	for _, r := range filter {

		if quote {
			if r == '\'' {
				quote = false
				continue
			}

			sb.WriteRune(r)
			continue
		}

		if r == '\'' {
			quote = true
			continue
		}

		if r == ' ' {
			if sb.Len() > 0 {
				tokens = append(tokens, sb.String())
				sb.Reset()
			}

			continue
		}

		if r == '(' || r == ')' {
			if sb.Len() > 0 {
				tokens = append(tokens, sb.String())
				sb.Reset()

				tokens = append(tokens, string(r))
			}
		}

		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			sb.WriteRune(r)
		}

		continue
	}

	if sb.Len() > 0 {
		tokens = append(tokens, sb.String())
	}

	return tokens
}
*/

package xgorm
