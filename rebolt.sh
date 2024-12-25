#!/bin/bash
cd 
for d in $HOME/live/*/ ; do
    cd $d
    go build -o $HOME/bin/$(basename ${d})
    $(basename ${d}) &
done
