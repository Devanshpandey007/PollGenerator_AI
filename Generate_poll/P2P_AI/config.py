import os

from dotenv import load_dotenv

load_dotenv()

class Config:
    
    USER_AGENTS = [
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
        "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0",
        "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36 Edge/120.0.0.0",
    ]
    
    # Google Trends settings (can be ignored for PunchNG if not needed)
    TRENDS_URL = "https://trends.google.com/trends/trendingsearches/daily?geo=NG"
    CACHE_KEY = "google_trending_topics_selenium"
    CACHE_TIMEOUT = 3600  
    TRENDING_TOPIC_COUNT = 10
    SLEEP_MIN = 4 
    SLEEP_MAX = 7
    ARTICLE_COUNT = 5  
    ELEMENT_CLASS_NAME = "mZ3RIc"


config = Config()  