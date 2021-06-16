#!/usr/bin/python

import random
import string

for i in range(10000):
    ran_str = ''
    for j in range(512):
        tmp_str = ''.join(random.sample(string.ascii_letters + string.digits, 10))
        ran_str = ran_str + tmp_str
    print 'set key'+str(i),ran_str
