package apachelogs

import (
	"sort"
	"time"
)

var DateTimeErrorFormat string = "Mon Jan 02 15:04:05 2006" // timestamp formatting in Apache err logs

type ErrorLine struct {
	DateTime     time.Time
	HasTimestamp bool // Sometimes log file entries don't have a timestamp
	Scope        []string
	Message      string
	FileName     string
}

type ErrorLog []ErrorLine

func (e ErrorLog) Remove(index int)   { e = append(e[:index], e[index+1:]...) }
func (e ErrorLog) SortByDateTime()    { sort.Sort(ErrorLog(e)) }
func (e ErrorLog) Len() int           { return len(e) }
func (e ErrorLog) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ErrorLog) Less(i, j int) bool { return e[i].DateTime.Before(e[j].DateTime) }
