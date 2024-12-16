default:
    @just --list

# build receipt-processor-challenge binary
build:
    @echo '{{ BOLD + CYAN }}Building receipt-processor-challenge!{{ NORMAL }}'
    go build -o receipt-processor-challenge ./cmd/

# update go packages
update:
    @cd ./cmd && go get -u

# run tests
test:
    go test -v ./... -covermode=atomic -coverprofile=coverage.out

# api load test using k6
k6: build
    ./receipt-processor-challenge > /dev/null 2>&1 &
    k6 run ./k6/script.js
