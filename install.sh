echo "changing directory to go-form-webserver"
cd $GOPATH/src/github.com/strick-j/go-form-webserver
echo "building the go binary"
go build -o go-form-webserver

echo "starting the binary"
./go-form-webserver