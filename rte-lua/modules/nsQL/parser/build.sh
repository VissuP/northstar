#!/bin/bash
go tool yacc -o parser.go -p nsQL parser.y
sed -i '.bak' 's/nsQLErrorVerbose = false/nsQLErrorVerbose = true/g' parser.go
sed -i '.bak' 's/nsQLErrorVerbose = false/nsQLErrorVerbose = true/g' parser.go
rm parser.go.bak y.output
go install