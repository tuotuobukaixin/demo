package algorithm

import (
	"container/list"
	"fmt"
)

//Statistics is a Testing struct
type Statistics struct {
	UUID     string
	EndPoint string
	CPUUsage float32
	MemUsage float32
	AppNum   int
}

//TC is the testing center struct
type TC struct {
	GSlist *list.List
	GTflag bool
	GSeq   int
}

//T is a global variable to get testing data
var T TC

func init() {
	T.GSlist = list.New()
	T.GTflag = false
	T.GSeq = 1
}

//GSlistReset is a interface to reset the T
func (t *TC) GSlistReset(f bool) {
	t.GSlist = list.New()
	t.GTflag = f
}

//GSeqIncrease is  increasing seq
func (t *TC) GSeqIncrease() {
	T.GSeq++
}

//GSeqReset is reset seq
func (t *TC) GSeqReset() {
	T.GSeq = 1
}

//MakeStatisticsNode is a interface to new a list node
func (t *TC) MakeStatisticsNode(uid string, ep string, c float32, m float32, an int) {
	var s Statistics
	s.UUID = uid
	s.EndPoint = ep
	s.CPUUsage = c
	s.MemUsage = m
	s.AppNum = an
	T.GSlist.PushBack(s)
}

//DisplayStatistics is a interface to display the static data
func (t *TC) DisplayStatistics() {
	fmt.Println("============== ", t.GSeq, " ===============")
	for e := T.GSlist.Front(); e != nil; e = e.Next() {
		s := e.Value.(Statistics)
		fmt.Println(s)
	}
}
