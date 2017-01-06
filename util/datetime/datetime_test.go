package datetime

import (
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

type timeTest struct {
	in string
	t  time.Time
}

type timeTests []timeTest

func (tt timeTests) run(t *testing.T, start time.Time) {
	for i, test := range tt {
		DPrintf("\nStarting parse of %q\n\n", test.in)
		ret, err := parse(test.in, start)
		DPrintf("\nEnding parse of %q\n", test.in)
		if err != nil || !ret.Equal(test.t) {
			t.Errorf("Unable to parse test %d\nin: %s\nexp: %s\ngot: %s (err=%v)",
				i+1, test.in, test.t, ret, err)
		}
	}

}

func TestParseTimeFormats(t *testing.T) {
	local = time.UTC
	// RFC822 doesn't specify seconds, and Stamp doesn't specify year
	ref := time.Date(2004, 6, 22, 13, 10, 0, 0, time.UTC)
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
		//		time.RFC3339Nano, // fails, nanosecs not supported
		//		time.Kitchen,     // fails, only contains HH and MM
		//		time.Stamp,       // fails, no year => assumed 2013
	}
	for i, f := range formats {
		in := ref.Format(f)
		ret, err := Parse(in)
		if err != nil || !ret.Equal(ref) {
			t.Errorf("Unable to parse format %d\nin: %s\ngot: %s (err=%v)", i, in, ret, err)
		}
	}
}

func TestParseTime(t *testing.T) {
	now := time.Now().UTC()
	mkt := func(h, m, s int, l ...string) time.Time {
		loc := time.UTC
		if len(l) > 0 {
			loc = zone(l[0])
		}
		return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
	}
	tests := timeTests{
		// T_INTEGER T_AMPM o_zone (also tests all possible zone permutations)
		{"11am", mkt(11, 0, 0)},
		{"11pm", mkt(23, 0, 0)},
		{"12am", mkt(0, 0, 0)},
		{"12pm", mkt(12, 0, 0)},
		{"920am", mkt(9, 20, 0)},
		{"1140pm", mkt(23, 40, 0)},
		{"3 am PDT", mkt(3, 0, 0, "US/Pacific")},
		{"3 pm PDT", mkt(15, 0, 0, "US/Pacific")},
		{"5AM -4:00", mkt(5, 0, 0, "Etc/GMT+4")},
		{"5PM -4:00", mkt(17, 0, 0, "Etc/GMT+4")},
		{"7 a.m. +0800", mkt(7, 0, 0, "Etc/GMT-8")},
		{"7 p.m. +0800", mkt(19, 0, 0, "Etc/GMT-8")},
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
		{"midday", mkt(12, 0, 0)},
		{"midnight", mkt(0, 0, 0)},
	}
	tests.run(t, now)
}

// Stolen from pkg/time/zoneinfo_read.go
// get4 returns the little-endian 32-bit value in b.
func get4(b []byte) int {
	if len(b) < 4 {
		return 0
	}
	return int(b[0]) | int(b[1])<<8 | int(b[2])<<16 | int(b[3])<<24
}

// get2 returns the little-endian 16-bit value in b.
func get2(b []byte) int {
	if len(b) < 2 {
		return 0
	}
	return int(b[0]) | int(b[1])<<8
}

func TestParseAllTimezonesInZoneinfo(t *testing.T) {
	// Ugh, half-parsing zip files, that's what the stdlib does.
	zip := path.Join(runtime.GOROOT(), "lib", "time", "zoneinfo.zip")
	fh, err := os.Open(zip)
	if err != nil {
		t.Fatalf("Failed to open zoneinfo.zip: %v", err)
		return
	}
	// Zip tail is 22 bytes at end of file.
	zTailSize := 22
	if _, err := fh.Seek(int64(-zTailSize), 2); err != nil {
		t.Fatalf("Seeking to tail of zoneinfo.zip: %v", err)
		return
	}
	buf := make([]byte, zTailSize)
	if n, err := fh.Read(buf); err != nil || n != zTailSize || get4(buf) != 0x06054b50 {
		t.Fatalf("Reading zoneinfo.zip tail: %v (%d of %d bytes)", err, n, zTailSize)
		return
	}
	n, size, offset := get2(buf[10:]), get4(buf[12:]), get4(buf[16:])
	buf = make([]byte, size)
	if _, err := fh.Seek(int64(offset), 0); err != nil {
		t.Fatalf("Seeking to data offset in zoneinfo.zip: %v", err)
		return
	}
	if n, err := fh.Read(buf); err != nil || n != size {
		t.Fatalf("Reading zoneinfo.zip data: %v (%d of %d bytes)", err, n, size)
		return
	}

	now := time.Now().UTC()
	mkt := func(h, m, s int, l ...string) time.Time {
		loc := time.Local
		if len(l) > 0 {
			loc = zone(l[0])
		}
		return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, loc)
	}
	tests := make(timeTests, 0, n)

	for i := 0; i < n; i++ {
		if get4(buf) != 0x02014b50 {
			break
		}
		ucsize := get4(buf[24:])
		namelen := get2(buf[28:])
		xlen := get2(buf[30:])
		fclen := get2(buf[32:])
		name := string(buf[46 : 46+namelen])
		buf = buf[46+namelen+xlen+fclen:]
		if ucsize == 0 || zone(name) == nil {
			// The zip header contains dirs and other things that
			// aren't real zones, skip them to avoid test fail.
			continue
		}
		tests = append(tests, timeTest{"6am " + name, mkt(6, 0, 0, name)})
	}
	tests.run(t, now)
}

