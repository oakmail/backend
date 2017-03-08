package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

const goTemplate = `package errors

var ((%range $key, $value := .%)
	(%$value.Name%) = &Error{Code: (%$value.ID%), Message: "(%$value.Encoded%)"}(%end%)
)
`

const mdTemplate = `# List of error codes

## Error listing
{{range $key, $value := .}}
### {{$value.ID}}: {{$value.Message}}
{{range $key2, $value2 := .Locations}} - [{{$value2.Path}}:{{$value2.Line}}]({{$value2.Link}})
{{end}}{{end}}`

var (
	nonalpha = regexp.MustCompile("[^A-Za-z0-9]+")
	comments = regexp.MustCompile("(?m)//.*$")
	errl     = regexp.MustCompile(`(errors\.[[:alnum:]]+)`)
)

type item struct {
	ID        int
	Name      string
	Message   string
	Encoded   string
	Locations []loc
}

type itemList []*item

func (p itemList) Len() int           { return len(p) }
func (p itemList) Less(i, j int) bool { return p[i].ID < p[j].ID }
func (p itemList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type loc struct {
	Path string
	Line int
	Link string
}

func main() {
	paths := []string{}
	if err := filepath.Walk("../pkg/api", func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Index(path, "_test") != -1 {
			return nil
		}

		if filepath.Ext(path) == ".go" {
			paths = append(paths, path)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	locations := map[string][]loc{}
	for _, path := range paths {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		path = strings.Replace(path[3:], string(os.PathSeparator), "/", -1)

		lines := strings.Split(string(contents), "\n")
		for index, line := range lines {
			matches := errl.FindAllString(string(line), -1)

			for _, match := range matches {
				if match == "errors.Abort" || match == "errors.Error" {
					continue
				}

				if _, ok := locations[match]; !ok {
					locations[match] = []loc{}
				}

				locations[match] = append(locations[match], loc{
					Path: path,
					Line: index + 1,
					Link: "../" + path + "#L" + strconv.Itoa(index+1),
				})
			}
		}
	}

	file, err := ioutil.ReadFile("./errors.txt")
	if err != nil {
		panic(err)
	}
	file = comments.ReplaceAll(file, []byte(""))

	errors := bytes.Split(file, []byte("\n"))

	var items itemList
	for id, bmsg := range errors {
		msg := string(bytes.Replace(bmsg, []byte("\r"), []byte(""), -1))
		id := id + 1

		key := strings.Replace(
			strings.Title(
				strings.Trim(
					nonalpha.ReplaceAllString(
						strings.Replace(
							msg, "-", "", -1,
						), " ",
					), " ",
				),
			), " ", "", -1,
		)

		locs := []loc{}
		if fl, ok := locations["errors."+key]; ok {
			locs = fl
		}

		items = append(items, &item{
			ID:        id,
			Name:      key,
			Message:   msg,
			Encoded:   strings.Replace(msg, "\"", "\\\"", -1),
			Locations: locs,
		})
	}

	sort.Sort(items)

	var buf bytes.Buffer
	if err := template.Must(template.New("golang").Delims("(%", "%)").Parse(goTemplate)).Execute(&buf, items); err != nil {
		panic(err)
	}
	source, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("../pkg/api/errors/errors.go", source, 0644); err != nil {
		panic(err)
	}
	buf.Reset()

	if err := template.Must(template.New("markdown").Parse(mdTemplate)).Execute(&buf, items); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile("../docs/api/errors.md", buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Generated %d errors definitions.\n", len(items))
}
