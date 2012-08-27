package datetime

import (
	"testing"
	"time"
)

func TestParseTimeFormats(t *testing.T) {
	// RFC822 doesn't specify seconds, and Stamp doesn't specify year
	ref := time.Date(time.Now().Year(), 6, 22, 13, 10, 0, 0, time.Local)
	formats := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
//		time.RFC3339Nano, // Nanosecs not supported
//		time.Kitchen,     // only contains HH and MM
		time.Stamp,
	}
	for i, f := range formats {
		in := ref.Format(f)
		ret, ok := Parse(in)
		if !ok || !ret.Equal(ref) {
			t.Errorf("Unable to parse format %d\nin: %s\ngot: %s", i, in, ret)
		}
	}
}

func TestParseTime(t *testing.T) {
	now := time.Now()
	mkt := func (h, m, s int, l ...string) time.Time {
		loc := time.Local
		if len(l) > 0 {
			loc = zone(l[0])
		}
		return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
	}
	tests := []struct{
		in string
		t  time.Time
	}{
		// T_INTEGER T_AMPM o_zone (also tests all possible zone permutations)
		{"11am", mkt(11, 0, 0)},
		{"11pm", mkt(23, 0, 0)},
		{"3 am PDT", mkt(3, 0, 0, "US/Pacific")},
		{"3 pm PDT", mkt(15, 0, 0, "US/Pacific")},
		{"5AM -4:30", mkt(5, 0, 0, "America/Caracas")},
		{"5PM -4:30", mkt(17, 0, 0, "America/Caracas")},
		{"7 a.m. +0800", mkt(7, 0, 0, "Asia/Shanghai")},
		{"7 p.m. +0800", mkt(19, 0, 0, "Asia/Shanghai")},
		{"9A.M. Africa/Nairobi", mkt(9, 0, 0, "Africa/Nairobi")},
		{"9P.M. Africa/Nairobi", mkt(21, 0, 0, "Africa/Nairobi")},
		// T_INTEGER : T_INTEGER T_AMPM o_zone
		{"11:23am", mkt(11, 23, 0)},
		{"11:23pm", mkt(23, 23, 0)},
		{"11:23am PDT", mkt(11, 23, 0, "US/Pacific")},
		{"11:23pm PDT", mkt(23, 23, 0, "US/Pacific")},
		// T_INTEGER : T_INTEGER : T_INTEGER T_AMPM o_zone
		{"11:23:45am", mkt(11, 23, 45)},
		{"11:23:45pm", mkt(23, 23, 45)},
		{"11:23:45am PDT", mkt(11, 23, 45, "US/Pacific")},
		{"11:23:45pm PDT", mkt(23, 23, 45, "US/Pacific")},
		// T_INTEGER zone
		{"03 PDT", mkt(3, 0, 0, "US/Pacific")},
		{"23 PDT", mkt(23, 0, 0, "US/Pacific")},
		{"0323 PDT", mkt(3, 23, 0, "US/Pacific")},
		{"2323 PDT", mkt(23, 23, 0, "US/Pacific")},
		{"032345 PDT", mkt(3, 23, 45, "US/Pacific")},
		{"232345 PDT", mkt(23, 23, 45, "US/Pacific")},
		// T_INTEGER : T_INTEGER o_zone
		{"11:23", mkt(11, 23, 0)},
		{"23:23", mkt(23, 23, 0)},
		{"11:23 PDT", mkt(11, 23, 0, "US/Pacific")},
		{"23:23 PDT", mkt(23, 23, 0, "US/Pacific")},
		// T_INTEGER : T_INTEGER : T_INTEGER o_zone
		{"11:23:45", mkt(11, 23, 45)},
		{"23:23:45", mkt(23, 23, 45)},
		{"11:23:45 PDT", mkt(11, 23, 45, "US/Pacific")},
		{"23:23:45 PDT", mkt(23, 23, 45, "US/Pacific")},
		// T_INTEGER (if len == 6)
		{"112345", mkt(11, 23, 45)},
		{"232345", mkt(23, 23, 45)},
		// These may be less-expected results... ;-)
		{"23am", mkt(23, 0, 0)},
		{"23pm", mkt(11, 0, 0)},
		{"11:63am", mkt(11, 03, 0)},
		{"11:83pm", mkt(23, 23, 0)},
		{"11:23:63am", mkt(11, 23, 03)},
		{"11:23:83pm", mkt(23, 23, 23)},
		{"27 PDT", mkt(3, 0, 0, "US/Pacific")},
		{"27:83", mkt(3, 23, 0)},
		{"27:73:83", mkt(3, 13, 23)},
	}
	for i, test := range tests {
		ret, ok := parse(test.in, now)
		if !ok || !ret.Equal(test.t) {
			t.Errorf("Unable to parse time %d\nin: %s\ngot: %s",
				i, test.in, ret)
		}
	}
}