func TestParseDate(t *testing.T) {
	h, n, s := 11, 22, 33
	mkt := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, h, n, s, 0, time.UTC)
	}
	rel := mkt(1, 2, 3)
	tests := timeTests{
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
		// T_THE T_INTEGER T_DAYQUAL
		{"the 1st", mkt(1, 2, 1)},
		{"the 2nd", mkt(1, 2, 2)},
		{"the 10th", mkt(1, 2, 10)},
		{"the 29th", mkt(1, 3, 1)},
		// T_THE T_INTEGER T_DAYQUAL T_OF T_MONTHNAME
		{"the 1st of January", mkt(1, 1, 1)},
		{"the 2nd of February", mkt(1, 2, 2)},
		{"the 10th of March", mkt(1, 3, 10)},
		{"the 31st of December", mkt(1, 12, 31)},
		// T_THE T_INTEGER T_DAYQUAL T_OF T_MONTHNAME o_comma T_INTEGER
		{"the 1st of January, 2001", mkt(2001, 1, 1)},
		{"the 2nd of February 2013", mkt(2013, 2, 2)},
		{"the 10th of March 67", mkt(2067, 3, 10)},
		{"the 31st of December, 69", mkt(1969, 12, 31)},
		// T_INTEGER o_dayqual o_of T_MONTHNAME
		{"2 Mar", mkt(1, 3, 2)},
		{"02 Mar", mkt(1, 3, 2)},
		{"2nd March", mkt(1, 3, 2)},
		{"3rd of March", mkt(1, 3, 3)},
		// T_MONTHNAME o_the T_INTEGER o_dayqual as Month YYYY
		{"March 2004", mkt(2004, 3, 1)},
		// T_MONTHNAME o_the T_INTEGER o_dayqual as Month [the] DD[th]
		{"March 2", mkt(1, 3, 2)},
		{"March 02", mkt(1, 3, 2)},
		{"Mar 2nd", mkt(1, 3, 2)},
		{"Mar the 3rd", mkt(1, 3, 3)},
		// T_INTEGER o_dayqual o_of T_MONTHNAME o_comma T_INTEGER
		{"2 Mar 2004", mkt(2004, 3, 2)},
		{"02 Mar 2004", mkt(2004, 3, 2)},
		{"2nd March, 2004", mkt(2004, 3, 2)},
		{"3rd of March 2004", mkt(2004, 3, 3)},
		{"2 March 4", mkt(2004, 3, 2)},
		{"2 March 68", mkt(2068, 3, 2)},
		{"2 March, 69", mkt(1969, 3, 2)},
		// T_INTEGER T_MINUS T_MONTHNAME T_MINUS T_INTEGER
		{"2-Mar-2004", mkt(2004, 3, 2)},
		{"02-Mar-2004", mkt(2004, 3, 2)},
		{"2-March-68", mkt(2068, 3, 2)},
		{"02-March-69", mkt(1969, 3, 2)},
		// T_MONTHNAME o_the T_INTEGER o_dayqual comma T_INTEGER
		{"March 2, 2004", mkt(2004, 3, 2)},
		{"Mar 02, 2004", mkt(2004, 3, 2)},
		{"March the 2nd, 04", mkt(2004, 3, 2)},
		{"March the 3rd, 04", mkt(2004, 3, 3)},
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
		{"2008W396", mkt(2008, 9, 27)},  // example from wikipedia
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
	tests.run(t, rel)
}

