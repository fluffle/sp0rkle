
// line 1 "prefixes.rl"
package util

/* 
 * This file requires version 6.7 of Ragel. Build the .go source with:
 *   $ ragel -Z -G2 -o prefixes.go prefixes.rl
 *
 * Hopefully it won't confuse the hell out of gb...
 *
 * This implements a simple state-machine parser for removing common 
 * prefixes for strings said on IRC. It's done here for the moment, until
 * the new Go exp/regex pkg is in release, because the regexp pkg sucks.
 *
 * We're replicating the following code from perlfu -- (c) jinzougen:
 *  qr/ o*k | see   | uh+m* | hey  | actually | ooo+
 * | we+ll+ | iirc  | but   | and  | or       | eh
 * | um+    | \.+   | like  | o+h+ | yea+h*
 * | yup    | lol   | wow   | hm+  | [ha]{2,} | [he]{3,}
 * | er+ /xi;
 *
 * Thank god ragel's syntax is pretty much regex/BNF ;-)
 */


// line 24 "prefixes.rl"

// line 29 "prefixes.go"
var config_start int = 47
var config_first_final int = 47
var config_error int = 0

var config_en_main int = 47


// line 25 "prefixes.rl"

// basic parser bits

// line 62 "prefixes.rl"


func RemovePrefixes(s string) string {
	data := make([]byte, len(s))
	copy(data, s)
	cs, p, m, pe, eof := 0, 0, 0, len(data), len(data)
	
	
// line 50 "prefixes.go"
	cs = config_start

// line 70 "prefixes.rl"
	
// line 55 "prefixes.go"
	{
	if p == pe { goto _test_eof }
	switch cs {
	case -666: // i am a hack D:
	fallthrough
case 47:
	switch data[p] {
		case 46: goto st1
		case 97: goto st2
		case 98: goto st4
		case 101: goto st7
		case 104: goto st12
		case 105: goto st16
		case 107: goto st6
		case 108: goto st19
		case 111: goto st23
		case 115: goto st27
		case 117: goto st28
		case 119: goto st30
		case 121: goto st35
	}
	goto st0
st0:
cs = 0;
	goto _out;
tr49:
// line 28 "prefixes.rl"
	{ m = p }
	goto st1
st1:
	p++
	if p == pe { goto _test_eof1 }
	fallthrough
case 1:
// line 90 "prefixes.go"
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 46: goto st1
	}
	goto st0
st48:
	p++
	if p == pe { goto _test_eof48 }
	fallthrough
case 48:
	switch data[p] {
		case 46: goto tr49
		case 97: goto tr50
		case 98: goto tr51
		case 101: goto tr52
		case 104: goto tr53
		case 105: goto tr54
		case 107: goto tr55
		case 108: goto tr56
		case 111: goto tr57
		case 115: goto tr58
		case 117: goto tr59
		case 119: goto tr60
		case 121: goto tr61
	}
	goto st0
tr50:
// line 28 "prefixes.rl"
	{ m = p }
	goto st2
st2:
	p++
	if p == pe { goto _test_eof2 }
	fallthrough
case 2:
// line 127 "prefixes.go"
	switch data[p] {
		case 97: goto st3
		case 99: goto st40
		case 104: goto st3
		case 110: goto st46
	}
	goto st0
st3:
	p++
	if p == pe { goto _test_eof3 }
	fallthrough
case 3:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 97: goto st3
		case 104: goto st3
	}
	goto st0
st49:
	p++
	if p == pe { goto _test_eof49 }
	fallthrough
case 49:
	switch data[p] {
		case 32: goto st48
		case 46: goto tr49
		case 97: goto tr50
		case 98: goto tr51
		case 101: goto tr52
		case 104: goto tr53
		case 105: goto tr54
		case 107: goto tr55
		case 108: goto tr56
		case 111: goto tr57
		case 115: goto tr58
		case 117: goto tr59
		case 119: goto tr60
		case 121: goto tr61
	}
	goto st0
tr51:
// line 28 "prefixes.rl"
	{ m = p }
	goto st4
st4:
	p++
	if p == pe { goto _test_eof4 }
	fallthrough
case 4:
// line 178 "prefixes.go"
	if data[p] == 117 { goto st5 }
	goto st0
st5:
	p++
	if p == pe { goto _test_eof5 }
	fallthrough
case 5:
	if data[p] == 116 { goto st6 }
	goto st0
tr55:
// line 28 "prefixes.rl"
	{ m = p }
	goto st6
st6:
	p++
	if p == pe { goto _test_eof6 }
	fallthrough
case 6:
// line 197 "prefixes.go"
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
	}
	goto st0
