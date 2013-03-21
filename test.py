#display out put line by line
import subprocess

def blah(foo,bar):
    print(foo)
    print(bar)

blah(bar="bar")
quit()

proc = subprocess.Popen(['/home/ubuntu/kafka/bin/kafka-console-consumer.sh','--zookeeper','localhost:2181','--topic','twitter'],stdout=subprocess.PIPE)
#works in python 3.0+
#for line in proc.stdout:
for line in iter(proc.stdout.readline,''):
    print "line: " + line
