#!/bin/sh

cd /home/matei/coalg/utils/hilite.me
python -c "from tools import *; print(hilite_me(\"$1\"))"
