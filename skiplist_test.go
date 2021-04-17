package go_ds

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSL(t *testing.T) {
	sl := NewSkipList()

	var keyArray []float64
	for i := 0; i < 5; i++ {
		rand.Seed(time.Now().UnixNano())
		lastKey := float64(rand.Intn(10000))
		sl.Insert(sl.CreateNode(lastKey, i))
		keyArray = append(keyArray, lastKey)
		//sl.Insert(sl.CreateNode(float64(rand.Intn(10000)), i))
		//sl.Insert(sl.CreateNode(float64(i), float64(rand.Intn(10000))))
	}

	levelNode := sl.head.level[0]
	for levelNode.forward.value != nil {
		fmt.Println(fmt.Sprintf("key: %f, value: %d, level: %d", levelNode.forward.key, levelNode.forward.value, len(levelNode.forward.level)))
		//fmt.Print(len(levelNode.forward.level), " ")
		levelNode = levelNode.forward.level[0]
	}

	for i := 0; i < len(keyArray); i++ {
		fmt.Println("====== Search key ======", keyArray[i])
		searchNode := sl.Search(keyArray[i])
		fmt.Println(searchNode.value)

		fmt.Println("====== After delete ======", keyArray[i])
		fmt.Println(sl.Delete(searchNode.key))

		levelNode = sl.head.level[0]
		for levelNode.forward.value != nil {
			fmt.Println(fmt.Sprintf("key: %f, value: %d, level: %d", levelNode.forward.key, levelNode.forward.value, len(levelNode.forward.level)))
			//fmt.Print(len(levelNode.forward.level), " ")
			levelNode = levelNode.forward.level[0]
		}
	}

	fmt.Println(sl.Delete(123))
	fmt.Println(sl.Search(123))
}
