package bintree

import "testing"

func TestAddGet(t *testing.T) {
	tree := testTree1()
	tree = tree.Add("basic_test", "basic_value")
	
	v, exists := tree.Get("basic_test")
	if v != "basic_value" {
		t.FailNow()
	}
	
	if !exists {
		t.FailNow()
	}
}

func TestRemove(t *testing.T) {
	tree := testTree1()
	tree = tree.Remove("lob")
	
	var keys []string
	for node := range tree.Iter() {
		keys = append(keys, node.Key)
	}
		
	shouldBe := []string{"bar", "foo", "gam", "gim", "jimmy", "lid", "purp", "pzz", "zim", "zom"}
	if len(keys) != len(shouldBe) {
		t.FailNow()
	}
	
	for i, key := range keys {
		if key != shouldBe[i] {
			t.FailNow()
		}
	}
}

func TestRange(t *testing.T) {
	tree := testTree1()
	expect := []string{"foo", "gam", "gim"}
	t.Log(expect)
	for node := range tree.Range("foo", "jimmy") {
		t.Log(node.Key)
		if len(expect) == 0 || node.Key != expect[0] {
			t.Fail()
		}
		expect = expect[1:]
	}
	
	if len(expect) > 0 {
		t.Fail()
	}
}

func testTree1() (tree *BinTree) {
	keys := []string{"foo", "bar", "zim", "purp", "lob", "gim", "jimmy", "lid", "gam", "zom", "pzz"}
	
	for i, key := range keys {
		tree = tree.Add(key, i)
	}
	
	return
}
