#!/bin/bash

if [ -z $GOPATH ]; then
    echo "FAIL: GOPATH environment variable is not set"
    exit 1
fi

go install github.com/mburman/hooli/prunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

go install github.com/mburman/hooli/arunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

go install github.com/mburman/hooli/tests
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Pick random ports between [10000, 20000).
PRUNNER_PORT1=$(((RANDOM % 10000) + 10000))
PRUNNER_PORT2=$(((RANDOM % 10000) + 10000))
PRUNNER_PORT3=$(((RANDOM % 10000) + 10000))

ARUNNER_PORT1=$(((RANDOM % 10000) + 10000))
ARUNNER_PORT2=$(((RANDOM % 10000) + 10000))
ARUNNER_PORT3=$(((RANDOM % 10000) + 10000))

ARUNNER=$GOPATH/bin/arunner
PRUNNER=$GOPATH/bin/prunner
TEST=$GOPATH/bin/tests
OUT_FILE1=./out1
OUT_FILE2=./out2
OUT_FILE3=./out3
OUT_FILE4=./out4
OUT_FILE5=./out5
OUT_FILE6=./out6
##################################################

# Start prunner 1
${ARUNNER} -aport=${ARUNNER_PORT1} &> ${OUT_FILE1} &
ARUNNER_PID1=$!
sleep 1

# Start prunner 2
${ARUNNER} -aport=${ARUNNER_PORT2} &> ${OUT_FILE2} &
ARUNNER_PID2=$!
sleep 1

# Start prunner 3
${ARUNNER} -aport=${ARUNNER_PORT3} &> ${OUT_FILE3} &
ARUNNER_PID3=$!
sleep 1

# Start prunner 1
${PRUNNER} -pport=${PRUNNER_PORT1} -ports=${ARUNNER_PORT1},${ARUNNER_PORT2},${ARUNNER_PORT3} &> ${OUT_FILE4} &
PRUNNER_PID1=$!
sleep 1

# Start prunner 2
${PRUNNER} -pport=${PRUNNER_PORT2} -ports=${ARUNNER_PORT2},${ARUNNER_PORT1},${ARUNNER_PORT3} &> ${OUT_FILE5} &
PRUNNER_PID2=$!
sleep 1

# Start prunner 3
${PRUNNER} -pport=${PRUNNER_PORT3} -ports=${ARUNNER_PORT3},${ARUNNER_PORT2},${ARUNNER_PORT1} &> ${OUT_FILE6} &
PRUNNER_PID3=$!
sleep 1

# Start test.
${TEST} -ports=${PRUNNER_PORT1},${PRUNNER_PORT2},${PRUNNER_PORT3}

# Kill storage server.
kill -9 ${PRUNNER_PID1}
kill -9 ${PRUNNER_PID2}
kill -9 ${PRUNNER_PID3}
kill -9 ${ARUNNER_PID1}
kill -9 ${ARUNNER_PID2}
kill -9 ${ARUNNER_PID3}
#wait ${STORAGE_SERVER_PID} 2> /dev/null

##################################################
