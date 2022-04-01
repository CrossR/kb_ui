#/bin/sh

if [ -z "$GOPATH" ]; then
    echo "GOPATH environment variable not set"
    exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
    echo "Installing 2goarray..."
    go get github.com/cratonica/2goarray
    if [ $? -ne 0 ]; then
        echo "Failure executing go get github.com/cratonica/2goarray"
        exit
    fi
fi

if [ -z "$1" ]; then
    echo "Please specify a PNG file"
    exit
fi

if [ ! -f "$1" ]; then
    echo $1 is not a valid file
    exit
fi

if [ -z "$2" ]; then
    echo "Please give a name for the icon"
    exit

OUTPUT=${1%.*}_unix.go
echo "Generating $OUTPUT"
echo "go:build linux || darwin" > $OUTPUT
echo >> $OUTPUT
cat "$1" | $GOPATH/bin/2goarray $2 icons >> $OUTPUT
if [ $? -ne 0 ]; then
    echo "Failure generating $OUTPUT"
    exit
fi
echo "Finished"

