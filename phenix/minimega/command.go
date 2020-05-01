// Taken (almost) as-is from minimega/miniweb.

package minimega

import (
	"fmt"
	"strings"
)

type Command struct {
	Command   string
	Columns   []string
	Filters   []string
	Namespace string
}

func NewCommand() *Command {
	return new(Command)
}

func NewNamespacedCommand(ns string) *Command {
	return &Command{Namespace: ns}
}

func (c *Command) String() string {
	cmd := c.Command

	// apply filters first so we don't need to worry about the columns not
	// including the filtered fields.
	for _, f := range c.Filters {
		cmd = fmt.Sprintf(".filter %v %v", f, cmd)
	}

	if len(c.Columns) > 0 {
		columns := make([]string, len(c.Columns))

		// quote all the columns in case there are spaces
		for i := range c.Columns {
			columns[i] = fmt.Sprintf("%q", c.Columns[i])
		}

		cmd = fmt.Sprintf(".columns %v %v", strings.Join(columns, ","), cmd)
	}

	// if there's a namespace, use it
	if c.Namespace != "" {
		cmd = fmt.Sprintf("namespace %q %v", c.Namespace, cmd)
	}

	// don't record command in history
	cmd = fmt.Sprintf(".record false %v", cmd)

	fmt.Printf("built command: `%v`\n", cmd)
	return cmd
}
