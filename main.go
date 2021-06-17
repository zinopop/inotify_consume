package main

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"inotify_consume/task"
	"os"
	"time"
)

func main(){
	g.Cfg().SetFileName("config.toml")
	initDb()

	taskChan := taskInit(
		task.StaticFile.AnalyticalFile,
		task.MysqlBinlog.AnalyticalFile,
	)

	// 用select模型阻塞住主线程
	for {
		select{
		case taskName := <- taskChan:
			fmt.Println(taskName)
		default:
			// fmt.Println("没有可执行的任务")
			time.Sleep(time.Second*1)

		}
	}
}

// 异步任务初始化
func taskInit(callback... func(ic chan string)) chan string {
	ch := make(chan string)
	for _,val := range callback{
		go val(ch)
	}
	return ch
}


// 数据库初始化
func initDb(){
	gdb.SetConfig(gdb.Config {
		"default" : gdb.ConfigGroup {
			gdb.ConfigNode {
				Host     : g.Cfg().GetString("mysql.conn.Host"),
				Port     : g.Cfg().GetString("mysql.conn.Port"),
				User     : g.Cfg().GetString("mysql.conn.User"),
				Pass     : g.Cfg().GetString("mysql.conn.Pass"),
				Name     : g.Cfg().GetString("mysql.conn.Name"),
				Type     : g.Cfg().GetString("mysql.conn.Type"),
				Charset  : g.Cfg().GetString("mysql.conn.Charset"),
			},
			//gdb.ConfigNode {
			//	Host     : connect.Host,
			//	Port     : connect.Port,
			//	User     : connect.User,
			//	Pass     : connect.Passwd,
			//	Name     : connect.DataBase,
			//	Type     : "mysql",
			//	Charset  : g.Cfg().GetString("mysql.Charset"),
			//},
		},
	})
	db ,err := gdb.New("default")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("正在测试连接数据库..请稍后(如果不提示`连接成功`说明程序没有运行)")
	_,err = db.Exec("show databases")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}
	fmt.Println("连接成功")
}