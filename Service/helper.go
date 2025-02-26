package main

// GetSecrets retrieves secrets from AWS Secrets Manager.
// func GetSecrets(path string, region string) (string, error) {
// 	Common.LogInfo("Reading secrets from AWS Secrets Manager", map[string]interface{}{"path": path, "region": region})

// 	// Get the current environment
// 	env := os.Getenv("ENV")
// 	if env == "" {
// 		env = "test" // Default to "test" if ENV is not set
// 		Common.LogInfo("ENV not set, defaulting to 'test'", nil)
// 	}

// 	// If in test environment, return a test API key
// 	if env == "test" {
// 		Common.LogInfo("Test environment detected, returning test API key", nil)
// 		return "test-api-key", nil
// 	}

// 	// Load AWS configuration with the specified region
// 	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
// 	if err != nil {
// 		Common.LogError("Failed to load AWS configuration", err, map[string]interface{}{"region": region})
// 		return "", fmt.Errorf("could not load AWS config: %w", err)
// 	}

// 	// Create Secrets Manager client
// 	svc := secretsmanager.NewFromConfig(cfg)
// 	Common.LogInfo("AWS Secrets Manager client initialized", nil)

// 	// Prepare the input for Secrets Manager
// 	input := &secretsmanager.GetSecretValueInput{
// 		SecretId:     aws.String(path),
// 		VersionStage: aws.String("AWSCURRENT"), // Use the current version of the secret
// 	}

// 	// Attempt to fetch the secret value
// 	Common.LogInfo("Fetching secret value from AWS Secrets Manager", map[string]interface{}{"path": path})
// 	result, err := svc.GetSecretValue(context.TODO(), input)
// 	if err != nil {
// 		Common.LogError("Failed to retrieve secret value", err, map[string]interface{}{"path": path})
// 		return "", fmt.Errorf("failed to retrieve secret value for %s: %w", path, err)
// 	}

// 	// Ensure the SecretString field is populated
// 	if result.SecretString == nil {
// 		errMsg := fmt.Sprintf("No SecretString found for secret %s", path)
// 		Common.LogError(errMsg, nil, nil)
// 		return "", fmt.Errorf(errMsg)
// 	}

// 	// Parse the JSON response to extract the value
// 	var secretMap map[string]string
// 	err = json.Unmarshal([]byte(*result.SecretString), &secretMap)
// 	if err != nil {
// 		Common.LogError("Failed to parse secret JSON", err, map[string]interface{}{"secret": *result.SecretString})
// 		return "", fmt.Errorf("failed to parse secret JSON: %w", err)
// 	}

// 	// Ensure there's only one key in the secret map
// 	if len(secretMap) != 1 {
// 		errMsg := fmt.Sprintf("Unexpected secret format for %s: expected one key-value pair, got %d", path, len(secretMap))
// 		Common.LogError(errMsg, nil, nil)
// 		return "", fmt.Errorf(errMsg)
// 	}

// 	// Return the value for the first key
// 	for _, value := range secretMap {
// 		Common.LogInfo("Successfully retrieved secret value", map[string]interface{}{"path": path})
// 		return value, nil
// 	}
// 	// This line should never be reached
// 	return "", fmt.Errorf("unexpected error occurred while retrieving secret value for %s", path)
// }
