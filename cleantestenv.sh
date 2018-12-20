#!/bin/bash

for i in `seq 1 3`
do
    rm -rf /tmp/$i
done

sudo ip addr del 10.29.1.1/16 dev ens33:1
sudo ip addr del 10.29.1.2/16 dev ens33:2
sudo ip addr del 10.29.1.3/16 dev ens33:3
sudo ip addr del 10.29.2.1/16 dev ens33:4
sudo ip addr del 10.29.2.2/16 dev ens33:5
