FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy app code
COPY . .

# Make the script executable
RUN chmod +x start.sh

# Run migrations, then start FastAPI
CMD ["./start.sh"]
