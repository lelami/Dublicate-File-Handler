package app

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type File struct {
	Size int64
	Hash string
	Path string
}

func Run() {
	if len(os.Args) != 2 {
		fmt.Println("Directory is not specified")
		return
	}
	dir := os.Args[1]

	fileExtension := readExtension()
	order := readSortingOrder()

	extFiles := readFilesByExt(fileExtension, dir)

	sizeMap := make(map[int64][]File)
	for _, file := range extFiles {
		sizeMap[file.Size] = append(sizeMap[file.Size], file)
	}

	var sizeKeys []int64
	for size, files := range sizeMap {
		if len(files) > 1 {
			sort.Slice(files, func(i, j int) bool {
				if order == 1 {
					return files[i].Path > files[j].Path
				}
				return files[i].Path < files[j].Path
			})
			sizeKeys = append(sizeKeys, size)
			sizeMap[size] = files
		} else {
			delete(sizeMap, size)
		}
	}

	sort.Slice(sizeKeys, func(i, j int) bool {
		if order == 1 {
			return sizeKeys[i] > sizeKeys[j]
		}
		return sizeKeys[i] < sizeKeys[j]
	})

	for _, size := range sizeKeys {

		fmt.Println()
		fmt.Println(size, "bytes")
		for _, file := range sizeMap[size] {
			fmt.Println(file.Path)
		}

	}

	if !readFindDublicates() {
		return
	}

	sizeHashMap := make(map[int64]map[string][]File)
	for size, files := range sizeMap {
		hashMap, ok := sizeHashMap[size]
		if !ok {
			hashMap = make(map[string][]File)
		}
		for _, file := range files {
			hash, err := getHash(file.Path)
			if err != nil {
				return
			}
			file.Hash = hash
			hashMap[hash] = append(hashMap[hash], file)
			sizeHashMap[size] = hashMap
		}
	}

	for _, hashMap := range sizeHashMap {
		for hash, files := range hashMap {
			if len(files) > 1 {
				sort.Slice(files, func(i, j int) bool {
					if order == 1 {
						return files[i].Path > files[j].Path
					}
					return files[i].Path < files[j].Path
				})

				hashMap[hash] = files

			} else {
				delete(hashMap, hash)
			}
		}
	}

	k := 0
	orderMap := make(map[int]File)
	for _, size := range sizeKeys {
		hashMap, ok := sizeHashMap[size]
		if !ok {
			continue
		}
		fmt.Println()
		fmt.Println(size, "bytes")

		for hash, files := range hashMap {
			fmt.Println("Hash:", hash)
			for _, file := range files {
				k++
				orderMap[k] = file
				fmt.Printf("%v. %v\n", k, file.Path)
			}
		}
	}

	inds := readDeleting(k)
	if len(inds) == 0 {
		return
	}

	var total int64
	for _, ind := range inds {
		file := orderMap[ind]
		if err := deleteFile(file.Path); err == nil {
			total += file.Size
		}
	}

	fmt.Println("\nTotal freed up space:", total, "bytes")
}

func deleteFile(fileName string) error {
	return os.Remove(fileName)
}

func getHash(fileName string) (string, error) {

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return "", err
	}

	hf := md5.New()
	_, err = io.Copy(hf, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hf.Sum(nil)), nil
}

func readExtension() string {
	fmt.Println("\nEnter file format:")
	var ext string
	fmt.Scanln(&ext)
	return ext
}

func readSortingOrder() int {

	fmt.Println("\nSize sorting options:\n1. Descending\n2. Ascending")
enterSorting:
	fmt.Println("\nEnter a sorting option:")
	var order int
	fmt.Scanln(&order)

	if order > 2 || order < 1 {
		fmt.Println("\nWrong option")
		goto enterSorting
	}
	return order
}

func readDeleting(maxInd int) []int {

enterChoise:
	fmt.Println("\nDelete files?")

	var ans string
	fmt.Scanln(&ans)

	switch ans {
	case "yes":
	case "no":
		return nil
	default:
		fmt.Println("\nWrong option")
		goto enterChoise
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

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

func readFindDublicates() bool {

enter:
	fmt.Println("\nCheck for duplicates?")

	var ans string
	fmt.Scanln(&ans)

	switch ans {
	case "yes":
		return true
	case "no":
		return false
	default:
		fmt.Println("\nWrong option")
		goto enter
	}
}

func readFilesByExt(ext, dir string) []File {
	var files []File

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if (ext == "" && !info.IsDir()) || filepath.Ext(path) == "."+ext {
			files = append(files, File{
				Size: info.Size(),
				Path: path,
			})
		}
		return nil

	}); err != nil {
		log.Fatal(err)
	}

	return files
}
