grammar m4l;

program: header? constantDef* statement+ ('loop:' statement+)?;

header: mapEntry (mapEntry)*;

mapEntry: ID ':' (INT | ANY);

constantDef: CONSTANTID ':=' (instrumentDef | tablature);
instrumentDef: MAPCLASS '{' mapEntry+ '}';

tablature: (
		ID
		| NOTE
		| SILENCE
		| OCTAVE
		| INCOCT
		| DECOCT
		| CONSTANTID
		| tuplet
		| '|'
	)+;

tuplet: '(' (NOTE | OCTAVE | INCOCT | DECOCT)+ ')' INT;

statement: channelFill | sync;

channelFill: CHANNELID '<-' tablature;

sync: '-' '-' '-'+;

CHANNELID: '@' ID;
CONSTANTID: '$' ID;
MAPCLASS: '#' ID;

NOTE: [a-gA-G](FLAT | SHARP)? INT? DOTS?;
SILENCE: [rR]DIGIT?;
OCTAVE: [oO]DIGIT?;
INCOCT: '>';
DECOCT: '<';

FLAT: '-';
SHARP: ('#' | '+');
DOTS: '.'+;

fragment DIGIT: [0-9];
fragment ALPHA: [a-zA-Z_.];

INT: DIGIT+;
ID: ALPHA (ALPHA | DIGIT)*;
ANY: (ALPHA | DIGIT)+;

COMMENT: ';' (~[\r\n])* -> skip;

NL: ('\r' | '\n') -> skip;
WS: (' ' | '\t') -> skip;
