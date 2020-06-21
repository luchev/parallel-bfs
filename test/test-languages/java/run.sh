#!/bin/bash

MAX_MEMORY=5G

java -Xmx${MAX_MEMORY} -classpath lib/commons-cli-1.4/commons-cli-1.4.jar:out/ Main $@
