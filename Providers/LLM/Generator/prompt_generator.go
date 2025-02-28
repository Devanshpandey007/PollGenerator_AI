package Generator

var PredictionMarketPrompt = `A trending topic is: %s
Context from news: %s
Generate a yes/no prediction poll question that is actionable (e.g., "Will [X] happen by [Y date]?"), 
make sure of these few things :
its currently 2025, so make sure any questions are relevant to that year.
make sure the question is clear and concise.
make sure the question is relevant to the news.
look up the trending topic and context to get more information.
`
