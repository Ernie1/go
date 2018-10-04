f = open("test", "w+")
for i in range(500):
	f.write("line %d\n" % i)
f.close()