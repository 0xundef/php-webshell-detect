package context

import (
	"github.com/0xundef/php-webshell-detect/internal/core/ir/stmt"
	"fmt"
	"strings"
)

type Context interface {
	GetLength() int
	GetElementAt(i int) interface{}
	String() string
	GetChild(element interface{}) Context
	GetElement() interface{}
	GetParent() Context
	GetRoot() Context
}
type TrieContext struct {
	parent   *TrieContext
	length   int
	element  interface{}
	children map[interface{}]*TrieContext
	root     Context
}

func (c *TrieContext) GetParent() Context {
	return c.parent
}
func (c *TrieContext) GetRoot() Context {
	if c.parent == nil {
		c.root = c
		return c
	}
	return c.root
}
func (c *TrieContext) GetChild(element interface{}) Context {
	if ctx, ok := c.children[element]; ok {
		return ctx
	} else {
		i := c.length
		c.children[element] = &TrieContext{
			parent:   c,
			length:   i + 1,
			element:  element,
			children: map[interface{}]*TrieContext{},
			root:     c.GetRoot(),
		}
		return c.children[element]
	}
}

func (c *TrieContext) GetLength() int {
	return c.length
}

func (c *TrieContext) GetElementAt(i int) interface{} {
	if i >= c.length {
		return nil
	}
	if i == c.length-1 {
		return c.element
	} else {
		return c.parent.GetElementAt(i)
	}
}

func (c *TrieContext) GetElement() interface{} {
	return c.element
}

func (c *TrieContext) String() string {
	var eles = make([]string, c.length)
	tmp := c
	for i := c.length - 1; i >= 0; i-- {
		if e, ok := tmp.element.(stmt.InvokeStmt); ok {
			eles[i] = e.String()
		}
		tmp = tmp.parent
	}
	return fmt.Sprint("[" + strings.Join(eles, ",") + "]")
}

// util function
func ContextEmpty() Context {
	return CreateEmptyContext()
}

func CreateEmptyContext() Context {
	return &TrieContext{
		parent:   nil,
		length:   0,
		element:  nil,
		children: map[interface{}]*TrieContext{},
	}
}

func ContextMakeLastK(context Context, limit int) Context {
	if limit == 0 {
		return ContextEmpty()
	}
	if limit >= context.GetLength() {
		return context
	}

	elements := make([]interface{}, limit)
	tmp := context
	for i := limit; i > 0; i-- {
		elements[i-1] = tmp.GetElement()
		tmp = tmp.GetParent()
	}
	return ContextMake(elements...)
}

func ContextMake(elemets ...interface{}) Context {
	root := ContextEmpty()
	ctx := root
	for _, element := range elemets {
		ctx = ctx.GetChild(element)
	}
	return ctx
}

func ContextAppend(parent Context, element interface{}, limit int) Context {
	if parent.GetLength() < limit {
		return parent.GetChild(element)
	} else {
		return ContextMakeLastK(parent, limit-1).GetChild(element)
	}
}
