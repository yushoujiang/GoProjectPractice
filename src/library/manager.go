package library

import (
	"errors"
)

type Music struct {
	Id     string
	Name   string
	Artist string
	Source string
	Type   string
}

type MusicManage struct {
	musics []Music
}

func NewMusicManage() (m *MusicManage) {

	music := MusicManage{musics: make([]Music, 0)}

	return &music
}

func (m MusicManage) Len() (mLen int) {
	if m.musics == nil {
		return 0
	}
	return len(m.musics)
}

func (m MusicManage) Get(index int) (retM *Music, err error) {
	if index < 0 || index > m.Len() {
		return nil, errors.New("index out of range")
	}

	return &m.musics[index], nil
}

func (m MusicManage) Find(name string) (retM *Music, err error) {
	if m.musics == nil || len(m.musics) == 0 {
		return nil, errors.New("array is empty")
	}

	for _, v := range m.musics {
		if v.Name == name {
			return &v, nil
		}
	}

	return nil, errors.New("not found in array")
}

func (m *MusicManage) Add(music Music) {
	m.musics = append(m.musics, music)
}

func (m *MusicManage) Remove(index int) *Music {

	if index < 0 || index > m.Len() {
		return nil
	}

	removedElement := &m.musics[index]

	if index < m.Len()-1 { //中间元素
		m.musics = append(m.musics[:index-1], m.musics[index+1:]...)
	} else if index == 0 { //第一个
		m.musics = make([]Music, 0)
	} else { //最后一个
		m.musics = m.musics[:index-1]
	}

	return removedElement
}
