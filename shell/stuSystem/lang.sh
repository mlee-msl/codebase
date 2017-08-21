#!/bin/bash

#print
print1=(学生成绩管理系统 增加学生记录 修改学生记录 删除学生记录 查询学生记录 排列学生记录 退出系统 '请输入操作编号> ' '无效的操作编号' '按任意键继续...')
print2=('Student Score Management System' 'Insert Student Record' 'Modify Student Record' 'Delete Student Record' 'Query Student Record' 'Sort Student Record' 'Quit' 'Please type number> ' 'Invalid number' 'Press any key to continue...')

#insert
insert1=('插入' '退出' '插入' '请输入学号' '学号为4位' '插入学生' '信息' '格式' '姓名' '性别' '班级' '举例' '张三' '男' '输入应有3个字段' '是否确定' '插入学生' '记录' '格式' '数据库' 'C语言' '数据结构' '百分制' '举例' '输入应有3个字段' '成绩为百分制' '成绩为非负数' '是否确定' '插入成功' '该学生' '已存在' '退出' '无效选择' '存在非整数')
insert2=('Insert' 'Quit' 'Insert' 'Please type number' 'The length of number should be 4' 'Insert student' 'information' 'Format' 'Name' 'Gender' 'Class' 'e.g.' 'zhangsan' 'male' 'Request three fields' 'Are you sure' 'Insert student' 'record' 'Format' 'Database' 'C-language' 'DataStructure' 'Hundred-mark System' 'e.g.' 'Request three fields' 'Hundred-mark System' 'Non-negative' 'Are you sure' 'Insert successfully!' 'The student' 'has existed' 'Quit' 'Invalid choice' 'Existed non-integer') 


#modify
modify1=('修改' '退出' '请输入学号' '学号为4位' '该学生' '不存在' '修改学生' '记录' '格式' '数据库' 'C语言' '数据结构' '百分制' '举例' '输入应有3个字段' '成绩为百分制' '成绩为非负数' '是否确定修改' '修改成功' '无效选择' '学生表为空' '所有学生' '存在非整数')
modify2=('Modify' 'Quit' 'Please type number' 'The length of number should be 4' 'The student' 'is not existed' 'Modify student' 'record' 'Format' 'Database' 'C-language' 'DataStructure' 'Hundred-mark System' 'e.g.' 'Request three fields' 'Hundred-m-ark System' 'Non-negative' 'Are you sure to modify' 'Modify successfully!' 'Invalid choice' 'None of students' 'All students' 'Existed non-integer')

#delete
delete1=('删除' '退出' '请输入学号' '学号为4位' '该学生' '不存在' '确定删除学生' '删除成功' '无效选择' '所有学生' '提示:All表示删除所有学生' '确定删除所有学生' '确定不后悔' '学生表为空')
delete2=('Delete' 'Quit' 'Please type number' 'The length of number should be 4' 'The student' 'is not existed' 'Are you sure to remove student' 'Remove successfully!' 'Invalid choice' 'All students' 'Note: "All" means to remove all students only once' 'Are you sure to remove all students' 'No regrets exactly' 'None of students')

#query
query1=('查询' '退出' '请输入学号' '学号为4位' '该学生' '不存在' '格式' '学号' '姓名' '数据库' 'C语言' '数据结构' '平均成绩' '修改时间' '是否查询学生信息' '格式' '学号' '姓名' '性别' '班级' '修改时间' '无效选择' '学生表为空' '所有学生')
query2=('Query' 'Quit' 'Please type number' 'The length of number should be 4' 'The student' 'is not existed' 'Format' 'SNo' 'Name' 'Database' 'C-language' 'DataStructure' 'Average' 'Modified-date' 'Whether to query student Information' 'Format' 'SNo' 'Name' 'Gender' 'Class' 'Modified-date' 'Invalid choice' 'None of students' 'All students')

#order
order1=('排序' '退出' '按照' '默认' '学号' '平均成绩' '请输入编号' '无效选择' '是否降序' '共计' '条记录' '格式' '学号' '姓名' '数据库' 'C语言' '数据结构' '平均成绩' '修改时间' '无效选择')
order2=('Sort' 'Quit' 'Sort by' 'Default' 'SNo' 'Average' 'Please type choice' 'Invalid choice' 'Whether to descend or not' 'Total' 'records' 'Format' 'SNo' 'Name' 'Database' 'C-language' 'DataStructure' 'Average' 'Modified-date' 'Invalid choice')

#quit
quit1=("是否确定退出?(y/n)> ")
quit2=("Are you sure to exit?(y/n)> ")

