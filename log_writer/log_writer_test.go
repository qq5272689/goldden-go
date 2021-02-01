package log_writer

import (
	"fmt"
	"testing"
	"time"
)

func TestCLose(t *testing.T) {
	t.Log("TestCLose")
	lw, err := NewLogWriter("test", "/tmp", "M")
	if err != nil {
		t.Error(err)
	}
	//defer func() {
	//	if err = lw.Close(); err != nil {
	//		t.Error(err)
	//	}
	//}()
	for i := 0; i < 1000; i++ {
		go func(d int) {
			n, err := lw.Write([]byte(fmt.Sprintf("%d\n", d)))
			if err != nil {
				t.Error(err)
			}
			t.Log(d, n)
		}(i)

	}
	time.Sleep(time.Minute)
	for i := 0; i < 1000; i++ {
		go func(d int) {
			n, err := lw.Write([]byte(fmt.Sprintf("%d\n", d)))
			if err != nil {
				t.Error(err)
			}
			t.Log(d, n)
		}(i)

	}
	time.Sleep(time.Minute)
	for i := 0; i < 1000; i++ {
		go func(d int) {
			n, err := lw.Write([]byte(fmt.Sprintf("%d\n", d)))
			if err != nil {
				t.Error(err)
			}
			t.Log(d, n)
		}(i)

	}
	if err = lw.Close(); err != nil {
		t.Error(err)
	}

}
