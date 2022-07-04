BINDIR=bin

#.PHONY: pbs

all: a i
#
#pbs:
#	cd pbs/ && $(MAKE)
#

a:
	 gomobile bind -v -o $(BINDIR)/simple.aar -target=android -ldflags=-s github.com/redeslab/simpleV2-go/mobile
i:
	gomobile bind -v -o $(BINDIR)/simple.xcframework -target=ios  -ldflags="-w" -ldflags=-s github.com/redeslab/simpleV2-go/mobile

clean:
	gomobile clean
	rm $(BINDIR)/*
