import time
import requests
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from bs4 import BeautifulSoup
from webdriver_manager.chrome import ChromeDriverManager
from newspaper import Article
import random
from datetime import datetime
import threading
import tempfile



USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
]

def get_random_user_agent():
    return random.choice(USER_AGENTS)

lock = threading.Lock()

def extract_article(url, extracted_articles):
    try:
        print(f"Processing: {url}")
        time.sleep(random.uniform(1.5, 3))

        article = Article(url)
        article.download()
        article.parse()
        if article.publish_date:
            if isinstance(article.publish_date, datetime):  
                date_last_updated = article.publish_date.date()  # Extract only the date part
            else:
                date_last_updated = None
        else:
            date_last_updated = None

        data = {
            "title": article.title,
            "url": article.url,
            "content": article.text,
            "date_last_updated": date_last_updated
        }

        with lock:
            extracted_articles.append(data)

    except Exception as e:
        print(f"Error processing article at {url}: {e}")

def extract_articles_from_punchng(topic):
    extracted_articles = []

    if not topic:
        print("Error: Topic cannot be empty.")
        return extracted_articles  
    
    user_data_dir = tempfile.mkdtemp() 
    
    chrome_options = Options()
    chrome_options.add_argument("--headless") 
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument(f"user-agent={get_random_user_agent()}")
    chrome_options.add_argument(f"--user-data-dir={user_data_dir}") 
    chrome_options.add_argument("--disable-blink-features=AutomationControlled")  
    chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])  
    chrome_options.add_experimental_option("useAutomationExtension", False)

    service = Service(ChromeDriverManager().install())
    driver = webdriver.Chrome(service=service, options=chrome_options)

    search_url = f"https://punchng.com/search-page/?q={topic}#gsc.tab=0&gsc.q={topic}"
    driver.get(search_url)

    try:
        WebDriverWait(driver, 15).until(
            EC.presence_of_all_elements_located((By.CSS_SELECTOR, "a.gs-title"))
        )

        headlines = driver.find_elements(By.CSS_SELECTOR, "a.gs-title")
        all_urls = list(set([headline.get_attribute("href") for headline in headlines if headline.get_attribute("href")]))[:5]

    except Exception as e:
        print(f"Error while fetching search results: {e}")
        driver.quit()
        return extracted_articles

    finally:
        driver.quit()

    
    threads = []
    for url in all_urls:
        thread = threading.Thread(target=extract_article, args=(url, extracted_articles))
        thread.start()
        threads.append(thread)

    
    for thread in threads:
        thread.join()

    return extracted_articles