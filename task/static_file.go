package task

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"inotify_consume/lib"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var StaticFile = new(static)

type static struct{}

func (s *static) AnalyticalFile(ch chan string) {
	listenDir := g.Cfg().GetString("static.dir.listenDir")
	fmt.Println("当前监听文件夹:", listenDir)
	dirExit, err := lib.Common.PathExists(listenDir)
	if err != nil {
		fmt.Println(err)
		ch <- err.Error()
	}
	if !dirExit {
		fmt.Println("目录:<" + listenDir + ">不存在")
		ch <- "目录:<" + listenDir + ">不存在"
	}
	for {
		filepathNamesarray, err := filepath.Glob(filepath.Join(listenDir, "*"))
		if len(filepathNamesarray) <= 0 {
			fmt.Println("static task sleep after 2 Minute")
			time.Sleep(time.Minute * 2)
			continue
		}
		if err != nil {
			ch <- err.Error()
			break
		}
		filepathNames := make([]string, 0)
		for _, val := range filepathNamesarray {
			_, fileName := filepath.Split(val)
			fileTmp := strings.Split(fileName, "_")
			//if len(fileTmp) == 2 && path.Ext(val) == ".zip"{
			//	if fileTmp[0] == "static" {
			//		filepathNames = append(filepathNames,val)
			//	}
			//}
			if path.Ext(val) == ".zip" {
				if fileTmp[0] == "static" {
					filepathNames = append(filepathNames, val)
				}
			}
		}

		if len(filepathNames) <= 0 {
			fmt.Println("static task sleep after 1 Minute")
			time.Sleep(time.Minute * 1)
			continue
		}

		for _, val := range filepathNames {
			fmt.Println("开始解压", val)
			//if _, err := lib.Zip.DeCompressByPath(val, g.Cfg().GetString("static.dir.targetDir")); err != nil {
			//	fmt.Println("解压失败", err)
			//	time.Sleep(time.Second * 1)
			//	continue
			//}

			// todo 密码进配置文件
			if _, err := lib.ZipPlus.UnZip(val, g.Cfg().GetString("zip.password"), g.Cfg().GetString("static.dir.targetDir")); err != nil {
				fmt.Println("解压失败", err)
				time.Sleep(time.Second * 1)
				continue
			}

			file := lib.Common.LoopHandelFile(val)
			fileInfo, _ := file.Stat()
			file.Close()
			fmt.Println("开始备份", val)
			//if err := os.Rename(val,g.Cfg().GetString("static.dir.localBakDir")+"\\"+fileInfo.Name()); err != nil {
			//	fmt.Println("remove",err)
			//}

			if _, err := lib.Common.CopyFile(val, g.Cfg().GetString("static.dir.localBakDir")+lib.Common.PathHandle()+fileInfo.Name()); err != nil {
				fmt.Println("copy fail", err)
			} else {
				if err := os.Remove(val); err != nil {
					fmt.Println("remove fail", err)
				}
			}
			fmt.Println("备份结束", val)
		}
	}
}
