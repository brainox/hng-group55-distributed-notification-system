# syntax=docker/dockerfile:1

FROM python:3.11-slim

WORKDIR /code

# Copy dependencies first for Docker layer caching
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy the rest of your project
COPY . .

# Expose the API port
EXPOSE 8000

# Default command (can be overridden by docker-compose)
CMD ["bash"]
