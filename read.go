package apachelogs

import (
	"bufio"
	"compress/gzip"
	"os"
)

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

func ReadAccessLog(filename string, callback func(access *AccessLog), errHandler func(err error)) {
	var (
		reader *bufio.Reader
		err    error
	)

	fi, err := os.Open(filename)
	if err != nil {
		errHandler(err)
		return
	}
	defer fi.Close()

	if len(filename) > 3 && filename[len(filename)-3:] == ".gz" {
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
		b, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				errHandler(err)
			}
			break
		}

		line, err, matched := ParseAccessLine(string(b))

		if !matched {
			continue
		}
		if err != nil {
			errHandler(err)
			continue
		}

		line.FileName = filename
		//*logs = append(*logs, line)
		callback(line)
	}

	return
}

/*
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
*/
