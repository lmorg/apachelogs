package apachelogs

import (
	"fmt"
	"sort"
)

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
		return (*sal.AccessLogs)[i].IP < (*sal.AccessLogs)[j].IP

	case FIELD_USER_ID:
		return (*sal.AccessLogs)[i].UserID < (*sal.AccessLogs)[j].UserID

	case FIELD_DATE_TIME, FIELD_DATE, FIELD_TIME:
		return (*sal.AccessLogs)[i].DateTime.Before((*sal.AccessLogs)[j].DateTime)

	case FIELD_METHOD:
		return (*sal.AccessLogs)[i].Method < (*sal.AccessLogs)[j].Method

	case FIELD_URI:
		return (*sal.AccessLogs)[i].URI < (*sal.AccessLogs)[j].URI

	case FIELD_QUERY_STRING:
		return (*sal.AccessLogs)[i].QueryString < (*sal.AccessLogs)[j].QueryString

	case FIELD_PROTOCOL:
		return (*sal.AccessLogs)[i].Protocol < (*sal.AccessLogs)[j].Protocol

	case FIELD_STATUS:
		return (*sal.AccessLogs)[i].Status.I < (*sal.AccessLogs)[j].Status.I

	case FIELD_SIZE:
		return (*sal.AccessLogs)[i].Size < (*sal.AccessLogs)[j].Size

	case FIELD_REFERRER:
		return (*sal.AccessLogs)[i].Referrer < (*sal.AccessLogs)[j].Referrer

	case FIELD_USER_AGENT:
		return (*sal.AccessLogs)[i].UserAgent < (*sal.AccessLogs)[j].UserAgent

	case FIELD_PROC_TIME:
		return (*sal.AccessLogs)[i].ProcTime < (*sal.AccessLogs)[j].ProcTime

	case FIELD_FILE_NAME:
		return (*sal.AccessLogs)[i].FileName < (*sal.AccessLogs)[j].FileName

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
