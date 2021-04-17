package go_ds

import (
	"math/rand"
	"time"
)

// Skip List
// Paper: https://www.cl.cam.ac.uk/teaching/0506/Algorithms/skiplists.pdf
// Time Complexity:
// 	Search: O (logN)
// 	Insert: O (logN)
// 	Delete: O (logN)

// 该实现与原著的区别：
// 每个节点有后退指针
// 每个索引有后退指针
// 索引直接指向节点

// Redis sort set与原著的区别：
// score值可重复
// 对比一个元素需要同时检查它的score和member
// 每个节点带有高度为1层的后退指针，用于从表尾方向向表头方向迭代

// TODO
// 1.索引指向索引
// 2.是否需要头尾节点
// 3.计算insert、delete、search的复杂度
// 4.是否需要双向链表

const (
	maxLevel int = 32
)

type SkipList struct {
	length int64
	level  int
	head   *Node
	tail   *Node
}

type Node struct {
	prev  *Node
	next  *Node
	key   float64
	value interface{}
	level []*nodeLevel
}

type nodeLevel struct {
	forward  *Node
	backward *Node
	span     int64 // TODO
}

func (sl *SkipList) randomLevel() int {
	rand.Seed(time.Now().UnixNano())

	level := 1
	p := rand.Float64()
	for p >= 0.5 && level <= sl.level && level < maxLevel {
		level++
		rand.Seed(time.Now().UnixNano())
		p = rand.Float64()
	}

	return level
}

func (sl *SkipList) upLevel(i int, n *Node) {
	sl.head.level = append(sl.head.level, &nodeLevel{
		forward: n,
		span:    1,
	})
	n.level[i] = &nodeLevel{
		backward: sl.head,
		forward:  sl.tail,
		span:     1,
	}
	sl.tail.level = append(sl.tail.level, &nodeLevel{
		backward: n,
		span:     1,
	})
}

func (sl *SkipList) searchInsertPlace(key float64) *Node {
	// 从左到右搜索，从上到下搜索
	//fmt.Println("Search node for insert place")
	i, startNode := sl.level-1, sl.head
	for i >= 0 {
		//fmt.Println(i)
		forWardNode := startNode.level[i].forward

		if forWardNode.key > key || forWardNode == sl.tail {
			i--
		} else if forWardNode.key < key {
			startNode = forWardNode
		} else {
			return nil
		}
	}

	return startNode
}

func (sl *SkipList) CreateNode(key float64, value interface{}) *Node {
	l := sl.randomLevel()
	//fmt.Println("New node level: ", l)
	return &Node{
		key:   key,
		value: value,
		level: make([]*nodeLevel, l),
	}
}

func (sl *SkipList) DeleteNode(n *Node) {
	sl.Search(n.key)
}

// 思路一：
// 先生成新节点层数 O (logN)
// 再搜索 O (logN)
// 再插入索引 O (level)

