package main

import (
	"errors"
	"fmt"
	// "math"
	"MyMath"
	. "Pratice"
)

func main() {
	// var d = 3
	// e := 4

	// d, e = e, d

	// fmt.Printf("d=%d,e=%d", d, e)

	// _, _, nickName := getName()
	// fmt.Printf("nickNmae=%s", nickName)

	// const (
	// 	pi   float64 = 3.1415926
	// 	a, b         = 1, 2
	// )

	//类似于别的语言中的枚举
	// const (
	// 	x1 = iota
	// 	x2
	// 	x3
	// 	x4 = 5 * iota
	// 	x5
	// )
	// fmt.Printf("x1=%d,x2=%d,x3=%d,x4=%d，x5=%d", x1, x2, x3, x4, x5)

	// var value1 complex64 = 3.2 + 1.2i
	// value2 := 3.2 + 12i
	// value3 := complex(3.2, 12) //value2 == value3

	// str := "hello world"
	// // ch := str[0]
	// for i := 0; i < len(str); i++ {
	// 	fmt.Printf("st1[%d]=%c\n", i, str[i])
	// }
	// fmt.Println("")
	// //两者等同
	// for i, v := range str {
	// 	fmt.Printf("str2[%d]=%c\n", i, v)
	// }

	// var array1 = [5]int{1, 2, 3, 4, 5}
	// // array2 := [5]int{6, 7, 8, 9, 10}
	// array2 := addArrayValue(array1)

	// //基于数组的slice
	// mySlice := array1[1:3] //first:last  first < x <= last
	// // mySlice := array1[:3]
	// // mySlice := array1[3:]
	// //直接创建切片数组
	// mySlice2 := make([]int, 5, 10)
	// mySlice3 := []int{1, 2, 3, 4, 5}

	// fmt.Println(cap(mySlice2), cap(mySlice3))

	// for _, v := range mySlice {
	// 	fmt.Println(v, " ")
	// }

	// for i, v := range array1 {
	// 	fmt.Println(i, v)
	// }

	// for i, v := range array2 {
	// 	fmt.Println(i, v)
	// }

	var myMap map[string]PersonInfo
	myMap = make(map[string]PersonInfo, 100)

	myMap2 := map[string]PersonInfo{
		"game":  PersonInfo{"1", "soso1", "hangz"},
		"game4": PersonInfo{"4", "soso2", "hangz"},
	}

	myMap["game1"] = PersonInfo{"2", "hehe", "hangzhou"}
	myMap2["game2"] = PersonInfo{"3", "heihei", "hangzhou"}

	// delete(myMap, "game1")
	findResult, ok := myMap["game1"]
	if ok {
		fmt.Println(findResult.id, findResult.name, findResult.address)
	} else {
		fmt.Println("can not find this")
	}

	result, err := addNum(1, -2)

	fmt.Println(result, err)

	MyMath.MySort(12, 25, 2, 8, 99, 16, 1)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("runtime error caught: %v", r)
			fmt.Println("")
		}
	}()

	// myDefer()

	defer fmt.Println("run here1?")

	// showIntger(1, 2, 3, 4, 5)

	myPrint(1, 2, "haha", 2.334)

	myFunc := func(x, y int) int {
		return x + y
	}
	defer fmt.Println("run here2?")
	fmt.Println(myFunc(1, 2))

	// switch findResult {
	// case true:
	// 	fmt.Println(ok)
	// 	fallthrough
	// 	break
	// case false:
	// 	fmt.Println(ok)
	// 	break
	// default:
	// 	fmt.Println("impossible")
	// 	break
	// }

	callShowArray := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	odd := showNeed(isOdd, callShowArray)

	fmt.Println("odd is", odd)

	even := showNeed(isEven, callShowArray)

	fmt.Println("even is", even)

	var soso PersonInfo
	soso.id, soso.address, soso.name = "1", "abc", "soso"

	yanyan := PersonInfo{id: "2", address: "bbc", name: "yanyan"}

	may := PersonInfo{"3", "may", "def"}

	fmt.Println(soso, yanyan, may)

	student := Student{Human{"soso", 20, "abc"}, "five high school"}
	employee := Employee{Human{"yu", 33, "bbc"}, "bianfeng"}

	student.SayHi(1)
	employee.SayHi(1)

	student.SayHi(1)
	employee.SayHi(1)

	// var box Box

	boxes := BoxList{
		Box{4, 4, 4, Red},
		Box{10, 10, 1, Yellow},
		Box{1, 1, 20, Black},
		Box{10, 10, 1, Blue},
		Box{10, 30, 1, White},
		Box{20, 20, 20, Yellow},
	}

	fmt.Println("getBiggestColor为", boxes.GetBiggestColor().ColorToString())
}

type Human struct {
	name  string
	age   int
	phone string
}

type Student struct {
	Human
	school string
}

type Employee struct {
	Human
	company string
}

func (h *Human) SayHi(value int) (ret int) {
	h.age += value
	fmt.Printf("human method name=%s,age=%d,phone=%s \n", h.name, h.age, h.phone)
	return h.age
}

func (e *Employee) SayHi(value int) (ret int) {
	e.age += value
	fmt.Printf("Employee method name=%s,age=%d,phone=%s \n", e.name, e.age, e.phone)
	return e.age
}

type myFunction func(int) bool

func showNeed(charge myFunction, lists []int) (ret []int) {

	var result []int

	for _, v := range lists {
		if charge(v) {
			result = append(result, v)
		}
	}
	return result

}

func isOdd(value int) bool {

	if value%2 == 0 {
		return false
	}

	return true
}

func isEven(value int) bool {
	if value%2 == 0 {
		return true
	}
	return false
}

type Integer int

func Less(a, b Integer) bool {
	return a < b
}

type PersonInfo struct {
	id      string
	name    string
	address string
}

func myPrint(printList ...interface{}) {

	var a = [5]int{1, 2, 3, 4, 5}
	b := a
	a[2] = 8

	fmt.Println(a, b)

	c := &a
	a[2] = 12
	fmt.Println(a, *c)

	for _, value := range printList {
		switch value.(type) {
		case int:
			fmt.Println(value, "是整型")
		case string:
			fmt.Println(value, "是字符串类型")
		case float64:
			fmt.Println(value, "是浮点型")
		default:
			fmt.Println("糟糕好像不认识呢")
		}
	}
	// panic("i need out here!")

}

func showIntger(values ...int) {
	for _, arg := range values {
		fmt.Println(arg)
	}
}

func addNum(a int, b int) (ret int, err error) {

	if a < 0 || b < 0 {
		err = errors.New("两个值都不能小于0")
		return
	}

	return a + b, nil

}

func addArrayValue(a1 [5]int) (a2 [5]int) {
	for i := 0; i < len(a1); i++ {
		a1[i] += 10
		fmt.Println(i, a1[i])
	}
	return a1
}

// func getName() (a, b, c string) {
// 	return "a", "b", "c"
// }
