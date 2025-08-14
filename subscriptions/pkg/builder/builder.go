package builder

import "github.com/Masterminds/squirrel"

func SetupBuilder() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}
