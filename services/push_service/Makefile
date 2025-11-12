.PHONY: help install run worker all logs rebuild stop clean

help:
	@echo "Usage:"
	@echo "  make install      Install Python dependencies"
	@echo "  make run          Run FastAPI server (Uvicorn)"
	@echo "  make worker       Run the push queue worker"
	@echo "  make all          Run both FastAPI + worker + RabbitMQ (Docker)"
	@echo "  make logs         View logs from running containers"
	@echo "  make rebuild      Rebuild Docker images"
	@echo "  make stop         Stop all Docker containers"
	@echo "  make clean        Remove containers and volumes"

install:
	pip install --upgrade pip && pip install -r requirements.txt

run:
	uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

worker:
	python app/consumer.py

all:
	sudo docker compose up --build

logs:
	sudo docker compose logs -f

rebuild:
	sudo docker compose up --build --force-recreate

stop:
	sudo docker compose down

clean:
	sudo docker compose down -v --remove-orphans
