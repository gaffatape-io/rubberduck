package header

import (
	"bufio"
	"io"
	"strings"

	errs "github.com/gaffatape-io/gopherrs"
)

type Output struct {
	Args []string
	Env []string
}

// trimWS removes space, tabs and newline chars from the beginning
// and end of the string.
func trimWS(s string) string {
	return strings.Trim(s, " \t\n")
}

// readSection reads a section from the buffered reader.
// A section starts with a header and ends with an empty newline.
// Ex.
// args:
//   flag1
//   flag2
//
// Here 'args' is the header and the flag1 and flag2 are the arguments.
func readSection(r *bufio.Reader) (string, []string, error) {
	header, err := r.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	header = trimWS(header)
	
	args := []string{}
	for {
		tmp, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		tmp = trimWS(tmp)
		if tmp == "" || err != nil {
			return header, args, err
		}
		args = append(args, tmp)
	}
	return header, args, nil
}

// Reads the default section from the reader.
func Read(r io.Reader) (Output, error) {
	reader := bufio.NewReader(r)

	wrongHeaderErr := func(header string) *errs.E {
		return errs.InvalidArgumentf("wrong header; got:%q", header)
	}
	
	argsHeader, args, err := readSection(reader)
	if err != nil {
		return Output{}, err
	} else if argsHeader != "args:" {
		return Output{}, wrongHeaderErr(argsHeader)
	}

	envHeader, env, err := readSection(reader)
	if err != nil {
		return Output{}, err
	} else if envHeader != "env:" {
		return Output{}, wrongHeaderErr(envHeader)
	}

	return Output{Args: args, Env: env}, nil
}
