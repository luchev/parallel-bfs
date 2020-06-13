#!/bin/bash

mkdir out 2>/dev/null
rm out/*.class

javac -classpath lib/commons-cli-1.4/commons-cli-1.4.jar:out/ src/* -d out
