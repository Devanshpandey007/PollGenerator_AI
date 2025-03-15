# Use official Python image
FROM python:3.10

# Set working directory
WORKDIR /app

# Copy and install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy project files
COPY . .

# Expose Django port
EXPOSE 8003

# Install cron
RUN apt-get update && apt-get install -y cron

# Copy crontab file inside container
COPY crontab /etc/cron.d/mycron

# Set permissions and install the cron job
RUN chmod 0644 /etc/cron.d/mycron && crontab /etc/cron.d/mycron

# Start cron in the background and run Django
CMD ["sh", "-c", "cron && python manage.py migrate && python manage.py runserver 0.0.0.0:8003"]
