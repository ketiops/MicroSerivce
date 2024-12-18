#!/bin/bash
# PodIP 입력
#host_ip=10.0.2.134

read -p 'Please enter stress module ip address: ' host_ip

read -p 'Please enter the 1st arguments(duration): ' duration

read -p 'Please enter the 2st arguments(mem_amount): ' mem_amount

curl "http://$host_ip:5000/memory_stress?duration=$duration&mem_amount=$mem_amount" &
