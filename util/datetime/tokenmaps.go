package datetime

import (
	"flag"
	"fmt"
	"time"
)

const TimeFormat = "15:04:05, Monday 2 January 2006 MST"

var timezone = flag.String("timezone", "Europe/London",
	"The timezone to display reminder times in.")

type tokenMap interface {
	Lookup(input string, lval *yySymType) (tokenType int, ok bool)
}

type wordMap map[string]int

var wordTokenMap = wordMap{
	"THE": T_THE,
	"OF":  T_OF,
	"IN":  T_IGNORE,
	"AT":  T_IGNORE,
	"ON":  T_IGNORE,
}

func (wtm wordMap) Lookup(input string, lval *yySymType) (int, bool) {
	if tok, ok := wtm[input]; ok {
		return tok, ok
	}
	return -1, false
}

type numMap map[string]struct {
	tokenType int
	tokenVal  int
}

var numTokenMap = numMap{
	"AGO":   {T_AGO, -1},
	"YEAR":  {T_OFFSET, int(O_YEAR)},
	"Y":     {T_ISOYD, int(O_YEAR)},
	"MONTH": {T_OFFSET, int(O_MONTH)},
	// *Many* ambiguity problems.
	//	"M":         {T_ISO, int(O_MONTH)},
	"FORTNIGHT": {T_DAYS, 14},
	"WEEK":      {T_DAYS, 7},
	// W is used as the week indicator in ISO 8601
	//	"W":         {T_DAYS, 7},
	"DAY":    {T_OFFSET, int(O_DAY)},
	"D":      {T_ISOYD, int(O_DAY)},
	"NIGHT":  {T_OFFSET, int(O_DAY)},
	"HOUR":   {T_OFFSET, int(O_HOUR)},
	"H":      {T_ISOHS, int(O_HOUR)},
	"MINUTE": {T_OFFSET, int(O_MIN)},
	"MIN":    {T_OFFSET, int(O_MIN)},
	//	"M":         {T_ISO, int(O_MIN)},
	"SECOND":    {T_OFFSET, int(O_SEC)},
	"SEC":       {T_OFFSET, int(O_SEC)},
	"S":         {T_ISOHS, int(O_SEC)},
	"TOMORROW":  {T_DAYSHIFT, 1},
	"YESTERDAY": {T_DAYSHIFT, -1},
	"TODAY":     {T_DAYSHIFT, 0},
	"NOW":       {T_DAYSHIFT, 0},
	"ST":        {T_DAYQUAL, 1},
	"ND":        {T_DAYQUAL, 2},
	"RD":        {T_DAYQUAL, 3},
	"TH":        {T_DAYQUAL, 4},
	"MIDDAY":    {T_MIDTIME, 12},
	"MIDNIGHT":  {T_MIDTIME, 0},
}

func (ntm numMap) Lookup(input string, lval *yySymType) (int, bool) {
	// allow plurals
	if last := len(input) - 1; len(input) != 1 && input[last] == 'S' {
		input = input[:last]
	}
	if tok, ok := ntm[input]; ok {
		lval.intval = tok.tokenVal
		return tok.tokenType, ok
	}
	return -1, false
}

type abbrMap map[string]struct {
	tokenType int
	tokenVal  int
}

var abbrTokenMap = abbrMap{
	"JAN": {T_MONTHNAME, 1},
	"FEB": {T_MONTHNAME, 2},
	"MAR": {T_MONTHNAME, 3},
	"APR": {T_MONTHNAME, 4},
	"MAY": {T_MONTHNAME, 5},
	"JUN": {T_MONTHNAME, 6},
	"JUL": {T_MONTHNAME, 7},
	"AUG": {T_MONTHNAME, 8},
	"SEP": {T_MONTHNAME, 9},
	"OCT": {T_MONTHNAME, 10},
	"NOV": {T_MONTHNAME, 11},
	"DEC": {T_MONTHNAME, 12},
	"MON": {T_DAYNAME, 1},
	"TUE": {T_DAYNAME, 2},
	"WED": {T_DAYNAME, 3},
	"THU": {T_DAYNAME, 4},
	"FRI": {T_DAYNAME, 5},
	"SAT": {T_DAYNAME, 6},
	"SUN": {T_DAYNAME, 0},
}

func (atm abbrMap) Lookup(input string, lval *yySymType) (int, bool) {
	if len(input) > 3 {
		input = input[:3]
	}
	if tok, ok := atm[input]; ok {
		lval.intval = tok.tokenVal
		return tok.tokenType, ok
	}
	return -1, false
}

type relMap map[string]int

var relTokenMap = relMap{
	"LAST":  -1,
	"THIS":  0,
	"NEXT":  1,
	"AN":    1,
	"FIRST": 1,
	"ONE":   1,
	//	"SECOND":   2,
	"TWO":      2,
	"THIRD":    3,
	"THREE":    3,
	"FOURTH":   4,
	"FOUR":     4,
	"FIFTH":    5,
	"FIVE":     5,
	"SIXTH":    6,
	"SIX":      6,
	"SEVENTH":  7,
	"SEVEN":    7,
	"EIGHTH":   8,
	"EIGHT":    8,
	"NINTH":    9,
	"NINE":     9,
	"TENTH":    10,
	"TEN":      10,
	"ELEVENTH": 11,
	"TWELFTH":  12,
}

