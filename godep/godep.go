package godep

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/anaskhan96/soup"
)

// Node struct
type Node struct {
	PkgName    string  `json:"pkgname"`
	IsRoot     bool    `json:"-"`
	Parent     *Node   `json:"parent"`
	ImportType string  `json:"import_type"` // internal/external/standard
	Imports    []*Node `json:"imports"`
}

// FindImports func
func (node *Node) FindImports() error {

	if node.ImportType == "standard" || node.ImportType == "external" {
		return nil
	}

	pkgImportURL := pkgImportURL(node.PkgName)
	imports, err := fetchImport(pkgImportURL, node.IsRoot)
	if err != nil && node.IsRoot {
		return errors.New("Package " + node.PkgName + " not found")
	}

	for key, val := range imports {
		pkgName := key
		childNode := &Node{
			PkgName:    pkgName,
			IsRoot:     false,
			Parent:     node,
			ImportType: val,
		}
		childNode.FindImports()
		node.Imports = append(node.Imports, childNode)
	}
	return nil
}

func (node *Node) graph(existEdges map[string]bool, nodes *[]string, edges *[]string) error {
	for _, imp := range node.Imports {
		edge := fmt.Sprintf("%s->%s", node.PkgName, imp.PkgName)
		// Check existEdges
		_, ok := existEdges[edge]
		if !ok {
			existEdges[edge] = true
			edge := fmt.Sprintf("	\"%s\" -> \"%s\";\n", node.PkgName, imp.PkgName)
			*edges = append(*edges, edge)
			if imp.ImportType == "standard" {
				node := fmt.Sprintf(" 	\"%s\"  [style=filled,color=steelblue1];\n", imp.PkgName)
				*nodes = append(*nodes, node)
			} else if imp.ImportType == "internal" {
				node := fmt.Sprintf(" 	\"%s\"  [style=filled,color=chocolate1];\n", imp.PkgName)
				*nodes = append(*nodes, node)
			} else {
				node := fmt.Sprintf(" 	\"%s\"  [style=filled,color=hotpink];\n", imp.PkgName)
				*nodes = append(*nodes, node)
			}
		}
		imp.graph(existEdges, nodes, edges)
	}
	return nil
}

// BuildGraph func
func (node *Node) BuildGraph() string {
	existEdges := make(map[string]bool)
	nodes := make([]string, 0)
	edges := make([]string, 0)

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("digraph G {\n")
	buf.WriteString("	 rankdir=\"LR\";\n")
	buf.WriteString("	 labelloc=\"t\";\n")
	buf.WriteString("	 label=\"Package: " + node.PkgName + "\";\n")
	buf.WriteString("    pad=.5;\n")
	buf.WriteString("    ratio=auto;\n")
	buf.WriteString("    dpi=360;\n")
	buf.WriteString("    graph [fontsize=16 fontname=\"Roboto Condensed, sans-serif\"];\n")
	buf.WriteString("    node [shape=box style=rounded fontname=\"Roboto Condensed, sans-serif\" fontsize=11 height=0 width=0 margin=.08];\n")
	buf.WriteString("    edge [fontsize=10, fontname=\"Roboto Condensed, sans-serif\" splines=\"polyline\"];\n")

	node.graph(existEdges, &nodes, &edges)

	sort.Strings(nodes)
	sort.Strings(edges)

	buf.WriteString(fmt.Sprintf(" 	\"%s\"  [style=filled,color=palegreen];\n", node.PkgName))

	buf.WriteString("// Nodes")
	for _, node := range nodes {
		buf.WriteString(node)
	}

	buf.WriteString("// Edges")
	for _, edge := range edges {
		buf.WriteString(edge)
	}

	buf.WriteString("}")
	return buf.String()
}

func pkgURL(pkgName string) string {
	return "https://pkg.go.dev/" + pkgName
}

func pkgImportURL(pkgName string) string {
	return "https://pkg.go.dev/" + pkgName + "?tab=imports"
}

func fetchImport(pkgImportURL string, isRoot bool) (map[string]string, error) {
	imports := make(map[string]string)
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
			return imports, errors.New("Not_Found")
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

	docs := root.FindAll("h2", "class", "Imports-heading")
	for _, doc := range docs {
		text := strings.ToLower(doc.Text())

		if strings.Contains(text, "standard library imports") {
			temp := doc.FindNextElementSibling()
			links := temp.FindAll("a")
			for _, link := range links {
				imports[link.FullText()] = "standard"
			}
		} else if strings.Contains(text, "imports in module") {
			temp := doc.FindNextElementSibling()
			links := temp.FindAll("a")
			for _, link := range links {
				imports[link.FullText()] = "internal"
			}
		} else if strings.Contains(text, "imports") {
			temp := doc.FindNextElementSibling()
			links := temp.FindAll("a")
			for _, link := range links {
				imports[link.FullText()] = "external"
			}
		}
	}
	return imports, nil
}
