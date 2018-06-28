# gensysops

Pre-built generic functions for exposure to [otto JS interpreter](https://godoc.org/github.com/robertkrimen/otto).

## Explanation

Otto is a package that allows you to run arbitrary JavaScript inside Go, and it supports exposing Go functions to the JavaScript interpreter. "gensysops" provides generic functions which you can expose to the otto runtime to perform common system operations such as file modification and even arbitrary system commands. 

## Operations

* File operations
    - Check if file exists
    - read in from file
    - write out to [new/existing] file
    - set file timestamps (Mod, Acc)
* System operations
    - arbitrary system commands

# TODO

* move/rename file
* remove file
* copy file
* create dir
* real error handling/exposure