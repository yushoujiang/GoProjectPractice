package library

import (
	"testing"
)

func TestOps(t *testing.T) {
	mm := NewMusicManage()
	if mm == nil {
		t.Error("gen music failed")
	}
	if mm.Len() != 0 {
		t.Error("")
	}

	music0 := Music{Id: "1", Name: "好音乐", Artist: "什么属性呀", Source: "电子版", Type: "avi"}
	mm.Add(music0)

	if mm.Len() != 1 {
		t.Error("add error")
	}

	m, err := mm.Find(music0.Name)

	if m == nil {
		t.Error("find error")
	}

	if err != nil {

	}

	m, err = mm.Get(0)
	if m == nil {
		t.Error("get error")
	}

	m = mm.Remove(0)
	if m == nil || mm.Len() != 0 {
		t.Error("remove error")
	}
}
