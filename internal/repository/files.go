package repository

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type File struct {
	Size int64
	Hash string
	Path string
}

func DeleteFile(fileName string) error {
	return os.Remove(fileName)
}

func ReadFilesByExt(ext, dir string) []File {
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

func GetHash(fileName string) (string, error) {

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
