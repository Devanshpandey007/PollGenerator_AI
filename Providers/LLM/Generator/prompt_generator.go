package Generator

var PredictionMarketPrompt = `A trending topic is: %s
Context from news: %s
Generate a yes/no prediction poll question that is actionable (e.g., "Will [X] happen by [Y date]?"), 
make sure of these few things :
make sure the news date isn't in the past, it needs to be a future event.
make sure the question is clear and concise.
make sure the question is relevant to the news.
you can look up the news to get more context.
`
