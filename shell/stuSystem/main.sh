#!/bin/bash

##########################
##########main############
##########################

. lang.sh

print()
{
	if [ $1 -eq 1 ]; then
#		PRINT=(${print1[*]})
		for((i = 0; i < ${#print1[*]}; i++)); do
			PRINT[$i]=${print1[$i]}
		done
	else
		for((i = 0; i < ${#print2[*]}; i++))
		do
			PRINT[$i]=${print2[$i]}
		done
	fi

	while true
	do
		clear
		if test $1 -eq 1
		then
			echo -e "*************\e[31;1m${PRINT[0]}\e[0m**************"
		else
			echo -e "*************\e[31;1m${PRINT[0]}\e[0m***************"
		fi
		echo -e "***\e[31;1m1.\e[32;1m${PRINT[1]}\t\e[31;1m2.\e[32;1m${PRINT[2]}\e[0m\t***"
		echo -e "***\e[31;1m3.\e[32;1m${PRINT[3]}\t\e[31;1m4.\e[32;1m${PRINT[4]}\e[0m\t***"
		if [ $1 -eq 1 ]; then
			echo -e "***\e[31;1m5.\e[32;1m${PRINT[5]}\t\e[31;1m0.${PRINT[6]}\e[0m\t***"
		else
			echo -e "***\e[31;1m5.\e[32;1m${PRINT[5]}\t\e[31;1m0.${PRINT[6]}\e[0m\t\t\t***"
		fi
		if test $1 -eq 1
		then
			echo "*******************************************"
		else
			echo "***********************************************************"
		fi
		read -p "${PRINT[7]}" num
		case $num in
			1)
				if [ $1 == 1 ]; then
					insert 1
				else
					insert 2
				fi
				;;
			2)
				if [ $1 == 1 ]; then
					modify 1
				else
					modify 2
				fi
				;;
			3)
				if [ $1 == 1 ]; then
					delete 1
				else
					delete 2
				fi
				;;
			4)
				if [ $1 == 1 ]; then
					query 1
				else
					query 2
				fi
				;;
			5)
				if [ $1 == 1 ]; then
					order 1
				else
				 	order 2
				fi
				;;
			0)
				if [ $1 == 1 ]; then
					quit 1
				else	
					quit 2
				fi
				;;
			*)
				echo -e "\e[31;1m${PRINT[8]}\e[0m"
				read -N 1 -p "${PRINT[9]}"
		esac
	done
}

function start()
{
	clear
	if [[ ! -e stuInfo.txt || ! -f stuInfo.txt ]]; then
		touch ./stuInfo.txt
	fi
	if [[ ! -e stuRecord.txt || ! -f stuRecord.txt ]]; then
		touch ./stuRecord.txt
	fi
	echo -e "**********\e[32;1m请选择语言(Please select language)\e[0m**********"
	select lang in '中文(Chinese)' '英文(English)' '退出(Exit)'
	do
		if [[ $lang == '中文(Chinese)' ]] #if [ $REPLY -eq 1 ]
		then
			print 1
			break
		elif [[ $lang == '英文(English)' ]]
		then
			print 2
			break
		elif [[ $lang = '退出(Exit)' ]]; then
			. ./startup.sh
		else
			echo -e "\e[31;1m无效选择(Invalid choice)\e[0m"
		fi
	done
}

