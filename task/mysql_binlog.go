package task

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"inotify_consume/lib"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var MysqlBinlog = new(mysqlBinlog)

type mysqlBinlog struct{}

func (m *mysqlBinlog) AnalyticalFile(ch chan string) {
	db := g.DB("default")
	listenDir := g.Cfg().GetString("mysql.dir.listenDir")
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
			fmt.Println("mysql task sleep after 2 Minute")
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
			//if len(fileTmp) == 7 && path.Ext(val) == ".zip"{
			//	if fileTmp[0] == "mysql" {
			//		filepathNames = append(filepathNames,val)
			//	}
			//}
			if path.Ext(val) == ".zip" {
				if fileTmp[0] == "mysql" && fileTmp[1] == "binlog" {
					filepathNames = append(filepathNames, val)
				}
			}
		}

		if len(filepathNames) <= 0 {
			fmt.Println("mysql task sleep after 1 Minute")
			time.Sleep(time.Minute * 1)
			continue
		}

		for _, val := range filepathNames {
			fmt.Println("开始解压", val)
			fileNames, err := lib.Zip.DeCompressByPath(val, g.Cfg().GetString("mysql.dir.localBakDir"))
			if err != nil {
				fmt.Println("解压失败", err)
				time.Sleep(time.Second * 1)
				continue
			}
			file := lib.Common.LoopHandelFile(val)
			fileInfo, _ := file.Stat()
			file.Close()
			fmt.Println("开始备份", val)

			if _, err := lib.Common.CopyFile(val, g.Cfg().GetString("mysql.dir.localBakDir")+lib.Common.PathHandle()+fileInfo.Name()); err != nil {
				fmt.Println("copy fail", err)
			} else {
				if err := os.Remove(val); err != nil {
					fmt.Println("remove fail", err)
				}
			}
			//if err := os.Rename(val,g.Cfg().GetString("mysql.dir.localBakDir")+"\\"+fileInfo.Name()); err != nil {
			//	fmt.Println("remove",err)
			//}
			fmt.Println("开始解析压缩包")
			for _, val := range fileNames {
				filearray := lib.Common.ReadFile(val)
				_, filenameall := filepath.Split(val)
				fileTmp := strings.Split(filenameall, "-")
				tableNameTmp := fileTmp[1]
				tableNameArray := strings.Split(tableNameTmp, ".")
				tableName := tableNameArray[0]
				//fmt.Println(filearray)
				flag, _ := db.HasTable(tableName)
				if !flag {
					_, err := db.Exec("CREATE TABLE `" + tableName + "`  (\n  `id` int(11) NOT NULL AUTO_INCREMENT,\n  `mid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '消息id',\n  `cid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '对话id',\n  `message_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '原消息id',\n  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL COMMENT '文字',\n  `media` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '图片hash',\n  `from_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '发送人',\n  `to_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '接收人',\n  `date` int(11) NULL DEFAULT 0 COMMENT '发送时间',\n  `msg_type` int(11) NULL DEFAULT 0 COMMENT '消息类型',\n  `app` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT 'telegram' COMMENT 'app',\n  PRIMARY KEY (`id`) USING BTREE,\n  UNIQUE INDEX `mid`(`mid`) USING BTREE\n) ENGINE = InnoDB AUTO_INCREMENT = 30154 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;")
					fmt.Println("创建表:", err)
				}
				// 批量插入
				fmt.Println("开始解析文件", val)
				inserts := g.List{}
				for _, val := range filearray {
					j := gjson.New(val)
					if j.Map() != nil {
						inserts = append(inserts, j.Map())
					}
				}
				fmt.Println("开始批量插入", tableName)
				_, err := db.Table(tableName).Data(inserts).Replace()
				if err != nil {
					fmt.Println("插入记录错误,改为单条插入", err)
					for _, val := range filearray {
						j := gjson.New(val)
						if j.Map() != nil {
							_, err := db.Table(tableName).Data(j.Map()).Replace()
							if err != nil {
								fmt.Println("单条插入失败", err)
							}
						}
					}
				}
			}
		}
	}
}
