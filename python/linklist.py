# -*- coding: utf-8 -*-
"""
Created on Wed Mar 27 20:53:53 2019

@author: MLee

Life is short, I use Python~_~
"""

class Node(object):
    def __init__(self, data):
        self.data = data
        self.next = None
       
#带有头结点的单链表
class LinkList(object):
    def __init__(self, seq):
        self.head = Node(0) # 头结点
        self.rear = None
        self.size = 0
        self.from_seq(seq)
        
    def size_(self):
        return self.size
    
    def adjust_index(self, index, action):
        if action == "insert":
            if index > self.size:
                index = self.size
            if index < -self.size:
                index = -self.size
            if index < 0:
                index = self.size - abs(index)
        elif action == "pop":
            if index >= self.size:
                raise Exception("pop index out of range")
            if index < -self.size:
                raise Exception("pop index out of range")
            if index < 0:
                index = self.size - abs(index)
        else:
            pass
        return index
        

    def is_empty(self):
        if self.size_() == 0:
            return True
        return False
    
    def clear(self):
        self.head.next = None;
        self.rear = None
        self.size = 0
    
    def append(self, data):
        if self.is_empty():
            self.insert(0, data)
        else:
            node = Node(data)
            self.rear.next = node
            self.rear = node
            self.size += 1
            
    def count(self, data):
        cnt = 0
        cur = self.head.next
        while cur is not None:
            if cur.data == data:
                cnt += 1
            cur = cur.next
        return cnt
    
    def extend(self, seq):
        if isinstance(seq, (str, tuple, list, dict)):
            if len(seq) == 0:
                return
            seq = LinkList(seq)
        if isinstance(seq, LinkList):
            if self.is_empty():
                self.head.next = seq.head.next
            else:
                self.rear.next = seq.head.next
            self.rear = seq.rear
            self.size += seq.size
        else:
            raise Exception("Invalid type", type(seq))
            
    def index(self, data):
        index = -1
        cur = self.head.next
        while cur is not None:
            index += 1
            if cur.data == data:
                return index
            cur = cur.next
        raise Exception("Not Found~_~")
            
   # 将data插入index处，其他节点依次向后移动
    def insert(self, index, data):
        index = self.adjust_index(index, "insert")
        pos = 0
        cur = self.head
        while pos != index:
            pos += 1
            cur = cur.next
        node = Node(data)
        node.next = cur.next
        cur.next = node
        if node.next is None: # 插入到最后的一个位置
            self.rear = node
        self.size += 1
    
    def pop(self, index=-1):
        index = self.adjust_index(index, "pop")
        pos = 0
        cur = self.head
        while pos != index:
            pos += 1
            cur = cur.next
        del_ = cur.next
        del_last = del_.next
        if del_last is None:
            self.rear = cur
        cur.next = del_last
        ret = del_.data
        del del_
        self.size -= 1
        return ret
    
    def remove(self, data):
        self.pop(self.index(data))
        
    def reverse(self):
        if self.size == 0:
            return
        first = self.head.next
        p = first.next
        first.next = None
        self.rear = first
        while p is not None:
            q = p
            p = p.next
            q.next = first
            first = q
        self.head.next = first
    
    def sort(self, reverse=False):
        lst = self.to_list()
        lst.sort(reverse=reverse)
        self.clear()
        self.from_seq(lst)
        
    def from_seq(self, seq):
        self.clear()
        for data in seq:
            self.append(data)
        
    def to_list(self):
        ret = []
        cur = self.head.next
        while cur is not None:
            ret.append(cur.data)
            cur = cur.next
        return ret
    
    def __len__(self): # 调用len方法时触发这个方法
        return self.size

    def __str__(self): # 调用str方法时触发这个方法
        if self.size == 0:
            return "<>" # str: "", tuple: (), list: [], dict: {}, set: set({}), linklist: <>  目前不知道怎么实现如何通过<>形式创建这个LinkList对象
        return "<" + ", ".join(map(str, self.to_list())) + ">"
    
    def __iter__(self):
        cur = self.head.next
        while cur is not None:
            yield cur.data
            cur = cur.next

    def __add__(self, link_list): # 两个LinkList进行+操作,self自身不发生改变
        return LinkList(self.to_list() + link_list.to_list())
    
    def __mul__(self, times): # LinkList进行*重复操作,self自身不发生改变
        return LinkList(self.to_list()*times)

if __name__ == "__main__":
    pass
