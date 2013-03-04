package bintree

import "fmt"

type BinTree struct {
	Key string
	Value interface{}
	left *BinTree
	right *BinTree
}

func (tree *BinTree) Get(key string) (interface{}, bool) {
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

func (tree *BinTree) First() (string, interface{}) {
	if tree == nil {
		panic("cannot find first key of empty BinTree")
	}
	
	if tree.left == nil {
		return tree.Key, tree.Value
	}
	
	return tree.left.First()
}

func (tree *BinTree) Last() (string, interface{}) {
	if tree == nil {
		panic("cannot find last key of empty BinTree")
	}
	
	if tree.right == nil {
		return tree.Key, tree.Value
	}
	
	return tree.right.Last()
}

func (tree *BinTree) Iter() <-chan *BinTree {
	ch := make(chan *BinTree)
	
	var visit func(*BinTree)
	visit = func(node *BinTree) {
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

func (tree *BinTree) Range(start string, end string) <-chan *BinTree {
	ch := make(chan *BinTree)
	
	started := false
	
	var visit func(*BinTree)
	visit = func(node *BinTree) {
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

func (tree *BinTree) Add(key string, value interface{}) *BinTree {
	if tree == nil {
		return &BinTree{key, value, nil, nil}
	} else if tree.Key == key {
		return &BinTree{key, value, tree.left, tree.right}
	} else if tree.Key < key {
		return &BinTree{tree.Key, tree.Value, tree.left, tree.right.Add(key, value)}
	}
	
	return &BinTree{tree.Key, tree.Value, tree.left.Add(key, value), tree.right}
}

func (tree *BinTree) Remove(key string) *BinTree {
	if tree == nil {
		return nil
	} else if tree.Key == key {
		if tree.left == nil {
			return tree.right
		} else if tree.right == nil {
			return tree.left
		} else {
			replaceKey, replaceValue := tree.left.Last()
			return &BinTree{replaceKey, replaceValue, tree.left.Remove(replaceKey), tree.right}
		}
	} else if tree.Key < key {
		return &BinTree{tree.Key, tree.Value, tree.left, tree.right.Remove(key)}
	}
	
	return &BinTree{tree.Key, tree.Value, tree.left.Remove(key), tree.right}
}

func (tree *BinTree) inspect() {
	var visit func(*BinTree, int) []string
	visit = func(node *BinTree, level int) []string {
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
