#!/usr/bin/env bash

ProgramName=test
ProgramVersion=1.0.0
ProgramBranch=`git rev-parse --abbrev-ref HEAD`
ProgramRevision=`git rev-parse HEAD`
CompilerVersion="`go version`"
BuildTime=`date -u '+%Y-%m-%d %H:%M:%S'`
Author=`whoami`@`hostname`

go build -ldflags "-X 'github.com/iTrellis/common/builder.ProgramName=$ProgramName' \
-X 'github.com/iTrellis/common/builder.ProgramVersion=$ProgramVersion' \
-X 'github.com/iTrellis/common/builder.ProgramBranch=$ProgramBranch' \
-X 'github.com/iTrellis/common/builder.ProgramRevision=$ProgramRevision' \
-X 'github.com/iTrellis/common/builder.CompilerVersion=${CompilerVersion}' \
-X 'github.com/iTrellis/common/builder.BuildTime=$BuildTime' \
-X 'github.com/iTrellis/common/builder.Author=$Author' \
" -o ${ProgramName} main.go

./${ProgramName}

rm ./${ProgramName}