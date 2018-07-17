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
	"encoding/csv"
	"time"
)
var gmylog string
var mybuffer = "C:\\Windows\\Temp\\md5utils"
var gto string
var result [][]string
func checkError(message string, err error){
	if err != nil {
		log.Fatal(message,err)
	}
}
func strs(s ...string)string{
	return strings.Join(s," ")
}
func inBuffer(s string) bool {
	return strings.Contains(s, mybuffer)
}
func drop(x string, sep string) string {
	ar := strings.Split(x, sep)
	return strings.Join(ar[0:len(ar)-1], sep)
}
func myexe(s ...string) {
	p(s...)
	app := s[0]
	args := s[1:len(s)]
	out, err := exec.Command(app, args...).Output()
	checkError("Error: myexe cannot run "+strs(s...),err)
	fmt.Printf("%s", out)
}
func myrmtree(imypath string) {
	os.RemoveAll(imypath)
}
func p(s ...string) string {
	s1:= strs(s...)
	fmt.Println(s1)
	gmylog=gmylog+s1+"\r\n"
	return s1
}
func mymd5(xfile string) string {
	f, err := os.Open(xfile)
	checkError("Error: mymd5 cannot open file "+xfile,err)
	defer f.Close()
	h := md5.New()
	_, err2 := io.Copy(h, f)
	checkError("Error: mymd5 cannot calculate md5 for file "+xfile,err2)
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}
func mytype(ipath string) string {
	var ytype string
	fi, err := os.Lstat(ipath)
	checkError("Error: mytype cannot os.Lstat("+ipath+")",err)
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
		checkError("Error: mytype cannot regexp.MatchString(regex, "+ipath+")",err2)
		if matched {
			ytype = "archive"
		} else {
			ytype = "file"
		}
	}
	return ytype
}
func myfiles(ipath string) []string {
	var result0 []string
	files, err := ioutil.ReadDir(ipath)
	checkError("Error: myfiles cannot ioutil.ReadDir("+ipath+")",err)
	for _, f := range files {
		result0 = append(result0, f.Name())
	}
	return result0
}
func isEmpty(s string) bool {
	return len(myfiles(s)) == 0
}
func readz(ipath string) {
	for _, i := range myfiles(ipath) {
		thisthing := ipath + "\\" + i
		p(thisthing)
		imytype := mytype(thisthing)
		switch imytype {
		case "file":
			short := strings.Replace(ipath, mybuffer+"\\", "", 1)
			m := mymd5(thisthing)
			result = append(result, []string{short, i, m})
			if inBuffer(thisthing) {
				myrmtree(thisthing)
			}
		case "dir":
			readz(thisthing)
		//	if(inBuffer(thisthing)){myrmtree(thisthing)}
		case "archive":
			newpath := drop(thisthing, ".")
			newpath = strings.Replace(newpath, os.Args[2], "", 1)
			newpath = mybuffer + newpath
			os.Mkdir(newpath, 0777)
			myexe("7z", "x", thisthing, "-o"+newpath, "-aou")
			readz(newpath)
		}
	}
}

func main() {
	myrmtree(mybuffer)
	os.Mkdir(mybuffer, 0777)
	if len(os.Args) <= 2 {
		p("md5scannerzgo.exe [myrmtree \"C:\\path\\to\\dir\\to\\remove\"|md5 \"C:\\path\\to\\file\"|mytype \"C:\\file\\of\\dir\\or\\archive\"|myfiles \"C:\\path\\to\\dir\"|readz \"C:\\dir\\from\" \"C:\\dir\\to\"|isEmpty \"C:\\some\\dir\"]")
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
		p(myfiles(os.Args[2])...)
	case "readz":
		gto = os.Args[3]
		myfrom := os.Args[2]
		result = append(result, []string{"path", "name", "md5"})
		readz(myfrom)
		fmt.Println(result)
		fmt.Println(len(result))
		amyfrom:=strings.Split(myfrom,"\\")
		mycsv:=gto+"\\"+amyfrom[len(amyfrom)-1]+" ("+time.Now().Format("02.01.2006")+").csv"
		file,err:=os.Create(mycsv)
		checkError("Error: main cannot os.Create("+mycsv+")",err)
		defer file.Close()
		writer:=csv.NewWriter(file)
		defer writer.Flush()
		for _,value:=range result {
			err:=writer.Write(value)
			checkError("Error: main cannot writer.Write("+strs(value...)+")",err)
		}
		mylog:=gto+"\\mylog.txt"
		d:=[]byte(gmylog)
		fmt.Println(gmylog)
		ioutil.WriteFile(mylog,d,0777)
	case "isEmpty":
		if isEmpty(os.Args[2]) {
			p("true")
		} else {
			p("false")
		}
	}
}
