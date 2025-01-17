# Structured Error

Structured Error is a little library to help build "structured" errors. For our purposes, a "structured error" is one that has a message, followed by a series of key/value pairs. Each key/value pair is in the form of "key: value", and each pair is seperated by a delimiter. The default delimiter is "-", but it can be changed.

## Example

```go
maker := structurederror.New()
err := maker("error message here", "key", "value")
// err == "error message here - key: value"
```

## Arguments

Arguments are either key/value pairs, where the key can be coerced to a string type, an ErrorArg, or a slice of ErrorArgs. An ErrorArg is just a structure with key and value properties.

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

You can also use a slice of ErrorArg as an argument to StructuredError. Here is a sample.

```go
func ParseTheThings(someImput string) []structurederror.ErrorArg {
  // parse the stuff
  return []structurederror.ErrorArg{
    {
      Key: "key1",
      Value: "value1",
    },
    {
      Key: "key2",
      Value: "value2",
    },
  }
}

err := maker("some error", ParseTheThings("some input"))
// err == "some error - key1: value1 - key2: value2

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

## Included Parsers

Below are some included parsers that can add context to your error.

### Http Client Error Parser
This provides a function that takes an `*http.Response` and will break it down into a slice of ErrorArgs. By default it only includes the status code, but with options, you can include the status line and body.

```go
maker := structurederror.New()

err := maker(
  "some error", 
  httpclientparser.Parse(
    resp,
    httpclientparser.WithStatus(),
    httpclientparser.WithBody(),
  ),
)

// err == "some error - body: HTTP response body goes here - status: 500 Internal Server Error - statusCode: 500"
```
