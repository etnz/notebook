package notebook

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"maps"
	"os"
	"path"
	"slices"
	"strings"
)

// very much like jupyter or colab, a notebook is made of cell that can hold different types of content (text, table, plots)
// unlike jupyter or colab this is not intended to be interactive. I don't need that.
// instead execution will generate a plain html that can be visualized live.

const HeaderCellStyle = "cell-style" // ID to set a style header to style cells. A default value is provided.

// cell is the main component of a Notebook.
type cell struct {
	content bytes.Buffer // plain html
	Title   string       // Cell Title
	// and that's all for now
}

func newCell() *cell { return &cell{} }

// Content returns the cell's content as HTML
func (c *cell) Content() template.HTML { return template.HTML(c.content.String()) }

// Notebook is a struct to receive all the contents of a notebook in memory.
type Notebook struct {
	Title  string // The notebook title.
	Output string // A filename to save the Notebook to.

	cells   []*cell           // The dynamic list of cells.
	headers map[string]string // id to header fragment

	console bytes.Buffer // to receive any fmt.Printf
}

// New creates a new Notebook.
//
// Output filename is initialized to be the os.Arg[0] +'.html'.
func New() *Notebook {
	name := path.Base(os.Args[0])
	return &Notebook{
		Output: name + ".html",
		Title:  strings.Title(name),
		headers: map[string]string{
			// default headers
			HeaderCellStyle: cellStyle,
		},
	}
}

// AddContent appends html content into a new Cell.
func (nb *Notebook) AddContent(title string, content string) {
	cell := newCell()
	cell.Title = title
	cell.content.WriteString(content)
	nb.cells = append(nb.cells, cell)
}

// AddHeader appends a header statement, once per ID.
// Calling AddHeader with the same ID will overwrite the existing value.
func (nb *Notebook) AddHeader(id string, content string) {
	if nb.headers == nil {
		nb.headers = make(map[string]string)
	}
	nb.headers[id] = content
}

// Print behave like fmt.Print but on the notebook console.
func (nb *Notebook) Print(a ...any) (n int, err error) { return fmt.Fprint(&nb.console, a...) }

// Printf behave like fmt.Printf but on the notebook console.
func (nb *Notebook) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(&nb.console, format, a...)
}

// Println behave like fmt.Println but on the notebook console.
func (nb *Notebook) Println(a ...any) (n int, err error) { return fmt.Fprintln(&nb.console, a...) }

// Close generates the notebook as HTML into the output file.
func (nb *Notebook) Close() error {
	f, err := os.Create(nb.Output)
	if err != nil {
		return fmt.Errorf("cannot create output file %q", nb.Output)
	}
	defer f.Close()
	if err := nb.Render(f); err != nil {
		return fmt.Errorf("cannot render notebook: %w", err)
	}
	return nil
}

// Render the notebook into a writer.
func (nb *Notebook) Render(w io.Writer) error {
	// pass nbView to the template to expose to it some private fields.
	return nbTemplate.Execute(w, nbView{nb})
}

type nbView struct{ *Notebook }

// Headers return the list of HTML fragment for the header
func (nbv nbView) Headers() []template.HTML {
	var result []template.HTML
	// retrieve by keys and sort they by ID so that the order is repeatable.
	keys := slices.AppendSeq([]string{}, maps.Keys(nbv.headers))
	slices.Sort(keys)
	for _, k := range keys {
		result = append(result, template.HTML(nbv.headers[k]))
	}
	return result
}

func (nbv nbView) Cells() []*cell  { return nbv.cells }
func (nbv nbView) Console() string { return nbv.console.String() }

var nbTemplate = template.Must(template.New("notebook").Parse(
	`<!DOCTYPE html>
<html>
	<head>
	{{- with .Title}}<title>{{.}}</title>{{end -}}
	{{- range .Headers}}{{.}}{{end -}}
	</head>
	<body>
		{{with .Title}}<h1>{{.}}</h1>{{end}}
		<div class="cell-container">
		{{- range .Cells}}
			<details open class="cell">
				<summary>{{.Title}}</summary>
				{{.Content}}
			</details>
		{{- end}}
			<details open class="cell">
				<summary>Console</summary>
				<pre>{{.Console}}</pre>
			</details>
		</div>
	</body>
</html>`))

const cellStyle = `<style>
	.cell-container {
		display: flex;
		flex-direction: column;
	}
	details.cell {
		border: 1px solid #aaa;
		border-radius: 4px;
		padding: 0.5em 0.5em 0;
		display: block;
	}
	details.cell > summary {
		font-weight: bold;
		margin: -0.5em -0.5em 0;
		padding: 0.5em;
	}
	details[open].cell {
		padding: 0.5em;
	}
	details[open].cell > summary {
		border-bottom: 1px solid #aaa;
		margin-bottom: 0.5em;
	}
</style>`