tr52:
// line 28 "prefixes.rl"
	{ m = p }
	goto st7
st7:
	p++
	if p == pe { goto _test_eof7 }
	fallthrough
case 7:
// line 212 "prefixes.go"
	switch data[p] {
		case 101: goto st8
		case 104: goto st10
		case 114: goto st11
	}
	goto st0
st8:
	p++
	if p == pe { goto _test_eof8 }
	fallthrough
case 8:
	switch data[p] {
		case 101: goto st9
		case 104: goto st10
		case 114: goto st11
	}
	goto st0
st9:
	p++
	if p == pe { goto _test_eof9 }
	fallthrough
case 9:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 101: goto st9
		case 104: goto st10
		case 114: goto st11
	}
	goto st0
st10:
	p++
	if p == pe { goto _test_eof10 }
	fallthrough
case 10:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 101: goto st10
		case 104: goto st10
	}
	goto st0
st11:
	p++
	if p == pe { goto _test_eof11 }
	fallthrough
case 11:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 114: goto st11
	}
	goto st0
tr53:
// line 28 "prefixes.rl"
	{ m = p }
	goto st12
st12:
	p++
	if p == pe { goto _test_eof12 }
	fallthrough
case 12:
// line 275 "prefixes.go"
	switch data[p] {
		case 97: goto st3
		case 101: goto st13
		case 104: goto st14
		case 109: goto st15
	}
	goto st0
st13:
	p++
	if p == pe { goto _test_eof13 }
	fallthrough
case 13:
	switch data[p] {
		case 101: goto st10
		case 104: goto st10
		case 121: goto st6
	}
	goto st0
st14:
	p++
	if p == pe { goto _test_eof14 }
	fallthrough
case 14:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 97: goto st3
		case 101: goto st10
		case 104: goto st14
		case 109: goto st15
	}
	goto st0
st15:
	p++
	if p == pe { goto _test_eof15 }
	fallthrough
case 15:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 109: goto st15
	}
	goto st0
tr54:
// line 28 "prefixes.rl"
	{ m = p }
	goto st16
st16:
	p++
	if p == pe { goto _test_eof16 }
	fallthrough
case 16:
// line 328 "prefixes.go"
	if data[p] == 105 { goto st17 }
	goto st0
st17:
	p++
	if p == pe { goto _test_eof17 }
	fallthrough
case 17:
	if data[p] == 114 { goto st18 }
	goto st0
st18:
	p++
	if p == pe { goto _test_eof18 }
	fallthrough
case 18:
	if data[p] == 99 { goto st6 }
	goto st0
tr56:
// line 28 "prefixes.rl"
	{ m = p }
	goto st19
st19:
	p++
	if p == pe { goto _test_eof19 }
	fallthrough
case 19:
// line 354 "prefixes.go"
	switch data[p] {
		case 105: goto st20
		case 111: goto st22
	}
	goto st0
st20:
	p++
	if p == pe { goto _test_eof20 }
	fallthrough
case 20:
	if data[p] == 107 { goto st21 }
	goto st0
st21:
	p++
	if p == pe { goto _test_eof21 }
	fallthrough
case 21:
	if data[p] == 101 { goto st6 }
	goto st0
st22:
	p++
	if p == pe { goto _test_eof22 }
	fallthrough
case 22:
	if data[p] == 108 { goto st6 }
	goto st0
tr57:
// line 28 "prefixes.rl"
	{ m = p }
	goto st23
