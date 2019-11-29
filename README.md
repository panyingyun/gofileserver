# gofileserver
go fileserver for study gin web framework

./gofileserver 

upload file:

curl -X POST http://127.0.0.1:6001/upload -F "file=@/gwclient_store/cae/output.txt" -H "Content-Type: multipart/form-data" 

upload multi files:

curl -X POST http://127.0.0.1:6001/upload -F "files=@/root/go.mod" -F "files=@/root/go.sum" -H "Content-Type: multipart/form-data"

pi upload:

curl -X POST http://127.0.0.1:6001/upload -F "file=@/home/pi/frp/frpc.ini" -H "Content-Type: multipart/form-data"

