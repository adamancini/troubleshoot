package redact

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"reflect"
	"regexp"

	"github.com/replicatedhq/troubleshoot/pkg/constants"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type MultiLineRedactor struct {
	re1        *regexp.Regexp
	re2        *regexp.Regexp
	maskText   string
	filePath   string
	redactName string
	isDefault  bool
}

func NewMultiLineRedactor(re1, re2 *regexp.Regexp, maskText, path, name string, isDefault bool) (*MultiLineRedactor, error) {
	return &MultiLineRedactor{re1: re1, re2: re2, maskText: maskText, filePath: path, redactName: name, isDefault: isDefault}, nil
}

func (r *MultiLineRedactor) Redact(input io.Reader, path string) io.Reader {
	out, writer := io.Pipe()

	_, span := otel.Tracer(constants.LIB_TRACER_NAME).Start(context.Background(), fmt.Sprintf("Redactor %s", r.redactName))
	span.SetAttributes(attribute.String("type", reflect.TypeOf(MultiLineRedactor{}).String()))

	go func() {
		var err error
		defer func() {
			writer.CloseWithError(err)
		}()

		substStr := getReplacementPattern(r.re2, r.maskText)

		reader := bufio.NewReader(input)
		line1, line2, err := getNextTwoLines(reader, nil)
		if err != nil {
			// this will print 2 blank lines for empty input...
			fmt.Fprintf(writer, "%s\n", line1)
			fmt.Fprintf(writer, "%s\n", line2)
			return
		}

		flushLastLine := false
		lineNum := 1
		for err == nil {
			lineNum++ // the first line that can be redacted is line 2

			// If line1 matches re1, then transform line2 using re2
			if !r.re1.MatchString(line1) {
				fmt.Fprintf(writer, "%s\n", line1)
				line1, line2, err = getNextTwoLines(reader, &line2)
				flushLastLine = true
				continue
			}
			flushLastLine = false

			clean := r.re2.ReplaceAllString(line2, substStr)

			// io.WriteString would be nicer, but reader strips new lines
			fmt.Fprintf(writer, "%s\n%s\n", line1, clean)
			if err != nil {
				return
			}

			// if clean is not equal to line2, a redaction was performed
			if clean != line2 {
				addRedaction(Redaction{
					RedactorName:      r.redactName,
					CharactersRemoved: len(line2) - len(clean),
					Line:              lineNum,
					File:              r.filePath,
					IsDefaultRedactor: r.isDefault,
				})
			}

			line1, line2, err = getNextTwoLines(reader, nil)
		}

		if flushLastLine {
			fmt.Fprintf(writer, "%s\n", line1)
		}
	}()
	span.End()
	return out
}

func getNextTwoLines(reader *bufio.Reader, curLine2 *string) (line1 string, line2 string, err error) {
	line1 = ""
	line2 = ""

	if curLine2 == nil {
		line1, err = readLine(reader)
		if err != nil {
			return
		}

		line2, err = readLine(reader)
		return
	}

	line1 = *curLine2
	line2, err = readLine(reader)
	if err != nil {
		return
	}

	return
}
