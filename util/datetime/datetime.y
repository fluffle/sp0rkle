%{
package datetime

// Based upon parse-datetime.y in GNU coreutils.
// also an exercise in learning goyacc in particular.
// This file contains the yacc grammar only.
// See lexer.go for the lexer and parse functions,
// and tokenmaps.go for the token maps.

import (
	"time"
)

type textint struct {
	i, l int
	s string
}

%}

%union
{
	tval    textint
	intval  int
	zoneval *time.Location
}

%token           T_OF T_THE T_IGNORE T_DAYQUAL
%token <tval>    T_INTEGER
%token <intval>  T_PLUS T_MINUS
%token <intval>  T_MONTHNAME T_DAYNAME T_DAYS T_DAYSHIFT
%token <intval>  T_OFFSET T_ISOYD T_ISOHS T_RELATIVE T_AGO
%token <zoneval> T_ZONE

%type <intval>   sign ampm
%type <tval>     o_sign_integer o_integer
%type <zoneval>  zone o_zone

%%

spec:
	unixtime | items;

comma: ',';

o_comma: /* empty */ | comma;

o_of: /* empty */ | T_OF;

o_the: /* empty */ | T_THE;

o_dayqual: /* empty */ | T_DAYQUAL;

o_integer: /* empty */ { $$ = textint{} } | T_INTEGER;

ampm:
	'A' 'M' {
		$$ = 0
	}
	| 'A' '.' 'M' '.' {
		$$ = 0
	}
	| 'P' 'M' {
		$$ = 12
	}
	| 'P' '.' 'M' '.' {
		$$ = 12
	};

timesep: ':' | '.';

sign:
	T_PLUS | T_MINUS; 

o_sign_integer:
	T_INTEGER
	| T_PLUS T_INTEGER {
		$2.s = "+" + $2.s
		$$ = $2
	}
	| T_MINUS T_INTEGER {
		$2.s = "-" + $2.s
		$2.i *= -1
		$$ = $2
	};

zone:
	T_ZONE
	| sign T_INTEGER {
		hrs, mins := $2.i, 0
		if ($2.l == 4) {
			hrs, mins = ($2.i / 100), ($2.i % 100)
		} else if ($2.l == 2) {
			hrs *= 100
		} else {
			yylex.Error("Invalid timezone offset " +$2.s)
		}
		$$ = time.FixedZone("WTF", $1 * (3600 * hrs + 60 * mins))
	}
	| sign T_INTEGER ':' T_INTEGER {
		$$ = time.FixedZone("WTF", $1 * (3600 * $2.i + 60 * $4.i))
	};

o_zone:
	/* empty */ { $$ = nil }
	| zone;
	
unixtime:
	'@' o_sign_integer {
		l := yylex.(*dateLexer)
		if ! l.state(HAVE_TIME, true) {
			l.time = time.Unix(int64($2.i), 0)
		}
	};

items:
	/* empty */
	| items item;

item:
	time
	| iso_8601_time
	| date
	| iso_8601_date
	| iso_8601_date_time
	| day_or_month
	| relative
	| iso_8601_duration
	| integer
	| T_IGNORE;

// ISO 8601 takes care of 24h time formats, so this deals with
// 12-hour HH, HH:MM or HH:MM:SS with am/pm and optional timezone
time:
	T_INTEGER ampm o_zone {
		l := yylex.(*dateLexer)
		// Hack to allow HHMMam to parse correctly, cos adie is a mong.
		if $1.l == 3 || $1.l == 4 {
			l.setTime(ampm($1.i / 100, $2), $1.i % 100, 0, $3)
		} else {
			l.setTime(ampm($1.i, $2), 0, 0, $3)
		}
	}
	| T_INTEGER timesep T_INTEGER ampm o_zone {
		yylex.(*dateLexer).setTime($1.i + $4, $3.i, 0, $5)
	}
	| T_INTEGER timesep T_INTEGER timesep T_INTEGER ampm o_zone {
		yylex.(*dateLexer).setTime($1.i + $6, $3.i, $5.i, $7)
	};

