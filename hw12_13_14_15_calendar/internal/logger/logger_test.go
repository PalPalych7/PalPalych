package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	q2 := New("", "debug")
	q2.Trace("trace")
	q2.Warn("warn")
	q2.Error("er")
	//
	//	q := New("6.log", "XZ")
	//	q.Trace("trace")
	//	q.Warn("warn")
	//	q.Fatal("fat")
	//
	//	fmt.Printf("val %v", q)
	//	q.Info("qqqq")
	//	q.Fatal("zzzz")
}
