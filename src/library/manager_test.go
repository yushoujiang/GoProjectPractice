package library

import (
	"test"
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

}
