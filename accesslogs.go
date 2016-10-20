package apachelogs

//go:generate stringer -type=FieldID

import (
	"regexp"
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
	//str_access_format = `^([\.0-9]+) (.*?) (.*?) \[(.*?) \+[0-9]{4}\] "(.*?)" ([\-0-9]+) ([\-0-9]+) "(.*?)" "(.*?)"`
	str_access_format = `^(.*?) (.*?) (.*?) \[(.*?) \+[0-9]{4}\] "(.*?)" ([\-0-9]+) ([\-0-9]+) "(.*?)" "(.*?)"` // some addresses aren't IPs

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

func (a *AccessLog) SetFieldID(id FieldID, val interface{}) {
	switch id {
	case FIELD_IP:
		a.IP = val.(string)
	case FIELD_USER_ID:
		a.UserID = val.(string)
	case FIELD_DATE_TIME, FIELD_DATE, FIELD_TIME:
		a.DateTime = val.(time.Time)
	case FIELD_METHOD:
		a.Method = val.(string)
	case FIELD_URI:
		a.URI = val.(string)
	case FIELD_QUERY_STRING:
		a.QueryString = val.(string)
	case FIELD_PROTOCOL:
		a.Protocol = val.(string)
	case FIELD_STATUS:
		a.Status = NewStatus(val.(string))
	case FIELD_SIZE:
		a.Size = val.(int)
	case FIELD_REFERRER:
		a.Referrer = val.(string)
	case FIELD_USER_AGENT:
		a.UserAgent = val.(string)
	case FIELD_PROC_TIME:
		a.ProcTime = val.(int)
	case FIELD_FILE_NAME:
		a.FileName = val.(string)
	}
}

type AccessLogs []*AccessLog

func (al AccessLogs) Remove(index int) { al = append(al[:index], al[index+1:]...) }
func (al AccessLogs) Len() int         { return len(al) }
func (al AccessLogs) Swap(i, j int)    { al[i], al[j] = al[j], al[i] }
