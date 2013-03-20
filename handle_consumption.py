import sys
import getopt
import tempfile
import os
import atexit
import subprocess
import datetime
import time

tmpFile = None
s3CmdDir = None
s3TargetDir = None
exiting = False

def bufferBlock(blockSize):
    global tmpFile, exiting

    fileSize = 0
    tmpFile = tempfile.NamedTemporaryFile(delete=False)

    while (exiting == False and fileSize < blockSize):
        s = sys.stdin.readline()
        tmpFile.write(s)
        fileSize += len(s)

    tmpFile.close()
    
def storeBlock():
    global tmpFile, s3CmdDir, s3TargetDir

    now = datetime.datetime.now()

    filePath = tmpFile.name
    tmpFileBaseName = os.path.basename(filePath)

    #target file name formatted to sort alphabetacally by newest file and be safe from duplicates
    targetFileName = "%04d-%02d-%02d-%d-%07d-%s" % (now.year, now.month, now.day, time.mktime(now.timetuple()), now.microsecond, tmpFileBaseName)
    targetFullPath = s3TargetDir + "/" + targetFileName

    s3CmdOutput = subprocess.check_output([s3CmdDir + "/s3cmd", "put", filePath, targetFullPath])
    #todo: error handling here

    os.unlink(filePath)

def beforeExit():
    global exiting, tmpFile

    exiting = True
    tmpFile.close()
    storeBlock()

def main(argv):
    global s3CmdDir, s3TargetDir

    #make sure we flush any temporary local data to s3 before exiting
    atexit.register(beforeExit)

    try:
        opts, args = getopt.getopt(argv[1:], "", ["blocksize=", "s3cmd-dir=", "s3target-dir="])
    except getopt.GetoptError:
        print("invalid option error")
        return 2
    
    for opt, arg in opts:
        if opt in ("--blocksize"):
            blockSize = int(arg)
        elif opt in ("--s3cmd-dir"):
            s3CmdDir = arg
        elif opt in ("--s3target-dir"):
            s3TargetDir = arg
    
    if (s3CmdDir == None or s3TargetDir == None or blockSize == None):
        print("one of the required options was missing")
        return 2

    while True:
        bufferBlock(blockSize)
        storeBlock()


if __name__ == "__main__":
    sys.exit(main(sys.argv))