func (rtm relMap) Lookup(input string, lval *yySymType) (int, bool) {
	if tok, ok := rtm[input]; ok {
		lval.intval = tok
		return T_RELATIVE, ok
	}
	return -1, false
}

type zoneMap map[string]string

var zoneCache = make(map[string]*time.Location)

func zone(loc string) *time.Location {
	if l, ok := zoneCache[loc]; ok {
		return l
	}
	l, err := time.LoadLocation(loc)
	if err != nil {
		return nil
	}
	zoneCache[loc] = l
	return l
}

func Zone(loc string) *time.Location {
	if _, ok := zoneTokenMap[loc]; ok {
		return zone(zoneTokenMap[loc])
	}
	return zone(loc)
}

func Format(t time.Time, format ...string) string {
	if len(format) == 1 {
		return t.In(Zone(*timezone)).Format(format[0])
	}
	return t.In(Zone(*timezone)).Format(TimeFormat)
}

var zoneTokenMap = zoneMap{
	"ADT":  "America/Barbados",
	"AFT":  "Asia/Kabul",
	"AKST": "US/Alaska",
	"AKDT": "US/Alaska",
	"AMT":  "America/Boa_Vista",
	"ANAT": "Asia/Anadyr",
	"ART":  "America/Argentina/Buenos_Aires",
	"AST":  "Asia/Qatar",
	"AZOT": "Atlantic/Azores",
	"BNT":  "Asia/Brunei",
	"BRT":  "Brazil/East",
	"BRST": "Brazil/East",
	"BST":  "GB",
	"CAT":  "Africa/Harare",
	"CCT":  "Indian/Cocos",
	"CDT":  "US/Central",
	"CET":  "Europe/Zurich",
	"CEST": "Europe/Zurich",
	"CLST": "Chile/Continental",
	"CST":  "Asia/Shanghai",
	"EAT":  "Africa/Nairobi",
	"EDT":  "US/Eastern",
	"EET":  "Europe/Athens",
	"EIT":  "Asia/Jayapura",
	"EEST": "Europe/Athens",
	"EST":  "Australia/Melbourne",
	"FET":  "Europe/Kaliningrad",
	"FJT":  "Pacific/Fiji",
	"FJST": "Pacific/Fiji",
	"GET":  "Asia/Tbilisi",
	"GMT":  "GMT",
	"GST":  "Asia/Dubai",
	"HADT": "US/Aleutian",
	"HAST": "US/Aleutian",
	"HKT":  "Hongkong",
	"HST":  "US/Hawaii",
	"ICT":  "Asia/Bangkok",
	"IDT":  "Asia/Tel_Aviv",
	"IDDT": "Asia/Tel_Aviv",
	"IRDT": "Iran",
	"IRST": "Iran",
	"IOT":  "Indian/Chagos",
	"IST":  "Asia/Kolkata",
	"JST":  "Asia/Tokyo",
	"KGT":  "Asia/Bishkek",
	"KST":  "Asia/Pyongyang",
	"MDT":  "US/Mountain",
	"MART": "Pacific/Marquesas",
	"MET":  "MET",
	"MEST": "MET",
	"MMT":  "Asia/Rangoon",
	"MST":  "US/Mountain",
	"MVT":  "Indian/Maldives",
	"MYT":  "Asia/Kuala_Lumpur",
	"NDT":  "Canada/Newfoundland",
	"NPT":  "Asia/Kathmandu",
	"NST":  "Canada/Newfoundland",
	"NZDT": "Pacific/Auckland",
	"NZST": "Pacific/Auckland",
	"PDT":  "US/Pacific",
	"PHT":  "Asia/Manila",
	"PKT":  "Asia/Karachi",
	"PST":  "US/Pacific",
	"PWT":  "Pacific/Palau",
	"RET":  "Indian/Reunion",
	"SAST": "Africa/Johannesburg",
	"SCT":  "Indian/Mahe",
	"SGT":  "Asia/Singapore",
	"SST":  "US/Samoa",
	"ULAT": "Asia/Ulaanbaatar",
	"UTC":  "UTC",
	"UZT":  "Asia/Tashkent",
	"WAT":  "Africa/Lagos",
	"WAST": "Africa/Lagos",
	"WET":  "WET",
	"WEST": "WET",
	"WIT":  "Asia/Jakarta",
	"WST":  "Australia/West",
	"VET":  "America/Caracas",
	"VLAT": "Asia/Vladivostok",
	"Z":    "UTC",
}

func (ztm zoneMap) Lookup(input string, lval *yySymType) (int, bool) {
	if tok, ok := ztm[input]; ok {
		lval.zoneval = zone(tok)
		return T_ZONE, ok
	}
	return -1, false
}

type tokenMapList []tokenMap

var tokenMaps = tokenMapList{wordTokenMap, numTokenMap, abbrTokenMap, relTokenMap, zoneTokenMap}

func (l tokenMapList) Lookup(input string, lval *yySymType) (int, bool) {
	if DEBUG {
		fmt.Printf("Map lookup: %s\n", input)
	}
	// These maps are defined in tokenmaps.go
	for _, m := range l {
		if tok, ok := m.Lookup(input, lval); ok {
			if DEBUG {
				fmt.Printf("Map got: %d %d\n", lval.intval, tok)
			}
			return tok, ok
		}
	}
	if DEBUG {
		fmt.Printf("Map lookup failed\n")
	}
	return 0, false
}
