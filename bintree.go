/*
Copyright (c) 2013 Will Conant, http://willconant.com/

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package bintree provides a trivial persistent binary search tree
package bintree

import "fmt"

// *Tree represents any node in a binary tree including a root node.
// For most methods, a nil-value of *Tree is a valid tree.
type Tree struct {
	Key string
	Value interface{}
	left *Tree
	right *Tree
}

// Get searches the tree for the given key. If the key is found, Get returns
// the associated value and true. If the key is not found, it returns nil
// and false.
func (tree *Tree) Get(key string) (interface{}, bool) {
	for tree != nil {
		if key == tree.Key {
			return tree.Value, true
		} else if key < tree.Key {
			tree = tree.left
		} else {
			tree = tree.right
		}
	}
	return nil, false
}

// First returns the left-most key and associated value.
// First panics if tree is nil.
func (tree *Tree) First() (string, interface{}) {
	if tree == nil {
		panic("cannot find first key of empty Tree")
	}
	
	if tree.left == nil {
		return tree.Key, tree.Value
	}
	
	return tree.left.First()
}

// Last returns the right-most key and associated value.
// Last panics if tree is nil.
func (tree *Tree) Last() (string, interface{}) {
	if tree == nil {
		panic("cannot find last key of empty Tree")
	}
	
	if tree.right == nil {
		return tree.Key, tree.Value
	}
	
	return tree.right.Last()
}

// Iter returns a channel that will produce every node in the tree in key-order.
func (tree *Tree) Iter() <-chan *Tree {
	ch := make(chan *Tree)
	
	var visit func(*Tree)
	visit = func(node *Tree) {
		if node.left != nil {
			visit(node.left)
		}
		
		ch <- node
		
		if node.right != nil {
			visit(node.right)
		}
	}
	
	go func() {
		if tree != nil {
			visit(tree)
		}
		close(ch)
	}()
	
	return ch
}

// Range returns a channel that will produce every node in the tree in key-order
// where start <= node.Key < end.
func (tree *Tree) Range(start string, end string) <-chan *Tree {
	ch := make(chan *Tree)
	
	started := false
	
	var visit func(*Tree)
	visit = func(node *Tree) {
		if node.left != nil {
			visit(node.left)
		}
		
		if !started {
			if node.Key >= start {
				started = true
				ch <- node
			}
		} else {
			if node.Key < end {
				ch <- node
			} else {
				return
			}
		}
		
		if node.right != nil {
			visit(node.right)
		}
	}
	
	go func() {
		if tree != nil {
			visit(tree)
		}
		close(ch)
	}()
	
	return ch
}

// Add returns a new tree with the given value associated with the given key.
// Add DOES NOT modify the original tree.
func (tree *Tree) Add(key string, value interface{}) *Tree {
	if tree == nil {
		return &Tree{key, value, nil, nil}
	} else if tree.Key == key {
		return &Tree{key, value, tree.left, tree.right}
	} else if tree.Key < key {
		return &Tree{tree.Key, tree.Value, tree.left, tree.right.Add(key, value)}
	}
	
	return &Tree{tree.Key, tree.Value, tree.left.Add(key, value), tree.right}
}

// Remove returns a new tree with the given key removed.
// Remove DOES NOT modify the original tree.
func (tree *Tree) Remove(key string) *Tree {
	if tree == nil {
		return nil
	} else if tree.Key == key {
		if tree.left == nil {
			return tree.right
		} else if tree.right == nil {
			return tree.left
		} else {
			replaceKey, replaceValue := tree.left.Last()
			return &Tree{replaceKey, replaceValue, tree.left.Remove(replaceKey), tree.right}
		}
	} else if tree.Key < key {
		return &Tree{tree.Key, tree.Value, tree.left, tree.right.Remove(key)}
	}
	
	return &Tree{tree.Key, tree.Value, tree.left.Remove(key), tree.right}
}

func (tree *Tree) inspect() {
	var visit func(*Tree, int) []string
	visit = func(node *Tree, level int) []string {
		if node == nil {
			return nil
		}
		
		line := ""
		for i := 0; i < level; i++ {
			line += ". "
		}
		
		line += fmt.Sprintf("%#v %#v", node.Key, node.Value)
		
		beforeLines := visit(node.left, level + 1)
		afterLines := visit(node.right, level + 1)
		
		lines := make([]string, 0, len(beforeLines) + len(afterLines) + 1)
		
		lines = append(lines, beforeLines...)
		lines = append(lines, line)
		lines = append(lines, afterLines...)
		
		return lines
	}
	
	lines := visit(tree, 0)
	for _, line := range lines {
		fmt.Print(line)
		fmt.Print("\n")
	}
}