insert()
{
	if [[ $1 == 1 ]]; then
		for((i=0;i<${#insert1[*]};i++)); do
			INSERT[$i]=${insert1[$i]}
		done
	else
		for((i=0;i<${#insert2[*]};i++)); do
			INSERT[$i]=${insert2[$i]}
		done
	fi
	clear
	select choose in "${INSERT[0]}" "${INSERT[1]}"; do
		if [[ $choose == "${INSERT[2]}" ]]; then
			while :; do
				read -p "${INSERT[3]}> " sNo
				if [[ ${#sNo} -ne 4 ]]; then
					echo -e "\e[31;1m${INSERT[4]}\e[0m"
					break
				fi
				res=$(cat stuInfo.txt | cut -f 1 | grep $sNo)
				if [[ ${#res} -ne 4 ]]; then # 该学生不存在
					echo -e "**********\e[32;1m${INSERT[5]}($sNo)${INSERT[6]}\e[0m**********"
					echo -e "\e[31;1m${INSERT[7]}:(\e[32;1m${INSERT[8]} ${INSERT[9]} ${INSERT[10]}\e[31;1m)\e[33;1m${INSERT[11]}: \e[32;1m${INSERT[12]} ${INSERT[13]} 006\e[0m"
					while true; do
						read -p "> " -a stuInfo
						if [[ ${#stuInfo[*]} < 3 ]]; then
							echo -e "\e[31;1m${INSERT[14]}\e[0m"
						else
							break
						fi
					done
					read -p "${INSERT[15]}?(y/n)> " flag1
					if [[ $flag1 = y ]]; then
						while :; do
							echo -e "******************\e[32;1m${INSERT[16]}($sNo)${INSERT[17]}\e[0m******************"
							echo -e "\e[31;1m${INSERT[18]}:(\e[32;1m${INSERT[19]} ${INSERT[20]} ${INSERT[21]}\e[31;1m)[${INSERT[22]}]\e[33;1m${INSERT[23]}: \e[32;1m80 92 85\e[0m"
							while :; do
								read -p "> " -a stuRecord
								if [[ ${#stuRecord[*]} < 3 ]]; then
									echo -e "\e[31;1m${INSERT[24]}\e[0m"
									continue
								elif [[ ${stuRecord[0]} == *[!0-9]* || ${stuRecord[1]} == *[!0-9]* || ${stuRecord[2]} == *[!0-9]* ]]; then
									echo -e "\e[31;1m${INSERT[33]}\e[0m"
									continue
								elif [[ ${stuRecord[0]} -gt 100 || ${stuRecord[1]} -gt 100 || ${stuRecord[2]} -gt 100 ]]; then
									echo -e "\e[31;1m${INSERT[25]}\e[0m"
									continue
								elif [[ ${stuRecord[0]} -lt 0 || ${stuRecord[1]} -lt 0 || ${stuRecord[2]} -lt 0 ]]; then
									echo -e "\e[31;1m${INSERT[26]}\e[0m"
									continue
								else
									break
								fi
							done
							read -p "${INSERT[27]}?(y/n)> " flag2
							if [[ $flag2 == y ]]; then
#								stuRecord[${#stuRecord[*]}]=$(date +%F)
								let "avg=(stuRecord[0]+stuRecord[1]+stuRecord[2])/3"
								d=`date +%F`/`date +%X`
								echo -e "$sNo\t${stuInfo[0]}\t${stuInfo[1]}\t${stuInfo[2]}\t$d" >> ./stuInfo.txt
								echo -e "$sNo\t${stuInfo[0]}\t${stuRecord[0]}\t${stuRecord[1]}\t${stuRecord[2]}\t$avg\t$d" >> ./stuRecord.txt
								echo -e "\e[33;1m${INSERT[28]}\e[0m"
								break
							else
								continue
							fi
						done
						break
					else
						continue
					fi	
				else # 该学生已存在
					# find special record from the file
					sName=`cat stuInfo.txt | grep "$res" | cut -f 2`
					echo -e "\e[31;1m${INSERT[29]}(\e[32;1m$sNo<-->${sName}\e[31;1m)${INSERT[30]}\e[0m"
					break
				fi
			done
		elif [[ $choose == "${INSERT[31]}" ]]
		then
			break
		else
			echo -e "\e[31;1m${INSERT[32]}\e[0m"
		fi
	done
}

modify()
{
	if [[ $1 == 1 ]];then
		for((i=0;i<${#modify1[*]};i++)); do
			MODIFY[$i]=${modify1[$i]}
		done
	else
		for((i=0;i<${#modify2[*]};i++)); do
			MODIFY[$i]=${modify2[$i]}
		done
	fi
	clear
	select choose in "${MODIFY[0]}" "${MODIFY[1]}"; do
		if [ $REPLY == 1 ]; then
			lines=$(wc -l stuInfo.txt | awk '{print $1}')
			if [[ $lines -eq 0 ]]; then
				echo -e "\e[33;1m${MODIFY[20]}\e[0m"
				continue
			else
				sNos=$(awk '{print $1}' stuInfo.txt)
				echo -n -e "\e[32;1m${MODIFY[21]}: (\e[0m"
				echo -n $sNos
				echo -e "\e[32;1m)\e[0m"
			fi
			while :; do
				read -p "${MODIFY[2]}> " sNo
				if [[ ${#sNo} -ne 4 ]]; then
					echo -e "\e[31;1m${MODIFY[3]}\e[0m"
					continue
				else
					break
				fi
			done
			res=$(cat stuInfo.txt | cut -f 1 | grep "$sNo")
			if [[ ${#res} != 4 ]]; then
				echo -e "\e[31;1m${MODIFY[4]}(\e[32;1m$sNo\e[31;1m)${MODIFY[5]}\e[0m"
				continue
			else
				grep -v "$sNo" stuRecord.txt > tmp.txt # 删除原始指定记录
				mv tmp.txt stuRecord.txt
				while :; do
					sName=$(cat stuInfo.txt | grep "$sNo" | cut -f 2)
					echo -e "******************\e[32;1m${MODIFY[6]}(\e[31;1m$sNo<-->$sName\e[32;1m)${MODIFY[7]}\e[0m******************"
					echo -e "\e[31;1m${MODIFY[8]}:(\e[32;1m${MODIFY[9]} ${MODIFY[10]} ${MODIFY[11]}\e[31;1m)[${MODIFY[12]}]\e[33;1m${MODIFY[13]}: \e[32;1m90 81 85\e[0m"
					while :; do
						read -p "> " -a stuRecord
						if [[ ${#stuRecord[*]} < 3 ]]; then
							echo -e "\e[31;1m${MODIFY[14]}\e[0m"
							continue
						elif [[ ${stuRecord[0]} == *[!0-9]* || ${stuRecord[1]} == *[!0-9]* || ${stuRecord[2]} == *[!0-9]* ]]; then
							echo -e "\e[31;1m${MODIFY[22]}\e[0m"
							continue
						elif [[ ${stuRecord[0]} -gt 100 || ${stuRecord[1]} -gt 100 || ${stuRecord[2]} -gt 100 ]]; then
							echo -e "\e[31;1m${MODIFY[15]}\e[0m"
							continue
						elif [[ ${stuRecord[0]} -lt 0 || ${stuRecord[1]} -lt 0 || ${stuRecord[2]} -lt 0 ]]; then
							echo -e "\e[31;1m${MODIFY[16]}\e[0m"
							continue
						else
							break
						fi
					done
					read -p "${MODIFY[17]}?(y/n)> " flag1
					if [[ $flag1 == y ]]; then
#						stuRecord[${#stuRecord[*]}]=$(date +%F)
						sName=$(cat stuInfo.txt | grep "$sNo" | cut -f 2)
						let "avg=(stuRecord[0]+stuRecord[1]+stuRecord[2])/3"
						d=`date +%F`/`date +%X`
						echo -e "$sNo\t${sName}\t${stuRecord[0]}\t${stuRecord[1]}\t${stuRecord[2]}\t$avg\t$d" >> ./stuRecord.txt
						echo -e "\e[33;1m${MODIFY[18]}\e[0m"
						break
					else
						continue
					fi
				done
			fi
		elif [ $REPLY == 2 ]; then
			break
		else
			echo -e "\e[31;1m${MODIFY[19]}\e[0m"
		fi
	done
}

delete(){
	if [[ $1 == 1 ]];then
		for((i=0;i<${#delete1[*]};i++)); do
			DELETE[$i]=${delete1[$i]}
		done
	else
		for((i=0;i<${#delete2[*]};i++)); do
			DELETE[$i]=${delete2[$i]}
		done
	fi
	clear
	select choice in "${DELETE[0]}" "${DELETE[1]}"; do
		if [ $REPLY == 1 ]; then
			lines=$(wc -l stuInfo.txt | awk '{print $1}')
			if [[ $lines -eq 0 ]]; then
				echo -e "\e[33;1m${DELETE[13]}\e[0m"
				continue
			else
				sNos=$(awk '{print $1}' stuInfo.txt)
				echo -n -e "\e[32;1m${DELETE[9]}: (\e[0m"
				echo -n $sNos
				echo -e " All\e[32;1m)\e[0m-->\e[31;1m(${DELETE[10]})\e[0m"
			fi
			while :; do
				read -p "${DELETE[2]}> " sNo
				if [[ ${#sNo} -ne 4 ]]; then
					if [[ ${sNo} = "All" || $sNo = 'all' ]]; then
						break
					else
						echo -e "\e[31;1m${DELETE[3]}\e[0m"
						continue
					fi
				else
					break
				fi
			done
			res=$(cat stuInfo.txt | cut -f 1 | grep "$sNo")
			if [ ${#res} != 4 ]; then
				if [[  $sNo = 'All' || $sNo = 'all' ]]; then
					read -p "${DELETE[11]}?(y/n)> " f1
					if [[ ${f1} == y ]]; then
						read -p "${DELETE[12]}?(y/n)> " f2
						if [[ $f2 = y ]]; then
							cat /dev/null > stuInfo.txt
							cat /dev/null > stuRecord.txt
							echo -e "\e[33;1m${DELETE[7]}\e[0m"
						else
							continue
						fi
					else
						continue
					fi
				else
					echo -e "\e[31;1m${DELETE[4]}(\e[32;1m$sNo\e[31;1m)${DELETE[5]}\e[0m"
					continue
				fi
			else
				sName=$(cat stuInfo.txt | grep "$sNo" | cut -f 2)
				read -p "${DELETE[6]}($sNo<-->$sName)?(y/n)> " flag
				if [[ $flag == y ]]; then
					cat stuInfo.txt | grep -v "$sNo" > tmp1.txt
					cat stuRecord.txt | grep -v "$sNo" > tmp2.txt
					mv tmp1.txt stuInfo.txt
					mv tmp2.txt stuRecord.txt
					echo -e "\e[33;1m${DELETE[7]}\e[0m"
				else
					continue
				fi
			fi
		elif [ $REPLY == 2 ]; then
			break
		else
			echo -e "\e[31;1m${DELETE[8]}\e[0m"
		fi
	done
}

query()
{
	if [[ $1 == 1 ]];then
		for((i=0;i<${#query1[*]};i++)); do
			QUERY[$i]=${query1[$i]}
		done
	else
		for((i=0;i<${#query2[*]};i++)); do
			QUERY[$i]=${query2[$i]}
		done
	fi
	clear
	select choice in "${QUERY[0]}" "${QUERY[1]}"; do
		if [ $REPLY == 1 ]; then
			lines=$(wc -l stuInfo.txt | awk '{print $1}')
			if [[ $lines -eq 0 ]]; then
				echo -e "\e[33;1m${QUERY[22]}\e[0m"
				continue
			else
				sNos=$(awk '{print $1}' stuInfo.txt)
				echo -n -e "\e[32;1m${QUERY[23]}: (\e[0m"
				echo -n $sNos
				echo -e "\e[32;1m)\e[0m"
			fi
			while :; do
				read -p "${QUERY[2]}> " sNo
				if [[ ${#sNo} -ne 4 ]]; then
					echo -e "\e[31;1m${QUERY[3]}\e[0m"
					continue
				else
					break
				fi
			done
			res=`cat stuInfo.txt | cut -f 1 | grep "$sNo"`
			if [ ${#res} != 4 ]; then
				echo -e "\e[31;1m${QUERY[4]}(\e[32;1m$sNo\e[31;1m)${QUERY[5]}\e[0m"
				continue
			else
#				printf "%12s%12s%12s%12s%16s%16s%16s\n" 学号 姓名 数据库 C语言 数据结构 平均成绩 修改时间
				echo -e "\e[31;1m${QUERY[6]}:(\e[32;1m${QUERY[7]} ${QUERY[8]} ${QUERY[9]} ${QUERY[10]} ${QUERY[11]} ${QUERY[12]} ${QUERY[13]}\e[31;1m)\e[0m"
				for((i = 0; i < 7; i++)); do
					stuRecord[$i]=$(cat stuRecord.txt | grep "$sNo" | cut -f `expr $i + 1`)
				done
				printf "%8s%12s%8s%8s%8s%8s%22s\n" ${stuRecord[0]} ${stuRecord[1]} ${stuRecord[2]} ${stuRecord[3]} ${stuRecord[4]} ${stuRecord[5]} ${stuRecord[6]}
				read -p "${QUERY[14]}?(y/n)> " flag
				if [[ $flag = "y" ]]; then
					echo -e "\e[31;1m${QUERY[15]}:(\e[32;1m${QUERY[16]} ${QUERY[17]} ${QUERY[18]} ${QUERY[19]} ${QUERY[20]}\e[31;1m)\e[0m"
					for((i=0;i<5;i++)); do
						stuInfo[$i]=`cat stuInfo.txt | grep "$sNo" | cut -f $(expr $i + 1)`
					done
					printf "%8s%12s%8s%8s%22s\n" ${stuInfo[0]} ${stuInfo[1]} ${stuInfo[2]} ${stuInfo[3]} ${stuInfo[4]}
				else
					continue
				fi
			fi
		elif [ $REPLY == 2 ]; then
			break
		else
			echo -e "\e[31;1m${QUERY[21]}\e[0m"
		fi
	done
}

order()
{
	if [[ $1 == 1 ]];then
		for((i=0;i<${#order1[*]};i++)); do
			ORDER[$i]=${order1[$i]}
		done
	else
		for((i=0;i<${#order2[*]};i++)); do
			ORDER[$i]=${order2[$i]}
		done
	fi
	clear
	select choice in "${ORDER[0]}" "${ORDER[1]}"; do
		if [ $REPLY == 1 ]; then
			echo -e "\e[31;1m${ORDER[2]}(\e[32;1m1\e[31;1m.${ORDER[3]}\e[32;1m2\e[31;1m.${ORDER[4]}\e[32;1m3\e[31;1m.${ORDER[5]})\e[0m"
			while :; do
				read -p "${ORDER[6]}> " num
				if [[ $num == 1 ]]; then
					field=2
				elif [[ $num == 2 ]]; then
					field=1
				elif [[ $num == 3 ]]; then
					field=6
				else
					echo -e "\e[31;1m${ORDER[7]}\e[0m"
					continue
				fi
				read -p "${ORDER[8]}?(y/n)> " flag
				if [[ $flag = "y" ]]; then
					option="-r"
				else
					option=""
				fi
				lines=`wc -l stuRecord.txt | awk '{print $1}'`
				echo -n -e "\e[31;1m${ORDER[9]}\e[32;1m$lines\e[31;1m${ORDER[10]}\e[0m"
				if [[ lines != 0 ]]; then
					echo -e "   \e[31;1m${ORDER[11]}:(\e[32;1m${ORDER[12]} ${ORDER[13]} ${ORDER[14]} ${ORDER[15]} ${ORDER[16]} ${ORDER[17]} ${ORDER[18]}\e[31;1m)\e[0m"
				fi
				sort $option -t $'\t' -k $field -n stuRecord.txt | awk '{printf("%-8s%-12s%-8s%-8s%-8s%-8s%-16s\n", $1,$2,$3,$4,$5,$6,$7)}'
				break
			done
		elif [ $REPLY == 2 ]; then
			break
		else
			echo -e "\e[31;1m${ORDER[19]}\e[0m"
		fi
	done
}

quit()
{
	if [ $1 == 1 ]; then
		for((i = 0; i < ${#quit1[*]}; i++)); do
			QUIT[$i]=${quit1[$i]}
		done
	else
		for((i = 0; i < ${#quit2[*]}; i++)); do
			QUIT[$i]=${quit2[$i]}
		done
	fi
	read -p "${QUIT[0]}" flag
	if [[ $flag == y ]]; then
		. ./main.sh
	fi
}

##################################StartUp######################################
start
##################################EndUp######################################
