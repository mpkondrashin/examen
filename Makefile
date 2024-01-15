

.PHONY: clean tidy

setup.exe: install.exe
	go build ./cmd/setup
install.exe: examen.exe examensvc.exe
	go build ./cmd/install
examen.exe:
	go build ./cmd/examen
examensvc.exe:
	go build ./cmd/examensvc

clean: tidy
	rm setup.exe

tidy:
	rm examen.exe examensvc.exe install.exe 
