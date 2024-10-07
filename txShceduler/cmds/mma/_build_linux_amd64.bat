echo

REM SYMBOL 바뀌는 경우 ROOT_DIR_HOST 와 SUB_DIR_HOST 도 수정해야함
set SUB_DIR_HOST=cmds\mma
set SUB_DIR_CONTAINER=%SUB_DIR_HOST:\=/%
set ROOT_DIR_HOST=%cd:\cmds\mma=%
set ROOT_DIR_CONTAINER=/usr/src/myapp

set BUILD_NAME=txmMma.app

set GOOS=linux
set GOARCH=amd64
set GOVER=1.18

TITLE TxScheduler_MMA

echo TXSCHEDULER_MMA compile try ...
docker run --rm -v %ROOT_DIR_HOST%:%ROOT_DIR_CONTAINER% -w %ROOT_DIR_CONTAINER% -e GOOS=%GOOS% -e GOARCH=%GOARCH% golang:%GOVER% go build -o %ROOT_DIR_CONTAINER%/%SUB_DIR_CONTAINER%/%BUILD_NAME% %ROOT_DIR_CONTAINER%/%SUB_DIR_CONTAINER%/main.go

pause
echo compile process is end ...




