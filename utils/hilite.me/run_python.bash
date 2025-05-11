#!/bin/sh

cd /home/ubuntu/Coalg/utils/hilite.me
python3 -c "from tools import *; print(hilite_me(\"$1\"))"
