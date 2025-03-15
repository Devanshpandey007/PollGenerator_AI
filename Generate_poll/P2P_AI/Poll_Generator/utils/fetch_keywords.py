

import time
import random
import requests
from django.core.cache import cache
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager

# Rotate User Agents
USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Safari/537.36 Edge/120.0.0.0",
]

def fetch_google_trends():
    """Fetch top 10 Google trending topics using Selenium safely."""
    cache_key = "google_trending_topics_selenium"
    cached_data = cache.get(cache_key)

    if cached_data:
        print("Serving from cache")
        return cached_data  

    
    options = Options()
    options.add_argument("--headless")  
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-blink-features=AutomationControlled")  
    options.add_argument(f"user-agent={random.choice(USER_AGENTS)}")  

    # Set up ChromeDriver
    service = Service(ChromeDriverManager().install())
    driver = webdriver.Chrome(service=service, options=options)

    url = "https://trends.google.com/trends/trendingsearches/daily?geo=NG"

    trending_topics = []

    try:
        print("Navigating to URL...")
        driver.get(url)

        time.sleep(random.uniform(4, 7))

        elements = driver.find_elements("class name", "mZ3RIc")

        if elements:
            print("Extracting trending topics...")
            for element in elements[:10]:  
                text = element.text.strip()
                if text:
                    trending_topics.append(text)

        else:
            print("No trending topics found.")

    except Exception as e:
        print(f"Error: {e}")

    finally:
        driver.quit()

    # Cache the results for 1 hour
    cache.set(cache_key, trending_topics, timeout=3600)

    return trending_topics
