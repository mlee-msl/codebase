#!/bin/bash

echo "****************************************"
echo -e "****\e[31;1mThe Youth, Signin more than life\e[0m****"
echo "****************************************"
if [[ ! -e passwd.txt || ! -f passwd.txt ]]; then
	touch ./passwd.txt
fi
select value in "Signin" "Signup"; do
	if [[ $REPLY == 1 ]]; then #Signin
		read -p "Please type username> " username
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
		echo -e "****\e[31;1mThe Youth, Signin more than life\e[0m****"
		echo "****************************************"
	else
		echo -e "\e[31;1mInvalid choice\e[0m"
	fi
done
