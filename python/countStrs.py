#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def preprocess(*data):
	s1, s2, len1, len2 = data
	l1 = list(s1)
	l2 = list(s2)
	m = max(len(l1), len(l2), len1, len2)
	for i in range(len(l1) - m):
		l1.append('0')
	s1 = ''.join(l1)
	for i in range(len(s2) - m):
		l2.append('0')
	s2 = ''.join(l2)
	return s1, s2, len1, len2

def count(*data):
	s1, s2, len1, len2 = preprocess(*data)
	print(data)
	c = abs(int(s1, 26) - int(s2, 26))
	return c
	

if __name__ == '__main__':
	count('ab', 'ce', 1, 4)
	
