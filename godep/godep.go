package godep

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/anaskhan96/soup"
)

// Node struct
type Node struct {
	PkgName string  `json:"pkgname"`
	IsRoot  bool    `json:-`
	Parent  *Node   `json:parent`
	Imports []*Node `json:imports`
}

// FindImports func
func (node *Node) FindImports() error {
	if node == nil {
		return nil
	}

	pkgImportURL := pkgImportURL(node.PkgName)
	imports, err := fetchImport(pkgImportURL, node.IsRoot)
	// fmt.Println("imports:", imports)
	if err != nil && node.IsRoot {
		log.Fatal(fmt.Sprintf("Package %s not found", node.PkgName))
	}

	for _, imp := range imports {
		pkgName := imp
		// fmt.Println(pkgName)
		childNode := &Node{
			PkgName: pkgName,
			IsRoot:  false,
			Parent:  node,
		}
		childNode.FindImports()
		node.Imports = append(node.Imports, childNode)
	}
	return nil
}

func (node *Node) graph(buf *bytes.Buffer, edges map[string]bool) error {
	for _, imp := range node.Imports {
		edge := fmt.Sprintf("%s->%s", node.PkgName, imp.PkgName)
		_, ok := edges[edge]
		if !ok {
			edges[edge] = true
			buf.WriteString(fmt.Sprintf("	\"%s\" -> \"%s\";\n", node.PkgName, imp.PkgName))
		}
		imp.graph(buf, edges)
	}
	return nil
}

// BuildGraph func
func (node *Node) BuildGraph() string {
	edges := make(map[string]bool)
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("digraph G {\n")
	buf.WriteString("	 rankdir=\"LR\";\n")
	buf.WriteString("    pad=.25;\n")
	buf.WriteString("    ratio=fill;\n")
	buf.WriteString("    dpi=360;\n")

	buf.WriteString("    node [shape=box];\n")

	node.graph(buf, edges)

	buf.WriteString(fmt.Sprintf(" 	\"%s\"  [style=filled];\n", node.PkgName))

	buf.WriteString("}")
	return buf.String()
}

func pkgURL(pkgName string) string {
	return "https://pkg.go.dev/" + pkgName
}

func pkgImportURL(pkgName string) string {
	return "https://pkg.go.dev/" + pkgName + "?tab=imports"
}

func fetchImport(pkgImportURL string, isRoot bool) ([]string, error) {
	imports := []string{}
	var htmlContent string
	if isRoot {
		resp, err := http.Get(pkgImportURL)
		if err != nil {
			return imports, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			htmlDoc, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return imports, err
			}
			htmlContent = string(htmlDoc)
		} else {
			return imports, errors.New("Not Found")
		}
	} else {
		resp, err := soup.Get(pkgImportURL)
		if err != nil {
			fmt.Println("Not Found", pkgImportURL)
			return imports, err
		}
		htmlContent = resp
	}

	root := soup.HTMLParse(htmlContent)
	doc := root.Find("h2", "class", "Imports-heading")
	if doc.Text() == "Imports" {
		doc = root.Find("ul", "class", "Imports-list")
		links := doc.FindAll("a")
		for _, link := range links {
			imports = append(imports, link.FullText())
			// fmt.Println("fetchImport:", link.FullText())
		}
	}
	return imports, nil
}
