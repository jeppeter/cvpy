package main

import (
	"fmt"
	"log"
)

const RB_BLACK int = 1
const RB_RED int = 2

type RBTreeData interface {
	Less(b RBTreeData) bool
	Equal(b RBTreeData) bool
	Stringer() string
	TypeName() string
}

type RBTreeElem struct {
	parent *RBTreeElem
	left   *RBTreeElem
	right  *RBTreeElem
	color  int
	Data   RBTreeData
}

type RBTree struct {
	root  *RBTreeElem
	count int
}

func NewRBTreeElem(data RBTreeData) *RBTreeElem {
	p := &RBTreeElem{}
	p.Data = data
	p.color = RB_RED
	p.parent = nil
	p.left = nil
	p.right = nil
	return p
}

func (elem *RBTreeElem) SetParent(parent *RBTreeElem) {
	elem.parent = parent
	return
}

func (elem *RBTreeElem) GetParent() *RBTreeElem {
	return elem.parent
}

func (elem *RBTreeElem) SetLeft(left *RBTreeElem) {
	elem.left = left
	return
}

func (elem *RBTreeElem) GetLeft() *RBTreeElem {
	return elem.left
}

func (elem *RBTreeElem) SetRight(right *RBTreeElem) {
	elem.right = right
	return
}

func (elem *RBTreeElem) GetRight() *RBTreeElem {
	return elem.right
}

func (elem *RBTreeElem) GetColor() int {
	return elem.color
}

func (elem *RBTreeElem) SetColor(color int) {
	elem.color = color
	return
}

func NewRBTree() *RBTree {
	p := &RBTree{}
	p.root = nil
	p.count = 0
	return p
}

func (rb *RBTree) __find_insert_parent(data RBTreeData, from *RBTreeElem) *RBTreeElem {
	if from.Data.Less(data) {
		if from.GetRight() != nil {
			return rb.__find_insert_parent(data, from.GetRight())
		} else {
			return from
		}
	} else {
		if from.GetLeft() != nil {
			return rb.__find_insert_parent(data, from.GetLeft())
		} else {
			return from
		}
	}
	return nil
}

func (rb *RBTree) find_insert_parent(data RBTreeData) *RBTreeElem {
	if rb.root.Data.Less(data) {
		if rb.root.GetRight() != nil {
			return rb.__find_insert_parent(data, rb.root.GetRight())
		} else {
			return rb.root
		}
	} else {
		if rb.root.GetLeft() != nil {
			return rb.__find_insert_parent(data, rb.root.GetLeft())
		} else {
			return rb.root
		}
	}
	return nil
}

func (rb *RBTree) __get_grandma(elem *RBTreeElem) *RBTreeElem {
	return elem.GetParent().GetParent()
}

func (rb *RBTree) __get_uncle(elem *RBTreeElem) *RBTreeElem {
	return rb.__get_sibling(elem.GetParent())
}

func (rb *RBTree) __is_leaf(elem *RBTreeElem) bool {
	if elem == nil {
		return true
	}
	if elem.GetLeft() != nil || elem.GetRight() != nil {
		return false
	}
	return true
}

func (rb *RBTree) __max_elem(from *RBTreeElem) *RBTreeElem {
	if from.GetRight() != nil {
		return rb.__max_elem(from.GetRight())
	}
	return from
}

func (rb *RBTree) __min_elem(from *RBTreeElem) *RBTreeElem {
	if from.GetLeft() != nil {
		return rb.__min_elem(from.GetLeft())
	}
	return from
}

func (rb *RBTree) __get_sibling(elem *RBTreeElem) *RBTreeElem {
	parent := elem.GetParent()
	if elem == parent.GetLeft() {
		return parent.GetRight()
	} else {
		return parent.GetLeft()
	}
	return nil

}

func (rb *RBTree) __rotate_left(insertp *RBTreeElem) {
	right := insertp.GetRight()
	rb.__replace_node(insertp, right)
	insertp.SetRight(right.GetLeft())
	if right.GetLeft() != nil {
		right.GetLeft().SetParent(insertp)
	}
	right.SetLeft(insertp)
	insertp.SetParent(right)
	return
}

func (rb *RBTree) __rotate_right(insertp *RBTreeElem) {
	left := insertp.GetLeft()
	rb.__replace_node(insertp, left)
	insertp.SetLeft(left.GetRight())
	if left.GetRight() != nil {
		left.GetRight().SetParent(insertp)
	}
	left.SetRight(insertp)
	insertp.SetParent(left)
	return
}

