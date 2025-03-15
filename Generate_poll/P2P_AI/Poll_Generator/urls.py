from django.urls import path
from . import views



urlpatterns = [
    path('topics/', views.TrendingTopicsAPIView.as_view(), name="trending-topics"),
    path('articles/', views.ArticlesAPIViews.as_view(), name="articles"),
    path('custom-topics/', views.AddArticleAPIView.as_view(), name="custom-topics")
]