setup:
	sudo mkdir /opt/scowl

start: setup
	docker-compose up -d
