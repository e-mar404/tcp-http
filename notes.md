# TCP -> HTTP

Will primarily use RFC9110 & RFC9112 specs. Even though this is just a guide for
implementing HTTP/1.1 it is ok since the underlying semantics are the same and
what varies is not fundamentally different.

Quick notes:

### UDP

- unordered
- "unreliable"
- not persistent connection
- really fast
- typical fail rate is 1%

### TCP

- ordered
- persistent connection
- not too fast

## Requests

At the center of a request is the http message ([RFC 9112
Sec2](https://datatracker.ietf.org/doc/html/rfc9112#name-message-format)).

A message is described like so:

```
start-line CRLF
*( field-line CRLF )
CRLF
[ message body ]
```

Note: CRLF -> Carriage Return Line Feed -> \r\n

Example of an http message:

```
Post /user/emar HTTP/1.1\r\n
Host: google.com
\r\n
{"name": "emar"}
```


