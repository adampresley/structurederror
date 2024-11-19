package httperrorparser

import (
	"io"
	"net/http"

	"github.com/adampresley/structurederror"
)

const (
	Key string = "http response"
)

type Parser struct {
	IncludeStatus       bool
	IncludeResponseBody bool
}

type Option func(*Parser)

func Parse(resp *http.Response, options ...Option) []structurederror.ErrorArg {
	parser := &Parser{
		IncludeStatus:       false,
		IncludeResponseBody: false,
	}

	for _, option := range options {
		option(parser)
	}

	return parser.Parse(resp)
}

func WithStatus() Option {
	return func(p *Parser) {
		p.IncludeStatus = true
	}
}

func WithResponseBody() Option {
	return func(p *Parser) {
		p.IncludeResponseBody = true
	}
}

func (p *Parser) Parse(resp *http.Response) []structurederror.ErrorArg {
	result := []structurederror.ErrorArg{}

	result = append(
		result,
		structurederror.ErrorArg{Key: "statusCode", Value: resp.StatusCode},
	)

	if p.IncludeStatus {
		result = append(
			result,
			structurederror.ErrorArg{Key: "status", Value: resp.Status},
		)
	}

	if p.IncludeResponseBody {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return result
		}

		result = append(
			result,
			structurederror.ErrorArg{Key: "body", Value: string(b)},
		)
	}

	return result
}
