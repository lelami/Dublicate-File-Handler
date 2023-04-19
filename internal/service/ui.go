package service

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadExtension() string {
	fmt.Println("\nEnter file format:")
	var ext string
	fmt.Scanln(&ext)
	return ext
}

func ReadSortingOrder() int {

	fmt.Println("\nSize sorting options:\n1. Descending\n2. Ascending")

	var order int

	for {

		fmt.Println("\nEnter a sorting option:")
		fmt.Scanln(&order)

		if order == 1 || order == 2 {
			return order
		}

		fmt.Println("\nWrong option")
	}
}

func ReadDeleting(maxInd int) []int {

	if !readYesOrNo("Delete files?") {
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines) // TODO придумать твою сплит функцию по словам или по байтам

enterInds:
	fmt.Println("\nEnter file numbers to delete:")

	var inds []int
	for scanner.Scan() {
		choiseStr := scanner.Text()

		strInds := strings.Split(choiseStr, " ")
		for _, sind := range strInds {
			ind, err := strconv.Atoi(sind)
			if err != nil || (ind < 1 || ind > maxInd) {
				fmt.Println("\nWrong format")
				goto enterInds
			}
			inds = append(inds, ind)
		}
		break
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("\nWrong format")
		goto enterInds
	}
	return inds

}

func ReadFindDuplicates() bool {
	return readYesOrNo("Check for duplicates?")
}

func readYesOrNo(question string) bool {

	for {
		fmt.Println("\n" + question)

		var ans string
		fmt.Scanln(&ans)

		switch ans {
		case "yes":
			return true
		case "no":
			return false
		default:
			fmt.Println("\nWrong option")
		}
	}

}
