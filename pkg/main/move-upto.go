package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

func listFiles(source string, olderFirst bool, filter string) ([]os.FileInfo, error) {
	// Check that the source is a folder
	srcInfo, err := os.Lstat(source)
	if err != nil {return nil, errors.Wrapf(err, "error while getting information about %q", source)}
	if !srcInfo.IsDir() {return nil, fmt.Errorf("The path %q is not a folder.", source)}
	files, err := ioutil.ReadDir(source)
	if err != nil {return nil, errors.Wrapf(err, "error while listing the content of %q", source)}
	if olderFirst {
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})
	}
	list := []os.FileInfo{}
	r := regexp.MustCompile(fmt.Sprintf(".*%s.*", filter))
	for _, file := range files {
		if file.IsDir() {continue}
		if r.MatchString(file.Name()) {list = append(list, file)}
	}
	return list, nil
}

func moveFile(source, filename, target string) error {
	sourceFile := fmt.Sprintf("%s/%s", source, filename)
	targetFile := fmt.Sprintf("%s/%s", target, filename)
	fmt.Printf("Moving file %s to %s\n", sourceFile, targetFile)
	content, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return errors.Wrapf(err, "failed to get the content of file %s", sourceFile)
	}
	err = ioutil.WriteFile(targetFile, content, 0644)
	if err != nil {
		os.Remove(targetFile)
		return errors.Wrapf(err, "failed to write content of file %s", targetFile)
	}
	os.Remove(sourceFile)
	return nil
}

func main() {
	var source = flag.String("source", "", "Defines the source folder from which the files are retrieved")
	var target = flag.String("target", "", "Defines the destination folder where the files will be moved")
	var olderFirst = flag.Bool("older-first", false, "Indicates that you want to move the older files first. Default is false.")
	var amount = flag.Int("amount", 0, "Define the number of gigabytes that you want to move. Default is 0.")
	var filter = flag.String("filter", "", "Set a string that the name must contain. This allows to only move files that are tainted.")

	flag.Parse()

	if *amount == 0 {fmt.Printf("Please supply a size to move (size will be interpreted as gigabytes). Only int values are accepted")}
	list, err := listFiles(*source, *olderFirst, *filter)
	if err != nil {
		fmt.Printf("Error while trying to list the files in folder %q: %+v\n", source, err)
		return
	}
	totalSize := int64(0)
	maximumSize := int64(*amount) * int64(1024 * 1024 * 1024)
	for _, file := range list {
		if file.Size() + totalSize > maximumSize {
			fmt.Printf("Files have been moved in the order supplied. File %q is too big to be moved within the specified space. Stopping here.\n", file.Name())
			return
		}
		totalSize += file.Size()
		err := moveFile(*source, file.Name(), *target)
		if err != nil {
			fmt.Printf("Error while moving file %q: %+v\n. Stopping execution.\n", file.Name(), err)
			os.Exit(1)
		}
	}
}
