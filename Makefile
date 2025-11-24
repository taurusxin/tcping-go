NAME=tcping
BINDIR=bin
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags '-w -s -buildid='
VERSION=1.2.0
# The -w and -s flags reduce binary sizes by excluding unnecessary symbols and debug info
# The -buildid= flag makes builds reproducible

all: releases

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

linux-arm64:
	GOARCH=arm64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

windows-386:
	GOARCH=386 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(NAME)-$@.exe

releases: linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-386
	chmod +x $(BINDIR)/$(NAME)-*
	tar zcf $(BINDIR)/$(NAME)-linux-amd64-$(VERSION).tar.gz -C $(BINDIR) $(NAME)-linux-amd64
	tar zcf $(BINDIR)/$(NAME)-linux-arm64-$(VERSION).tar.gz -C $(BINDIR) $(NAME)-linux-arm64
	tar zcf $(BINDIR)/$(NAME)-darwin-amd64-$(VERSION).tar.gz -C $(BINDIR) $(NAME)-darwin-amd64
	tar zcf $(BINDIR)/$(NAME)-darwin-arm64-$(VERSION).tar.gz -C $(BINDIR) $(NAME)-darwin-arm64
	zip -j $(BINDIR)/$(NAME)-windows-386-$(VERSION).zip $(BINDIR)/$(NAME)-windows-386.exe
	zip -j $(BINDIR)/$(NAME)-windows-amd64-$(VERSION).zip $(BINDIR)/$(NAME)-windows-amd64.exe
	rm -f $(BINDIR)/$(NAME)-darwin-amd64 $(BINDIR)/$(NAME)-darwin-arm64 $(BINDIR)/$(NAME)-linux-amd64 $(BINDIR)/$(NAME)-linux-arm64 $(BINDIR)/$(NAME)-windows-386.exe $(BINDIR)/$(NAME)-windows-amd64.exe

clean:
	rm $(BINDIR)/*-$(VERSION).tar.gz $(BINDIR)/*-$(VERSION).zip
