package main

import (
	"fmt"
	"github.com/qq5272689/goutils/logrus-hooks/TimedRotatingFileHook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {
	log := logrus.New()
	fmt.Println("start")
	log.Warnln("start")
	hook, err := TimedRotatingFileHook.NewTRFileHook("/tmp/logs", "test.log", "M")
	hook.SetFilePrefix("/home/hj/gopaths/goutils_gopath/src/goutils")
	defer hook.CloseWrites()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	log.AddHook(hook)
	for i := 0; i <= 1; i++ {
		time.Sleep(time.Second * 1)
		log.Errorln(i)

	}

}
