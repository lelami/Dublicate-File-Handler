package app

import (
	"dfg/internal/repository"
	"dfg/internal/service"
	"fmt"
	"os"
	"sort"
)

func Run() {
	if len(os.Args) != 2 {
		fmt.Println("Directory is not specified")
		return
	}
	dir := os.Args[1]

	fileExtension := service.ReadExtension()
	order := service.ReadSortingOrder()

	extFiles := repository.ReadFilesByExt(fileExtension, dir)

	sizeMap := make(map[int64][]repository.File)
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

	if !service.ReadFindDuplicates() {
		return
	}

	sizeHashMap := make(map[int64]map[string][]repository.File)
	for size, files := range sizeMap {
		hashMap, ok := sizeHashMap[size]
		if !ok {
			hashMap = make(map[string][]repository.File)
		}
		for _, file := range files {
			hash, err := repository.GetHash(file.Path)
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
	orderMap := make(map[int]repository.File)
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

	inds := service.ReadDeleting(k)
	if len(inds) == 0 {
		return
	}

	var total int64
	for _, ind := range inds {
		file := orderMap[ind]
		if err := repository.DeleteFile(file.Path); err == nil {
			total += file.Size
		}
	}

	fmt.Println("\nTotal freed up space:", total, "bytes")
}
