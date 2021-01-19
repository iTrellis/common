#!/usr/bin/env bash

ProgramName=test
ProgramVersion=1.0.0
CompilerVersion="`go version`"
BuildTime=`date -u '+%Y-%m-%d %H:%M:%S'`
Author=`whoami`

go build -ldflags "-X 'github.com/iTrellis/go-common/builder.ProgramName=$ProgramName' \
-X 'github.com/iTrellis/go-common/builder.ProgramVersion=$ProgramVersion' \
-X 'github.com/iTrellis/go-common/builder.CompilerVersion=${CompilerVersion}' \
-X 'github.com/iTrellis/go-common/builder.BuildTime=$BuildTime' \
-X 'github.com/iTrellis/go-common/builder.Author=$Author' \
" -o ${ProgramName} main.go

./${ProgramName}

rm ./${ProgramName}