st23:
	p++
	if p == pe { goto _test_eof23 }
	fallthrough
case 23:
// line 390 "prefixes.go"
	switch data[p] {
		case 104: goto st24
		case 107: goto st6
		case 111: goto st25
		case 114: goto st6
	}
	goto st0
st24:
	p++
	if p == pe { goto _test_eof24 }
	fallthrough
case 24:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 104: goto st24
	}
	goto st0
st25:
	p++
	if p == pe { goto _test_eof25 }
	fallthrough
case 25:
	switch data[p] {
		case 104: goto st24
		case 107: goto st6
		case 111: goto st26
	}
	goto st0
st26:
	p++
	if p == pe { goto _test_eof26 }
	fallthrough
case 26:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 104: goto st24
		case 107: goto st6
		case 111: goto st26
	}
	goto st0
tr58:
// line 28 "prefixes.rl"
	{ m = p }
	goto st27
st27:
	p++
	if p == pe { goto _test_eof27 }
	fallthrough
case 27:
// line 442 "prefixes.go"
	if data[p] == 101 { goto st21 }
	goto st0
tr59:
// line 28 "prefixes.rl"
	{ m = p }
	goto st28
st28:
	p++
	if p == pe { goto _test_eof28 }
	fallthrough
case 28:
// line 454 "prefixes.go"
	switch data[p] {
		case 104: goto st29
		case 109: goto st15
	}
	goto st0
st29:
	p++
	if p == pe { goto _test_eof29 }
	fallthrough
case 29:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 104: goto st29
		case 109: goto st15
	}
	goto st0
tr60:
// line 28 "prefixes.rl"
	{ m = p }
	goto st30
st30:
	p++
	if p == pe { goto _test_eof30 }
	fallthrough
case 30:
// line 481 "prefixes.go"
	switch data[p] {
		case 101: goto st31
		case 111: goto st34
	}
	goto st0
st31:
	p++
	if p == pe { goto _test_eof31 }
	fallthrough
case 31:
	switch data[p] {
		case 101: goto st31
		case 108: goto st32
	}
	goto st0
st32:
	p++
	if p == pe { goto _test_eof32 }
	fallthrough
case 32:
	if data[p] == 108 { goto st33 }
	goto st0
st33:
	p++
	if p == pe { goto _test_eof33 }
	fallthrough
case 33:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 108: goto st33
	}
	goto st0
st34:
	p++
	if p == pe { goto _test_eof34 }
	fallthrough
case 34:
	if data[p] == 119 { goto st6 }
	goto st0
tr61:
// line 28 "prefixes.rl"
	{ m = p }
	goto st35
st35:
	p++
	if p == pe { goto _test_eof35 }
	fallthrough
case 35:
// line 531 "prefixes.go"
	switch data[p] {
		case 97: goto st36
		case 101: goto st37
		case 117: goto st39
	}
	goto st0
st36:
	p++
	if p == pe { goto _test_eof36 }
	fallthrough
case 36:
	switch data[p] {
		case 97: goto st36
		case 104: goto st24
	}
	goto st0
st37:
	p++
	if p == pe { goto _test_eof37 }
	fallthrough
case 37:
	switch data[p] {
		case 97: goto st38
		case 101: goto st37
		case 104: goto st24
	}
	goto st0
st38:
	p++
	if p == pe { goto _test_eof38 }
	fallthrough
case 38:
	switch data[p] {
		case 32: goto st48
		case 44: goto st49
		case 97: goto st38
		case 104: goto st24
	}
	goto st0
st39:
	p++
	if p == pe { goto _test_eof39 }
	fallthrough
case 39:
	if data[p] == 112 { goto st6 }
	goto st0
st40:
	p++
	if p == pe { goto _test_eof40 }
	fallthrough
case 40:
	if data[p] == 116 { goto st41 }
	goto st0
st41:
	p++
	if p == pe { goto _test_eof41 }
	fallthrough
case 41:
	if data[p] == 117 { goto st42 }
	goto st0
