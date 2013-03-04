package bintree

import "fmt"

type Tree struct {
	Key string
	Value interface{}
	left *Tree
	right *Tree
}

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

func (tree *Tree) First() (string, interface{}) {
	if tree == nil {
		panic("cannot find first key of empty Tree")
	}
	
	if tree.left == nil {
		return tree.Key, tree.Value
	}
	
	return tree.left.First()
}

func (tree *Tree) Last() (string, interface{}) {
	if tree == nil {
		panic("cannot find last key of empty Tree")
	}
	
	if tree.right == nil {
		return tree.Key, tree.Value
	}
	
	return tree.right.Last()
}

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
