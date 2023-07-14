package outputparser

import (
	"fmt"
	"regexp"

	"github.com/tmc/langchaingo/schema"
)

// RegexParser is an output parser used to parse the output of an llm as a map.
type RegexDict struct {
	OutputKeyToFormat map[string]string
	NoUpdateValue     string
}

const (
	REGEX_DICT_PATTERN = `(?:%s):\s?(?P<value>(?:[^.'\n']*)\.?)`
)

// NewRegexParser returns a new RegexParser.
func NewRegexDict(outputKeyToFormat map[string]string, noUpdateValue string) RegexDict {
	return RegexDict{
		OutputKeyToFormat: outputKeyToFormat,
		NoUpdateValue:     noUpdateValue,
	}
}

// Statically assert that RegexParser implements the OutputParser interface.
var _ schema.OutputParser[any] = RegexDict{}

// GetFormatInstructions returns instructions on the expected output format.
func (p RegexDict) GetFormatInstructions() string {
	instructions := "Your output should be a map of strings. e.g.:\n"
	instructions += "map[string]string{\"key1\": \"value1\", \"key2\": \"value2\"}\n"

	return instructions
}

func (p RegexDict) parse(text string) (map[string]string, error) {
	results := make(map[string]string, len(p.OutputKeyToFormat))

	for key, format := range p.OutputKeyToFormat {
		expression := regexp.MustCompile(fmt.Sprintf(REGEX_DICT_PATTERN, format))
		matches := expression.FindStringSubmatch(text)

		if len(matches) < 2 {
			return nil, ParseError{
				Text:   text,
				Reason: fmt.Sprintf("No match found for expression %s", expression),
			}
		}

		if len(matches) > 2 {
			return nil, ParseError{
				Text:   text,
				Reason: fmt.Sprintf("Multiple matches found for expression %s", expression),
			}
		}

		match := matches[1]

		if match == p.NoUpdateValue {
			continue
		}

		results[key] = match
	}

	return results, nil
}

// Parse parses the output of an llm into a map of strings.
func (p RegexDict) Parse(text string) (any, error) {
	return p.parse(text)
}

// ParseWithPrompt does the same as Parse.
func (p RegexDict) ParseWithPrompt(text string, _ schema.PromptValue) (any, error) {
	return p.parse(text)
}

// Type returns the type of the parser.
func (p RegexDict) Type() string {
	return "regex_dict_parser"
}
