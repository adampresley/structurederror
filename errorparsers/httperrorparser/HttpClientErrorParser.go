package httperrorparser

import (
	"io"
	"net/http"
	"strconv"
	"strings"

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

func Parse(resp *http.Response, options ...Option) structurederror.ErrorArg {
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

func (p *Parser) Parse(resp *http.Response) structurederror.ErrorArg {
	var (
		message strings.Builder = strings.Builder{}
	)

	result := structurederror.ErrorArg{
		Key: Key,
	}

	message.WriteString("Status Code: ")
	message.WriteString(strconv.Itoa(resp.StatusCode))
	result.Value = message.String()

	if p.IncludeStatus {
		message.WriteString(", Status: ")
		message.WriteString(resp.Status)

		result.Value = message.String()
	}

	if p.IncludeResponseBody {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return result
		}

		message.WriteString(", Body: ")
		message.WriteString(string(b))
		result.Value = message.String()
	}

	return result
}
