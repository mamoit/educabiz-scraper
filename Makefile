all: linux windows

linux:
	go build -trimpath
	# strip educabiz-scraper
	# upx --best educabiz-scraper

windows:
	CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go build -trimpath
	# strip educabiz-scraper.exe
	# upx --best educabiz-scraper.exe
