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
	Field      AccessFieldId
	Operator   OperatorID
	Comparison interface{}
	regExp     *regexp.Regexp
	datetime   uint64
}

var rxRegexSubMatch *regexp.Regexp

func init() {
	rxRegexSubMatch, _ = regexp.Compile(`^{(.*?)}{(.*?)}$`)
}

func NewPattern(fieldId AccessFieldId, operator OperatorID, comparison string) (p Pattern, err error) {
	a := new(AccessLine)

	switch v := a.ByFieldID(fieldId).(type) {
	default:
		err = errors.New(fmt.Sprintf("Unexpected type %T", v))
		return

	case string:
		p.Comparison = strings.ToLower(comparison)
		if operator == OpRegexEqual || operator == OpRegexNotEqual || operator == OpRegexSubstitute {

			var (
				rx      *regexp.Regexp
				replace string
			)

			if operator == OpRegexSubstitute {
				match := rxRegexSubMatch.FindAllStringSubmatch(comparison, 1)
				if len(match) != 1 || len(match[0]) != 3 {
					err = errors.New(fmt.Sprintf("Cannot match {search}{replace} with '%s'", comparison))
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
			err = errors.New(fmt.Sprintf("Not a number '%s'", comparison))
			return
		}
		p.Comparison = i

	case time.Time:
		var t time.Time

		parse := map[AccessFieldId]string{
			AccFieldDate:     "01-02-2006",
			AccFieldTime:     "15:04",
			AccFieldDateTime: "01-02-2006 15:04",
		}
		t, err = time.Parse(parse[fieldId], comparison)
		if err != nil {
			//err = errors.New(fmt.Sprintf(`time.Parse("01-02-06 15:04",%s\n)`, comparison))
			return
		}

		switch fieldId {
		case AccFieldDateTime:
			p.datetime, _ = strconv.ParseUint(t.Format("200602011504"), 10, 64)
		case AccFieldDate:
			p.datetime, _ = strconv.ParseUint(t.Format("20060201"), 10, 64)
		case AccFieldTime:
			p.datetime, _ = strconv.ParseUint(t.Format("1504"), 10, 64)
		}

		p.Comparison = t
	}

	p.Field = fieldId
	p.Operator = operator
	return
}

var Patterns []Pattern

////////////////////////////////////////
// field id
////////////////////////////////////////

type OperatorID byte

const (
	OpLessThan OperatorID = iota + 1
	OpGreaterThan
	OpEqualTo
	OpNotEqual
	OpRegexEqual
	OpRegexNotEqual
	OpContains
	OpDoesNotContain
	OpRegexSubstitute
	OpRoundDown
	OpRoundUp
	OpDivide
	OpMultiply
)

func roundDown(val, round int) int { return int(val/round) * round }
func roundUp(val, round int) int   { return (int(val/round) * round) + round }

func PatternMatch(a *AccessLine) (r bool, err error) {
	if len(Patterns) == 0 {
		return true, nil
	}

	for _, p := range Patterns {
		switch v := p.Comparison.(type) {
		default:
			err = errors.New(fmt.Sprintf("Unexpected type %T", v))
			return

		case string:
			switch p.Operator {
			default:
				err = errors.New(fmt.Sprintf("Unexpected operator %s for %T", p.Operator, v))
				return
			case OpEqualTo:
				r = strings.ToLower(a.ByFieldID(p.Field).(string)) == p.Comparison.(string)
			case OpNotEqual:
				r = strings.ToLower(a.ByFieldID(p.Field).(string)) != p.Comparison.(string)
			case OpContains:
				r = strings.Contains(strings.ToLower(a.ByFieldID(p.Field).(string)), p.Comparison.(string))
			case OpDoesNotContain:
				r = !strings.Contains(strings.ToLower(a.ByFieldID(p.Field).(string)), p.Comparison.(string))
			case OpRegexEqual:
				r = p.regExp.MatchString(a.ByFieldID(p.Field).(string))
			case OpRegexNotEqual:
				r = !p.regExp.MatchString(a.ByFieldID(p.Field).(string))
			case OpRegexSubstitute:
				a.SetFieldID(p.Field, p.regExp.ReplaceAllString(a.ByFieldID(p.Field).(string), p.Comparison.(string)))
				r = true
			}

		case int:
			switch p.Operator {
			default:
				err = errors.New(fmt.Sprintf("Unexpected operator %s for %T", p.Operator, v))
				return
			case OpEqualTo:
				r = a.ByFieldID(p.Field).(int) == p.Comparison.(int)
			case OpNotEqual:
				r = a.ByFieldID(p.Field).(int) != p.Comparison.(int)
			case OpLessThan:
				r = a.ByFieldID(p.Field).(int) < p.Comparison.(int)
			case OpGreaterThan:
				r = a.ByFieldID(p.Field).(int) > p.Comparison.(int)
			case OpRoundDown:
				a.SetFieldID(p.Field, roundDown(a.ByFieldID(p.Field).(int), p.Comparison.(int)))
				r = true
			case OpRoundUp:
				a.SetFieldID(p.Field, roundUp(a.ByFieldID(p.Field).(int), p.Comparison.(int)))
				r = true
			case OpDivide:
				a.SetFieldID(p.Field, a.ByFieldID(p.Field).(int)/p.Comparison.(int))
				r = true
			case OpMultiply:
				a.SetFieldID(p.Field, a.ByFieldID(p.Field).(int)/p.Comparison.(int))
				r = true
			}

		case time.Time:
			switch p.Field {
			default:
				err = errors.New(fmt.Sprintf("Unexpected type %T for %s", v, p.Field))
				return
			case AccFieldDateTime:
				switch p.Operator {
				default:
					err = errors.New(fmt.Sprintf("Unexpected operator %s for %T", p.Operator, v))
					return
				case OpEqualTo:
					i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("200602011504"), 10, 64)
					r = i == p.datetime
				case OpNotEqual:
					i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("200602011504"), 10, 64)
					r = i != p.datetime
				case OpLessThan:
					r = a.ByFieldID(p.Field).(time.Time).Before(p.Comparison.(time.Time))
				case OpGreaterThan:
					r = a.ByFieldID(p.Field).(time.Time).After(p.Comparison.(time.Time))
				}

			case AccFieldDate:
				i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("200602011504"), 10, 64)
				switch p.Operator {
				default:
					err = errors.New(fmt.Sprintf("Unexpected operator id %s for %T", p.Operator, v))
					return
				case OpEqualTo:
					r = i == p.datetime
				case OpNotEqual:
					r = i != p.datetime
				case OpLessThan:
					r = i < p.datetime
				case OpGreaterThan:
					r = i > p.datetime
				}

			case AccFieldTime:
				i, _ := strconv.ParseUint(a.ByFieldID(p.Field).(time.Time).Format("1504"), 10, 64)
				switch p.Operator {
				default:
					err = errors.New(fmt.Sprintf("Unexpected operator id %s for %T", p.Operator, v))
					return
				case OpEqualTo:
					r = i == p.datetime
				case OpNotEqual:
					r = i != p.datetime
				case OpLessThan:
					r = i < p.datetime
				case OpGreaterThan:
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
