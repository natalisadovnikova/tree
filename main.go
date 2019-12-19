package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	//"path/filepath"
	//"strings"
)

func getDir(output io.Writer, dirName string, full bool, recursionNumber int) error {
	dir, err := os.Open(dirName)

	//defer dir.Close()

	if err != nil {
		return err
	}

	fileInfos, err := dir.Readdir(-1)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})
	if err != nil {
		return err
	}
	countInSlice := len(fileInfos)
	for key, fi := range fileInfos {
		fName := fi.Name()
		isDir := fi.IsDir()
		fPath := fmt.Sprintf("%s%s%s", dirName, string(os.PathSeparator), fName)
		size := fi.Size()
		//fmt.Println("├───", fName, isDir, "(", size, ")")

		if isDir || (!isDir && full) {
			prefix := "├───"
			if key+1 == countInSlice {
				prefix = "└───"
			}
			if recursionNumber != 1 {
				for i := 1; i < recursionNumber; i++ {
					prefix = fmt.Sprintf("	%s", prefix)
					//fmt.Fprintln(output, "печатаем пробелы", i)
				}
				prefix = fmt.Sprintf("│%s", prefix)
			}
			fmt.Fprintf(output, "%s%s ( %vb)\n", prefix, fName, size)
		}

		if isDir {
			// fmt.Println("in is dir")
			err = getDir(output, fPath, full, recursionNumber+1)
		}

	}
	if err != nil {
		return err
	}

	return nil
}
func dirTree(output io.Writer, dirName string, full bool) error {

	err := getDir(output, dirName, full, 1)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	//дирректория
	path := os.Args[1]
	//флаг вывода файлов
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