// The "basic" ISO 8601 format (without a timezone) is lexed as
// an integer and handled in 'integer' below
iso_8601_time:
	T_INTEGER zone {
		yylex.(*dateLexer).setHMS($1.i, $1.l, $2)
	}
	| T_INTEGER timesep T_INTEGER o_zone {
		yylex.(*dateLexer).setTime($1.i, $3.i, 0, $4)
	}
	| T_INTEGER timesep T_INTEGER timesep T_INTEGER o_zone o_integer {
		yylex.(*dateLexer).setTime($1.i, $3.i, $5.i, $6)
		// Hack to make time.ANSIC, time.UnixDate and time.RubyDate parse
		if $7.l == 4 {
			yylex.(*dateLexer).setYear($7.i)
		}
	};

// ISO 8601 takes care of dash-separated big-endian date formats,
// so this deals with /-separated little-endian formats (dd/mm/yyyy)
// and more "english" ones like "20th of March 2012"
date:
	T_INTEGER '/' T_INTEGER {
		l := yylex.(*dateLexer)
		if $3.l == 4 {
			// assume we have MM/YYYY
			l.setDate($3.i, $1.i, 1)
		} else {
			// assume we have DD/MM (too bad, americans)
			l.setDate(0, $3.i, $1.i)
		}
	}
	| T_INTEGER '/' T_INTEGER '/' T_INTEGER {
		l := yylex.(*dateLexer)
		if $5.l == 4 {
			// assume we have DD/MM/YYYY
			l.setDate($5.i, $3.i, $1.i)
		} else if $5.i > 68 {
			// assume we have DD/MM/YY, add 1900 if YY > 68
			l.setDate($5.i + 1900, $3.i, $1.i)
		} else {
			// assume we have DD/MM/YY, add 2000 otherwise
			l.setDate($5.i + 2000, $3.i, $1.i)
		}
	}
	| T_THE T_INTEGER T_DAYQUAL {
		// the DDth
		yylex.(*dateLexer).setDay($2.i)
	}
	| T_THE T_INTEGER T_DAYQUAL T_OF T_MONTHNAME {
		// the DDth of Month
		yylex.(*dateLexer).setDate(0, $5, $2.i)
	}
	| T_THE T_INTEGER T_DAYQUAL T_OF T_MONTHNAME o_comma T_INTEGER {
		l := yylex.(*dateLexer)
		if $7.l == 4 {
			// the DDth of Month[,] YYYY
			l.setDate($7.i, $5, $2.i)
		} else if $7.i > 68 {
			// the DDth of Month[,] YY, add 1900 if YY > 68
			l.setDate($7.i + 1900, $5, $2.i)
		} else {
			// the DDth of Month[,] YY, add 2000 otherwise
			l.setDate($7.i + 2000, $5, $2.i)
		}
	}
	| T_INTEGER o_dayqual o_of T_MONTHNAME {
		// DD[th] [of] Month
		yylex.(*dateLexer).setDate(0, $4, $1.i)
	}
	| T_MONTHNAME o_the T_INTEGER o_dayqual {
		l := yylex.(*dateLexer)
		if $3.l == 4 {
			// assume Month YYYY
			l.setDate($3.i, $1, 1)
		} else {
			// assume Month [the] DD[th]
			l.setDate(0, $1, $3.i)
		}
	}
	| T_INTEGER o_dayqual o_of T_MONTHNAME o_comma T_INTEGER {
		l := yylex.(*dateLexer)
		if $6.l == 4 {
			// assume DD[th] [of] Month[,] YYYY
			l.setDate($6.i, $4, $1.i)
		} else if $6.i > 68 {
			// assume DD[th] [of] Month[,] YY, add 1900 if YY > 68
			l.setDate($6.i + 1900, $4, $1.i)
		} else {
			// assume DD[th] [of] Month[,] YY, add 2000 otherwise
			l.setDate($6.i + 2000, $4, $1.i)
		}
	}
	| T_INTEGER T_MINUS T_MONTHNAME T_MINUS T_INTEGER {
		// RFC 850, srsly :(
		l := yylex.(*dateLexer)
		if $5.l == 4 {
			// assume DD-Mon-YYYY
			l.setDate($5.i, $3, $1.i)
		} else if $5.i > 68 {
			// assume DD-Mon-YY, add 1900 if YY > 68
			l.setDate($5.i + 1900, $3, $1.i)
		} else {
			// assume DD-Mon-YY, add 2000 otherwise
			l.setDate($5.i + 2000, $3, $1.i)
		}
	}
	| T_MONTHNAME o_the T_INTEGER o_dayqual comma T_INTEGER {
		// comma cannot be optional here; T_MONTHNAME T_INTEGER T_INTEGER
		// can easily match [March 02 10]:30:00 and break parsing.
		l := yylex.(*dateLexer)
		if $6.l == 4 {
			// assume Month [the] DD[th], YYYY
			l.setDate($6.i, $1, $3.i)
		} else if $6.i > 68 {
			// assume Month [the] DD[th], YY, add 1900 if YY > 68
			l.setDate($6.i + 1900, $1, $3.i)
		} else {
			// assume Month [the] DD[th], YY, add 2000 otherwise
			l.setDate($6.i + 2000, $1, $3.i)
		}
	};

