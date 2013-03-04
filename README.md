## Trivial Binary Tree in Go ##

### Why? ###

For some reasons, this is how I've started to spend my Sundays. I think of some bit of comp-sci homework
I would have hated to do when I was in college, and I do it in Go. I had been poking around with some
other more interesting data structures when I got the idea to publish this work on GitHub. I figured it
would be best to start with something trivial, so here it is: A Trivial Binary Tree in Go.

### Features ###

This package implements a simple binary search tree:

    var tree *bintree.Tree
    tree = tree.Add("some_key", "some_value")
    tree = tree.Add("some_other_key", "some_other_value")
    
    fmt.Printf("%#v\n", tree.Get("some_other_key")
    
    >> "some_other_value"

The `Add()` and `Remove()` functions do not modify the original tree. Rather, they build a new tree that
reuses pieces of the old tree as possible, thus the tree is a persistent data structure that requires no
read or write locks for concurrent access. One goroutine may be traversing the tree while another modifies
it and neitehr will explode.

Speaking of traversal, `*bintree.Tree` supports two methods for in-order traversal of nodes. The first is
`Iter()` which is used to traverse all nodes in key-order. The second is `Range(start, end)` which is used
to traverse all nodes where `start <= node.key < end`. Both produce a channel and can, therefore, be used
with Go's idiomatic for/range construct. There are, however, certain disadvantages to this bit of trickery
that are mentioned below. First, an example:
    
    var tree *bintree.Tree
    
    // make a simple tree
    keys := []string{"foo", "bar", "zim", "pow"}
    for i, key := range keys {
        tree = tree.Add(key, i)
    }

    // will print "bar" "foo" "pow" "zim"
    for node := range tree.Iter() {
    	fmt.Printf("%#v ", node.Key)
    }
    
    // will print "foo" "pow"
    for node := range tree.Range("foo", "zim") {
    	fmt.Printf("%#v ", node.Key)
    }
    
*Important Caveat:* Using channels for iterators in Go is a pretty nifty trick, but it has some potential
downsides. First, and most importantly, if you stop reading from the iterator channel before its end, you
will end up leaking memory. The goroutine tasked with traversing the binary tree and sending elements
across the channel will hang forever waiting for the next element to be read. For more information about
this issue, check out the following thread on the golang-nuts group:

https://groups.google.com/forum/?fromgroups=#!topic/golang-nuts/bfuwmhbZHYw

The other issue with channel-based iterators is performance. Currently, `Iter()` and `Range()` use an
unbuffered channel which means Go has to switch between goroutines for every single item being traversed.
It would definitely be more efficient to use a buffered channel, but any guess for the size of the channel
would be fairly wild. Long story short: it sure would be nice if Go supported `range` on some built-in
`iterable` interface.

### Interface ###

#### *bintree.Tree ####

Represents any node in a binary tree including a root node. A nil-value is a valid empty tree.

##### Key string #####

The key associated with this node of the tree.

##### Value interface{} #####

The value associated with this node of the tree.

##### Get(key string) (value interface{}, exists bool) #####

Retrieves the value associated with the provided key from the tree. If the key is found, `exists` will be
true, otherwise, `value` will be nil and `exists` will be false.

##### First() (key string, value interface{}) #####

Finds the first (smallest) key and associated value in the tree. Panics on nil tree.

##### Last() (key string, value interface{}) #####

Finds the last (largest) key and associated value in the tree. Panics on nil tree.

##### Iter() <-chan *bintree.Tree #####

Returns a channel that will emit every node in the tree in key-order.

##### Range(start, end) <-chan *bintree.Tree #####

Returns a channel that will emit every node in the tree in key-order where `start <= node.key < end`

##### Add(key string, value interface{}) *bintree.Tree #####

Returns a new tree with the provided value associated with the provided key. DOES NOT MODIFY existing
tree.

##### Remove(key string) *bintree.Tree #####

Returns a new tree with the provided key removed. DOES NOT MODIFY existing tree.


### Status of This Project ###

This package is currently just some fiddling around. It includes some very basic tests, but it hasn't been
deeply validated. Furthermore, there are some practical performance enhancements that should be made
before anyone takes it seriously. For instance, promotion of child in removal case shouldn't always be
the right-most of the left tree.
