from rest_framework import serializers
from .models import Topics, Articles

class TopicSerializer(serializers.ModelSerializer):
    class Meta:
        model = Topics
        fields = ['id', 'topic_name', 'is_poll_generated', 'date_created', 'date_last_updated']  

class ArticleSerializer(serializers.ModelSerializer):
    topic_name = serializers.CharField(source="topic.topic_name", read_only=True)  

    class Meta:
        model = Articles
        fields = ['topic_name', 'title', 'content', 'url', 'date_created', 'date_last_updated']  



