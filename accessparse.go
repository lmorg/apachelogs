package apachelogs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse access log entry
func ParseAccessLine(s string) (line *AccessLog, err error, matched bool) {
	// Quick strings.Split parser
	line, err = parserStringSplit(s)

	// Quick parse failed, falling back to slower regexp parser
	if err != nil || line.Status.I == 0 {
		line, err = parserRegexp(s)
	}

	if err == nil {
		matched, err = PatternMatch(line)
	}

	return
}

// Internal function: strings.Split the access log. It's faster than regexp.
func parserStringSplit(s string) (line *AccessLog, err error) {
	line = new(AccessLog)
	defer func() (line *AccessLog, err error) {
		// Catch any unforeseen errors that might cause a panic.
		// Mostly just to catch any out-of-bounds errors when working with slices,
		// all of which should already be accounted for but this should
		// cover any human error on my part.
		if r := recover(); r != nil {
			err = errors.New("panic caught in string split parser")
		}
		return
	}()

	split := strings.Split(s, " ")
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
		err = errors.New("empty request (typically 408)")
		return

	} else {
		line.IP = split[_S_IP]
		line.Method = split[_S_METHOD][1:]
		uri := strings.SplitN(split[_S_URI], "?", 2)
		line.URI = uri[0]
		if len(uri) == 2 {
			line.QueryString = "?" + uri[1]
		}
		line.Protocol = split[_S_PROTOCOL][:len(split[_S_PROTOCOL])-1]
		line.Size, _ = strconv.Atoi(split[_S_SIZE])
		line.Referrer = split[_S_REFERRER][1 : len(split[_S_REFERRER])-1]
		line.UserID = split[_S_USER_ID]

		pos := len(split) - 1
		for ; pos >= _S_USER_AGENT; pos-- {
			if split[pos][len(split[pos])-1:] == `"` {
				break
			}
		}
		if _S_USER_AGENT != pos {
			line.UserAgent = strings.Join(split[_S_USER_AGENT:pos+1], " ")
		} else {
			line.UserAgent = split[_S_USER_AGENT]
		}
		line.UserAgent = line.UserAgent[1 : len(line.UserAgent)-1]

		// Get processing time. This isn't part of the combined log format
		// but it is used in by Level 10 Fireball and Bronze Dagger
		if pos+1 < len(split) {
			//line.ProcTime, _ = strconv.Atoi(split[len(split)-2])
			line.ProcTime, _ = strconv.Atoi(split[pos+1])
		}

		line.Status = NewStatus(split[_S_STATUS])
	}

	return
}

// Internal function: regexp.FindStringSubmatch. Accurate but slooooow.
func parserRegexp(s string) (line *AccessLog, err error) {
	line = new(AccessLog)
	defer func() (line *AccessLog, err error) {
		if r := recover(); r != nil {
			err = errors.New("panic caught in regexp parser")
		}
		return
	}()

	split := rx_access_format.FindStringSubmatch(s)

	if len(split) < _SRX_DATE_TIME {
		//err = errors.New("Too few indexes in", filename, "line", i, split)
		err = errors.New(fmt.Sprintf("len(split){%d} < _SRX_DATE_TIME: %s", len(split), s))
		return
	}

	line.DateTime, err = time.Parse(APACHE_DATE_TIME, split[_SRX_DATE_TIME])
	if err != nil {
		return
	}

	if len(split) < _SRX_USER_AGENT {
		//err = errors.New("Too few indexes in", filename, "line", i, split)
		err = errors.New(fmt.Sprintf("len(split){%d} < _SRX_USER_AGENT: %s", len(split), s))
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
