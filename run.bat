mkdir bin
cd bin && del quail-gui.exe && cd ..
rsrc -ico quail-gui.ico -manifest quail-gui.exe.manifest
copy /y quail-gui.exe.manifest bin\quail-gui.exe.manifest
go build -buildmode=pie -ldflags="-s -w" -o quail-gui.exe main.go
move quail-gui.exe bin/quail-gui.exe
cd bin && quail-gui.exe c:\games\eq\rebuildeq\rkp.eqg
rem rkp.eqg