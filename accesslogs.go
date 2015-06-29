package apachelogs

//go:generate stringer -type=FieldID

import (
	"fmt"
	"regexp"
	"sort"
	"time"
)

const (
	// Timestamp formatting in Apache logs, best not to change this.
	APACHE_DATE_TIME = "02/Jan/2006:15:04:05"

	// Slice indexes after the log has been strings.Split.
	// This method is faster than using regexp matches.
	_S_IP         = 0
	_S_USER_ID    = 2
	_S_DATE_TIME  = 3
	_S_METHOD     = 5
	_S_URI        = 6
	_S_PROTOCOL   = 7
	_S_STATUS     = 8
	_S_SIZE       = 9
	_S_REFERRER   = 10
	_S_USER_AGENT = 11

	// regexp match string and slice indexes.
	// This method is slower but more accurate, so this
	// is only used as a fallback if the above fails
	str_access_format = `^([\.0-9]+) (.*?) (.*?) \[(.*?) \+[0-9]{4}\] "(.*?)" ([\-0-9]+) ([\-0-9]+) "(.*?)" "(.*?)"`

	_SRX_IP         = 1
	_SRX_USER_ID    = 3
	_SRX_DATE_TIME  = 4
	_SRX_REQUEST    = 5
	_SRX_STATUS     = 6
	_SRX_SIZE       = 7
	_SRX_REFERRER   = 8
	_SRX_USER_AGENT = 9
)

var rx_access_format *regexp.Regexp

func init() {
	rx_access_format, _ = regexp.Compile(str_access_format)
}

type FieldID byte

const (
	FIELD_IP FieldID = iota + 1
	FIELD_USER_ID
	FIELD_DATE_TIME
	FIELD_DATE
	FIELD_TIME
	FIELD_METHOD
	FIELD_URI
	FIELD_QUERY_STRING
	FIELD_PROTOCOL
	FIELD_STATUS
	FIELD_SIZE
	FIELD_REFERRER
	FIELD_USER_AGENT
	FIELD_PROC_TIME
	FIELD_FILE_NAME
)

type AccessLog struct {
	IP          string
	UserID      string
	DateTime    time.Time
	Method      string
	URI         string
	QueryString string
	Protocol    string
	Status      Status
	Size        int
	Referrer    string
	UserAgent   string
	ProcTime    int
	FileName    string
}

func (a AccessLog) ByFieldID(id FieldID) interface{} {
	switch id {
	case FIELD_IP:
		return a.IP
	case FIELD_USER_ID:
		return a.UserID
	case FIELD_DATE_TIME, FIELD_DATE, FIELD_TIME:
		return a.DateTime
	case FIELD_METHOD:
		return a.Method
	case FIELD_URI:
		return a.URI
	case FIELD_QUERY_STRING:
		return a.QueryString
	case FIELD_PROTOCOL:
		return a.Protocol
	case FIELD_STATUS:
		return a.Status.A
	case FIELD_SIZE:
		return a.Size
	case FIELD_REFERRER:
		return a.Referrer
	case FIELD_USER_AGENT:
		return a.UserAgent
	case FIELD_PROC_TIME:
		return a.ProcTime
	case FIELD_FILE_NAME:
		return a.FileName
	default:
		return nil
	}
}

type AccessLogs []*AccessLog

func (al AccessLogs) Remove(index int) { al = append(al[:index], al[index+1:]...) }
func (al AccessLogs) Len() int         { return len(al) }
func (al AccessLogs) Swap(i, j int)    { al[i], al[j] = al[j], al[i] }

type Sort struct {
	AccessLogs *AccessLogs
	Key        FieldID
}

// !!! the following methods are untested !!!
func (sal Sort) Remove(index int) { sal.AccessLogs.Remove(index) }
func (sal Sort) Len() int         { return sal.AccessLogs.Len() }
func (sal Sort) Swap(i, j int)    { sal.AccessLogs.Swap(i, j) }
func (sal Sort) Sort()            { sort.Sort(sal) }
func (sal Sort) Less(i, j int) bool {
	switch sal.Key {
	case FIELD_IP:
		return (*sal.AccessLogs)[i].IP < (*sal.AccessLogs)[i].IP

	case FIELD_USER_ID:
		return (*sal.AccessLogs)[i].UserID < (*sal.AccessLogs)[i].UserID

	case FIELD_DATE_TIME, FIELD_DATE, FIELD_TIME:
		return (*sal.AccessLogs)[i].DateTime.Before((*sal.AccessLogs)[j].DateTime)

	case FIELD_METHOD:
		return (*sal.AccessLogs)[i].Method < (*sal.AccessLogs)[i].Method

	case FIELD_URI:
		return (*sal.AccessLogs)[i].URI < (*sal.AccessLogs)[i].URI

	case FIELD_QUERY_STRING:
		return (*sal.AccessLogs)[i].QueryString < (*sal.AccessLogs)[i].QueryString

	case FIELD_PROTOCOL:
		return (*sal.AccessLogs)[i].Protocol < (*sal.AccessLogs)[i].Protocol

	case FIELD_STATUS:
		return (*sal.AccessLogs)[i].Status.I < (*sal.AccessLogs)[i].Status.I

	case FIELD_SIZE:
		return (*sal.AccessLogs)[i].Size < (*sal.AccessLogs)[i].Size

	case FIELD_REFERRER:
		return (*sal.AccessLogs)[i].Referrer < (*sal.AccessLogs)[i].Referrer

	case FIELD_USER_AGENT:
		return (*sal.AccessLogs)[i].UserAgent < (*sal.AccessLogs)[i].UserAgent

	case FIELD_PROC_TIME:
		return (*sal.AccessLogs)[i].ProcTime < (*sal.AccessLogs)[i].ProcTime

	case FIELD_FILE_NAME:
		return (*sal.AccessLogs)[i].FileName < (*sal.AccessLogs)[i].FileName

	case 0:
		panic("Key unset on sort")
	}

	panic(fmt.Sprintf("%s is not a valid sort key", sal.Key))
}

/*
	// Example usage of experimental sort:

	a := new(apachelogs.AccessLogs)
	sort := apachelogs.Sort{
		AccessLogs: a,
		Key:        apachelogs.FIELD_DATE,
	}
	sort.Sort()

*/
