package apachelogs

import "time"

// Parse error log entry. Input is a byte slice rather than string (as used in `ParseApacheLine`) because we need to
// inspect each character and Go's Reader interface returns byte slices anyway.
// `last` is the previous lines timestamp and is used only if the current line doesn't have a timestamp (sometimes that
// happens - it's annoying when it does!)
//
// This code is very new so there's scope for a great deal of optimisation still. However I expect the function
// parameters and returns to remain as is because it follows the same design as `ParseApacheLine` which has been
// stable for a long time now.
func ParseErrorLine(line []byte, last time.Time) (errLog *ErrorLine, err error) {
	errLog = new(ErrorLine)

	for i := range line {
		var (
			matchBrace bool
			start      int
		)

		switch {

		case line[i] == ' ':
			continue

		case line[i] != ']' && matchBrace == true:
			continue

		case line[i] != '[' && start == 0:
			errLog.Message = string(line)
			errLog.DateTime = last
			return

		case line[i] == '[':
			matchBrace = true
			start = i
			continue

		case line[i] == ']':
			matchBrace = false

			if start < 3 {
				errLog.DateTime, err = time.Parse(DateTimeErrorFormat, string(line[1:i-1]))

				if err == nil {
					errLog.HasTimestamp = true
				} else {
					errLog.Scope = append(errLog.Scope, string(line[1:i-1]))
					errLog.DateTime = last
				}

				start = i + 1
			}

		default:
			errLog.Message = string(line[start:])
		}
	}

	return
}
