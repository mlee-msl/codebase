#!/bin/bash

while true # while :
do
	read -p 'Username> ' username
	read -s -p 'Password> ' password
	if [ $username == 'MLee' -a $password == 'MLee' ]
	then
		echo -e '\n***************' # echo * --> show content of the current directory
		echo 'Welcome -^_^-'
		echo '***************'
		break
	else
		echo -e "\n***************"
		echo 'incorrect(username Or password) _~_~_'
		echo '***************'
		read -p 'Press any key to continue.'
		clear
	fi
done

echo -n -e '\e[32;1m'
echo '[1]-->shanghai'
echo '[2]-->chengdu'
echo '[3]-->wuhan'
echo '[4]-->beijing'
echo -n -e '\e[31;1m'
echo '[0]-->exit'
echo -n -e '\e[0m'

while true
do
	read -t 10 -p 'Please type place[1/2/3/4/0]: ' place

	if [ $place ] # judge place whether or not null
	then
		case $place in
			1)
				echo -e '\e[32;1mshanghai\e[0m'
				;;
			2)
				echo -e '\e[32;1mchengdu\e[0m'
				;;
			3)
				echo -e '\e[32;1mwuhan\e[0m'
				;;
			4)
				echo -e '\e[32;1mbeijing\e[0m'
				;;
			0)
				echo -e '\e[31;1mexit\e[0m'
				break
				;;
			*)
				echo -e '\e[31;1merror\e[0m'
				;;
		esac
	else
		echo -e "\e[31;1mYour input timed out.\e[0m"
		break
	fi
done





################################################
# touch file{1..10}
# touch file{1,2,3,4} # can't space character between comma and digit(e.g. 1,2) or else wrongly think that will create several files(e.g. touch file{1, 2} --> file1: file{1,file2: 2})
