package apachelogs

import (
	"regexp"
	"sort"
	"time"
)

const (
	APACHE_DATE_TIME = "02/Jan/2006:15:04:05"     // timestamp formatting in Apache logs, best not to change this
	ERR_DATE_TIME    = "Mon Jan 02 15:04:05 2006" // timestamp formatting in Apache err logs, best not to change this
	//	DISPLAY_DATE_TIME = "02-01-2006 15:04:05" // timestamp formatting in for greater / less filters

	_S_IP         = 0 // array indexes after the log has been split
	_S_USER_ID    = 2
	_S_DATE_TIME  = 3
	_S_METHOD     = 5
	_S_URI        = 6
	_S_PROTO      = 7
	_S_STATUS     = 8
	_S_SIZE       = 9
	_S_REFERRER   = 10
	_S_USER_AGENT = 11

	_S2_IP        = 0 // array indexes after the log has been split
	_S2_USER_ID   = 2
	_S2_DATE_TIME = 3
	_S2_STATUS    = 6
	_S2_SIZE      = 7
	_S2_REFERRER  = 8

	_SRX_IP         = 1 // array indexes for regex splits
	_SRX_USER_ID    = 3
	_SRX_DATE_TIME  = 4
	_SRX_REQUEST    = 5
	_SRX_STATUS     = 6
	_SRX_SIZE       = 7
	_SRX_REFERRER   = 8
	_SRX_USER_AGENT = 9

	regexp_log = `^([\.0-9]+) (.*?) (.*?) \[(.*?) \+[0-9]{4}\] "(.*?)" ([\-0-9]+) ([\-0-9]+) "(.*?)" "(.*?)"`
	//regexp_log = `^([\.0-9]+) (.*?) (.*?) \[(.*?) \+[0-9]{4}\] "(.*?)" ([\-0-9]+) ([\-0-9]+) "(.*?)"`
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
}

func (a AccessLog) ByFieldID(id byte) interface{} {
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
	default:
		return nil
	}
}

type AccessLogs []AccessLog

////////////////////////////////////////////////////////////////////////////////

type ErrorLog struct {
	DateTime time.Time
	Message  string
}

type ErrorLogs []ErrorLog

func (e ErrorLogs) Remove(index int)   { e = append(e[:index], e[index+1:]...) }
func (e ErrorLogs) SortByDateTime()    { sort.Sort(ErrorLogs(e)) }
func (e ErrorLogs) Len() int           { return len(e) }
func (e ErrorLogs) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e ErrorLogs) Less(i, j int) bool { return e[i].DateTime.Before(e[j].DateTime) }

var (
	rx_404_err    *regexp.Regexp
	rx_log_format *regexp.Regexp
)

func init() {
	rx_404_err, _ = regexp.Compile(`^\[error\] \[client [\.0-9]+\] File does not exist: `)
	rx_log_format, _ = regexp.Compile(regexp_log) // more precise, but slower.
}
