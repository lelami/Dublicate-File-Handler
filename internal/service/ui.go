package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
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

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		advance, token, err = bufio.ScanWords(data, atEOF)
		if err == nil {
			if token != nil {
				ind, er := strconv.ParseInt(string(token), 10, 0)
				if er != nil {
					return advance, token, er
				}

				if ind < 0 || int(ind) > maxInd {
					return advance, token, errors.New("wrong index")
				}
			}

			if data[len(token)] == '\n' {
				return 0, token, bufio.ErrFinalToken
			}
		}

		return advance, token, err
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(split)

		fmt.Println("\nEnter file numbers to delete:")

		var inds []int
		for scanner.Scan() {
			indStr := scanner.Text()
			ind, err := strconv.Atoi(indStr)
			if err == nil {
				inds = append(inds, ind)
			}
		}

		if err := scanner.Err(); err == nil && len(inds) != 0 {
			return inds
		}

		fmt.Println("\nWrong format")

	}
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
