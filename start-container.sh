#!/usr/bin/env bash

alias go="/usr/lib/go-1.13/bin/go"

# Welcome script
WELCOME="
/ ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ \\
|  /~~\                                                                                       /~~\  |
|\ \   |                                                                                     |   / /|
| \   /|                                                                                     |\   / |
|  ~~  |  Here, we create a linux container by moving bytes!                                 |  ~~  |
|      |      42 65 74 74 65 72 20 6b 6e 6f 77 20 79 6f 75 72 20 68 65 78 65 73 20 3b 2d 29  |      |
|      |                                                                                     |      |
|      |                                                                                     |      |
|      |                                                                                     |      |
\     |~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~|     /
\   /                                                                                       \   /
~~~                                                                                         ~~~
"

echo -e "\033[1;36m${WELCOME}\033[0m"
go build cruntime.go
./cruntime run bash