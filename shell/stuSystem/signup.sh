#!/bin/bash

echo "*****************************************"
echo -e "******\e[32;1mSuccess, To Start With Signup\e[0m******"
echo "*****************************************"

while :; do
	read -p "Please type username> " username
	res=`grep "$username" passwd.txt | cut -f 1`
	if [[ $res == $username ]]; then #A used USERNAME
		echo -e "\e[31;1mThe username($username) has used.\e[0m"
		continue
	else
		while :; do
			read -s -p "Please type password> " password
			echo
			read -s -p "Please type it again> " PASSWORD
			echo
			if [[ $password != $PASSWORD ]]; then
				echo -e "\e[31;1mBoth the type is inconsistent.\e[0m"
				continue
			else
				echo -e "$username\t$password" >> passwd.txt
				echo -e "\e[32;1mSignup Successfully."
				read -n 1 -p "Press any key to continue..."
				break
			fi
		done
	fi
	break
done
