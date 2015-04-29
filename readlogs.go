package apachelogs

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	//	"io/ioutil"
	"os"
	"strconv"
	"strings"
	//	"sync"
	"time"
)

var Debug bool

/*
func ScanDirectories() {
	var wg sync.WaitGroup

	logs = make(map[string][]Logs)
	errs = make(map[string][]Errs)
	init_log = ""

	items, _ := ioutil.ReadDir(LOGS_PATH)
	for _, d := range items {
		if !d.IsDir() {
			continue
		}

		wg.Add(1)
		go scanForLogs(LOGS_PATH+d.Name()+"/", d.Name(), &wg)
	}
	wg.Wait()
}
*/
/*
func scanForLogs(path string, site string, wg *sync.WaitGroup) {
	defer wg.Done()
	items, _ := ioutil.ReadDir(path)
	for _, f := range items {
		if f.IsDir() {
			continue
		}

		if strings.Contains(f.Name(), "log") && (LOGS_FILTER == "" || strings.Contains(f.Name(), LOGS_FILTER)) {
			readLog(path+f.Name(), site)
		}
		if strings.Contains(f.Name(), "err") && (LOGS_FILTER == "" || strings.Contains(f.Name(), LOGS_FILTER)) {
			readErr(path+f.Name(), site)
		}
	}
}
*/
func ReadAccessLog(filename string, logs *AccessLogs, errHandler func(err error)) {
	var (
		reader *bufio.Reader
		i, j   int
		err    error
	)

	fi, err := os.Open(filename)
	if err != nil {
		//fmt.Println("ERROR:", err)
		errHandler(err)
		return
	}
	defer fi.Close()

	if filename[len(filename)-3:] == ".gz" {
		fz, err := gzip.NewReader(fi)
		if err != nil {
			//fmt.Println("ERROR:", err)
			errHandler(err)
			return
		}
		reader = bufio.NewReader(fz)
	} else {
		reader = bufio.NewReader(fi)
	}

	for {
		b, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				//fmt.Println("ERROR:", err)
				errHandler(err)
			}
			break
		}

		i++
		//var line AccessLog
		line, err, matched := ParseAccessLine(&b, errHandler)

		if err != nil || !matched {
			//fmt.Println("ERROR:", err)
			errHandler(err)
			continue
		}

		*logs = append(*logs, line)
		j++

	}

	//fmt.Printf("%-65s: %6d lines read (%6d used)\n", filename, i, j)

	return
}

func ParseAccessLine(b *[]byte, errHandler func(err error)) (line AccessLog, err error, matched bool) {
	//line.ProcTime = "Not logged."

	line, err = parserStringSplit(b)
	if err != nil || line.Status.I == 0 {
		// Quick parse failed, falling back to regex
		line, err = parserRegEx(*b)
	}

	if err == nil {
		matched = PatternMatch(&line)
	} else {
		errHandler(err)
	}

	return
}

func parserStringSplit(b *[]byte) (line AccessLog, err error) {
	/*************************
	 ** string split parser **
	 *************************/
	//if !Debug {
	defer func() (line AccessLog, err error) {
		if r := recover(); r != nil {
			err = errors.New("panic caught in string split parser")
		}
		return
	}()
	//} else {
	//	fmt.Println(string(b))
	//}

	split := strings.Split(string(*b), " ")
	if len(split) < 3 {
		err = errors.New(fmt.Sprint("len(split) < 3", split))
		return
	}

	if len(split) >= 9 {
		line.DateTime, err = time.Parse(APACHE_DATE_TIME, split[_S_DATE_TIME][1:])

		if err != nil {
			return
		}

	} else {
		err = errors.New(fmt.Sprint("len(split) < 9", split))
		return
	}

	if len(split) < _S_STATUS {
		err = errors.New(fmt.Sprint("len(split) < _S_STATUS", split))
		return
	}
	if len(split) < _S_REFERRER {
		err = errors.New(fmt.Sprint("len(split) < _S_REFERRER", split))
		return
	}
	if split[_S_STATUS] == `"-"` {
		//err = errors.New("408s bug")
		err = errors.New("empty request bug (typically 408)")
		return

	} else {
		line.IP = split[_S_IP]
		line.Method = split[_S_METHOD][1:]
		uri := strings.SplitN(split[_S_URI], "?", 2)
		line.URI = uri[0]
		if len(uri) == 2 {
			line.QueryString = "?" + uri[1]
		}
		line.Protocol = split[_S_PROTO][:len(split[_S_PROTO])-1]
		line.Size, _ = strconv.Atoi(split[_S_SIZE])
		line.Referrer = split[_S_REFERRER][1 : len(split[_S_REFERRER])-1]
		line.UserID = split[_S_USER_ID]

		k := len(split) - 1
		for ; k >= _S_USER_AGENT; k-- {
			if split[k][len(split[k])-1:] == `"` {
				break
			}
		}
		if _S_USER_AGENT != k {
			line.UserAgent = strings.Join(split[_S_USER_AGENT:k+1], " ")
		} else {
			line.UserAgent = split[_S_USER_AGENT]
		}
		line.UserAgent = line.UserAgent[1 : len(line.UserAgent)-1]

		if k+1 < len(split) { // TODO: eh!?!
			//line.ProcTime = split[len(split)-2]
			line.ProcTime, _ = strconv.Atoi(split[len(split)-2])
		}

		line.Status = NewStatus(split[_S_STATUS])
	}

	return
}

