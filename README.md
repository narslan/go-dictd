A golang implementation of dictd server. İt is at pre-alpha stage.

#### Usage:

Install
```
go install github.com/narslan/go-dictd


```
### Usage
Get it:
```
go install github.com
```
Start:
```
$ dictdserver
```
Connect via the dict client.
```
$ dict -D # list DBs
$ dict word
```

There is lots of room for improvements.: 
[-] Configuration parser should be removed, as it parse nothing. A small hook to lexer would be fine
[-] Speed improvements: 
