#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from enum import Enum

class Month(Enum):
	Jan = 0
	Feb = 1
	Mar = 2
	Apr = 3
	May = 4
	Jun = 5
	Jul = 6
	Aug = 7
	Sep = 8
	Oct = 9
	Nov = 10
	Dec = 11

if __name__ == '__main__':
	keys = Month.__members__.keys()
	members = Month.__members__.values()
	items = Month.__members__.items()
	for key, member in items:
		print('%s(%s)-->%s' % (key, member, member.value))
