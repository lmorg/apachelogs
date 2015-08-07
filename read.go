package apachelogs

import (
	"bufio"
	"compress/gzip"
	"os"
)

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

		if err != nil {
			errHandler(err)
			continue
		}

		if !matched {
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
