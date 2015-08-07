# apachelogs
Go package for parsing Apache logs.

### Version 2.2.2

### Change Log

##### 2.2.2

_New pattern operators:_
* `OP_DIVIDE      ` Divide integer by defined integer
* `OP_MULTIPLY    ` Multiply integer by defined integer

##### 2.2.0

_Minor breaking API changes:_
* `NewPattern()` has new additional return for error handling.
* `PatternMatch()` has new additional return for error handling.

_New pattern operators:_
* `OP_REGEX_SUB   ` Regexp substitution: `{search}{replace}`
* `OP_ROUND_DOWN  ` Round integer down to nearest defined integer
* `OP_ROUND_UP    ` Round integer up to nearest defined integer

### Warning:
This package is still very much alpha. While I will try to impliment new features without disrupting existing features, I cannot make any guarantees that there wouldn't be any API changes in the future.
