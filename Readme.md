# GoRedis

GoRedis is supposed to be a project made up of two parts:

- a high performance, stable, and light-weight Redis server
- a fast, high-performance, and easy-to-use Go client library

## Quick Start

```bash
	________      __________           .___.__        
 /  _____/  ____\______   \ ____   __| _/|__| ______
/   \  ___ /  _ \|       _// __ \ / __ | |  |/  ___/
\    \_\  (  <_> )    |   \  ___// /_/ | |  |\___ \ 
 \______  /\____/|____|_  /\___  >____ | |__/____  >
				\/              \/     \/     \/         \/ 
```

## Architecture

```bash
.
├─aof
├─cluster
├─config
├─database
├─data_struct
│  ├─bitmap
│  ├─dict
│  ├─list
│  ├─lock
│  ├─set
│  └─sorted_set
├─interface
│  ├─database
│  ├─redis
│  └─tcp
├─lib
│  ├─log
│  ├─sync
│  └─util
├─pubsub
├─redis
│  ├─client
│  ├─conn
│  ├─parser
│  ├─protocol
│  └─server
└─tcp
```

## Redis protocol part

The Redis protocol is a binary-safe protocol that is used to communicate with Redis, which is called "RESP", short for "REdis Serialization Protocol".

RESP works on TCP connections.

In RESP, messages use CRLF (\r\n) as a line separator.

Binary-safe means that the protocol is not text-based, you can put any characters in your messages without causing problems.

There are five types of messages in RESP:

- Simple Strings (S) Not binary safe, CRLF not allowed.
- Error (E) Simple error message. Not binary safe. CRLF not allowed.
- Integers (I) uint64, as return values from commands like llen, scard, etc.
- Bulk Strings (B) Binary safe string. Used as return values from commands like get, set, etc.
- Arrays (A) (Also called Multi-Bulk-Strings) A list of bulk strings.

The RESP protocol uses the first character to indicate the type of the message.

- "+" for Simple Strings
- "-" for Error
- ":" for Integers
- "$" for Bulk Strings
- "*" for Arrays

### Bulk Strings

Bulk string is combination of `$` + length and a string as the value.

For example, the following RESP message:

```redis
$3\r\nSET\r\n
```

is a bulk string with length 3 and value "SET".

### nil

`$-1` is used to indicate that the value of the key is nil.

### Arrays

`*` is used to indicate that the value is an array, containing the number of elements in the array. The elements are separated by CRLF.

Here is an example of an array: `["foo", "bar", "baz"]`

```redis
*3\r\n$3\r\nfoo\r\n$3\r\nbar\r\n$3\r\nbaz\r\n
```

And also, the client can send an array of strings to the server:

```redis
*3\r\n$3SET\r\n$3key\r\n$5value\r\n
```

### Unmarshal

Notice that the RESP protocol is binary-safe, it allows you to send data like this: `*3\r\n$3\r\nSET\r\n$4\r\na\vr\nb\r\n$5hello\r\n`.

So you cannot simply use `ReadBytes('\n')` to read the next line, you shall use `io.Readfull(reader, msg)` to read  a given number of bytes.

```golang
msg := make([]byte, 6) // abcd/r/n
_, err := io.ReadFull(reader, msg)
```

So we can build a parser like this:

<!-- markdownlint-disable MD010 -->

```golang
type Payload struct{
	Data redis.Reply
	Err error
}

func ParseStream ( reader io.Reader) <-chan *Payload {
	ch := make(chan Payload)
	go func() {
		readingMultiLine := false
		expectedArgsCount := 0
		var args [][]byte
		var bulkLen int64
		for {
			line, err = readLine(reader, bulkLen)
			if err != nil {
				ch <- Payload{Err: err} // deal with error
				break
			}
			// We classify the line into two categories:
			// - Single line: Status reply, integer reply, or error reply.
			// - Multi-line: Bulk string reply and array reply.
			if !readingMultiLine {
				if isMultiBulkHeader(line) {
					// Got the first line of a multi-line reply.
					// fetch the expected number of arguments.
					expectedArgsCount = parseMultiBulkHeader(line)
					// wait for the rest of the multi-line reply.
					readingMultiLine = true
				} else if isBulkStringHeader(line) {
					// Got the first line of a bulk string reply.
					// fetch the length of the second line and save it in bulkLen.
					bulkLen = parseBulkStringHeader(line)
					// 1 line in the reply.
					expectedArgsCount = 1
					readingMultiLine = true
				} else {
					// Got a single line reply.
					reply, err := parseReply(line)
					emit(ch, reply, err)
				}
			} else {
				// We are awaiting the rest of the multi bulk reply or bulk string reply.
				// There are two cases when we get a multi bulk reply:
				// A BulkHeader
				if isBulkHeader(Line){
					bulkLen = parseBulkStringHeader(line)
				} else {
					// or a BulkString.
					args = append(args, line)
				}
			}
			// If we have read all the arguments, we can emit the reply.
			if len(args) == expectedArgsCount {
				reply, err := parseReply(args)
				emit(ch, reply, err)
				args = nil
				readingMultiLine = false
				expectedArgsCount = 0
				bulkLen = 0
			}
		}
	}()
	return ch
}
```
## Thread-safe HashMap

GoRedis works in a multi-threaded way, so it is important to make sure that the hashmap is thread-safe.

A common way to make a hashmap thread-safe is to use `sync.Map`, but `m.dirty` will duplicate the `m.read` into 

## Licensing and Disclaimers

This project, except the 3-rd party libraries, is written by mizumoto-cn. No GPL-like licensed codes are used.

This project is licensed under the [Mizumoto General Public License version 1.2](https://github.com/mizumoto-cn/TRPcG/blob/master/License/Mizumoto%20General%20Public%20License%20v1.2.md), which is a Mozilla Public License version 2.0 based license but with following restrictions:

By using any part of this project, you are deemed to have fully understanding and acceptance of the following terms：

1. You must conspicuously display, without modification, this License and the notice on each redistributed or derivative copy of the License Covered Work.
2. Any non-independent developers companies/groups/legal entities or other organizations should ensure that employees are not oppressed or exploited, and that employees can always receive a reasonable salary for their legal working hours.
3. Any independent or non-independent developers/companies/groups/legal entities or other organizations, shall ensure that it has a clear conscience, including and not limited to opposition to any form of Nazi or Neo-Nazism organization(s).

Otherwise these Individuals / Companies / Groups / Legal-entities will not have the right to copy / modify / redistribute any code / file / algorithm governed by MGPL v1.2.
