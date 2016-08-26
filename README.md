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

Any key may have **TTL** specified by seconds. After TTL key will be expired and will be removed from storage by GC. TTL equal to 0 means unlimited TTL.

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

####EXPIRE
Command updates key ttl. It works for **all** value types. It returns error if key doesn't exist.

	EXPIRE <key> <ttl>\r\n
	OK\r\n

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
Command sets hash field value. It is allowed to change existing field. If hash doesn't exist yet, it will be created with ttl=0.

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
Command adds new value to the ending (right) of the list. If list doesn't exist yet, it will be created with ttl=0.

	LRPUSH <key> <value_length>\r\n<value>\r\n
	OK\r\n

####LLPUSH
Command adds new value to the beginning (left) of the list. If list doesn't exist yet, it will be created with ttl=0.

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

Example:

	--> HCREATE hash 0\r\n
	<-- OK\r\n
	--> GET hash
	<-- ERROR Key type is not string

## Server
###Storage types
There are 3 implemented types of storages: memory, multi_memory and bolt. All storages have "garbage collector" (GC) to remove expired values from storage. Interval of GC running is defined by `storage_gc_interval` option.

####Memory storage
Memory storage is a simple in-memory storage with limited count of stored keys. Maximum size is defined by `storage_memory_size` option. Memory storage uses LRU algorithm, so less recent key will be removed in case of adding new key to full storage. 

####Multi-memory storage
It's the same in-memory storage but separated on several buckets. Distribution by buckets is normal and made by key check sum. Number of buckets is defined by `storage_multi_memory_count` option.

####Bolt
This storage has underlying [Bolt](https://github.com/boltdb/bolt) file storage. Path to Bolt file is defined by `storage_boltdb_path` option. If file doesn't exist it will be created.

###Authentication
If you want server supports authentication, just pass path to .htpasswd file with `htpasswd` option. If server is running with `htpasswd` option then it requires `AUTH` command with valid credentials after connection is open. All other commands will work only after valid authentication.

###How to build

	git clone git@github.com:Barberrrry/jcache.git ./
	make

It will install vendor dependencies and build `jcache` file.

###Benchmarks
Run benchmarks to see some storages and server performance:

	make bench

###How to run
Just run server with default parameters:

	./jcache
	
Run options:

	./jcache --help
	Usage of ./jcache:
        -htpasswd string
            Path to .htpasswd file for authentication. Leave blank to disable authentication.
        -listen string
            Host and port to listen connection (default ":9999")
        -storage_bolt_path string
            Path to Bolt file
        -storage_gc_interval duration
            Storage GC interval (default 1m0s)
        -storage_memory_size uint
            Max number of stored elements (default 10000)
        -storage_multi_memory_count uint
            Number of storages inside multi memory storage (default 1)
        -storage_type value
            Type of storage (memory, multi_memory, bolt) (default memory)

Example:

	./jcache -listen=127.0.0.1:9999 -storage_type=bolt -storage_boltdb_path=bold.db -storage_gc_interval=5m

## Client
Import client package:
		
	import "github.com/Barberrrry/jcache/client"

Client package support [glide](github.com/Masterminds/glide), so you can just run `glide up` if you use glide in your project. Also you may download client dependencies manually:

	go get "gopkg.in/fatih/pool.v2"

Example of client usage:

	client, clientErr := client.New("127.0.0.1:9999", "admin", "admin", 5*time.Second, 5)
	setErr := client.Set("key", "value1", 3600)
