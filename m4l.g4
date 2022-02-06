grammar m4l;

program: header? constantDef* statement+ ('loop:' statement+)?;

header: mapEntry (mapEntry)*;

mapEntry: ID ':' (NUM | ANY);

constantDef: CONSTANTID ':=' (instrumentDef | tablature);
instrumentDef: ID '{' mapEntry+ '}';


tablature: (ID | NOTE | SILENCE | OCTAVE | INCOCT | DECOCT | CONSTANTID | tuplet | '|')+;

tuplet : '(' (NOTE|OCTAVE|INCOCT|DECOCT) + ')' NUM;

statement : channelFill | sync;

channelFill : CHANNELID '<-' tablature;

sync: '-''-''-'+;

CHANNELID: '@' ID;
CONSTANTID: '$' ID;

NOTE: [a-gA-G]DIGIT?; // TODO: dots and bemol
SILENCE: [rR]DIGIT?;
OCTAVE: [oO]DIGIT?;
INCOCT: '>';
DECOCT: '<';

fragment DIGIT : [0-9] ;
fragment ALPHA : [a-zA-Z_.];

ID: ALPHA (ALPHA | DIGIT)*;
ANY: (ALPHA | DIGIT)+; 
NUM: DIGIT+ ('.' DIGIT+)?;


COMMENT: ';' (~[\r\n])* -> skip;

NL: ('\r' | '\n') -> skip;
WS: (' ' | '\t') -> skip;
