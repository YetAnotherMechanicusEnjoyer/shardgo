package network

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrInvalidRequest = errors.New("invalid request format")

type Request struct {
	Command string
	Args    []string
}

func ParseRequest(reader *bufio.Reader) (*Request, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, ErrInvalidRequest
	}
	line = strings.TrimSuffix(line, "\r\n")

	if len(line) < 2 || line[0] != '*' {
		return nil, ErrInvalidRequest
	}

	numArgs, err := strconv.Atoi(line[1:])
	if err != nil || numArgs < 1 {
		return nil, ErrInvalidRequest
	}

	req := &Request{
		Args: make([]string, numArgs),
	}

	for i := 0; i < numArgs; i++ {
		argLine, err := reader.ReadString('\n')
		if err != nil {
			return nil, ErrInvalidRequest
		}
		argLine = strings.TrimSuffix(argLine, "\r\n")
		if len(argLine) < 2 || argLine[0] != '$' {
			return nil, ErrInvalidRequest
		}

		argLen, err := strconv.Atoi(argLine[1:])
		if err != nil {
			return nil, err
		}

		arg := make([]byte, argLen)
		_, err = io.ReadFull(reader, arg)
		if err != nil {
			return nil, err
		}

		_, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		req.Args[i] = string(arg)
	}

	req.Command = strings.ToUpper(req.Args[0])
	return req, nil
}

func ParseResponse(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSuffix(line, "\r\n")

	if len(line) == 0 {
		return "", errors.New("empty response")
	}

	switch line[0] {
	case '+':
		return line[1:], nil
	case '-':
		return line, nil
	case ':':
		return line, nil
	case '$':
		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return "", err
		}
		if length == -1 {
			return "", nil
		}
		val := make([]byte, length)
		_, err = io.ReadFull(reader, val)
		if err != nil {
			return "", err
		}
		_, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		return string(val), nil
	case '*':
		return "", errors.New("arrays not supported in responses yet")
	default:
		return "", errors.New("unknown response type")
	}
}

func FormatResponse(resp string) string {
	return fmt.Sprintf("+%s\r\n", resp)
}

func FormatBulkString(s string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

func FormatError(err error) string {
	return fmt.Sprintf("-ERR %v\r\n", err)
}
