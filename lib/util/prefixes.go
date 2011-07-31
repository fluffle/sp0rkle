
// line 1 "prefixes.rl"
package util

import "fmt"

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


// line 26 "prefixes.rl"

// line 31 "prefixes.go"
var config_start int = 1
var config_first_final int = 1
var config_error int = 0

var config_en_main int = 1


// line 27 "prefixes.rl"

// basic parser bits

// line 64 "prefixes.rl"


func RemovePrefixes(s string) string {
	data := make([]byte, len(s))
	copy(data, s)
	cs, p, m, pe, eof := 0, 0, 0, len(data), len(data)
	
	
// line 52 "prefixes.go"
	cs = config_start

// line 72 "prefixes.rl"
	
// line 57 "prefixes.go"
	{
	if p == pe { goto _test_eof }
	switch cs {
	case -666: // i am a hack D:
	fallthrough
case 1:
	switch data[p] {
		case 46: goto st3
		case 97: goto st5
		case 98: goto st8
		case 101: goto st11
		case 104: goto st16
		case 105: goto st20
		case 107: goto st10
		case 108: goto st23
		case 111: goto st27
		case 115: goto st31
		case 117: goto st32
		case 119: goto st34
		case 121: goto st39
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr17:
// line 30 "prefixes.rl"
	{ m = p }
	goto st2
st2:
	p++
	if p == pe { goto _test_eof2 }
	fallthrough
case 2:
// line 90 "prefixes.go"
	if data[p] <= 127 { goto st2 }
	goto st0
st0:
cs = 0;
	goto _out;
tr18:
// line 30 "prefixes.rl"
	{ m = p }
	goto st3
st3:
	p++
	if p == pe { goto _test_eof3 }
	fallthrough
case 3:
// line 105 "prefixes.go"
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 46: goto st3
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr34:
// line 30 "prefixes.rl"
	{ m = p }
	goto st4
st4:
	p++
	if p == pe { goto _test_eof4 }
	fallthrough
case 4:
// line 122 "prefixes.go"
	switch data[p] {
		case 46: goto tr18
		case 97: goto tr19
		case 98: goto tr20
		case 101: goto tr21
		case 104: goto tr22
		case 105: goto tr23
		case 107: goto tr24
		case 108: goto tr25
		case 111: goto tr26
		case 115: goto tr27
		case 117: goto tr28
		case 119: goto tr29
		case 121: goto tr30
	}
	if data[p] <= 127 { goto tr17 }
	goto st0
tr19:
// line 30 "prefixes.rl"
	{ m = p }
	goto st5
st5:
	p++
	if p == pe { goto _test_eof5 }
	fallthrough
case 5:
// line 149 "prefixes.go"
	switch data[p] {
		case 97: goto st6
		case 99: goto st44
		case 104: goto st6
		case 110: goto st50
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st6:
	p++
	if p == pe { goto _test_eof6 }
	fallthrough
case 6:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 97: goto st6
		case 104: goto st6
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st7:
	p++
	if p == pe { goto _test_eof7 }
	fallthrough
case 7:
	switch data[p] {
		case 32: goto tr34
		case 46: goto tr18
		case 97: goto tr19
		case 98: goto tr20
		case 101: goto tr21
		case 104: goto tr22
		case 105: goto tr23
		case 107: goto tr24
		case 108: goto tr25
		case 111: goto tr26
		case 115: goto tr27
		case 117: goto tr28
		case 119: goto tr29
		case 121: goto tr30
	}
	if data[p] <= 127 { goto tr17 }
	goto st0
tr20:
// line 30 "prefixes.rl"
	{ m = p }
	goto st8
st8:
	p++
	if p == pe { goto _test_eof8 }
	fallthrough
case 8:
// line 203 "prefixes.go"
	if data[p] == 117 { goto st9 }
	if data[p] <= 127 { goto st2 }
	goto st0
st9:
	p++
	if p == pe { goto _test_eof9 }
	fallthrough
case 9:
	if data[p] == 116 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
tr24:
// line 30 "prefixes.rl"
	{ m = p }
	goto st10
st10:
	p++
	if p == pe { goto _test_eof10 }
	fallthrough
case 10:
// line 224 "prefixes.go"
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr21:
// line 30 "prefixes.rl"
	{ m = p }
	goto st11
st11:
	p++
	if p == pe { goto _test_eof11 }
	fallthrough
case 11:
// line 240 "prefixes.go"
	switch data[p] {
		case 101: goto st12
		case 104: goto st14
		case 114: goto st15
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st12:
	p++
	if p == pe { goto _test_eof12 }
	fallthrough
case 12:
	switch data[p] {
		case 101: goto st13
		case 104: goto st14
		case 114: goto st15
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st13:
	p++
	if p == pe { goto _test_eof13 }
	fallthrough
case 13:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 101: goto st13
		case 104: goto st14
		case 114: goto st15
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st14:
	p++
	if p == pe { goto _test_eof14 }
	fallthrough
case 14:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 101: goto st14
		case 104: goto st14
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st15:
	p++
	if p == pe { goto _test_eof15 }
	fallthrough
case 15:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 114: goto st15
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr22:
// line 30 "prefixes.rl"
	{ m = p }
	goto st16
st16:
	p++
	if p == pe { goto _test_eof16 }
	fallthrough
case 16:
// line 308 "prefixes.go"
	switch data[p] {
		case 97: goto st6
		case 101: goto st17
		case 104: goto st18
		case 109: goto st19
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st17:
	p++
	if p == pe { goto _test_eof17 }
	fallthrough
case 17:
	switch data[p] {
		case 101: goto st14
		case 104: goto st14
		case 121: goto st10
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st18:
	p++
	if p == pe { goto _test_eof18 }
	fallthrough
case 18:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 97: goto st6
		case 101: goto st14
		case 104: goto st18
		case 109: goto st19
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st19:
	p++
	if p == pe { goto _test_eof19 }
	fallthrough
case 19:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 109: goto st19
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr23:
// line 30 "prefixes.rl"
	{ m = p }
	goto st20
st20:
	p++
	if p == pe { goto _test_eof20 }
	fallthrough
case 20:
// line 365 "prefixes.go"
	if data[p] == 105 { goto st21 }
	if data[p] <= 127 { goto st2 }
	goto st0
st21:
	p++
	if p == pe { goto _test_eof21 }
	fallthrough
case 21:
	if data[p] == 114 { goto st22 }
	if data[p] <= 127 { goto st2 }
	goto st0
st22:
	p++
	if p == pe { goto _test_eof22 }
	fallthrough
case 22:
	if data[p] == 99 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
tr25:
// line 30 "prefixes.rl"
	{ m = p }
	goto st23
st23:
	p++
	if p == pe { goto _test_eof23 }
	fallthrough
case 23:
// line 394 "prefixes.go"
	switch data[p] {
		case 105: goto st24
		case 111: goto st26
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st24:
	p++
	if p == pe { goto _test_eof24 }
	fallthrough
case 24:
	if data[p] == 107 { goto st25 }
	if data[p] <= 127 { goto st2 }
	goto st0
st25:
	p++
	if p == pe { goto _test_eof25 }
	fallthrough
case 25:
	if data[p] == 101 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
st26:
	p++
	if p == pe { goto _test_eof26 }
	fallthrough
case 26:
	if data[p] == 108 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
tr26:
// line 30 "prefixes.rl"
	{ m = p }
	goto st27
st27:
	p++
	if p == pe { goto _test_eof27 }
	fallthrough
case 27:
// line 434 "prefixes.go"
	switch data[p] {
		case 104: goto st28
		case 107: goto st10
		case 111: goto st29
		case 114: goto st10
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st28:
	p++
	if p == pe { goto _test_eof28 }
	fallthrough
case 28:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 104: goto st28
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st29:
	p++
	if p == pe { goto _test_eof29 }
	fallthrough
case 29:
	switch data[p] {
		case 104: goto st28
		case 107: goto st10
		case 111: goto st30
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st30:
	p++
	if p == pe { goto _test_eof30 }
	fallthrough
case 30:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 104: goto st28
		case 107: goto st10
		case 111: goto st30
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr27:
// line 30 "prefixes.rl"
	{ m = p }
	goto st31
st31:
	p++
	if p == pe { goto _test_eof31 }
	fallthrough
case 31:
// line 490 "prefixes.go"
	if data[p] == 101 { goto st25 }
	if data[p] <= 127 { goto st2 }
	goto st0
tr28:
// line 30 "prefixes.rl"
	{ m = p }
	goto st32
st32:
	p++
	if p == pe { goto _test_eof32 }
	fallthrough
case 32:
// line 503 "prefixes.go"
	switch data[p] {
		case 104: goto st33
		case 109: goto st19
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st33:
	p++
	if p == pe { goto _test_eof33 }
	fallthrough
case 33:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 104: goto st33
		case 109: goto st19
	}
	if data[p] <= 127 { goto st2 }
	goto st0
tr29:
// line 30 "prefixes.rl"
	{ m = p }
	goto st34
st34:
	p++
	if p == pe { goto _test_eof34 }
	fallthrough
case 34:
// line 532 "prefixes.go"
	switch data[p] {
		case 101: goto st35
		case 111: goto st38
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st35:
	p++
	if p == pe { goto _test_eof35 }
	fallthrough
case 35:
	switch data[p] {
		case 101: goto st35
		case 108: goto st36
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st36:
	p++
	if p == pe { goto _test_eof36 }
	fallthrough
case 36:
	if data[p] == 108 { goto st37 }
	if data[p] <= 127 { goto st2 }
	goto st0
st37:
	p++
	if p == pe { goto _test_eof37 }
	fallthrough
case 37:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 108: goto st37
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st38:
	p++
	if p == pe { goto _test_eof38 }
	fallthrough
case 38:
	if data[p] == 119 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
tr30:
// line 30 "prefixes.rl"
	{ m = p }
	goto st39
st39:
	p++
	if p == pe { goto _test_eof39 }
	fallthrough
case 39:
// line 587 "prefixes.go"
	switch data[p] {
		case 97: goto st40
		case 101: goto st41
		case 117: goto st43
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st40:
	p++
	if p == pe { goto _test_eof40 }
	fallthrough
case 40:
	switch data[p] {
		case 97: goto st40
		case 104: goto st28
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st41:
	p++
	if p == pe { goto _test_eof41 }
	fallthrough
case 41:
	switch data[p] {
		case 97: goto st42
		case 101: goto st41
		case 104: goto st28
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st42:
	p++
	if p == pe { goto _test_eof42 }
	fallthrough
case 42:
	switch data[p] {
		case 32: goto st4
		case 44: goto st7
		case 97: goto st42
		case 104: goto st28
	}
	if data[p] <= 127 { goto st2 }
	goto st0
st43:
	p++
	if p == pe { goto _test_eof43 }
	fallthrough
case 43:
	if data[p] == 112 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
st44:
	p++
	if p == pe { goto _test_eof44 }
	fallthrough
case 44:
	if data[p] == 116 { goto st45 }
	if data[p] <= 127 { goto st2 }
	goto st0
st45:
	p++
	if p == pe { goto _test_eof45 }
	fallthrough
case 45:
	if data[p] == 117 { goto st46 }
	if data[p] <= 127 { goto st2 }
	goto st0
st46:
	p++
	if p == pe { goto _test_eof46 }
	fallthrough
case 46:
	if data[p] == 97 { goto st47 }
	if data[p] <= 127 { goto st2 }
	goto st0
st47:
	p++
	if p == pe { goto _test_eof47 }
	fallthrough
case 47:
	if data[p] == 108 { goto st48 }
	if data[p] <= 127 { goto st2 }
	goto st0
st48:
	p++
	if p == pe { goto _test_eof48 }
	fallthrough
case 48:
	if data[p] == 108 { goto st49 }
	if data[p] <= 127 { goto st2 }
	goto st0
st49:
	p++
	if p == pe { goto _test_eof49 }
	fallthrough
case 49:
	if data[p] == 121 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
st50:
	p++
	if p == pe { goto _test_eof50 }
	fallthrough
case 50:
	if data[p] == 100 { goto st10 }
	if data[p] <= 127 { goto st2 }
	goto st0
	}
	_test_eof2: cs = 2; goto _test_eof; 
	_test_eof3: cs = 3; goto _test_eof; 
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
	_test_eof47: cs = 47; goto _test_eof; 
	_test_eof48: cs = 48; goto _test_eof; 
	_test_eof49: cs = 49; goto _test_eof; 
	_test_eof50: cs = 50; goto _test_eof; 

	_test_eof: {}
	if p == eof {
	switch cs {
	case 4, 7:
// line 30 "prefixes.rl"
	{ m = p }
	break
// line 753 "prefixes.go"
	}
	}

	_out: {}
	}

// line 73 "prefixes.rl"

	if cs < config_first_final {
		fmt.Printf("Parse error at %d\n", p)
		fmt.Printf("%s <- HERE -> %s", data[:p], data[p:])
	}

	if m > 0 {
		return s[m:]
	}
	return s
}






