setup:
	sudo mkdir -p /opt/scowl

start: setup
	docker-compose up -d

getscowl:
	wget 'https://sourceforge.net/projects/wordlist/files/SCOWL/2020.12.07/scowl-2020.12.07.tar.gz' 

unzip:
	tar -xvf scowl-2020.12.07.tar.gz 

build:
	docker build -t fuzzylemma/scowldb:latest .

push:
	docker push fuzzylemma/scowldb:latest
