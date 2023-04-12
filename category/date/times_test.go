package date

import (
	"fmt"
	"testing"
	"time"
)

func TestMinutes1(t *testing.T) {
	t1 := time.Now()
	//设置期间经历了50秒时间
	t2 := time.Now().Add(time.Second * 150)
	fmt.Println("t2与t1相差：", t2.Sub(t1)/60) //t2与t1相差： 50s
}

func TestMinutes(t *testing.T) {
	n := Minutes()
	fmt.Println(n)
}
