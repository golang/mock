include $(GOROOT)/src/Make.inc

all:	install

# Order matters!
DIRS=\
	mockgen\
	gomock\

install clean nuke:
	for dir in $(DIRS); do \
		make -C $$dir $@ || break; \
	done
