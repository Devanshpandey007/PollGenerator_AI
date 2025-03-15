# P2P_AI

## Description
P2P_AI is a project aimed at fetching the top trending topics from Google Trends and extracting relevant articles from popular Nigerian websites like PunchNG. It then generates polls using LLMs (Large Language Models) based on the extracted data.

## Technologies Used
- **Docker**
- **Django**
- **SQLite3**

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
pip install -r requirements.txt
```

### **4. Set Up Environment Variables**
Create a `.env` file and add necessary environment variables required by the project.

### **5. Run Migrations**
```sh
python manage.py migrate
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
| `/topics/` | GET/POST | Fetch trending topics and view topics |
| `/articles/` | GET/POST | Extract trending articles and view articles |
| `/generate-polls/` | POST | Generate polls using LLMs |
| `/custom-topics/` | POST | Manually post topic, title,content and url |
| `/get-polls/<str:topic>/` | GET | Fetch polls for a particular topic |

## Database
The project uses **SQLite3** as its database.

## Docker Support
If you want to use Docker, follow these steps:
1. **Build the Docker Image:**
   ```sh
   docker build -t p2p_ai .
   ```
2. **Run the Docker Container:**
   ```sh
   docker-compose up -d
   ```
This will start the project inside a container.

## Additional Notes
- Ensure you have the correct `.env` file before running the project.
- SQLite3 is used for local development; for production, consider PostgreSQL or another database.

## Contributors
- **Devansh Pandey** 

---
This README serves as documentation for setting up and running the P2P_AI project. 

