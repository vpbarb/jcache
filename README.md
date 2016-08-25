# jcache key-value store
**jcache** is a key-value cache store. It provides own communication protocol over TCP connection. Server supports different storage types, including simple in-memory storage and persistent storage based on key-value file-based store [Bolt](https://github.com/boltdb/bolt).

## Protocol
### Terms: key, value and TTL
**Key** may contain only alphabetical symbols (a-z) in any case and numbers.

Supported **value** types:
- string
- hash (key-value subset)
- list

Hash field key limitation is similar to key limitation.

Any key may have **TTL** specified by seconds. When after TTL key will be expired and will be removed from storage. TTL equal to 0 forces to keep key forever.

###Commands
**jcache** protocol provides simple human-readable commands and responses format. 

Following list of command contains each command description and command/responses format.

Some of command descriptions provide examples in format:

    --> example of data sent to server
    <-- example of data sent to client

####KEYS
Command returns all keys from store. If storage is empty, then the command returns `COUNT 0`.

    KEYS\r\n
    COUNT <number_of_keys>\r\n[KEY <key>\r\n...]

Example:

    --> KEYS\r\n
    <-- COUNT 3\r\nKEY some_key1\r\nKEY some_key2\r\nKEY some_key3\r\n

####GET
Command returns string value by key. It works only for string value type. Command responses `VALUE N` where N is a length of following value. It returns error if key doesn't exist.

	GET <key>\r\n
	VALUE <value_length>\r\n<value>\r\n

Example:

    --> GET some_key\r\n
    <-- VALUE 10\r\nsome_value\r\n

####SET
Command sets new key-value pair. It works only for string value type. It returns error if key already exists.

    SET <key> <ttl> <value length>\r\n<value>\r\n
	OK\r\n

Example:

    --> SET some_key 60 10\r\nsome_value\r\n
    <-- OK\r\n

####UPD
Command updates existing key string value. It works only for string value type. It returns error if key doesn't exist.

    UPD <key> <value_length>\r\n<value>\r\n
    OK\r\n

####DEL
Command deletes key value. It works for **all** value types. It returns error if key doesn't exist.

    DEL <key>\r\n
    OK\r\n

####HCREATE
Command creates new hash. It returns error if key already exists.

	HCREATE <key> <ttl>\r\n
	OK\r\n

####HGET
Command returns hash field value. It returns error if key or field doesn't exist or key type is not hash.

	HGET <key> <field>\r\n
	VALUE <value_length>\r\n<value>\r\n

####HSET
Command sets hash field value. It is allowed to change existing field. It returns error if key already exists.

	HSET <key> <field> <value_length>\r\n<value>\r\n
	OK\r\n

####HDEL
Command deletes hash field. It returns error if key or field doesn't exist.

	HDEL <key> <field>\r\n
	OK\r\n

####HKEYS
Command returns list of all fields in hash. It returns error if key doesn't exist.

	HKEYS <key>\r\n
	COUNT <number_of_fields>\r\n[KEY <field>\r\n...]

####HLEN
Command returns number of fields in hash. It returns error if key doesn't exist.

	HLEN <key>\r\n
	LEN <number_of_fields>\r\n

####LCREATE
Command creates new list. It returns error if key already exists.

	HLIST <key> <ttl>\r\n
	OK\r\n

####LRPUSH
Command adds new value to the ending (right) of the list. It returns error if key doesn't exist.

	LRPUSH <key> <value_length>\r\n<value>\r\n
	OK\r\n

####LLPUSH
Command adds new value to the beginning (left) of the list. It returns error if key doesn't exist.

	LLPUSH <key> <value_length>\r\n<value>\r\n
	OK\r\n

####LRPOP
Command returns and removes value from the ending of the list. It returns error if key doesn't exist or if the list is empty.

	LRPOP <key>\r\n
	VALUE <value_length>\r\n<value>\r\n

####LLPOP
Command returns and removes value from the beginning of the list. It returns error if key doesn't exist or if the list is empty.

	LLPOP <key>\r\n
	VALUE <value_length>\r\n<value>\r\n

####LLEN
Command returns number of values in the list. It returns error if key doesn't exist.

	LLEN <key>\r\n
	LEN <number_of_values>\r\n

####LRANGE
Command returns sublist of values from and to specified indexes. If specified indexes are out of range it doesn't cause an error. It returns error if key doesn't exist.

	LRANGE <key> <start> <stop>\r\n
	COUNT <number_of_values>\r\n[VALUE <value_length>\r\n<value>\r\n...]

Example:

	LRANGE some_list 0 2\r\n
	COUNT 3\r\nVALUE 10\r\nsome_value\r\nVALUE 13\r\nanother_value\r\nVALUE 0\r\n\r\n

####AUTH
Command authenticate user within the opened connection. If server is started with authentication support, then AUTH command must be first after connection open. If authentication is not passed, then all commands will return error.

	AUTH <user> <password>\r\n
	OK\r\n

###Errors
Server may return protocol-related errors if it could not parse incoming request. Also, most of commands may return command-related error as a response instead of normal response. In both cases error response format will be:

	ERROR <description>\r\n

Examples:

    --> NONEXISTINGCOMMAND\r\n
    <-- ERROR Unknown command\r\n
    
    --> GET\r\n
    <-- ERROR Invalid command format\r\n

All commands related to specific value type return error if client tries to work with key of another type (except of DEL command which is universal).

## Server
###Storage types
TODO
###How to run
TODO
## Client
TODO
