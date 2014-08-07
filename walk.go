package viewmaster

import (
	"errors"
	"fmt"
	"text/template/parse"
)

// walk recurses the node tree starting at n and calls the
// walkFn for each node found. Errors returned by walkFn
// cause iteration to stop and the error is returned by walk.
func walk(n parse.Node, walkFn func(n parse.Node) error) error {
	// Inspired by https://codereview.appspot.com/4918041/patch/13001/11004

	if n == nil || walkFn == nil {
		return nil
	}

	var err error

	err = walkFn(n)
	if err != nil {
		return err
	}

	switch n := n.(type) {
	case nil:
	case *parse.ActionNode:
		if n.Pipe != nil {
			err = walk(n.Pipe, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.BoolNode:
	case *parse.CommandNode:
		for _, arg := range n.Args {
			err = walk(arg, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.DotNode:
	case *parse.FieldNode:
	case *parse.IdentifierNode:
	case *parse.IfNode:
		if n.Pipe != nil {
			err = walk(n.Pipe, walkFn)
			if err != nil {
				return err
			}
		}
		if n.List != nil {
			err = walk(n.List, walkFn)
			if err != nil {
				return err
			}
		}
		if n.ElseList != nil {
			err = walk(n.ElseList, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.ListNode:
		for _, node := range n.Nodes {
			err = walk(node, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.NumberNode:
	case *parse.PipeNode:
		for _, decl := range n.Decl {
			err = walk(decl, walkFn)
			if err != nil {
				return err
			}
		}
		for _, cmd := range n.Cmds {
			err = walk(cmd, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.RangeNode:
		if n.Pipe != nil {
			err = walk(n.Pipe, walkFn)
			if err != nil {
				return err
			}
		}
		if n.List != nil {
			err = walk(n.List, walkFn)
			if err != nil {
				return err
			}
		}
		if n.ElseList != nil {
			err = walk(n.ElseList, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.StringNode:
	case *parse.TemplateNode:
		if n.Pipe != nil {
			err = walk(n.Pipe, walkFn)
			if err != nil {
				return err
			}
		}
	case *parse.TextNode:
	case *parse.VariableNode:
	case *parse.WithNode:
		if n.Pipe != nil {
			err = walk(n.Pipe, walkFn)
			if err != nil {
				return err
			}
		}
		if n.List != nil {
			err = walk(n.List, walkFn)
			if err != nil {
				return err
			}
		}
		if n.ElseList != nil {
			err = walk(n.ElseList, walkFn)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unknown node of type " + fmt.Sprintf("%T", n))
	}

	return nil
}