func TestParseDate(t *testing.T) {
	h, n, s := 11, 22, 33
	mkt := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, h, n, s, 0, time.UTC)
	}
	rel := mkt(1,1,1)
	tests := []struct{
		in string
		t  time.Time
	}{
		// T_INTEGER / T_INTEGER as MM/YYYY
		{"3/2004", mkt(2004, 3, 1)},
		{"03/2004", mkt(2004, 3, 1)},
		{"12/2004", mkt(2004, 12, 1)},
		// T_INTEGER / T_INTEGER as DD/MM
		{"12/4", mkt(1, 4, 12)},
		{"30/4", mkt(1, 4, 30)},
		{"31/12", mkt(1, 12, 31)},
		// T_INTEGER / T_INTEGER / T_INTEGER as DD/MM/YYYY
		{"2/3/2004", mkt(2004, 3, 2)},
		{"02/03/2004", mkt(2004, 3, 2)},
		{"31/12/2004", mkt(2004, 12, 31)},
		// T_INTEGER / T_INTEGER / T_INTEGER as DD/MM/YY
		{"2/3/4", mkt(2004, 3, 2)},
		{"02/03/04", mkt(2004, 3, 2)},
		{"31/12/04", mkt(2004, 12, 31)},
		{"2/3/68", mkt(2068, 3, 2)},
		{"2/3/69", mkt(1969, 3, 2)},
		// T_INTEGER o_dayqual o_of T_MONTHNAME
		{"2 Mar", mkt(1, 3, 2)},
		{"02 Mar", mkt(1, 3, 2)},
		{"2nd March", mkt(1, 3, 2)},
		{"3rd of March", mkt(1, 3, 3)},
		// T_MONTHNAME T_INTEGER o_dayqual as Mon YYYY
		{"March 2004", mkt(2004, 3, 1)},
		// T_MONTHNAME T_INTEGER o_dayqual as Mon DDth
		{"March 2", mkt(1, 3, 2)},
		{"March 02", mkt(1, 3, 2)},
		{"Mar 2nd", mkt(1, 3, 2)},
		{"Mar 3rd", mkt(1, 3, 3)},
		// T_INTEGER o_dayqual o_of T_MONTHNAME T_INTEGER
		{"2 Mar 2004", mkt(2004, 3, 2)},
		{"02 Mar 2004", mkt(2004, 3, 2)},
		{"2nd March 2004", mkt(2004, 3, 2)},
		{"3rd of March 2004", mkt(2004, 3, 3)},
		{"2 March 4", mkt(2004, 3, 2)},
		{"2 March 68", mkt(2068, 3, 2)},
		{"2 March 69", mkt(1969, 3, 2)},
		// T_INTEGER T_MINUS T_MONTHNAME T_MINUS T_INTEGER
		{"2-Mar-2004", mkt(2004, 3, 2)},
		{"02-Mar-2004", mkt(2004, 3, 2)},
		{"2-March-68", mkt(2068, 3, 2)},
		{"02-March-69", mkt(1969, 3, 2)},
		// T_MONTHNAME T_INTEGER o_dayqual comma T_INTEGER
		{"March 2, 2004", mkt(2004, 3, 2)},
		{"Mar 02, 2004", mkt(2004, 3, 2)},
		{"March 2nd, 04", mkt(2004, 3, 2)},
		{"March 3rd, 04", mkt(2004, 3, 3)},
		{"March 2, 68", mkt(2068, 3, 2)},
		{"March 2, 69", mkt(1969, 3, 2)},
		// T_INTEGER T_MINUS T_INTEGER as YYYY-DDD
		{"2004-062", mkt(2004, 3, 2)}, // 2004 is a leap year
		{"2003-062", mkt(2003, 3, 3)},
		{"2004-001", mkt(2004, 1, 1)},
		{"2004-366", mkt(2004, 12, 31)},
		// T_INTEGER T_MINUS T_INTEGER as YYYY-MM
		{"2004-03", mkt(2004, 3, 1)},
		{"2004-3", mkt(2004, 3, 1)},
		{"2004-12", mkt(2004, 12, 1)},
		// T_INTEGER T_MINUS T_INTEGER as MM-DD
		{"3-2", mkt(1, 3, 2)},
		{"03-02", mkt(1, 3, 2)},
		{"12-31", mkt(1, 12, 31)},
		// T_INTEGER T_MINUS T_INTEGER T_MINUS T_INTEGER
		{"2004-03-02", mkt(2004, 3, 2)},
		{"2004-3-2", mkt(2004, 3, 2)},
		{"4-3-2", mkt(2004, 3, 2)},
		{"68-03-02", mkt(2068, 3, 2)},
		{"69-03-02", mkt(1969, 3, 2)},
		// T_INTEGER W T_INTEGER as YYYYWwwD
		{"2004W102", mkt(2004, 3, 2)},
		{"2003W097", mkt(2003, 3, 2)},
		{"2008W396", mkt(2008, 9, 27)}, // example from wikipedia
		{"2009W011", mkt(2008, 12, 29)}, // 2009 is special!
		{"2009W537", mkt(2010, 1, 3)},
		// T_INTEGER W T_INTEGER as YYYYWww
		{"2004W01", mkt(2003, 12, 29)},
		{"2004W02", mkt(2004, 1, 5)},
		{"2004W52", mkt(2004, 12, 20)},
		{"2004W53", mkt(2004, 12, 27)},
		// T_INTEGER T_MINUS 'W' T_INTEGER as YYYY-Www
		{"2004-W01", mkt(2003, 12, 29)},
		{"2004-W02", mkt(2004, 1, 5)},
		{"2004-W52", mkt(2004, 12, 20)},
		{"2004-W53", mkt(2004, 12, 27)},
		// T_INTEGER T_MINUS 'W' T_INTEGER T_MINUS T_INTEGER as YYYY-Www-D
		{"2004-W01-1", mkt(2003, 12, 29)},
		{"2004-W02-2", mkt(2004, 1, 6)},
		{"2004-W52-6", mkt(2004, 12, 25)},
		{"2004-W53-7", mkt(2005, 1, 2)},
	}
	for i, test := range tests {
		ret, ok := parse(test.in, rel)
		if !ok || !ret.Equal(test.t) {
			t.Errorf("Unable to parse date %d\nin: %s\ngot: %s",
				i, test.in, ret)
		}
	}
}
