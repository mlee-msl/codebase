#!/usr/bin/env python3
# -*- coding: utf-8 -*-

class Fibonacci(object):
	def __init__(self):
		self.__first, self.__second = 0, 1

	def __iter__(self):
		return self

	def __next__(self):
		self.__first, self.__second = self.__second, self.__first + self.__second
		if self.__first > 1000:
			raise StopIteration()
		return self.__first

	def __getitem__(self, n):
		'unsuported negative index owing to unknown length'
		if isinstance(n ,int):
#			for i in range(n): # error, have changed value of the object
#				self.__first, self.__second = self.__second, self.__first + self.__second
			first, second = 0, 1
			for i in range(n):
				first, second = second, first + second
			return second
		if isinstance(n, slice):
			start = n.start
			stop = n.stop
			step = n.step
			if start is None:
				start = 0
			if step is None:
				step = 1
			indices = list(range(start, stop, step))
			l = [] # save the list of result by slice operation
			for i in indices:
				l.append(self.__getitem__(i))
			return l

	def __call__(self, n=0): # if you define this method, you can treat object as a function who can be callable, e.g. Fibonacci()(3)
		return self.__getitem__(n)
