include $(GOROOT)/src/Make.inc

TARG=mybot
GOFILES=\
	ants.go\
	map.go\
	main.go\
	debugging.go\
	mybot.go\

include $(GOROOT)/src/Make.cmd
