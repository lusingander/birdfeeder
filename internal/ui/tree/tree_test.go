package tree

import (
	"fmt"
	"strings"
	"testing"

	"github.com/lusingander/birdfeeder/internal/domain"
)

func TestNew(t *testing.T) {
	posts := []*domain.Post{
		{
			Title:      "p1",
			Categories: []string{},
		},
		{
			Title:      "p2",
			Categories: []string{"a"},
		},
		{
			Title:      "p3",
			Categories: []string{"a"},
		},
		{
			Title:      "p4",
			Categories: []string{"a", "b"},
		},
		{
			Title:      "p5",
			Categories: []string{"a", "c"},
		},
		{
			Title:      "p6",
			Categories: []string{"b", "b"},
		},
		{
			Title:      "p7",
			Categories: []string{"c", "d", "e"},
		},
	}

	New(posts)
}

func (n *node) String() string {
	return n.toString(0)
}

func (n *node) toString(depth int) string {
	indent := strings.Repeat("+", depth)
	str := fmt.Sprintf("%s%s\n", indent, n.name)
	for _, child := range n.children {
		str += fmt.Sprintf("%s\n", child.toString(depth+1))
	}
	return str
}
