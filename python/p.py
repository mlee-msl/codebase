#!/usr/bin/env python
# -*- coding: utf-8 -*-

'''
test some features...
'''

import urllib
import re

def getList():
	html = urllib.urlopen('http://www.dytt8.net').read()
	print(html)
	reg = '<a href="(.*?)">'
	return re.findall(reg, html)

if __name__ == '__main__':
	print('*'*6 + 'begin' + '*'*6)
	getList()
	print('*'*6 + 'end' + '*'*6)
