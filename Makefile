include $(GOROOT)/src/Make.$(GOARCH)

TARG=file_splitter
GOFILES=\
				file_splitter.go

include $(GOROOT)/src/Make.pkg

