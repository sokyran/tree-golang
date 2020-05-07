package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

// myFolder : myFolder describes folder and contains it's path, counter of folders being last, level of depth, files list, and parent folder
type myFolder struct {
	path        string
	lastCount   int
	level       int
	prev        *myFolder
	files       map[string]bool
	levelSymbol string
}

func (f *myFolder) setFiles(filesInfo []os.FileInfo) {
	newFiles := make(map[string]bool, 0)
	for _, file := range filesInfo {
		newFiles[file.Name()] = file.IsDir()
	}
	f.files = newFiles
}

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
	files, err := ioutil.ReadDir(path)
	startFolder := myFolder{path: path}
	startFolder.setFiles(files)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if printFiles {
		err := showTreeNoFiles(out, startFolder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	} else {
		err := showTreeNoFiles(out, startFolder)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
	}
	return nil
}

func showTreeNoFiles(out io.Writer, Folder myFolder) error {
	fileNotLastSymbol := "├───"
	fileLastSymbol := "└───"
	files := Folder.files
	lastCount := Folder.lastCount

	folders := make([]string, 0)
	for file, isDir := range files {
		if isDir {
			folders = append(folders, file)
		}
	}
	sort.Strings(folders)

	for idx, folder := range folders {
		fullPath := Folder.path + string(os.PathSeparator) + folder
		folderIsLast := idx == len(folders)-1
		fileLevel := strings.Count(fullPath, string(os.PathSeparator))

		beforeFile := ""
		if folderIsLast {
			beforeFile = fileLastSymbol
			lastCount++
		} else {
			beforeFile = fileNotLastSymbol
		}
		levelSymbol := ""
		if Folder.levelSymbol == "" {
			for i := 1; i < fileLevel; i++ {
				levelSymbol += "│\t"
			}
		} else {
			levelSymbol = Folder.levelSymbol + "│\t"
		}

		if Folder.lastCount > 0 {
			for i := 0; i < Folder.lastCount; i++ {
				parentCount := strings.Count(Folder.levelSymbol, "│")
				childCount := strings.Count(levelSymbol, "│")
				if parentCount == childCount {
					break
				}
				levelReversed := reverse(levelSymbol)
				levelReversed = strings.Replace(levelReversed, "│", "", 1)
				levelSymbol = reverse(levelReversed)
			}
		}

		if !folderIsLast {
			lastCount = 0
		}
		fmt.Fprintf(out, "%v%v%v \n", levelSymbol, beforeFile, folder)

		allFiles, _ := ioutil.ReadDir(fullPath)
		newFolder := myFolder{path: fullPath, lastCount: lastCount, level: fileLevel, prev: &Folder, levelSymbol: levelSymbol}
		newFolder.setFiles(allFiles)

		showTreeNoFiles(out, newFolder)
	}
	return nil
}

// func showTreeWithFiles(out io.Writer, folder myFolder) error {
// 	fileNotLastSymbol := "├───"
// 	fileLastSymbol := "└───"
// 	files := folder.files
// 	usualFiles := make([]string, 0)
// 	for _, file := range files {
// 		if file.IsDir() {
// 			folders = append(folders, file.Name())
// 		} else {
// 			usualFiles = append(usualFiles, file.Name())
// 		}
// 	}

// 	sort.Strings(usualFiles)
// 	folders = append(folders, usualFiles...)
// 	sort.Strings(folders)

// 	for idx, folder := range folders {
// 		fullPath := path + string(os.PathSeparator) + folder
// 		fileIsLast := idx == len(folders)-1
// 		fileLevel := strings.Count(fullPath, string(os.PathSeparator))

// 		beforeFile := ""
// 		if fileIsLast {
// 			beforeFile = fileLastSymbol
// 		} else {
// 			beforeFile = fileNotLastSymbol
// 		}

// 		levelSymbol := ""
// 		for i := 1; i < fileLevel; i++ {
// 			levelSymbol += "│\t"
// 		}

// 		if lastFolderCount > 0 {
// 			levelReversed := reverse(levelSymbol)
// 			levelReversed = strings.Replace(levelReversed, "│", "", lastFolderCount)
// 			levelSymbol = reverse(levelReversed)
// 		}

// 		if !find(usualFiles, folder) {
// 			if fileIsLast {
// 				lastFolderCount++
// 			} else {
// 				lastFolderCount = 0
// 			}
// 		}

// 		if foundInFiles := find(usualFiles, folder); foundInFiles {
// 			fileInfo, _ := os.Stat(fullPath)
// 			size := strconv.FormatInt(fileInfo.Size(), 10) + "b"
// 			if size == "0b" {
// 				size = "empty"
// 			}
// 			fmt.Fprintf(out, "%v%v%v (%v)\n", levelSymbol, beforeFile, folder, size)
// 		} else {
// 			fmt.Fprintf(out, "%v%v%v\n", levelSymbol, beforeFile, folder)
// 		}
// 		showTreeWithFiles(out, fullPath, lastFolderCount)
// 	}
// 	return nil
// }

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
