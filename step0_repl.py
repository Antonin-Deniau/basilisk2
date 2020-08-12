import sys

def read():
    return input()

def evl(e):
    return e

def prnt(e):
    sys.stdout.write(e)
    sys.stdout.write("\n")

def rep():
    while True:
        sys.stdout.write("user>")
        b = read()
        c = evl(b)
        prnt(c)

rep()