// 思路二：
//新节点和各层索引节点逐一比较，确定原链表的插入位置。O（logN）
//把索引插入到原链表。O（1）
//利用抛硬币的随机方式，决定新节点是否提升为上一级索引。结果为“正”则提升并继续抛硬币，结果为“负”则停止。O（logN）
//总体上，跳跃表插入操作的时间复杂度是O（logN），而这种数据结构所占空间是2N，既空间复杂度是 O（N）。
func (sl *SkipList) Insert(n *Node) {
	if sl.length == 0 {
		// 第一个节点插入
		//fmt.Println("【first】")
		// 维护level关系
		for i := 0; i < len(n.level); i++ {
			sl.upLevel(i, n)
		}

		// 维护node关系
		sl.tail.prev = n
		n.next = sl.tail
		n.prev = sl.head
		sl.head.next = n
	} else {
		//fmt.Println("【not first】")
		if sl.head.level[0].forward.key > n.key {
			// 头部插入
			//fmt.Println("【head insert】")
			// 维护level
			for i := 0; i < len(n.level); i++ {
				//fmt.Println(i)
				if i < sl.level {
					// 高度小于等于sl最大高度，从右到左
					sl.head.level[i].forward.level[i].backward = n
					n.level[i] = &nodeLevel{
						backward: sl.head,
						forward:  sl.head.level[i].forward,
						span:     1,
					}
					sl.head.level[i].forward = n
				} else {
					// 高度大于sl最大高度
					sl.upLevel(i, n)
				}
			}

			// 维护node： 从右到左
			sl.head.next.prev = n
			n.next = sl.head.next
			n.prev = sl.head
			sl.head.next = n

		} else if sl.tail.prev.key < n.key {
			// 尾部插入
			//fmt.Println("【tail insert】")
			// 维护level
			for i := 0; i < len(n.level); i++ {
				//fmt.Println(i)
				if i < sl.level {
					// 高度小于等于sl最大高度
					sl.tail.level[i].backward.level[i].forward = n
					n.level[i] = &nodeLevel{
						backward: sl.tail.level[i].backward,
						forward:  sl.tail,
						span:     1,
					}
					sl.tail.level[i].backward = n
				} else {
					// 高度大于sl最大高度
					sl.upLevel(i, n)
				}
			}
			// 维护node： 从左到右
			sl.tail.prev.next = n
			n.prev = sl.tail.prev
			n.next = sl.tail
			sl.tail.prev = n
		} else {
			// 中间插入
			//fmt.Println("【middle insert】")
			// 从左到右搜索，从上到下搜索
			i := sl.level - 1
			startNode := sl.head
			//fmt.Println("Search")
			for i >= 0 {
				//fmt.Println(i)
				forWardNode := startNode.level[i].forward

				if forWardNode.key > n.key || forWardNode == sl.tail {
					i--
				} else if forWardNode.key < n.key {
					startNode = forWardNode
				} else {
					//fmt.Println("Exist!Exist!Exist!Exist!Exist!Exist!Exist!")
					forWardNode.value = n.value
					return
				}
			}

			// 维护 level，startNode，nextNode
			//fmt.Println("Level")
			for i := 0; i < len(n.level); i++ {
				//fmt.Println(i)
				if i < len(startNode.level) {
					// 高度小于等于前节点最大高度
					startNode.level[i].forward.level[i].backward = n
					n.level[i] = &nodeLevel{
						backward: startNode,
						forward:  startNode.level[i].forward,
						span:     1,
					}
					startNode.level[i].forward = n
				} else {
					// 高度大于前节点最大高度，往前寻找更高层节点
					if i < sl.level {
						backWardSearchLevel := i - 1
						backWardNode := startNode.level[backWardSearchLevel].backward
						for backWardNode != sl.head && len(backWardNode.level) <= len(startNode.level) {
							backWardNode = backWardNode.level[backWardSearchLevel].backward
						}

						startNode = backWardNode
						i--
					} else {
						// 大于最大值的情况
						sl.upLevel(i, n)
					}
				}
			}

			// 维护 node，从右到左
			startNode.next.prev = n
			n.next = startNode.next
			n.prev = startNode
			startNode.next = n
		}
	}

	// 维护最大层数
	if len(n.level) > sl.level {
		sl.level = len(n.level)
	}

	// 维护跳表长度
	sl.length++

	return
}

func (sl *SkipList) Delete(key float64) bool {
	if sl.length > 0 {
		node := sl.Search(key)
		if node != nil {
			// 前接连后，自身弹出
			node.prev.next = node.next
			// 维护 level
			for i := 0; i < len(node.level); i++ {
				levelBackwardNode := node.level[i].backward
				levelForwardNode := node.level[i].forward

				levelBackwardNode.level[i].forward = levelForwardNode
				levelForwardNode.level[i].backward = levelBackwardNode
				node.level[i].forward = nil
				node.level[i].backward = nil
				node.level[i] = nil
			}

			node = nil
			sl.length--

			return true
		}
	}

	return false
}

func (sl *SkipList) Search(key float64) *Node {
	// 从左到右搜索，从上到下搜索
	//fmt.Println("Search node by key")
	for i, startNode := sl.level-1, sl.head; i >= 0; {
		//fmt.Println(i)
		forWardNode := startNode.level[i].forward

		if forWardNode.key > key || forWardNode == sl.tail {
			i--
		} else if forWardNode.key < key {
			startNode = forWardNode
		} else {
			return forWardNode
		}
	}
	return nil
}

func NewSkipList() *SkipList {
	return &SkipList{
		length: 0,
		level:  0,
		head: &Node{
			value: nil,
			prev:  nil,
			next:  nil,
			level: make([]*nodeLevel, 0, maxLevel),
		},
		tail: &Node{
			value: nil,
			prev:  nil,
			next:  nil,
			level: make([]*nodeLevel, 0, maxLevel),
		},
	}
}

func FreeSkipList(sl *SkipList) {
	sl = nil
}
