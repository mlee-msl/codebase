#!/bin/bash

read -p 'please input a number: ' num
if [ $num -lt 60 ]
then
	echo 'bad..._~_~_'
elif [ $num -ge 60 -a $num -lt 80 ]; then
	echo 'good'
else
	echo 'better'
fi
