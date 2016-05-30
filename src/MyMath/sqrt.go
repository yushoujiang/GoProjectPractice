package MyMath

import (
	"errors"
	"fmt"
)

// func sqrt(x float64) float64 {

// }

func MySort(numList ...int) (err error) {

	if numList == nil {
		err = errors.New("numList 位 空")
		return
	}

	for i := 0; i < len(numList)-1; i++ {
		for j := i + 1; j < len(numList); j++ {
			if numList[i] > numList[j] {
				numList[i], numList[j] = numList[j], numList[i]
			}
		}
	}
	fmt.Println("")

	for i, v := range numList {
		fmt.Printf("numList[%d]=%d \n", i, v)
	}

	return nil
}
