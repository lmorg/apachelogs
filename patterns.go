package apachelogs

//go:generate stringer -type=OperatorID

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Pattern struct {
	Field      FieldID
	Operator   OperatorID
	Comparison interface{}
	regExp     *regexp.Regexp
	datetime   uint64
}

var rx_regex_sub_match *regexp.Regexp

func init() {
	rx_regex_sub_match, _ = regexp.Compile(`^{(.*?)}{(.*?)}$`)
}

func NewPattern(field_id FieldID, operator OperatorID, comparison string) (p Pattern, err error) {
	a := new(AccessLog)

	switch v := a.ByFieldID(field_id).(type) {
	default:
		err = errors.New(fmt.Sprintf("unexpected type %T\n", v))
		return

	case string:
		p.Comparison = strings.ToLower(comparison)
		if operator == OP_REGEX_EQ || operator == OP_REGEX_NE || operator == OP_REGEX_SUB {

			var (
				rx      *regexp.Regexp
				replace string
			)

			if operator == OP_REGEX_SUB {
				match := rx_regex_sub_match.FindAllStringSubmatch(comparison, 1)
				if len(match[0]) != 3 {
					err = errors.New(fmt.Sprintf("Cannot match {}{} with", comparison))
					return
				}

				comparison, replace = match[0][1], match[0][2]
			}

			rx, err = regexp.Compile(`(?i)` + comparison)
			if err != nil {
				//err = errors.New(fmt.Sprintf("regexp.Compile(`(?i)%s`)\n", comparison))
				return
			}

			p.regExp = rx
			p.Comparison = replace

		}

	case int:
		var i int

		i, err = strconv.Atoi(comparison)
		if err != nil {
			//err = errors.New(fmt.SpPrintf("strconv.Atoi(%s) fails\n", comparison))
			return
		}
		p.Comparison = i

	case time.Time:
		var t time.Time

		parse := map[FieldID]string{
			FIELD_DATE:      "01-02-2006",
			FIELD_TIME:      "15:04",
			FIELD_DATE_TIME: "01-02-2006 15:04",
		}
		t, err = time.Parse(parse[field_id], comparison)
		if err != nil {
			//err = errors.New(fmt.Sprintf(`time.Parse("01-02-06 15:04",%s\n)`, comparison))
			return
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

////////////////////////////////////////
// field id
////////////////////////////////////////

type OperatorID byte

const (
	OP_LESS_THAN OperatorID = iota + 1
	OP_GREATER_THAN
	OP_EQUAL_TO
	OP_NOT_EQUAL
	OP_REGEX_EQ
	OP_REGEX_NE
	OP_CONTAINS
	OP_NOT_CONTAIN
	OP_REGEX_SUB
	OP_ROUND_DOWN
	OP_ROUND_UP
)

func roundDown(val, round int) int { return int(val/round) * round }
func roundUp(val, round int) int   { return (int(val/round) * round) + round }

func PatternMatch(a *AccessLog) (r bool, err error) {
	if len(Patterns) == 0 {
		return true, nil
	}

	for _, p := range Patterns {
		switch v := p.Comparison.(type) {
		default:
			err = errors.New(fmt.Sprintf("unexpected type %T\n", v))
			return

		case string:
			switch p.Operator {
			default:
				err = errors.New(fmt.Sprintf("unexpected operator %s for %T\n", p.Operator, v))
				return
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
			case OP_REGEX_SUB:
				a.SetFieldID(p.Field, p.regExp.ReplaceAllString(a.ByFieldID(p.Field).(string), p.Comparison.(string)))
				r = true
			}

		case int:
			switch p.Operator {
			default:
				err = errors.New(fmt.Sprintf("unexpected operator %s for %T\n", p.Operator, v))
				return
			case OP_EQUAL_TO:
				r = a.ByFieldID(p.Field).(int) == p.Comparison.(int)
			case OP_NOT_EQUAL:
				r = a.ByFieldID(p.Field).(int) != p.Comparison.(int)
			case OP_LESS_THAN:
				r = a.ByFieldID(p.Field).(int) < p.Comparison.(int)
			case OP_GREATER_THAN:
				r = a.ByFieldID(p.Field).(int) > p.Comparison.(int)
			case OP_ROUND_DOWN:
				a.SetFieldID(p.Field, roundDown(a.ByFieldID(p.Field).(int), p.Comparison.(int)))
				r = true
			case OP_ROUND_UP:
				a.SetFieldID(p.Field, roundUp(a.ByFieldID(p.Field).(int), p.Comparison.(int)))
				r = true
			}

		case time.Time:
			switch p.Field {
			default:
				err = errors.New(fmt.Sprintf("unexpected type %T for %s\n", v, p.Field))
				return
			case FIELD_DATE_TIME:
				switch p.Operator {
				default:
					err = errors.New(fmt.Sprintf("unexpected operator %s for %T\n", p.Operator, v))
					return
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
					err = errors.New(fmt.Sprintf("unexpected operator id %s for %T\n", p.Operator, v))
					return
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
					err = errors.New(fmt.Sprintf("unexpected operator id %s for %T\n", p.Operator, v))
					return
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
