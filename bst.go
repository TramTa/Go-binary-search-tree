package main

import ( 
	"os"
	"bufio"
	"fmt"

	"flag"
	"time"

	"strings"
	"strconv"

	"sync" 
	// "reflect"
	// io/ioutil 
)

// ================================================ Args

type Args struct {
	hash_workers  *int  // flag returns ptr 
	data_workers *int 
	comp_workers *int 
	input_file *string 
}

func process_args() Args {
	ff := flag.String("input", "", "input")
	hw := flag.Int("hash-workers", 1, "hash-workers")
	dw := flag.Int("data-workers", 0, "data-workers")
	cw := flag.Int("comp-workers", 0, "comp-workers")
	flag.Parse()

	return Args{ hash_workers: hw, data_workers: dw, comp_workers: cw, input_file: ff }
}

// ================================================ Node 

type Node struct {
	val int
	leftchild *Node
	rightchild *Node
}


func (node *Node) in_order_traversal(in_order_list []int) []int {
	if node != nil {
		in_order_list = node.leftchild.in_order_traversal(in_order_list)
		in_order_list = append(in_order_list, node.val)
		in_order_list = node.rightchild.in_order_traversal(in_order_list)
	}
	return in_order_list
}


func (node *Node) get_hash() int {
	hash := 1 // var hash int = 1 
	var in_order_list []int 
	in_order_list = node.in_order_traversal( in_order_list ) // node.in_order_traversal([]int)
	var new_value int

	for _, value := range in_order_list {  
		new_value = value + 2
		hash = (hash*new_value + new_value) % 1000
	}
	return hash 
}


func (node *Node) insert_value (v int){
	if v < node.val {
		if node.leftchild == nil {
			node.leftchild = &Node{ val: v }
		} else {
			node.leftchild.insert_value(v)
		}
	} else {
		// v > node.val 
		if node.rightchild == nil {
			node.rightchild = &Node{ val: v }
		} else {
			node.rightchild.insert_value(v)
		}
	}
}


func compare_2_tree(t1 *Node , t2 *Node) bool {  
	var t1_list, t2_list  []int 
	t1_list = t1.in_order_traversal(t1_list)
	t2_list = t2.in_order_traversal(t2_list)

	// if len(t1_list) != len(t2_list) {
	// 	return false 
	// }

	for i:= 0; i < len(t1_list) ; i++ {
		if t1_list[i] != t2_list[i] {
			return false
		}
	}
	return true 
}


func get_all_trees(filePath *string, trees *[]*Node ) {  
	file, err := os.Open(*filePath)

	if err != nil {
		fmt.Printf("Cannot open the file. Error: %s \n", err)
		os.Exit(1) 
	}
	// defer file.Close() 

	fileScanner := bufio.NewScanner(file)
    fileScanner.Split(bufio.ScanLines)

    var word_list []string
    var tree_id int = 0  // didn't use 

   	for fileScanner.Scan() { 
      	word_list = strings.Fields( fileScanner.Text() ) // fileScanner.Text() is line of 1 tree 
 
	   		var tree *Node

      	for i, word := range word_list {
          	value, _ := strconv.Atoi(word) // _ is err 

          	if i == 0 {
          		tree = &Node{val: value} // create new root 
          		tree_id++

          	} else {
          		tree.insert_value(value) // insert into root 
          	}
          	// fmt.Println(val, reflect.TypeOf(val))
      	}

      	*trees = append(*trees, tree)    
   	}

   	if err = file.Close(); err != nil {
   		fmt.Printf("Cannot close the file. Error: %s \n", err)
   	}
}
  

type IDHashPair struct {
	tree_id int
	hash int 
}


