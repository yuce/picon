@echo off

REM default target is build
if "%1" == "" (
    goto :build
)

2>NUL call :%1
if errorlevel 1 (
    echo Unknown target: %1
)

goto :end

:build
    go get github.com/yuce/picon/cmd/picon
	go build github.com/yuce/picon/cmd/picon
    goto :end

:end
