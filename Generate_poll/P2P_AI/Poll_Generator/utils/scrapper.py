

import time
import tempfile
import random
import concurrent.futures
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from webdriver_manager.chrome import ChromeDriverManager
from newspaper import Article
from datetime import datetime
from tenacity import retry, stop_after_attempt, wait_fixed, retry_if_exception_type
from config import config

def get_random_user_agent():
    return random.choice(config.USER_AGENTS)

def setup_driver(headless=True):
    """Set up Selenium WebDriver with options."""
    user_data_dir = tempfile.mkdtemp()
    chrome_options = Options()
    if headless:
        chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument(f"user-agent={get_random_user_agent()}")
    chrome_options.add_argument(f"--user-data-dir={user_data_dir}")
    chrome_options.add_argument("--disable-blink-features=AutomationControlled")
    chrome_options.add_experimental_option("excludeSwitches", ["enable-automation"])
    chrome_options.add_experimental_option("useAutomationExtension", False)

    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=chrome_options)

@retry(stop=stop_after_attempt(3), wait=wait_fixed(5), retry=retry_if_exception_type(Exception))
def extract_from_url(url):
    """Extract article data from a single URL with retries."""
    extracted_articles = []
    article = Article(url, request_timeout=30)
    article.download()
    article.parse()
    date_last_updated = article.publish_date.date() if article.publish_date else None

    data = {
        "title": article.title,
        "url": article.url,
        "content": article.text,
        "date_last_updated": date_last_updated
    }
    extracted_articles.append(data)
    print(f"Extracted from: {url}")
    return extracted_articles

def extract_articles_from_punchng(topic):
    """Fetch and extract articles from PunchNG for a given topic."""
    if not topic:
        print("Error: Topic cannot be empty.")
        return []

    hub_url = f"https://punchng.com/search-page/?q={topic}#gsc.tab=0&gsc.q={topic}"
    driver = setup_driver(headless=False)

    try:
        driver.get(hub_url)
        WebDriverWait(driver, 15).until(
            EC.presence_of_all_elements_located((By.CSS_SELECTOR, "a.gs-title"))
        )

        headlines = driver.find_elements(By.CSS_SELECTOR, "a.gs-title")
        all_urls = list(set([headline.get_attribute("href") for headline in headlines if headline.get_attribute("href")]))[:config.ARTICLE_COUNT]

    except Exception as e:
        print(f"Failed to load search page: {e}")
        return []

    finally:
        driver.quit()

    extracted_articles = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        results = list(executor.map(extract_from_url, all_urls))

    for result in results:
        extracted_articles.extend(result)

    return extracted_articles



