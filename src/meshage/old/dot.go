package meshage

import "fmt"

// Dot returns a graphviz 'dotfile' string representing the topology known to a node
func (n *Node) Dot() string {
	var ret string

	ret = fmt.Sprintf("graph %s {\n", n.name)
	ret += "size=\"8,11\";\n"
	ret += fmt.Sprintf("Legend [shape=box, shape=plaintext, label=\"total=%d\"];\n", len(n.mesh))

	// we avoid listing a connection twice by maintaining a list of visited nodes.
	// when emitting edges, we don't list those when the node has already been visited.
	visited := make(map[string]bool)

	for k, v := range n.mesh {
		var color string
		if k == n.name {
			color = "red"
		} else {
			color = "green"
		}
		ret += fmt.Sprintf("%s [style=filled, color=%s];\n", k, color)
		for _, c := range v {
			if _, ok := visited[c]; !ok {
				ret += fmt.Sprintf("%s -- %s\n", k, c)
			}
		}
		ret += "\n"
		visited[k] = true
	}
	ret += "}"
	return ret
}
