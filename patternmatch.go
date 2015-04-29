package apachelogs

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Pattern struct {
	Field      byte
	Operator   byte
	Comparison interface{}
	regExp     *regexp.Regexp
	datetime   uint64
}

func NewPattern(field_id, operator byte, comparison string) (p Pattern) {
	a := new(AccessLog)

	switch v := a.ByFieldID(field_id).(type) {
	default:
		fmt.Printf("unexpected type %T\n", v)
		os.Exit(1)

	case string:
		p.Comparison = strings.ToLower(comparison)
		if operator == OP_REGEX_EQ || operator == OP_REGEX_NE {
			rx, err := regexp.Compile(`(?i)` + comparison)
			if err != nil {
				fmt.Printf("regexp.Compile(`(?i)%s`)", comparison)
				os.Exit(1)
			}
			p.regExp = rx
		}

	case int:
		i, err := strconv.Atoi(comparison)
		if err != nil {
			fmt.Printf("strconv.Atoi(%s) fails", comparison)
			os.Exit(1)
		}
		p.Comparison = i

	case time.Time:
		parse := map[byte]string{
			FIELD_DATE:      "01-02-2006",
			FIELD_TIME:      "15:04",
			FIELD_DATE_TIME: "01-02-2006 15:04",
		}
		t, err := time.Parse(parse[field_id], comparison)
		if err != nil {
			fmt.Printf(`time.Parse("01-02-06 15:04",%s)`, comparison)
			os.Exit(1)
		}

		switch field_id {
		case FIELD_DATE_TIME:
			p.datetime, _ = strconv.ParseUint(t.Format("200602011504"), 10, 64)
		case FIELD_DATE:
			p.datetime, _ = strconv.ParseUint(t.Format("20060201"), 10, 64)
		case FIELD_TIME:
			p.datetime, _ = strconv.ParseUint(t.Format("1504"), 10, 64)
		}

		p.Comparison = t
	}

	p.Field = field_id
	p.Operator = operator
	return
}

var Patterns []Pattern

const (
	OP_LESS_THAN byte = iota
	OP_GREATER_THAN
	OP_EQUAL_TO
	OP_NOT_EQUAL
	OP_REGEX_EQ
	OP_REGEX_NE
	OP_CONTAINS
	OP_NOT_CONTAIN
)

const (
	FIELD_IP byte = iota
	FIELD_USER_ID
	FIELD_DATE_TIME
	FIELD_METHOD
	FIELD_URI
	FIELD_QUERY_STRING
	FIELD_PROTOCOL
	FIELD_STATUS
	FIELD_SIZE
	FIELD_REFERRER
	FIELD_USER_AGENT
	FIELD_PROC_TIME
	FIELD_DATE
	FIELD_TIME
)

func PatternMatch(a *AccessLog) (r bool) {
	if len(Patterns) == 0 {
		return true
	}

	for _, p := range Patterns {
		switch v := p.Comparison.(type) {
		default:
			fmt.Printf("unexpected type %T\n", v)
			os.Exit(1)

		case string:
			switch p.Operator {
			default:
				fmt.Printf("unexpected operator id %d for %T", p.Operator, v)
				os.Exit(1)
			case OP_EQUAL_TO:
				r = strings.ToLower(a.ByFieldID(p.Field).(string)) == p.Comparison.(string)
			case OP_NOT_EQUAL:
				r = strings.ToLower(a.ByFieldID(p.Field).(string)) != p.Comparison.(string)
			case OP_CONTAINS:
				r = strings.Contains(strings.ToLower(a.ByFieldID(p.Field).(string)), p.Comparison.(string))
			case OP_NOT_CONTAIN:
				r = !strings.Contains(strings.ToLower(a.ByFieldID(p.Field).(string)), p.Comparison.(string))
			case OP_REGEX_EQ:
				r = p.regExp.MatchString(a.ByFieldID(p.Field).(string))
			case OP_REGEX_NE:
				r = !p.regExp.MatchString(a.ByFieldID(p.Field).(string))
			}

		case int:
			switch p.Operator {
			default:
				fmt.Printf("unexpected operator id %d for %T", p.Operator, v)
				os.Exit(1)
			case OP_EQUAL_TO:
				r = a.ByFieldID(p.Field).(int) == p.Comparison.(int)
			case OP_NOT_EQUAL:
				r = a.ByFieldID(p.Field).(int) != p.Comparison.(int)
			case OP_LESS_THAN:
				r = a.ByFieldID(p.Field).(int) < p.Comparison.(int)
			case OP_GREATER_THAN:
				r = a.ByFieldID(p.Field).(int) > p.Comparison.(int)
			}

		case time.Time:
			switch p.Field {
			default:
				fmt.Printf("unexpected type %T for field id %d\n", v, p.Field)
				os.Exit(1)
			case FIELD_DATE_TIME:
				switch p.Operator {
				default:
					fmt.Printf("unexpected operator id %d for %T", p.Operator, v)
					os.Exit(1)
				case OP_EQUAL_TO:
					i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("200602011504"), 10, 64)
					r = i == p.datetime
				case OP_NOT_EQUAL:
					i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("200602011504"), 10, 64)
					r = i != p.datetime
				case OP_LESS_THAN:
					r = a.ByFieldID(p.Field).(time.Time).Before(p.Comparison.(time.Time))
				case OP_GREATER_THAN:
					r = a.ByFieldID(p.Field).(time.Time).After(p.Comparison.(time.Time))
				}

			case FIELD_DATE:
				i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("200602011504"), 10, 64)
				switch p.Operator {
				default:
					fmt.Printf("unexpected operator id %d for %T", p.Operator, v)
					os.Exit(1)
				case OP_EQUAL_TO:
					r = i == p.datetime
				case OP_NOT_EQUAL:
					r = i != p.datetime
				case OP_LESS_THAN:
					r = i < p.datetime
				case OP_GREATER_THAN:
					r = i > p.datetime
				}

			case FIELD_TIME:
				i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("1504"), 10, 64)
				switch p.Operator {
				default:
					fmt.Printf("unexpected operator id %d for %T", p.Operator, v)
					os.Exit(1)
				case OP_EQUAL_TO:
					r = i == p.datetime
				case OP_NOT_EQUAL:
					r = i != p.datetime
				case OP_LESS_THAN:
					r = i < p.datetime
				case OP_GREATER_THAN:
					r = i > p.datetime
				}
			}
		}

		if !r {
			return
		}
	}

	return
}
