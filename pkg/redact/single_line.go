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

type SingleLineRedactor struct {
	re         *regexp.Regexp
	maskText   string
	filePath   string
	redactName string
	isDefault  bool
}

func NewSingleLineRedactor(re *regexp.Regexp, maskText, path, name string, isDefault bool) (*SingleLineRedactor, error) {
	return &SingleLineRedactor{re: re, maskText: maskText, filePath: path, redactName: name, isDefault: isDefault}, nil
}

func (r *SingleLineRedactor) Redact(input io.Reader, path string) io.Reader {
	out, writer := io.Pipe()

	_, span := otel.Tracer(constants.LIB_TRACER_NAME).Start(context.Background(), fmt.Sprintf("Redactor %s", r.redactName))
	span.SetAttributes(attribute.String("type", reflect.TypeOf(SingleLineRedactor{}).String()))

	go func() {
		var err error
		defer func() {
			if err == io.EOF {
				writer.Close()
			} else {
				writer.CloseWithError(err)
			}
		}()

		substStr := getReplacementPattern(r.re, r.maskText)

		reader := bufio.NewReader(input)
		lineNum := 0
		for {
			lineNum++
			var line string
			line, err = readLine(reader)
			if err != nil {
				return
			}

			if !r.re.MatchString(line) {
				fmt.Fprintf(writer, "%s\n", line)
				continue
			}

			clean := r.re.ReplaceAllString(line, substStr)

			// io.WriteString would be nicer, but scanner strips new lines
			fmt.Fprintf(writer, "%s\n", clean)
			if err != nil {
				return
			}

			// if clean is not equal to line, a redaction was performed
			if clean != line {
				addRedaction(Redaction{
					RedactorName:      r.redactName,
					CharactersRemoved: len(line) - len(clean),
					Line:              lineNum,
					File:              r.filePath,
					IsDefaultRedactor: r.isDefault,
				})
			}
		}
	}()
	span.End()
	return out
}