func TestParseIsoDateTime(t *testing.T) {
	tests := timeTests{
		// some random iso_8601_date 'T' iso_8601_time tests
		{"2004-03-02T13:14:15",
			time.Date(2004, 3, 2, 13, 14, 15, 0, time.Local)},
		{"2004-062T13:14Z",
			time.Date(2004, 3, 2, 13, 14, 0, 0, time.UTC)},
		{"2004W102T13+0400",
			time.Date(2004, 3, 2, 13, 0, 0, 0, zone("Etc/GMT-4"))},
		{"2004-W10-2T13:14:15-08:00",
			time.Date(2004, 3, 2, 13, 14, 15, 0, zone("Etc/GMT+8"))},
		// T_INTEGER 'T' T_INTEGER o_zone
		{"20040302T131415",
			time.Date(2004, 3, 2, 13, 14, 15, 0, time.Local)},
		{"2004062T1314Z",
			time.Date(2004, 3, 2, 13, 14, 0, 0, time.UTC)},
		{"20040302T13+0400",
			time.Date(2004, 3, 2, 13, 0, 0, 0, zone("Etc/GMT-4"))},
	}
	tests.run(t, time.Time{})
}

func TestParseRelativeDays(t *testing.T) {
	mkt := func(off int) time.Time {
		// return offset from Wed 22nd Jan 2014
		return time.Date(2014, 1, 22+off, 0, 0, 0, 0, time.UTC)
	}
	rel := mkt(0)
	tests := timeTests{
		{"wednesday", mkt(0)},
		{"this wednesday", mkt(0)},
		{"next wednesday", mkt(7)},
		{"last wednesday", mkt(-7)},
		{"this day", mkt(0)},
		{"first wednesday", mkt(7)},
		{"second wednesday", mkt(14)},
		{"thursday", mkt(1)},
		{"this thursday", mkt(1)},
		{"next thursday", mkt(1)},
		{"last thursday", mkt(-6)},
		{"next day", mkt(1)},
		{"first day", mkt(1)},
		{"friday", mkt(2)},
		{"this friday", mkt(2)},
		{"next friday", mkt(2)},
		{"last friday", mkt(-5)},
		{"second day", mkt(2)},
		{"saturday", mkt(3)},
		{"this saturday", mkt(3)},
		{"next saturday", mkt(3)},
		{"last saturday", mkt(-4)},
		{"third day", mkt(3)},
		{"sunday", mkt(4)},
		{"this sunday", mkt(4)},
		{"next sunday", mkt(4)},
		{"last sunday", mkt(-3)},
		{"monday", mkt(5)},
		{"this monday", mkt(5)},
		{"next monday", mkt(5)},
		{"last monday", mkt(-2)},
		{"tuesday", mkt(6)},
		{"this tuesday", mkt(6)},
		{"next tuesday", mkt(6)},
		{"last tuesday", mkt(-1)},
		{"2 wednesdays", mkt(14)},
		{"-3 wednesdays", mkt(-21)},
		{"+4 wednesdays", mkt(28)},
		{"yesterday", mkt(-1)},
		{"tomorrow", mkt(1)},
		{"today", mkt(0)},
		{"now", mkt(0)},
	}
	tests.run(t, rel)
}

func TestParseRelativeMonths(t *testing.T) {
	mkt := func(off int) time.Time {
		// return offset from Wed 22nd June 2014
		return time.Date(2014, time.Month(6+off), 22, 0, 0, 0, 0, time.UTC)
	}
	rel := mkt(0)
	tests := timeTests{
		{"june", mkt(0)},
		{"this june", mkt(0)},
		{"next june", mkt(12)},
		{"last june", mkt(-12)},
		{"this month", mkt(0)},
		{"second june", mkt(24)},
		{"july", mkt(1)},
		{"this july", mkt(1)},
		{"next july", mkt(1)},
		{"last july", mkt(-11)},
		{"next month", mkt(1)},
		{"first month", mkt(1)},
		{"august", mkt(2)},
		{"this august", mkt(2)},
		{"next august", mkt(2)},
		{"last august", mkt(-10)},
		{"second month", mkt(2)},
		{"september", mkt(3)},
		{"this september", mkt(3)},
		{"next september", mkt(3)},
		{"last september", mkt(-9)},
		{"october", mkt(4)},
		{"this october", mkt(4)},
		{"next october", mkt(4)},
		{"last october", mkt(-8)},
		{"november", mkt(5)},
		{"this november", mkt(5)},
		{"next november", mkt(5)},
		{"last november", mkt(-7)},
		{"december", mkt(6)},
		{"this december", mkt(6)},
		{"next december", mkt(6)},
		{"last december", mkt(-6)},
		{"january", mkt(-5)},
		{"this january", mkt(-5)},
		{"next january", mkt(7)},
		{"last january", mkt(-5)},
		{"february", mkt(-4)},
		{"this february", mkt(-4)},
		{"next february", mkt(8)},
		{"last february", mkt(-4)},
		{"march", mkt(-3)},
		{"this march", mkt(-3)},
		{"next march", mkt(9)},
		{"last march", mkt(-3)},
		{"april", mkt(-2)},
		{"this april", mkt(-2)},
		{"next april", mkt(10)},
		{"last april", mkt(-2)},
		{"may", mkt(-1)},
		{"this may", mkt(-1)},
		{"next may", mkt(11)},
		{"last may", mkt(-1)},
	}
	tests.run(t, rel)
}

