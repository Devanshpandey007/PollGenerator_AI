import time
import concurrent.futures
from django.shortcuts import render
from .models import *
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from .serializers import *
import requests
import random
from .utils.fetch_keywords import fetch_google_trends
from .utils.scrapper import extract_articles_from_provider
from django.db import transaction
from rest_framework import generics
import os
from dotenv import load_dotenv
load_dotenv()


class TrendingTopicsAPIView(APIView):
    def get(self, request):
        try:
            trending_topics = Topics.objects.all().order_by('-date_created')

            serializer = TopicSerializer(trending_topics, many=True)

            return Response(serializer.data, status=status.HTTP_200_OK)

        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

    

    def post(self, request):
        try:
            trending_topics = fetch_google_trends()
            saved_topics = []
            for topic in trending_topics:
                obj, created = Topics.objects.get_or_create(topic_name=topic)
                if created:  
                    saved_topics.append(obj)

            if saved_topics:
                serializer = TopicSerializer(saved_topics, many=True)
                return Response(serializer.data, status=status.HTTP_201_CREATED)
            else:
                return Response({"message": "No new topics were added."}, status=status.HTTP_200_OK)

        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)




class ArticlesAPIView(APIView):
    """
    API view for managing articles: retrieving all articles (GET) and adding new ones (POST).
    """

    def get(self, request):
        try:
            articles = Articles.objects.all()
            serializer = ArticleSerializer(articles, many=True)
            if serializer.data:
                return Response(serializer.data, status=status.HTTP_200_OK)
            return Response({"error": "No articles found"}, status=status.HTTP_404_NOT_FOUND)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

    def post(self, request):
        try:
            added_articles = []
            latest_topics = Topics.objects.order_by('-id')[:os.getenv('EXTRACTED_TOPICS_COUNT')]

            new_topics = [
                topic for topic in latest_topics
                if not Articles.objects.filter(topic=topic).exists()
            ]

            if not new_topics:
                return Response({"error": "No new topics to process."}, status=status.HTTP_404_NOT_FOUND)

           
            with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
                articles_list = list(executor.map(extract_articles_from_provider, [topic.topic_name for topic in new_topics]))

      
            for topic, articles in zip(new_topics, articles_list):
                if not articles:
                    continue

                topic_instance = Topics.objects.get(topic_name=topic.topic_name)

                for article in articles:
                    if Articles.objects.filter(title=article['title']).exists() or Articles.objects.filter(url=article['url']).exists():
                        continue

                    with transaction.atomic():  
                        new_article = Articles.objects.create(
                            topic=topic_instance,
                            title=article['title'],
                            url=article['url'],
                            content=article['content'],
                            date_last_updated=article['date_last_updated']
                        )

                    article_data = ArticleSerializer(new_article).data
                    added_articles.append(article_data)

            return Response({"message": added_articles}, status=status.HTTP_201_CREATED)

        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)



class AddArticleAPIView(APIView):
    def post(self, request):
        try:
            topic_name = request.data.get("topic_name")
            title = request.data.get("title")
            url = request.data.get("url")
            content = request.data.get("content")

            if not all([topic_name, title, url, content]):
                return Response({"error": "All fields are required"}, status=status.HTTP_400_BAD_REQUEST)

           
            topic, _ = Topics.objects.get_or_create(topic_name=topic_name)

            if Articles.objects.filter(url=url).exists():
                return Response({"error": "Article with this URL already exists"}, status=status.HTTP_400_BAD_REQUEST)

            if Articles.objects.filter(title=title, topic=topic).exists():
                return Response({"error": "Article with this title already exists under this topic"}, status=status.HTTP_400_BAD_REQUEST)

            article = Articles.objects.create(topic=topic, title=title, url=url, content=content)

            serializer = ArticleSerializer(article)
            return Response(serializer.data, status=status.HTTP_201_CREATED)

        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_500_INTERNAL_SERVER_ERROR)

        




