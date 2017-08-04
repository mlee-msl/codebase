#!/usr/bin/env python3
# -*- coding: utf-8 -*-

'''
the functions or attributes like __XXX__, to represent special features
'''

class Student(object):
	__slot__ = ('__name', '__gender', '__age') # only enable to set those attrs, or AttributeError
	__classAttr1 = 'Student1' # private (interpreter realizes it by changing the name of attributes)
	classAttr2 = 'Student2' # public

	def __init__(self, name, gender, age):
		'''initial an object values'''
		self.__name = name
		self.__gender = gender
		self.__age = age

	def show(self):
		'''show an object all values'''
		print('(%s, %s, %s)' % (self.__name, self.__gender, self.__age))

	def __len__(self):
		'''enable to call len() to get special length of this object'''
		return len(self.__name)
	
	def __str__(self): # call it when print(object)
#		print('__str__')
		return '(%s, %s, %s)' % (self.__name, self.__gender, self.__age)

#	__repr__ = __str__ # call it when object directly
	def __repr__(self):
#		print('__repr__')
		return self.__str__()

if __name__ == '__main__':
	s = Student('MLee', 'male', 30)
	hasattr(student.Student, '__classAttr1') # hasattr only find attr whose modifier(access permittion) is the public, getattr & setattr both are the same as it
	hasattr(student.Student, 'classAttr2')