func (rb *RBTree) __rebalanced_case5(insertp *RBTreeElem) {
	grandma := rb.__get_grandma(insertp)
	parent := insertp.GetParent()
	parent.SetColor(RB_BLACK)
	grandma.SetColor(RB_RED)
	if insertp == parent.GetLeft() && parent == grandma.GetLeft() {
		rb.__rotate_right(grandma)
	} else {
		if insertp != parent.GetRight() || parent != grandma.GetRight() {
			panic(fmt.Sprintf("insert (%s) not right of (%s) or parent (%s) not right of grandma (%s)", insertp.Data.Stringer(), parent.Data.Stringer(), parent.Data.Stringer(), grandma.Data.Stringer()))
		}
		rb.__rotate_left(grandma)
	}
}

func (rb *RBTree) __rebalanced_case4(insertp *RBTreeElem) {
	var next *RBTreeElem
	grandma := rb.__get_grandma(insertp)
	parent := insertp.GetParent()
	next = insertp
	if insertp == parent.GetRight() && parent == grandma.GetLeft() {
		rb.__rotate_left(parent)
		next = insertp.GetLeft()
	} else if insertp == parent.GetLeft() && parent == grandma.GetRight() {
		rb.__rotate_right(parent)
		next = insertp.GetRight()
	}
	rb.__rebalanced_case5(next)
	return
}

func (rb *RBTree) __rebalanced_case3(insertp *RBTreeElem) {
	uncle := rb.__get_uncle(insertp)
	if uncle != nil && uncle.GetColor() == RB_RED {
		insertp.GetParent().SetColor(RB_BLACK)
		uncle.SetColor(RB_BLACK)
		grandma := rb.__get_grandma(insertp)
		grandma.SetColor(RB_RED)
		rb.__rebalanced_case1(grandma)
	} else {
		rb.__rebalanced_case4(insertp)
	}
}

func (rb *RBTree) __rebalanced_case2(insertp *RBTreeElem) {
	parent := insertp.GetParent()
	if parent.GetColor() == RB_BLACK {
		return
	} else {
		rb.__rebalanced_case3(insertp)
	}
	return
}

func (rb *RBTree) __rebalanced_case1(insertp *RBTreeElem) {
	parent := insertp.GetParent()
	if parent == nil {
		insertp.SetColor(RB_BLACK)
		return
	} else {
		rb.__rebalanced_case2(insertp)
	}
	return
}

func (rb *RBTree) Insert(data RBTreeData) error {
	var err error
	err = nil
	if rb.count == 0 {
		rb.root = NewRBTreeElem(data)
		rb.root.SetColor(RB_BLACK)
		rb.count++
	} else {
		parent := rb.find_insert_parent(data)
		if parent != nil {
			insertp := NewRBTreeElem(data)
			if parent.Data.Less(data) {
				parent.SetRight(insertp)
			} else {
				parent.SetLeft(insertp)
			}
			insertp.SetParent(parent)
			rb.__rebalanced_case1(insertp)
			rb.count++
		} else {
			err = fmt.Errorf("can not find parent for (%s)", data.Stringer())
		}

	}
	return err
}

func (rb *RBTree) __find_data_from(data RBTreeData, from *RBTreeElem) *RBTreeElem {
	if from.Data.Equal(data) {
		return from
	} else {
		if from.Data.Less(data) && from.GetLeft() != nil {
			return rb.__find_data_from(data, from.GetLeft())
		} else if from.GetRight() != nil {
			return rb.__find_data_from(data, from.GetRight())
		}
	}
	return nil
}

func (rb *RBTree) __find_data(data RBTreeData) *RBTreeElem {
	if rb.count == 0 {
		return nil
	}
	return rb.__find_data_from(data, rb.root)
}

func (rb *RBTree) __get_color(elem *RBTreeElem) int {
	if elem == nil {
		return RB_BLACK
	}
	return elem.GetColor()
}

func (rb *RBTree) __delete_elem_case6(elem *RBTreeElem) {
	sibling := rb.__get_sibling(elem)
	parent := elem.GetParent()
	sibling.SetColor(parent.GetColor())
	parent.SetColor(RB_BLACK)
	if elem == parent.GetLeft() {
		if sibling.GetRight().GetColor() != RB_RED {
			log.Fatalf("(%s) not red", sibling.GetRight().Data.Stringer())
		}
		sibling.GetRight().SetColor(RB_BLACK)
		rb.__rotate_left(parent)
	} else {
		if sibling.GetLeft().GetColor() != RB_RED {
			log.Fatalf("(%s) not red", sibling.GetLeft().Data.Stringer())
		}
		sibling.GetLeft().SetColor(RB_BLACK)
		rb.__rotate_right(parent)
	}
	return
}

