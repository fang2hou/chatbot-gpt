#!/bin/sh

APP_NAME=$1
SRC_DIR=./cmd/$1
BIN_DIR=./bin/$1

if [ -n "$2" ]; then
  GOOS=$2
else
  GOOS=$(go env GOOS)
fi

if [ -n "$3" ]; then
  GOARCH=$3
else
  GOARCH=$(go env GOARCH)
fi

if [ -n "$4" ]; then
  GOARM=$4;
else
  GOARM=$(go env GOARM)
fi

ARCH=$GOARCH
if [ -n "$GOARM" ]; then
    ARCH="$ARCH"v"$GOARM"
fi

START_MSG="--- 🛠 START BUILDING $APP_NAME FOR $GOOS/$ARCH ---"

echo "$START_MSG"

FILENAME="$APP_NAME"_"$GOOS"_"$ARCH"
if [ "$GOOS" = "windows" ]; then
  FILENAME="$FILENAME.exe";
fi

mkdir -p "$BIN_DIR"

CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" GOARM="$GOARM" go build -o "$BIN_DIR"/"$FILENAME" "$SRC_DIR"

END_MSG="--- ✅ DONE ($APP_NAME $GOOS/$ARCH)"
END_MSG_LENGTH="${#END_MSG}"
END_MSG_LENGTH=$((END_MSG_LENGTH + 2))
DASHES=$(printf "%0.s-" $(seq $END_MSG_LENGTH "${#START_MSG}"))
echo "$END_MSG $DASHES"