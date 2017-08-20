#!/bin/bash

echo "****************************************"
echo -e "****\e[31;1mThe Youth, Signin More Than Life\e[0m****"
echo "****************************************"
if [[ ! -e passwd.txt || ! -f passwd.txt ]]; then
	touch ./passwd.txt
fi
select value in "Signin" "Signup" "Cancel"; do
	if [[ $REPLY == 1 ]]; then #Signin
		while :; do
			read -p "Please type username> " username
			if [[ -z $username ]]; then
				echo -e "\e[31;1mUSERNAME can't be empty.\e[0m"
			else
				break
			fi
		done
		read -s -p "Please type password> " password
		u=`grep "$username" passwd.txt | cut -f 1` #可以匹配用户名或密码的部分
		p=`grep "$username" passwd.txt | cut -f 2`
		if [[ $u != $username || $p != $password ]]; then
			echo -e "\n\e[31;1mUsername or Password is incorrect.\e[0m"
			continue
		else
			clear
			. ./main.sh
		fi
	elif [[ $REPLY == 2 ]]; then #Signup
		clear
		. ./signup.sh
		clear
		echo "****************************************"
		echo -e "****\e[31;1mThe Youth, Signin More Than Life\e[0m****"
		echo "****************************************"
	elif [[ $REPLY == 3 ]]; then
		exit 0
	else
		echo -e "\e[31;1mInvalid choice\e[0m"
	fi
done