func (rb *RBTree) __delete_elem_case5(elem *RBTreeElem) {
	sibling := rb.__get_sibling(elem)
	parent := elem.GetParent()
	if elem == parent.GetLeft() && rb.__get_color(sibling) == RB_BLACK && rb.__get_color(sibling.GetLeft()) == RB_RED && rb.__get_color(sibling.GetRight()) == RB_BLACK {
		sibling.SetColor(RB_RED)
		sibling.GetLeft().SetColor(RB_BLACK)
		rb.__rotate_right(sibling)
	} else if elem == parent.GetRight() && rb.__get_color(sibling) == RB_BLACK && rb.__get_color(sibling.GetRight()) == RB_RED && rb.__get_color(sibling.GetLeft()) == RB_BLACK {
		sibling.SetColor(RB_RED)
		sibling.GetRight().SetColor(RB_BLACK)
		rb.__rotate_left(sibling)
	}

	rb.__delete_elem_case6(elem)
	return
}

func (rb *RBTree) __delete_elem_case4(elem *RBTreeElem) {
	sibling := rb.__get_sibling(elem)
	parent := elem.GetParent()
	if rb.__get_color(parent) == RB_RED && rb.__get_color(sibling) == RB_BLACK && rb.__get_color(sibling.GetLeft()) == RB_BLACK && rb.__get_color(sibling.GetRight()) == RB_BLACK {
		sibling.SetColor(RB_RED)
		parent.SetColor(RB_BLACK)
	} else {
		rb.__delete_elem_case5(elem)
	}
	return
}

func (rb *RBTree) __delete_elem_case3(elem *RBTreeElem) {
	sibling := rb.__get_sibling(elem)
	if rb.__get_color(elem.GetParent()) == RB_BLACK && rb.__get_color(sibling) == RB_BLACK && rb.__get_color(sibling.GetLeft()) == RB_BLACK && rb.__get_color(sibling.GetRight()) == RB_BLACK {
		sibling.SetColor(RB_RED)
		rb.__delete_elem_case1(elem.GetParent())
	} else {
		rb.__delete_elem_case4(elem)
	}
	return
}

func (rb *RBTree) __delete_elem_case2(elem *RBTreeElem) {
	sibling := rb.__get_sibling(elem)
	if rb.__get_color(sibling) == RB_RED {
		parent := elem.GetParent()
		parent.SetColor(RB_RED)
		sibling.SetColor(RB_BLACK)
		if elem == parent.GetLeft() {
			rb.__rotate_left(parent)
		} else {
			rb.__rotate_right(parent)
		}
	}

	rb.__delete_elem_case3(elem)

}

func (rb *RBTree) __delete_elem_case1(elem *RBTreeElem) {
	if elem.GetParent() == nil {
		return
	} else {
		rb.__delete_elem_case2(elem)
	}
}

func (rb *RBTree) __replace_node(oldelem *RBTreeElem, newelem *RBTreeElem) {
	var parent *RBTreeElem
	parent = nil
	if oldelem.GetParent() == nil {
		rb.root = newelem
	} else {
		parent = oldelem.GetParent()
		if oldelem == parent.GetLeft() {
			parent.SetLeft(newelem)
		} else {
			parent.SetRight(newelem)
		}
	}
	if newelem != nil {
		newelem.SetParent(parent)
	}
	return
}

func (rb *RBTree) __delete_one(elem *RBTreeElem) (cnt int, err error) {
	var chld *RBTreeElem
	deleteone := elem
	if elem.GetLeft() != nil && elem.GetRight() != nil {
		pred := rb.__max_elem(elem.GetLeft())
		elem.Data = pred.Data
		deleteone = pred
		//log.Printf("change dete to (%s)", deleteone.Data.Stringer())
	}

	if deleteone.GetRight() != nil && deleteone.GetLeft() != nil {
		log.Fatalf("not set (%s) has both left and right", deleteone.Data.Stringer())
	}

	if deleteone.GetRight() == nil {
		chld = deleteone.GetLeft()
	} else {
		chld = deleteone.GetRight()
	}

	if deleteone.GetColor() == RB_BLACK {
		deleteone.SetColor(rb.__get_color(chld))
		rb.__delete_elem_case1(deleteone)
	}

	rb.__replace_node(deleteone, chld)

	if deleteone.GetParent() == nil && chld != nil {
		/*root should be black*/
		chld.SetColor(RB_BLACK)
	}

	rb.count--
	if false {
		err = rb.__verify()
		if err != nil {
			return 0, err
		}
	}

	return rb.count, nil
}

