# Use official Python image
FROM python:3.10

# Set working directory
WORKDIR /app

# Install system dependencies for Chrome & ChromeDriver
RUN apt-get update && apt-get install -y \
    wget unzip curl \
    libnss3 libgconf-2-4 libxi6 libxrandr2 libatk1.0-0 \
    libxcomposite1 libxcursor1 libxdamage1 libxfixes3 libasound2 \
    libpangocairo-1.0-0 fonts-liberation libappindicator3-1 \
    xdg-utils

# Install Google Chrome
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    dpkg -i google-chrome-stable_current_amd64.deb || apt-get -fy install

# Install ChromeDriver (must match Chrome version)
RUN wget -q https://chromedriver.storage.googleapis.com/114.0.5735.90/chromedriver_linux64.zip && \
    unzip chromedriver_linux64.zip -d /usr/local/bin/ && \
    chmod +x /usr/local/bin/chromedriver

# Copy and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy project files
COPY . .

# Expose Django port
EXPOSE 8003

# Start Django server
CMD ["sh", "-c", "python manage.py migrate && python manage.py runserver 0.0.0.0:8003"]