// The "basic" ISO 8601 format is lexed as an integer and handled in "integer"
iso_8601_date:
	T_INTEGER T_MINUS T_INTEGER {
		l := yylex.(*dateLexer)
		if $1.l == 4 && $3.l == 3 {
			// assume we have YYYY-DDD
			l.setDate($1.i, 1, $3.i)
		} else if $1.l == 4 {
			// assume we have YYYY-MM
			l.setDate($1.i, $3.i, 1)
		} else {
			// assume we have MM-DD (not strictly ISO compliant)
			// this is for americans, because of DD/MM above ;-)
			l.setDate(0, $1.i, $3.i)
		}
	}
	| T_INTEGER T_MINUS T_INTEGER T_MINUS T_INTEGER {
		l := yylex.(*dateLexer)
		if $1.l == 4 {
			// assume we have YYYY-MM-DD
			l.setDate($1.i, $3.i, $5.i)
		} else if $1.i > 68 {
			// assume we have YY-MM-DD, add 1900 if YY > 68
			l.setDate($1.i + 1900, $3.i, $5.i)
		} else {
			// assume we have YY-MM-DD, add 2000 otherwise
			l.setDate($1.i + 2000, $3.i, $5.i)
		}
	}
	| T_INTEGER 'W' T_INTEGER {
		l := yylex.(*dateLexer)
		wday, week := 1, $3.i
		if $3.l == 3 {
			// assume YYYY'W'WWD
			wday = week % 10
			week = week / 10
		}
		l.setWeek($1.i, week, wday)
	}
	| T_INTEGER T_MINUS 'W' T_INTEGER {
		// assume YYYY-'W'WW
		yylex.(*dateLexer).setWeek($1.i, $4.i, 1)
	}
	| T_INTEGER T_MINUS 'W' T_INTEGER T_MINUS T_INTEGER {
		// assume YYYY-'W'WW-D
		yylex.(*dateLexer).setWeek($1.i, $4.i, $6.i)
	};

// NOTE: this doesn't enforce that the date is complete.
iso_8601_date_time:
	iso_8601_date 'T' iso_8601_time
	| T_INTEGER 'T' T_INTEGER o_zone {
		// this goes here because the YYYYMMDD and HHMMSS forms of the
		// ISO 8601 format date and time are handled by 'integer' below.
		l := yylex.(*dateLexer)
		if $1.l == 8 {
			// assume ISO 8601 YYYYMMDD
			l.setYMD($1.i, $1.l)
		} else if $1.l == 7 {
			// assume ISO 8601 ordinal YYYYDDD
			l.setDate($1.i / 1000, 1, $1.i % 1000)
		}
		l.setHMS($3.i, $3.l, $4)
	};

