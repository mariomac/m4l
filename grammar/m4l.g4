grammar m4l;

program: header? constantDef* statement+ ('loop:' statement+)?;

header: MAPENTRY (MAPENTRY)*;
MAPENTRY: ALPHANUMERIC WS* ':' WS* ALPHANUMERIC;

constantDef: CONSTANTID ':=' ( tablature | instrumentDef);
instrumentDef: OPENINSTRUMENT MAPENTRY+ '}';
OPENINSTRUMENT: ALPHANUMERIC WS* '{';

tablature:
	(
		NOTE
		| SILENCE
		| OCTAVE
		| INCOCT
		| DECOCT
		| VOLUME
		| '|'
		| CONSTANTID
		| tuplet
		| NL
	)+;

tuplet: '(' (NOTE | OCTAVE | INCOCT | DECOCT)+ ')' INT;

statement: channelFill | sync;

channelFill: CHANNELID '<-' tablature;

sync: '-' '-' '-'+;

CHANNELID: '@' ALPHANUMERIC;
CONSTANTID: '$' ALPHANUMERIC;

NOTE: [a-gA-G](FLAT | SHARP)? INT? DOTS?;
SILENCE: [rR]INT? DOTS?;
OCTAVE: [oO]INT?;
INCOCT: '>';
DECOCT: '<';
VOLUME: [Vv]INT;

FLAT: '-';
SHARP: ('#' | '+');
DOTS: '.'+;

fragment DIGIT: [0-9];
fragment ALPHA: [a-zA-Z_.];
fragment ALPHANUMERIC: (DIGIT | ALPHA)+;

INT: DIGIT+;

COMMENT: ';' (~[\r\n])* -> skip;

NL: ('\r' | '\n') -> skip;
WS: (' ' | '\t') -> skip;
