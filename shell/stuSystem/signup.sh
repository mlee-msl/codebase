#!/bin/bash

##########################
##########Signup##########
##########################

echo "*****************************************"
echo -e "******\e[32;1mSuccess, To Start With Signup\e[0m******"
echo "*****************************************"

select choice in 'Signup' 'Cancel'; do
	if [[ $REPLY == 1 ]]; then
		while true; do
			read -p "Please type username> " username
			if [[ -z $username ]]; then
				echo -e "\e[31;1mUSERNAME can't be empty.\e[0m"
			else
				break
			fi
		done
		res=`grep "$username" passwd.txt | cut -f 1`
		if [[ $res == $username ]]; then #A used USERNAME
			echo -e "\e[31;1mThe username($username) has used.\e[0m"
			continue
		else
			while :; do
				read -s -p "Please type password> " password
				echo
				if [[ -z $password ]]; then
					echo -e "\e[31;1mPASSWORD can't be empty.\e[0m"
					continue
				fi
				read -s -p "Please type it again> " PASSWORD
				echo
				if [[ $password != $PASSWORD ]]; then
					echo -e "\e[31;1mBoth the type is inconsistent.\e[0m"
					continue
				else
					echo -e "$username\t$password" >> passwd.txt
					echo -e "\e[32;1mSignup Successfully.\e[0m"
					break
				fi
			done
		fi
	elif [[ $REPLY == 2 ]]; then
		read -n 1 -p "Press any key to continue..."
		break
	else
		echo -e "\e[31;1mInvalid choice\e[0m"
	fi
done
