#!/usr/bin/python
# -*- coding: utf-8 -*-

def preprocess(*data):
	s1, s2, len1, len2 = data
	s1 = transfer(fill(s1, len1-1, 'a'))
	s2 = transfer(fill(s2, len1-1, 'z'))
	l1 = transfer(fill(s1, len2, 'a'))
	l2 = transfer(fill(s2, len2, 'z'))
	return s1, s2, l1, l2

def fill(s, num, char):
	l = list(s)
	if num <= len(l): # truncate
		l = l[:num]
	else: # append
		l.extend(char*(num-len(l)))
	return ''.join(l)

def transfer(s):
	d = {}
	for i in range(26):
		d[chr(ord('a') + i)] = i if i < 10 else chr(ord('a') + (i - 10))
	l = list(s)
	for i, v in enumerate(l):
		l[i] = str(d[v])
	return ''.join(l)

def count(s1, s2, len1, len2):
	s1, s2 = (min(s1, s2), max(s1, s2))
	len1, len2 = (min(len1, len2), max(len1, len2))
	f = 1 if len(s1) <= len2 else 0
	s1, s2, l1, l2 = preprocess(s1, s2, len1, len2)
	print(s1, s2, l1, l2)
	# return (int(l2, 26) - int(l1, 26) + f) - (int(s2, 26) - int(s1, 26) + 1)

def main(s1, s2, len1, len2):
	return count(s1, s2, len1, len2)

if __name__ == '__main__':
	s1, s2, len1, len2 = input().strip().split()
	print(main(s1, s2, int(len1), int(len2)))
