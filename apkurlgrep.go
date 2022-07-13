/*
 made with love by @ndelphit 5/2020
*/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/ndelphit/apkurlgrep/command/apktool"
	dependency "github.com/ndelphit/apkurlgrep/command/dependency"
	"github.com/ndelphit/apkurlgrep/directory"
	"github.com/ndelphit/apkurlgrep/extractor"
)

func main() {

	parser := argparse.NewParser("apkurlgrep", "ApkUrlGrep - Extract endpoints from APK files")
	apk := parser.String("a", "apk", &argparse.Options{Required: true, Help: "Input a path to APK file."})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(-1)
	}

	var baseApk = *apk
	var tempDir = directory.CreateTempDir()

	// get dir from file
	p := filepath.Join(baseApk)
	root := filepath.Dir(p)

	// get all files in file dir
	files := getFilesInDirectory(root)

	for _, file := range files {
		doTheWalk(file, tempDir, getFilePath(file))
		fmt.Println("\n" + file + "\n ")
	}

}

func getFilePath(file string) string {
	var filename = file
	// Create path for new .txt file
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext) + ".txt"

}

func doTheWalk(baseApk string, tempDir string, path string) {
	dependency.AreAllReady()
	apktool.RunApktool(baseApk, tempDir)
	extractor.Extract(tempDir, path)
	directory.RemoveTempDir(tempDir)
}

func getFilesInDirectory(root string) []string {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {

			fmt.Println(err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".apk" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return files

}
