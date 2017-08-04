#!/usr/bin/python3
# -*- coding: utf-8 -*-

from urllib import request
import re

def getList():
	html = request.urlopen(r'http://www.dytt8.net').read().decode('gb2312') # bytes-like to str by decode
#	html = request.urlopen(r'http://www.dytt8.net').read() # html is a bytes object
	f = open(r'/home/mlee/Desktop/html.txt', 'w+') # only enable to write by bytes(wb+)
	f.write(html)
	reg = r'<a href=[\'"](.*?)[\'"]>'
	txt = f.read()
	f.close()
	return re.findall(reg, txt)

if __name__ == '__main__':
	print('*'*6 + 'begin' + '*'*6)
	getList()
	print('*'*6 + 'end' + '*'*6)
