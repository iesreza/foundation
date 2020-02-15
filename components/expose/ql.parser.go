package expose

import (
	"strings"
	"unicode"
)

type Node struct {
	Name    string
	Sub     []Node
	Params  map[string]string
	EndNode bool
}

func ParseQueryLanguage(query string) *Node {
	root := &Node{}

	parseQL(query, root)

	return root
}

func parseQL(input string, node *Node) int {
	outer := true
	subCandid := true
	node.Sub = []Node{}
	node.Params = map[string]string{}
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '(' {
			i += parseParameters(input[i:], node)
		}
		if c == '#' {
			for input[i] != '\n' && i < len(input) {
				i++

			}
			continue
		}
		if c == '{' {
			outer = false
			continue
		}
		if !outer && c == '}' {
			node.Name = strings.TrimSpace(node.Name)

			return i
		}
		if outer && (c == '}' || unicode.IsSpace(rune(c))) {
			node.Name = strings.TrimSpace(node.Name)
			node.EndNode = true
			return i
		}

		if outer {
			node.Name += string(c)
		}
		if !outer {
			//fmt.Print(string(c))
			if unicode.IsSpace(rune(c)) {
				subCandid = true

			} else {

				if subCandid {
					subCandid = false
					node.Sub = append(node.Sub, Node{})
					node.EndNode = false
					i += parseQL(input[i:], &node.Sub[len(node.Sub)-1])
					if node.Sub[len(node.Sub)-1].EndNode == true {
						subCandid = true
					}
				}
			}
		}
	}
	return 0
}

func parseParameters(input string, node *Node) int {
	buff := ""
	inside := false
	var char uint8
	for i := 0; i < len(input); i++ {
		if input[i] == '\'' || input[i] == '"' {
			if !inside {
				char = input[i]
				inside = true
				//buff += string(char)
				continue
			}
			if inside && char == input[i] {
				inside = !inside
				//buff += string(char)
				continue
			}
		}
		if !inside && (input[i] == ')' || input[i] == ',') {
			parseParam(buff[1:], node)
			buff = ""
		}
		if !inside && input[i] == ')' {
			return i
		}
		buff += string(input[i])
	}
	return 0
}

func parseParam(input string, node *Node) {
	chunks := strings.SplitN(input, ":", 2)
	if len(chunks) == 2 {
		node.Params[strings.TrimSpace(chunks[0])] = strings.TrimSpace(chunks[1])
	}
}
