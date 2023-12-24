package main

import (
	"fmt"
	"github.com/zombie-k/algorithms/tree/rbtree"
	"github.com/zombie-k/algorithms/tree/util"
	"time"
)

func main() {
	TestRedBlackTreePut()
}

func TestRedBlackTreePut() {
	tree := rbtree.NewTree(util.IntComparator)
	tree.Insert(5, "e")
	tree.Insert(6, "f")
	tree.Insert(7, "g")
	tree.Insert(3, "c")
	tree.Insert(4, "d")
	tree.Insert(1, "x")
	tree.Insert(2, "b")
	tree.Insert(1, "a") //overwrite
	tree.Insert(10, "b")
	tree.Insert(20, "b")
	tree.Insert(23, "b")
	tree.Insert(11, "b")
	tree.Insert(18, "b")




	//fmt.Println(tree)

	key := 9
	k, v, _:= tree.Floor(key)
	fmt.Printf("floor query_key=%d, key=%v, value=%v\n", key, k, v)
	k, v, _= tree.Ceil(key)
	fmt.Printf("ceil query_key=%d, key=%v, value=%v\n", key, k, v)

	for i := 0 ; i<15; i++ {
		tree.Delete(i)
	}

	num := 10000
	for i:=0;i<num;i++{
		tree.Insert(i, fmt.Sprintf("value%d",i))
	}
	start := time.Now()

	minKey, _ := tree.Minimum()
	maxKey, _ := tree.Maximum()
	fmt.Printf("max:%v, min:%v\n", maxKey, minKey)

	for i:=0;i<num+10;i++{
		tree.Delete(i)
	}
	fmt.Printf("delete size:%d, cost:%v\n", num, time.Since(start))

	//fmt.Println(tree)
	fmt.Println(tree.Size())

}

type HashTable struct{
	vv map[int]interface{}

	keys []int
	mapping map[int]interface{}
}
