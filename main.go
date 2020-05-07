package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	if printFiles {
		err := showTreeWithFiles(out, path, 0)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	} else {
		err := showTreeNoFiles(out, path, 0)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}

func showTreeWithFiles(out io.Writer, path string, lastFolderCount int) error {
	fileNotLastSymbol := "├───"
	fileLastSymbol := "└───"
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	folders := make([]string, 0)
	usualFiles := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		} else {
			usualFiles = append(usualFiles, file.Name())
		}
	}

	sort.Strings(usualFiles)
	folders = append(folders, usualFiles...)
	sort.Strings(folders)

	for idx, folder := range folders {
		fullPath := path + string(os.PathSeparator) + folder
		fileIsLast := idx == len(folders)-1
		fileLevel := strings.Count(fullPath, string(os.PathSeparator))

		beforeFile := ""
		if fileIsLast {
			beforeFile = fileLastSymbol
		} else {
			beforeFile = fileNotLastSymbol
		}

		levelSymbol := ""
		for i := 1; i < fileLevel; i++ {
			levelSymbol += "│\t"
		}

		if lastFolderCount > 0 {
			levelReversed := reverse(levelSymbol)
			levelReversed = strings.Replace(levelReversed, "│", "", lastFolderCount)
			levelSymbol = reverse(levelReversed)
		}

		if !find(usualFiles, folder) {
			if fileIsLast {
				lastFolderCount++
			} else {
				lastFolderCount = 0
			}
		}

		if foundInFiles := find(usualFiles, folder); foundInFiles {
			fileInfo, _ := os.Stat(fullPath)
			size := strconv.FormatInt(fileInfo.Size(), 10) + "b"
			if size == "0b" {
				size = "empty"
			}
			fmt.Fprintf(out, "%v%v%v (%v)\n", levelSymbol, beforeFile, folder, size)
		} else {
			fmt.Fprintf(out, "%v%v%v\n", levelSymbol, beforeFile, folder)
		}
		showTreeWithFiles(out, fullPath, lastFolderCount)
	}
	return nil
}

func showTreeNoFiles(out io.Writer, path string, lastFolderCount int) error {
	fileNotLastSymbol := "├───"
	fileLastSymbol := "└───"
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	folders := make([]string, 0)
	for _, file := range files {
		if file.IsDir() {
			folders = append(folders, file.Name())
		}
	}
	sort.Strings(folders)

	for idx, folder := range folders {
		fullPath := path + string(os.PathSeparator) + folder
		fileIsLast := idx == len(folders)-1
		fileLevel := strings.Count(fullPath, string(os.PathSeparator))

		beforeFile := ""
		if fileIsLast {
			beforeFile = fileLastSymbol
		} else {
			beforeFile = fileNotLastSymbol
		}

		levelSymbol := ""
		for i := 1; i < fileLevel; i++ {
			levelSymbol += "│\t"
		}

		if lastFolderCount > 0 {
			levelReversed := reverse(levelSymbol)
			levelReversed = strings.Replace(levelReversed, "│", "", lastFolderCount)
			levelSymbol = reverse(levelReversed)
		}

		if fileIsLast {
			lastFolderCount++
		} else {
			lastFolderCount = 0
		}
		fmt.Fprintf(out, "%v%v%v\n", levelSymbol, beforeFile, folder)
		showTreeNoFiles(out, fullPath, lastFolderCount)
	}
	return nil
}

func find(list []string, a string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func reverse(s string) (result string) {
	for _, v := range s {
		result = string(v) + result
	}
	return
}