st42:
	p++
	if p == pe { goto _test_eof42 }
	fallthrough
case 42:
	if data[p] == 97 { goto st43 }
	goto st0
st43:
	p++
	if p == pe { goto _test_eof43 }
	fallthrough
case 43:
	if data[p] == 108 { goto st44 }
	goto st0
st44:
	p++
	if p == pe { goto _test_eof44 }
	fallthrough
case 44:
	if data[p] == 108 { goto st45 }
	goto st0
st45:
	p++
	if p == pe { goto _test_eof45 }
	fallthrough
case 45:
	if data[p] == 121 { goto st6 }
	goto st0
st46:
	p++
	if p == pe { goto _test_eof46 }
	fallthrough
case 46:
	if data[p] == 100 { goto st6 }
	goto st0
	}
	_test_eof1: cs = 1; goto _test_eof; 
	_test_eof48: cs = 48; goto _test_eof; 
	_test_eof2: cs = 2; goto _test_eof; 
	_test_eof3: cs = 3; goto _test_eof; 
	_test_eof49: cs = 49; goto _test_eof; 
	_test_eof4: cs = 4; goto _test_eof; 
	_test_eof5: cs = 5; goto _test_eof; 
	_test_eof6: cs = 6; goto _test_eof; 
	_test_eof7: cs = 7; goto _test_eof; 
	_test_eof8: cs = 8; goto _test_eof; 
	_test_eof9: cs = 9; goto _test_eof; 
	_test_eof10: cs = 10; goto _test_eof; 
	_test_eof11: cs = 11; goto _test_eof; 
	_test_eof12: cs = 12; goto _test_eof; 
	_test_eof13: cs = 13; goto _test_eof; 
	_test_eof14: cs = 14; goto _test_eof; 
	_test_eof15: cs = 15; goto _test_eof; 
	_test_eof16: cs = 16; goto _test_eof; 
	_test_eof17: cs = 17; goto _test_eof; 
	_test_eof18: cs = 18; goto _test_eof; 
	_test_eof19: cs = 19; goto _test_eof; 
	_test_eof20: cs = 20; goto _test_eof; 
	_test_eof21: cs = 21; goto _test_eof; 
	_test_eof22: cs = 22; goto _test_eof; 
	_test_eof23: cs = 23; goto _test_eof; 
	_test_eof24: cs = 24; goto _test_eof; 
	_test_eof25: cs = 25; goto _test_eof; 
	_test_eof26: cs = 26; goto _test_eof; 
	_test_eof27: cs = 27; goto _test_eof; 
	_test_eof28: cs = 28; goto _test_eof; 
	_test_eof29: cs = 29; goto _test_eof; 
	_test_eof30: cs = 30; goto _test_eof; 
	_test_eof31: cs = 31; goto _test_eof; 
	_test_eof32: cs = 32; goto _test_eof; 
	_test_eof33: cs = 33; goto _test_eof; 
	_test_eof34: cs = 34; goto _test_eof; 
	_test_eof35: cs = 35; goto _test_eof; 
	_test_eof36: cs = 36; goto _test_eof; 
	_test_eof37: cs = 37; goto _test_eof; 
	_test_eof38: cs = 38; goto _test_eof; 
	_test_eof39: cs = 39; goto _test_eof; 
	_test_eof40: cs = 40; goto _test_eof; 
	_test_eof41: cs = 41; goto _test_eof; 
	_test_eof42: cs = 42; goto _test_eof; 
	_test_eof43: cs = 43; goto _test_eof; 
	_test_eof44: cs = 44; goto _test_eof; 
	_test_eof45: cs = 45; goto _test_eof; 
	_test_eof46: cs = 46; goto _test_eof; 

	_test_eof: {}
	if p == eof {
	switch cs {
	case 48, 49:
// line 28 "prefixes.rl"
	{ m = p }
	break
// line 684 "prefixes.go"
	}
	}

	_out: {}
	}

// line 71 "prefixes.rl"

	if m > 0 {
		return s[m:]
	}
	return s
}