func (rb *RBTree) Delete(data RBTreeData) (cnt int, err error) {
	elem := rb.__find_data(data)
	if elem == nil {
		cnt = rb.count
		err = fmt.Errorf("can not find (%s)", data.Stringer())
		return
	}

	return rb.__delete_one(elem)
}

func (rb *RBTree) GetMin() RBTreeData {
	if rb.root == nil {
		return nil
	}
	elem := rb.__min_elem(rb.root)
	if elem == nil {
		return nil
	}

	_, err := rb.__delete_one(elem)
	if err != nil {
		log.Fatal("%s", err.Error())
	}
	return elem.Data
}

func (rb *RBTree) GetMax() RBTreeData {
	if rb.root == nil {
		return nil
	}
	elem := rb.__max_elem(rb.root)
	if elem == nil {
		return nil
	}
	_, err := rb.__delete_one(elem)
	if err != nil {
		log.Fatal("%s", err.Error())
	}
	return elem.Data
}

func (rb *RBTree) __verify_property1(elem *RBTreeElem) error {
	var err error
	if elem == nil {
		return nil
	}

	if elem.GetColor() == RB_BLACK || elem.GetColor() == RB_RED {
		err = rb.__verify_property1(elem.GetLeft())
		if err != nil {
			return err
		}

		return rb.__verify_property1(elem.GetRight())
	}
	log.Fatalf("(%s) color is (%d)", elem.Data.Stringer(), elem.GetColor())
	return fmt.Errorf("(%s) color is (%d)", elem.Data.Stringer(), elem.GetColor())
}

func (rb *RBTree) __verify_property2(elem *RBTreeElem) error {
	if elem == nil {
		return nil
	}

	if elem.GetColor() != RB_BLACK {
		log.Fatalf("(%s) color not black (%d)", elem.Data.Stringer(), elem.GetColor())
		return fmt.Errorf("(%s) color not black (%d)", elem.Data.Stringer(), elem.GetColor())
	}
	return nil
}

func (rb *RBTree) __verify_property4(elem *RBTreeElem) error {
	if elem == nil {
		return nil
	}

	if elem.GetColor() == RB_RED {
		if rb.__get_color(elem.GetLeft()) != RB_BLACK {
			log.Fatalf("(%s).left(%s) = (%d) not black", elem.Data.Stringer(), elem.GetLeft().Data.Stringer(), elem.GetLeft().GetColor())
			return fmt.Errorf("(%s).left(%s) = (%d) not black", elem.Data.Stringer(), elem.GetLeft().Data.Stringer(), elem.GetLeft().GetColor())
		}

		if rb.__get_color(elem.GetRight()) != RB_BLACK {
			log.Fatalf("(%s).right(%s) = (%d) not black", elem.Data.Stringer(), elem.GetRight().Data.Stringer(), elem.GetRight().GetColor())
			return fmt.Errorf("(%s).right(%s) = (%d) not black", elem.Data.Stringer(), elem.GetRight().Data.Stringer(), elem.GetRight().GetColor())
		}

	}
	err := rb.__verify_property4(elem.GetLeft())
	if err != nil {
		return err
	}
	return rb.__verify_property4(elem.GetRight())
}

func (rb *RBTree) __verify_property5_recursive(elem *RBTreeElem, cnt int, setcnt *int) error {
	if elem == nil {
		if *setcnt == -1 {
			*setcnt = cnt
		}
		if *setcnt != cnt {
			log.Fatalf("(%s) black count (%d) != setcnt (%d)", elem.Data.Stringer(), cnt, *setcnt)
			return fmt.Errorf("(%s) black count (%d) != setcnt (%d)", elem.Data.Stringer(), cnt, *setcnt)
		}
		return nil
	}

	if elem.GetColor() == RB_BLACK {
		cnt++
	}

	err := rb.__verify_property5_recursive(elem.GetLeft(), cnt, setcnt)
	if err != nil {
		return err
	}
	return rb.__verify_property5_recursive(elem.GetRight(), cnt, setcnt)
}

func (rb *RBTree) __verify_property5(elem *RBTreeElem) error {
	setcnt := -1
	return rb.__verify_property5_recursive(elem, 0, &setcnt)
}

func (rb *RBTree) __verify() error {
	var err error
	err = rb.__verify_property1(rb.root)
	if err != nil {
		return err
	}

	err = rb.__verify_property2(rb.root)
	if err != nil {
		return err
	}
	err = rb.__verify_property4(rb.root)
	if err != nil {
		return err
	}
	err = rb.__verify_property5(rb.root)
	if err != nil {
		return err
	}
	return nil
}
