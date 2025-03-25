# P2P_AI

## Description
P2P_AI is a project aimed at fetching the top trending topics from Google Trends and extracting relevant articles from popular Nigerian websites like PunchNG. It then generates polls using LLMs (Large Language Models) based on the extracted data.

## Prerequisites
- **Docker**
- **Django**
- **SQLite3**
- **Python 3.10**

## Features
- Fetches trending topics from Google Trends.
- Extracts trending articles from Nigerian news sources.
- Uses LLMs to generate polls based on extracted content.
- Allows fetching the latest polls related to specific topics.

## Installation & Setup

### **1. Clone the Repository**
```sh
git clone <repo-url>
cd P2P_AI
```

### **2. Create & Activate Virtual Environment**
```sh
python3 -m venv venv
source venv/bin/activate  # On MacOS/Linux
venv\Scripts\activate  # On Windows
```

### **3. Install Dependencies**
```sh
pip install --upgrade pip setuptools wheel  
pip install -r requirements.txt --use-pep517
```

### **4. Set Up Environment Variables**
Create a `.env` file in the project's root directory and add the following required environment variables:

```ini
# Django Secret Key (Used for security and cryptographic signing)
SECRET_KEY=your-secure-random-key  

# GPT API Key (Used for interacting with GPT models)
GPT_KEY=your-gpt-api-key  

# Django Debug Mode (Set to False in Production)
DJANGO_DEBUG=True  

# Django Allowed Hosts (Comma-separated, e.g., "localhost,127.0.0.1,example.com")
DJANGO_ALLOWED_HOSTS=localhost,127.0.0.1

### **5. Run Migrations**
```sh 
python manage.py migrate # Apply the changes to the database
```

### **6. Start the Server**
```sh
python manage.py runserver
```
The application should now be running at `http://127.0.0.1:8000/`.

## API Endpoints
The project exposes the following API endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `api/topics/` | GET/POST | Fetch trending topics and view topics |
| `api/articles/` | GET/POST | Extract trending articles and view articles |
| `api/generate-polls/` | POST | Generate polls using LLMs |
| `api/custom-topics/` | POST | Manually post topic, title,content and url |
| `api/get-polls/<str:topic>/` | GET | Fetch polls for a particular topic |

## Database
The project uses **SQLite3** as its database.

## Docker Support
If you want to use Docker, follow these steps:

1. **Build and run the Docker Container:**
   ```sh
   docker-compose up -d --build
   ```
This will start the project inside a container.

2. **Stop the Docker Container:**
   ```sh
   docker-compose down
   ```
3. **Check docker logs:**
   ```sh
   docker logs -f <container_id>
   ```

## Additional Notes
- Ensure you have the correct `.env` file before running the project.
- SQLite3 is used for local development; for production, consider PostgreSQL or another database.

## Contributors
- **Devansh Pandey** 

---
This README serves as documentation for setting up and running the P2P_AI project. 