func TestAbsDayMonth(t *testing.T) {
	h, n, s := 11, 22, 33
	mkt := func(y, m, d int) time.Time {
		return time.Date(y, time.Month(m), d, h, n, s, 0, time.UTC)
	}
	rel := mkt(2001, 2, 3)
	tests := timeTests{
		// ... of implicitly this month
		{"1st Monday", mkt(2001, 2, 5)},
		{"1st Wednesday", mkt(2001, 2, 7)},
		{"1st Thursday", mkt(2001, 2, 1)},
		{"1st Sunday", mkt(2001, 2, 4)},
		{"2nd Monday", mkt(2001, 2, 12)},
		{"2nd Wednesday", mkt(2001, 2, 14)},
		{"2nd Thursday", mkt(2001, 2, 8)},
		{"2nd Sunday", mkt(2001, 2, 11)},
		// ... of explicit month
		{"3rd Saturday of December", mkt(2000, 12, 16)},
		{"2nd Sunday of January", mkt(2001, 1, 14)},
		{"4th Thursday of February", mkt(2001, 2, 22)},
		{"1st Tuesday of March", mkt(2001, 3, 6)},
		{"2nd Wednesday of August", mkt(2001, 8, 8)},
		{"2nd Wednesday of September", mkt(2000, 9, 13)},
		// ... of explicit month of year
		{"3rd Tuesday of January 2014", mkt(2014, 1, 21)},
		{"3rd Friday of January 2014", mkt(2014, 1, 17)},
		{"1st Monday of April 2014", mkt(2014, 4, 7)},
		{"1st Wednesday of April 2014", mkt(2014, 4, 2)},
		// ... of explicit year
		{"1st Tuesday of 2014", mkt(2014, 1, 7)},
		{"1st Wednesday of 2014", mkt(2014, 1, 1)},
		{"1st Friday of 2014", mkt(2014, 1, 3)},
		{"1st Sunday of 2014", mkt(2014, 1, 5)},
		{"2nd Tuesday of 2014", mkt(2014, 1, 14)},
		{"2nd Wednesday of 2014", mkt(2014, 1, 8)},
		{"2nd Friday of 2014", mkt(2014, 1, 10)},
		{"2nd Sunday of 2014", mkt(2014, 1, 12)},
	}
	tests.run(t, rel)
}

/* I'm not smart enough for this :-(
func TestRelDayMonth(t *testing.T) {
	mkt := func(day, month int) time.Time {
		// Absolute day of month relative to June 2014
		return time.Date(2014, time.Month(6+month), day, 0, 0, 0, 0, time.UTC)
	}
	// Our base date is 22nd June 2014
	rel := mkt(22, 0)
	tests := timeTests{
		{"this day of the month", mkt(22, 0)},
		{"first day of this month", mkt(1, 0)},
		{"first sunday of this month", mkt(1, 0)},
		{"first tuesday of this month", mkt(3, 0)},
		{"first saturday of this month", mkt(7, 0)},
		{"first day of next month", mkt(1, 1)},
		{"first tuesday of next month", mkt(1, 1)},
		{"first thursday of next month", mkt(3, 1)},
		{"first monday of next month", mkt(7, 1)},
		{"third tuesday of next month", mkt(15, 1)},
		{"second day of last month", mkt(2, -1)},
		{"second wednesday of last month", mkt(14, -1)},
		// Testing the behaviour of the time package too.
		{"last day of this month", mkt(30, 0)},
		{"last day of this month", mkt(0, 1)},
	}
	tests.run(t, rel)
}
*/

func TestUnixTime(t *testing.T) {
	rel := time.Now().UTC()
	tests := timeTests{
		{"@1234567890", time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC)},
	}
	tests.run(t, rel)
}
