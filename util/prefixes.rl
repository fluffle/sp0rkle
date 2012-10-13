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

%% machine config;
%% write data;

// basic parser bits
%%{
	action mark { m = p }

	ha = ( "h" | "a" );
	he = ( "h" | "e" );
	sp = ( "," " " | " " | "," );

	prefix = 
	      "o"* "k"
		| "s" "e" "e"
		| "u" ( "h"+ "m"* | "m"+ )
		| "h" "e" "y"
		| "a" "c" "t" "u" "a" "l" "l" "y"
		| "o" "o" "o"+
		| "w" "e"+ "l" "l"+
		| "i" "i" "r" "c"
		| "b" "u" "t"
		| "a" "n" "d"
		| "o" "r"
		| "e" "h"
		| "."+
		| "l" "i" "k" "e"
		| "o"+ "h"+
		| "y" ( "e"+ "a"+ "h"* | "e"+ "h"+ | "a"+ "h"+ )
		| "y" "u" "p"
		| "l" "o" "l"
		| "w" "o" "w"
		| "h"+ "m"+
		| "e"+ "r"+
		| ha{2,}
		| he{3,} 
	;
	prefixes = prefix sp %mark ;

	main := prefixes* ascii*;
}%%

func RemovePrefixes(s string) string {
	data := make([]byte, len(s))
	copy(data, s)
	cs, p, m, pe, eof := 0, 0, 0, len(data), len(data)
	
	%% write init;
	%% write exec;

	if cs < config_first_final {
		fmt.Printf("Parse error at %d\n", p)
		fmt.Printf("%s <- HERE -> %s", data[:p], data[p:])
	}

	if m > 0 {
		return s[m:]
	}
	return s
}






