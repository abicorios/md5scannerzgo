if exist C:\t rmdir C:\t /s /q
mkdir C:\t
del md5scannerzgo.exe
go build md5scannerzgo.exe
md5scannerzgo.exe readz "C:\Atari Jaguar by genres" "C:\t"
