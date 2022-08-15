package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

// Build time variables
var (
	version = "main"
)

// Command line flags
var (
	input        string
	output       string
	printVersion bool
	printHelp    bool
)

func init() {
	flag.StringVarP(&input, "file", "f", "-", "Input file containing helm values.yaml")
	flag.StringVarP(&output, "output", "o", "-", "Output file to write Terraform code")
	flag.BoolVarP(&printVersion, "version", "v", false, "Show version")
	flag.BoolVarP(&printHelp, "help", "h", false, "Show help")

	flag.Parse()
}

func main() {
	if printHelp {
		flag.Usage()
		os.Exit(0)
	}

	if printVersion {
		fmt.Printf("helm-to-hcl %s\n", version)
		os.Exit(0)
	}

	var file *os.File
	if input == "-" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprint(os.Stderr, "error: standard input empty\n")
			os.Exit(1)
		}
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	inputBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	hcl, err := helmToHcl(inputBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}

	if output == "-" {
		fmt.Print(hcl)
	} else {
		err = ioutil.WriteFile(output, []byte(hcl), 0o644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func helmToHcl(input []byte) (string, error) {
	m := yaml.MapSlice{}
	err := yaml.Unmarshal(input, &m)
	if err != nil {
		return "", err
	}

	var buf strings.Builder

	buf.WriteString(`resource "helm_release" "default" {`)
	buf.WriteByte('\n')
	buf.WriteString(stringWithIndent(`values = [`, 2))
	buf.WriteByte('\n')
	buf.WriteString(stringWithIndent(`yamlencode(`, 4))
	buf.WriteByte('\n')
	buf.WriteString(stringWithIndent(`{`, 6))
	buf.WriteString(convert(m, 8))
	buf.WriteByte('\n')
	buf.WriteString(stringWithIndent(`}`, 6))
	buf.WriteByte('\n')
	buf.WriteString(stringWithIndent(`)`, 4))
	buf.WriteByte('\n')
	buf.WriteString(stringWithIndent(`]`, 2))
	buf.WriteByte('\n')
	buf.WriteString(`}`)
	buf.WriteByte('\n')

	return buf.String(), nil
}

func convert(m yaml.MapSlice, indent int) string {
	var buf strings.Builder
	for _, item := range m {

		key := formatKey(item.Key.(string))

		switch v := item.Value.(type) {

		case bool, string, int:
			buf.WriteByte('\n')
			buf.WriteString(stringWithIndent(fmt.Sprintf(`%s = %s`, key, convertValue(v)), indent))
		case []interface{}:
			buf.WriteByte('\n')
			buf.WriteString(stringWithIndent(fmt.Sprintf(`%s = [`, key), indent))
			if len(v) > 0 {
				for _, item := range v {
					switch v := item.(type) {
					case yaml.MapSlice:
						buf.WriteByte('\n')
						buf.WriteString(stringWithIndent(`{`, indent+2))
						buf.WriteString(convert(v, indent+4))
						buf.WriteByte('\n')
						buf.WriteString(stringWithIndent(`},`, indent+2))
					default:
						buf.WriteByte('\n')
						buf.WriteString(stringWithIndent(convertValue(v), indent+2))
					}
				}
				buf.WriteByte('\n')
				buf.WriteString(strings.Repeat(" ", indent))
			}
			buf.WriteByte(']')

		case yaml.MapSlice:
			buf.WriteByte('\n')
			buf.WriteString(stringWithIndent(fmt.Sprintf(`%s = {`, key), indent))

			if len(v) > 0 {
				buf.WriteString(convert(v, indent+2))
				buf.WriteByte('\n')
				buf.WriteString(strings.Repeat(" ", indent))
			}
			buf.WriteByte('}')

		case nil:

		default:
			fmt.Printf("Type missed: %v\n", reflect.TypeOf(v))
		}

	}

	return buf.String()
}

func convertValue(value interface{}) string {
	switch v := value.(type) {
	case bool, int:
		return fmt.Sprintf(`%v`, v)
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case nil:
		return ""
	default:
		fmt.Printf("Value Type missed: %v\n", reflect.TypeOf(v))

		return ""
	}
}

func stringWithIndent(s string, indent int) string {
	return strings.Repeat(" ", indent) + s
}

func formatKey(key string) string {
	m := regexp.MustCompile(`^[A-Za-z][0-9A-Za-z-_]+$`)
	if m.MatchString(key) {
		return key
	} else {
		return `"` + key + `"`
	}
}