func parserRegEx(b []byte) (line AccessLog, err error) {
	/******************
	 ** regex parser **   (slower, but more precise)
	 ******************/
	//if !Debug {
	defer func() (line AccessLog, err error) {
		if r := recover(); r != nil {
			err = errors.New("panic caught in regexp parser")
		}
		return
	}()
	//} else {
	//	fmt.Println(string(b))
	//}

	split := rx_log_format.FindStringSubmatch(string(b))

	if len(split) < _SRX_DATE_TIME {
		//err = errors.New("Too few indexes in", filename, "line", i, split)
		err = errors.New(fmt.Sprintf("len(split){%d} < _SRX_DATE_TIME: %s", len(split), b))
		return
	}

	line.DateTime, err = time.Parse(APACHE_DATE_TIME, split[_SRX_DATE_TIME])
	if err != nil {
		return
	}

	if len(split) < _SRX_USER_AGENT {
		//err = errors.New("Too few indexes in", filename, "line", i, split)
		err = errors.New(fmt.Sprintf("len(split){%d} < _SRX_USER_AGENT: %s", len(split), b))
		return
	}

	line.IP = split[_SRX_IP]
	line.Status = NewStatus(split[_SRX_STATUS])
	line.Size, _ = strconv.Atoi(split[_SRX_SIZE])
	line.Referrer = split[_SRX_REFERRER]
	line.UserAgent = split[_SRX_USER_AGENT]
	line.UserID = split[_SRX_USER_ID]

	request := strings.Split(split[_SRX_REQUEST], " ")

	switch len(request) {
	case 0:
		line.Method = "???"
		line.URI = "???"
		line.Protocol = "???"
	case 1:
		line.Method = "???"
		line.Protocol = "???"

		// left out of function for micro-optimisation reasons
		uri := strings.SplitN(request[0], "?", 2)
		line.URI = uri[0]
		if len(uri) == 2 {
			line.QueryString = "?" + uri[1]
		}
	case 2:
		line.Method = request[0]
		line.Protocol = "???"

		// left out of function for micro-optimisation reasons
		uri := strings.SplitN(request[1], "?", 2)
		line.URI = uri[0]
		if len(uri) == 2 {
			line.QueryString = "?" + uri[1]
		}
	case 3:
		line.Method = request[0]
		line.Protocol = request[2]

		// left out of function for micro-optimisation reasons
		uri := strings.SplitN(request[1], "?", 2)
		line.URI = uri[0]
		if len(uri) == 2 {
			line.QueryString = "?" + uri[1]
		}
	default:
		line.Method = request[0]
		line.Protocol = request[len(request)-1]

		s := strings.Join(request[1:len(request)-1], " ")
		// left out of function for micro-optimisation reasons
		uri := strings.SplitN(s, "?", 2)
		line.URI = uri[0]
		if len(uri) == 2 {
			line.QueryString = "?" + uri[1]
		}
	}

	return
}

func ReadErrorLog(filename string, logs *ErrorLogs, errHandler func(err error)) {
	var (
		reader *bufio.Reader
		i, j   int
		err    error
		line   ErrorLog
	)

	fi, err := os.Open(filename)
	if err != nil {
		errHandler(err)
		return
	}
	defer fi.Close()

	if filename[len(filename)-3:] == ".gz" {
		fz, err := gzip.NewReader(fi)
		if err != nil {
			errHandler(err)
			return
		}
		reader = bufio.NewReader(fz)
	} else {
		reader = bufio.NewReader(fi)
	}

	for {
		s_line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				errHandler(err)
			}
			break
		}
		i++

		if len(s_line) < 28 {
			continue
		}
		line.DateTime, err = time.Parse(ERR_DATE_TIME, string(s_line)[1:25])
		if err != nil {
			continue
		}

		line.Message = string(s_line[27:])

		*logs = append(*logs, line)
		j++

	}

	//fmt.Printf("%-65s: %6d lines read (%6d used)\n", filename, i, j)

	return
}
