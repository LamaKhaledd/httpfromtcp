package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.SplitN(string(data), "\r\n", 2)
	if len(lines) == 0 {
		return nil, errors.New("empty request")
	}

	reqLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(line string) (*RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line, expected 3 parts")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	for _, r := range method {
		if !unicode.IsUpper(r) {
			return nil, errors.New("invalid method, must be uppercase")
		}
	}

	if version != "HTTP/1.1" {
		return nil, errors.New("unsupported HTTP version, only 1.1 allowed")
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}
