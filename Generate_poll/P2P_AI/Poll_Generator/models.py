from django.db import models

class Topics(models.Model):
    topic_name = models.CharField(max_length=255)
    is_poll_generated = models.BooleanField(default=False) 
    date_last_updated = models.DateTimeField(blank=True, null=True)
    date_created = models.DateTimeField(auto_now_add=True)



class Articles(models.Model):
    topic = models.ForeignKey(Topics, on_delete=models.CASCADE, related_name="articles")  
    title = models.TextField(unique=True) 
    url = models.URLField(unique=True)
    content = models.TextField()
    date_created = models.DateTimeField(auto_now_add=True)
    date_last_updated = models.DateField(blank=True, null=True)

    def __str__(self):
        return self.title