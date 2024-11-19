# Structured Error

Structured Error is a little library to help build "structured" errors. For our purposes, a "structured error" is one that has a message, followed by a series of key/value pairs. Each key/value pair is in the form of "key: value", and each pair is seperated by a delimiter. The default delimiter is "-", but it can be changed.

## Example

```go
maker := structurederror.New()
err := maker("error message here", "key", "value", someFuncThatReturnsAnErrorArg())
```

## Arguments

Arguments are either key/value pairs, where the key can be coerced to a string type, or an ErrorArg. An ErrorArg is just a structure with key and value properties.

```go
type ErrorArg struct {
  Key   string
  Value any
}
```

The ErrorArg struct allows you to do interesting things, like build functions that can parse data and return an ErrorArg. For example:

```go
func ParseTheThing(someImput string) structurederror.ErrorArg {
    // parse the stuff
    return structurederror.ErrorArg{
        Key: "parserOutput",
        Value: "some value",
    }
}

err := maker("some error", ParseTheThing("some input"))
// err == "some error - parserOutput: some input"
```

## Options

The structured error maker can be configured with options. Below are the available options.

### WithDelimiter

**WithDelimiter** allows you to change the delimiter used between key/value pairs. The default value is "-".

```go
maker := structurederror.New(
  structurederror.WithDelimiter(";;"),
)
```

### WithSlog

**WithSlog** will call the **Error()** method on a provided logger when the error is created. Here is an example.

```go
logger := slog.New(slog.Default().Handler())

maker := structurederror.New(
  structurederror.WithSlog(logger),
)

err := maker("some error", "key", "value")
// This will output to stdout:
// 2024/11/19 16:22:00 ERROR some error key=value
```
