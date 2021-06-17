package lib

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var Common = new(common)

type common struct {}

func (c *common) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (c *common) LoopHandelFile(file string) (*os.File){
	//var filesize int64
	var loop = func(f string) *os.File{
		file, err := os.Open(f)
		if err != nil {
			return nil
		}
		//for  {
		//
		//	fileinfo, err := file.Stat()
		//	if err != nil {
		//		return nil
		//	}
		//
		//	if fileinfo.Size() != filesize{
		//		filesize = fileinfo.Size()
		//		time.Sleep(time.Millisecond*1001)
		//		continue
		//	}else{
		//		break
		//	}
		//	//if !fileinfo.Mode().IsRegular(){
		//	//	continue
		//	//}else{
		//	//	break
		//	//}
		//}
		return file
	}(file)
	return loop
}

func (c *common)ReadFile(filename string) []string{
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	file.Close()
	return strings.Split(string(buffer), "\n")
}

func (c *common) CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)

	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.CopyN(destination, source, sourceFileStat.Size())
	return nBytes, err
}