day_or_month:
	T_DAYNAME o_comma {
		// Tuesday
		yylex.(*dateLexer).setDays($1, 0)
	}
	| T_MONTHNAME {
		// March
		yylex.(*dateLexer).setMonths($1, 0)
	}
	| T_RELATIVE T_DAYNAME {
		// Next tuesday
		yylex.(*dateLexer).setDays($2, $1)
	}
	| T_RELATIVE T_MONTHNAME {
		// Next march
		yylex.(*dateLexer).setMonths($2, $1)
	}
	| o_sign_integer T_DAYNAME {
		// +-N Tuesdays
		yylex.(*dateLexer).setDays($2, $1.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME {
		// 3rd Tuesday (of implicit this month)
		l := yylex.(*dateLexer)
		l.setDays($3, $1.i)
		l.setMonths(0, 0)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME T_OF T_MONTHNAME {
		// 3rd Tuesday of (implicit this) March
		l := yylex.(*dateLexer)
		l.setDays($3, $1.i)
		l.setMonths($5, 0)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME T_OF T_INTEGER {
		// 3rd Tuesday of 2012
		yylex.(*dateLexer).setDays($3, $1.i, $5.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME T_OF T_MONTHNAME T_INTEGER {
		// 3rd Tuesday of March 2012
		l := yylex.(*dateLexer)
		l.setDays($3, $1.i)
		l.setMonths($5, 0, $6.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME T_OF T_RELATIVE T_MONTHNAME {
		// 3rd Tuesday of next March
		l := yylex.(*dateLexer)
		l.setDays($3, $1.i)
		l.setMonths($6, $5)
	};

relative:
	relunits 
	| relunits T_AGO {
		yylex.(*dateLexer).setAgo()
	};

relunits:
	relunit
	| relunit relunits;

relunit:
	o_sign_integer T_OFFSET {
		yylex.(*dateLexer).addOffset(offset($2), $1.i)
	}
	| T_RELATIVE T_OFFSET {
		yylex.(*dateLexer).addOffset(offset($2), $1)
	} 
	| 'A' T_OFFSET {
		yylex.(*dateLexer).addOffset(offset($2), 1)
	} 
	| o_sign_integer T_DAYS {
		// Special-case to handle "week" and "fortnight"
		yylex.(*dateLexer).addOffset(O_DAY, $1.i * $2)
	}
	| T_RELATIVE T_DAYS {
		yylex.(*dateLexer).addOffset(O_DAY, $1 * $2)
	}
	| 'A' T_DAYS {
		yylex.(*dateLexer).addOffset(O_DAY, $2)
	}
	| o_sign_integer T_ISOYD {
		// As we need to be able to separate out YD from HS in ISO durations
		// this becomes a fair bit messier than if Y D H S were just T_OFFSET
		// Because writing "next y" or "two h" would be odd, disallow
		// T_RELATIVE tokens from being used with ISO single-letter notation
		yylex.(*dateLexer).addOffset(offset($2), $1.i)
	}
	| o_sign_integer T_ISOHS {
		yylex.(*dateLexer).addOffset(offset($2), $1.i)
	}
	| o_sign_integer 'M' {
		// Resolve 'm' ambiguity in favour of minutes outside ISO duration
		yylex.(*dateLexer).addOffset(O_MIN, $1.i)
	}
	| o_sign_integer 'W' {
		yylex.(*dateLexer).addOffset(O_DAY, $1.i * 7)
	}
	| T_DAYSHIFT {
		// yesterday or tomorrow
		yylex.(*dateLexer).addOffset(O_DAY, $1)
	};

/* date/time based durations not yet supported */
iso_8601_duration:
	'P' ymd_units o_t_hms_units
	| 'P' t_hms_units
	| 'P' T_INTEGER 'W' {
		yylex.(*dateLexer).addOffset(O_DAY, 7 * $2.i)
	};

/* This is a bit lazy compared to specifying the combinations of nYnMnS */
ymd_units:
	ymd_unit
	| ymd_units ymd_unit;

ymd_unit:
	T_INTEGER T_ISOYD {
		// takes care of Y and D
		yylex.(*dateLexer).addOffset(offset($2), $1.i)
	}
	| T_INTEGER 'M' {
		yylex.(*dateLexer).addOffset(O_MONTH, $1.i)
	};

hms_units:
	hms_unit
	| hms_units hms_unit;

hms_unit:
	T_INTEGER T_ISOHS {
		// takes care of H and S
		yylex.(*dateLexer).addOffset(offset($2), $1.i)
	}
	| T_INTEGER 'M' {
		yylex.(*dateLexer).addOffset(O_MIN, $1.i)
	};

t_hms_units:
	'T' hms_units;

o_t_hms_units:
	/* empty */
	| t_hms_units;

integer:
	T_INTEGER {
		l := yylex.(*dateLexer)
		if $1.l == 8 {
			// assume ISO 8601 YYYYMMDD
			l.setYMD($1.i, $1.l)
		} else if $1.l == 7 {
			// assume ISO 8601 ordinal YYYYDDD
			l.setDate($1.i / 1000, 1, $1.i % 1000)
		} else if $1.l == 6 {
			// assume ISO 8601 HHMMSS with no zone
			l.setHMS($1.i, $1.l, nil)
		} else if $1.l == 4 {
			// Assuming HHMM because that's more useful on IRC.
			l.setHMS($1.i, $1.l, nil)
		} else if $1.l == 2 {
			// assume HH with no zone
			l.setHMS($1.i, $1.l, nil)
		}
	};
%%
