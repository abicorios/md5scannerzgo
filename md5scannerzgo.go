package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var mybuffer = "C:\\Windows\\Temp\\md5utils"

func myexe(s ...string) {
	pa(s)
	app := s[0]
	args := s[1:len(s)]
	out, err := exec.Command(app, args...).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out)
}
func myrmtree(imypath string) {
	os.RemoveAll(imypath)
}
func p(s string) string {
	fmt.Println(s)
	return s
}
func pa(s []string) []string {
	fmt.Println(s)
	return s
}
func mymd5(xfile string) string {
	f, err := os.Open(xfile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
func mytype(ipath string) string {
	var ytype string
	fi, err := os.Lstat(ipath)
	if err != nil {
		log.Fatal(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsRegular():
		ytype = "afile"
	case mode.IsDir():
		ytype = "dir"
	default:
		ytype = "it is not file and not dir"
	}
	if ytype == "afile" {
		matched, err2 := regexp.MatchString(".*\\.(7z|zip|rar)$", ipath)
		if err2 != nil {
			log.Fatal(err2)
		}
		if matched {
			ytype = "archive"
		} else {
			ytype = "file"
		}
	}
	return ytype
}
func myfiles(ipath string) []string {
	var result []string
	files, err := ioutil.ReadDir(ipath)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		result = append(result, f.Name())
	}
	return result
}
func readz(ipath string) {
	for _, i := range myfiles(ipath) {
		thisthing := ipath + "\\" + i
		p(thisthing)
		switch mytype(thisthing) {
		case "file":
			p(mymd5(thisthing))
		case "dir":
			readz(thisthing)
		}
	}
}

func main() {
	if len(os.Args) <= 2 {
		p("longpathgo.exe [myrmtree \"C:\\path\\to\\dir\\to\\remove\"|md5 \"C:\\path\\to\\file\"|mytype \"C:\\file\\of\\dir\\or\\archive\"|myfiles \"C:\\path\\to\\dir\"|test \"C:\\some\\path\"]")
		os.Exit(0)
	}
	if os.Args[1] == "myrmtree" {
		myrmtree(os.Args[2])
		if _, err := os.Stat(os.Args[2]); os.IsNotExist(err) {
			p(os.Args[2] + " is deleted")
		}
	}
	switch os.Args[1] {
	case "md5":
		p(mymd5(os.Args[2]))
	case "mytype":
		p(mytype(os.Args[2]))
	case "myfiles":
		pa(myfiles(os.Args[2]))
	case "test":
		readz(os.Args[2])
	}
}
