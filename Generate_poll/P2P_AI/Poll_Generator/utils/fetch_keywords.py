
import time
import random
from django.core.cache import cache
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager
from config import config

def get_random_user_agent():
    return random.choice(config.USER_AGENTS)

def setup_driver():
    """Set up Selenium WebDriver with options."""
    options = Options()
    options.add_argument("--headless")
    options.add_argument("--no-sandbox")
    options.add_argument("--disable-dev-shm-usage")
    options.add_argument("--disable-blink-features=AutomationControlled")
    options.add_argument(f"user-agent={get_random_user_agent()}")

    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=options)

def fetch_google_trends():
    """Fetch top trending topics from Google Trends using Selenium."""
    cache_key = config.CACHE_KEY
    cached_data = cache.get(cache_key)
    if cached_data:
        print("Serving from cache")
        return cached_data

    driver = setup_driver()
    trending_topics = []

    try:
        print("Navigating to URL...")
        driver.get(config.TRENDS_URL)
        time.sleep(random.uniform(config.SLEEP_MIN, config.SLEEP_MAX))

        elements = driver.find_elements("class name", config.ELEMENT_CLASS_NAME)

        if elements:
            print("Extracting trending topics...")
            for element in elements[:config.TRENDING_TOPIC_COUNT]:
                text = element.text.strip()
                if text:
                    trending_topics.append(text)
        else:
            print("No trending topics found.")

    except Exception as e:
        print(f"Error: {e}")

    finally:
        driver.quit()

    cache.set(cache_key, trending_topics, timeout=config.CACHE_TIMEOUT)
    return trending_topics

