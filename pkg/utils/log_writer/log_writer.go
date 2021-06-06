package log_writer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogWriter struct {
	file     *os.File
	when     string
	path     string
	service  string
	filename string
	mu       *sync.Mutex
	wg       *sync.WaitGroup
	nowTime  time.Time
}

func NewLogWriter(service, path, when string) (lw *LogWriter, err error) {
	lw = &LogWriter{when: when, path: path, service: service, mu: new(sync.Mutex), wg: new(sync.WaitGroup), nowTime: time.Now()}
	err = lw.NewFile()
	return
}

type FileIsDirErr struct {
	file string
}

func (e FileIsDirErr) Error() string {
	return e.file + " 是一个目录不是文件"
}

func (lw *LogWriter) NewFile() (err error) {
	lw.path, _ = filepath.Abs(lw.path)
	fmt.Println("base path", lw.path)
	if fi, err := os.Stat(lw.path); err != nil {
		if err = os.MkdirAll(lw.path, os.ModePerm); err != nil {
			return errors.New("目录:" + lw.path + "创建失败！！！")
		}
	} else {
		if !fi.IsDir() {
			return errors.New("路径:" + lw.path + "非目录！！！")
		}
	}
	lw.filename = filepath.Join(lw.path, lw.service) + ".log"
	flag := os.O_CREATE | os.O_RDWR | os.O_APPEND
	if fi, err := os.Stat(lw.filename); err == nil {
		if fi.IsDir() {
			return FileIsDirErr{lw.filename}
		}
		mod_logstr := lw.logstr(fi.ModTime())
		if mod_logstr != lw.logstr(lw.nowTime) {
			os.Rename(lw.filename, lw.filename+"-"+mod_logstr)
		}
	}
	if lw.file, err = os.OpenFile(lw.filename, flag, 0664); err != nil {
		return err
	}
	return nil
}

func (lw *LogWriter) updateWrite(nt time.Time) (err error) {
	if lw.file != nil {
		if err = lw.file.Close(); err != nil {
			return err
		}
		mod_logstr := lw.logstr(lw.nowTime)
		os.Rename(lw.filename, lw.filename+"-"+mod_logstr)
		lw.nowTime = nt
	}
	return lw.NewFile()
}

func (lw *LogWriter) Write(p []byte) (n int, err error) {
	lw.wg.Add(1)
	defer lw.wg.Done()
	lw.mu.Lock()
	defer lw.mu.Unlock()
	nt := time.Now()
	if lw.logstr(lw.nowTime) != lw.logstr(nt) {
		if err = lw.updateWrite(nt); err != nil {
			return 0, err
		}
	}
	if lw.file == nil {
		if err = lw.NewFile(); err != nil {
			return 0, err
		}
	}
	return lw.file.Write(p)
}

func (lw *LogWriter) Close() error {
	fmt.Println("close logWriter")
	lw.wg.Wait()
	if err := lw.file.Close(); err != nil {
		return err
	}
	lw.file = nil
	return nil
}

func (lw *LogWriter) Sync() error {
	return lw.Close()
}

func (lw *LogWriter) logstr(t time.Time) string {
	switch lw.when {
	case "H":
		return t.Format("2006-01-02-15")
	case "M":
		return t.Format("2006-01-02-15-04")
	case "D":
		return t.Format("2006-01-02")
	default:
		return t.Format("2006-01-02-15")
	}
}
