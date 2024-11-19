package structurederror

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

/*
StructuredError is a way to craft errors using similar semantics to slog.
*/
type StructuredError struct {
	Message string
	Args    map[string]any

	Delimiter string
	Logger    *slog.Logger
}

type ErrorArg struct {
	Key   string
	Value any
}

type ErrorMaker func(message string, args ...any) *StructuredError
type ErrorMakerOption func(*StructuredError)

/*
New creates a new error maker. An error maker is a function that accepts
an error message and a series of arguments. Args can be either key/value
pairs, or an ErrorArg struct. "Key" is expected to be "stringy".
*/
func New(options ...ErrorMakerOption) ErrorMaker {
	return func(message string, args ...any) *StructuredError {
		var (
			key   string
			value any
		)

		/*
		 * Create the structured error object, altering it with
		 * any options provided.
		 */
		result := &StructuredError{
			Message: message,
			Args:    make(map[string]any),

			Delimiter: "-",
		}

		for _, option := range options {
			option(result)
		}

		/*
		 * Parse arguments. If the arg is an ErrorArg struct
		 * we'll use it as is. Otherwise we are expecting
		 * key/value pairs. Keys should be "stringy".
		 */
		for index := 0; index < len(args); index++ {
			if arg, ok := isErrorArg(args[index]); ok {
				result.Args[arg.Key] = arg.Value
				continue
			}

			if multipleArgs, ok := isErrorArgSlice(args[index]); ok {
				for _, arg := range multipleArgs {
					result.Args[arg.Key] = arg.Value
				}

				continue
			}

			if key == "" {
				key = fmt.Sprintf("%v", args[index])
				continue
			}

			value = args[index]

			arg := makeErrorArg(key, value)
			result.Args[arg.Key] = arg.Value

			// Reset the key for the next pair.
			key = ""
		}

		if result.Logger != nil {
			result.writeLog()
		}

		return result
	}
}

/*
Error implements the Error value interface.
*/
func (se *StructuredError) Error() string {
	result := strings.Builder{}

	result.WriteString(se.Message)

	argKeys := []string{}

	for key := range se.Args {
		argKeys = append(argKeys, key)
	}

	sort.Strings(argKeys)

	for _, key := range argKeys {
		value := se.Args[key]

		result.WriteString(" " + se.Delimiter + " ")
		result.WriteString(key + ": ")
		result.WriteString(fmt.Sprintf("%v", value))
	}

	return result.String()
}

func (se *StructuredError) writeLog() {
	args := []any{}

	for key, value := range se.Args {
		args = append(args, key)
		args = append(args, value)
	}

	se.Logger.Error(se.Message, args...)
}

/*
WithDelimiter allows you to configure a StructuredError's arg delimiter.
*/
func WithDelimiter(delimiter string) ErrorMakerOption {
	return func(s *StructuredError) {
		s.Delimiter = delimiter
	}
}

/*
WithSlog tells the maker to call the provided logger Error method, using Args
as key value pairs in the log. This happens when you create the error.
*/
func WithSlog(logger *slog.Logger) ErrorMakerOption {
	return func(se *StructuredError) {
		se.Logger = logger
	}
}

func isErrorArg(item any) (ErrorArg, bool) {
	switch v := item.(type) {
	case ErrorArg:
		return v, true

	case *ErrorArg:
		return *v, true

	default:
		return ErrorArg{}, false
	}
}

func isErrorArgSlice(item any) ([]ErrorArg, bool) {
	switch v := item.(type) {
	case []ErrorArg:
		return v, true

	default:
		return []ErrorArg{}, false
	}
}

func makeErrorArg(item1 any, item2 any) ErrorArg {
	key := fmt.Sprintf("%v", item1)
	return ErrorArg{
		Key:   key,
		Value: item2,
	}
}
