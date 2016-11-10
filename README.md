# apachelogs
Go package for parsing Apache logs.

## Change Log

### Version 3.0.0

_STABLE!_

This is now a stable release. Unfortunately it's also quite a rewrite to made the code more idiomatic Go. A lot of that is in the form of variables and constants being renamed - which will break compatibility. I have also changed the naming of the AccessLog / AccessLogs (for slices) to be more meaningful: `AccessLine` is a single entity, `AccessLog` is a slice of AccessLine's (ie `[]AccessLine`).

_New in version 3:_

Support for error logs. Currently untested but the API follows the same format as for access log parsing so I cannot envisage any breaking changes as that code matures.  

## Previous versions:

### Version 2.2.2

_New pattern operators:_
* `OP_DIVIDE      ` Divide integer by defined integer
* `OP_MULTIPLY    ` Multiply integer by defined integer

### Version 2.2.0

_Minor breaking API changes:_
* `NewPattern()` has new additional return for error handling.
* `PatternMatch()` has new additional return for error handling.

_New pattern operators:_
* `OP_REGEX_SUB   ` Regexp substitution: `{search}{replace}`
* `OP_ROUND_DOWN  ` Round integer down to nearest defined integer
* `OP_ROUND_UP    ` Round integer up to nearest defined integer
