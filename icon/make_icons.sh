#!/bin/bash

cd "$(dirname "$0")" || (echo "Failed to change into script directory" >&2 && exit)

# Validate parameters
if [[ $# -ne 4 ]]; then
    echo "Usage: $0 --ext <jpg|png|ico> --target <unix|win>"
    exit 1
fi

while [[ $# -gt 0 ]]; do
    case "$1" in
        --ext)
            shift
            case "$1" in
                jpg|png|ico)
                    EXTENSION="$1"
                    ;;
                *)
                    echo "Invalid --ext parameter. Supported values: jpg, png, ico"
                    exit 1
                    ;;
            esac
            ;;
        --target)
            shift
            case "$1" in
                unix|win)
                    BUILD_OS="$1"
                    ;;
                *)
                    echo "Invalid --target parameter. Supported values: unix, win"
                    exit 1
                    ;;
            esac
            ;;
        *)
            echo "Unknown parameter: $1"
            exit 1
            ;;
    esac
    shift
done

if [ -z "$GOPATH" ]; then
    echo GOPATH environment variable not set
    exit
fi

if [ ! -e "$GOPATH/bin/2goarray" ]; then
    echo "Installing 2goarray..."
    go get github.com/cratonica/2goarray
    if [ $? -ne 0 ]; then
        echo Failure executing go get github.com/cratonica/2goarray
        exit
    fi
fi

if [[ $BUILD_OS == 'unix' ]]; then
  OUTPUT=icon-unix.go
  echo "//+build !windows" > $OUTPUT
else
  OUTPUT=icon-win.go
  echo "//+build windows" > $OUTPUT
fi

# Loop through files with specified extension
for FILE in *."$EXTENSION"; do
  FILENAME=$(basename "$FILE")
  FILENAME_NO_EXT="${FILENAME%.*}"

  echo >> $OUTPUT

  if ! < "$FILENAME" $GOPATH/bin/2goarray Data_"$FILENAME_NO_EXT" icon >> $OUTPUT
  then
      echo Failure writing to $OUTPUT
      exit
  fi
done

