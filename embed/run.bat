mkdir bin
rsrc -ico embed.ico -manifest embed.exe.manifest
copy /y embed.exe.manifest bin\embed.exe.manifest
go build -buildmode=pie -ldflags="-s -w" -o embed.exe embed.go
move embed.exe bin/embed.exe
cd bin && embed.exe c:\games\eq\rebuildeq\rkp.eqg