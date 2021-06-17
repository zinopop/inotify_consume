package lib

import (
	"fmt"
	"io"
	"os"
	"archive/zip"
)

var Zip = new(zipLib)

type zipLib struct {}

/**
@tarFile：压缩文件路径
@dest：解压文件夹
*/
func (z *zipLib)DeCompressByPath(tarFile, dest string) ([]string,error) {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return nil,err
	}
	defer srcFile.Close()
	return DeCompress(srcFile, dest)
}

/**
@zipFile：压缩文件
@dest：解压之后文件保存路径
*/
func DeCompress(srcFile *os.File, dest string) ([]string,error) {
	zipFile, err := zip.OpenReader(srcFile.Name())
	files := make([]string,0)
	if err != nil {
		fmt.Println("Unzip File Error：", err.Error())
		return nil,err
	}
	defer zipFile.Close()
	for _, innerFile := range zipFile.File {
		info := innerFile.FileInfo()
		if info.IsDir() {
			err = os.MkdirAll(innerFile.Name, os.ModePerm)
			if err != nil {
				fmt.Println("Unzip File Error : " + err.Error())
				return nil,err
			}
			continue
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			fmt.Println("Unzip File Error : " + err.Error())
			continue
		}
		defer srcFile.Close()
		newFile, err := os.Create(dest+"\\"+innerFile.Name)
		if err != nil {
			fmt.Println("Unzip File Error : " + err.Error())
			continue
		}
		io.Copy(newFile, srcFile)
		files = append(files, dest+"\\"+innerFile.Name)
		newFile.Close()
	}
	return files,nil
}
