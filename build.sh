OUTPUT_DIR=bin

# mac 64
GOARCH=amd64 GOOS=darwin CGO_ENABLED=0 \
    go build -o $OUTPUT_DIR/macosx64/netkits -v netkits

# win 32
GOARCH=386 GOOS=windows CGO_ENABLED=0 \
    go build -o $OUTPUT_DIR/win32/netkits.exe  -v netkits

# win 64
GOARCH=amd64 GOOS=windows CGO_ENABLED=0 \
    go build -o $OUTPUT_DIR/win64/netkits.exe  -v netkits

# linux 32
GOARCH=386 GOOS=linux CGO_ENABLED=0 \
    go build -o $OUTPUT_DIR/linux32/netkits  -v netkits

# linux 64
GOARCH=amd64 GOOS=linux CGO_ENABLED=1 \
    go build -o $OUTPUT_DIR/linux64/netkits  -v netkits
