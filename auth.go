package main

import "slices"

var authCache map[string][]string // {"collection-action" : "access token"}

func buildAuthCache(config Config) map[string][]string {
	var tokenCache = make(map[string]string)
	for _, accessToken := range config.AccessTokens {
		tokenCache[accessToken.Name] = accessToken.Token
	}
	var authCache = make(map[string][]string)
	for _, collection := range config.Collections {
		baseTokenNames := collection.Auth[ActionAll]
		authCache[collection.Name+"-"+ActionCreate] = tokensFromTokenNames(baseTokenNames, collection.Auth[ActionCreate], tokenCache)
		authCache[collection.Name+"-"+ActionRead] = tokensFromTokenNames(baseTokenNames, collection.Auth[ActionRead], tokenCache)
		authCache[collection.Name+"-"+ActionList] = tokensFromTokenNames(baseTokenNames, collection.Auth[ActionList], tokenCache)
		authCache[collection.Name+"-"+ActionReplace] = tokensFromTokenNames(baseTokenNames, collection.Auth[ActionReplace], tokenCache)
		authCache[collection.Name+"-"+ActionPatch] = tokensFromTokenNames(baseTokenNames, collection.Auth[ActionPatch], tokenCache)
		authCache[collection.Name+"-"+ActionDelete] = tokensFromTokenNames(baseTokenNames, collection.Auth[ActionDelete], tokenCache)
	}
	return authCache
}

func tokensFromTokenNames(baseTokenNames []string, tokenNames []string, tokenCache map[string]string) []string {
	mergedTokens := []string{}
	for _, tokenName := range baseTokenNames {
		mergedTokens = append(mergedTokens, tokenCache[tokenName])
	}
	for _, tokenName := range tokenNames {
		mergedTokens = append(mergedTokens, tokenCache[tokenName])
	}
	return mergedTokens
}

func isAuthTokenValid(token string, collectionName string, action string) bool {
	validTokens := authCache[collectionName+"-"+action]
	return slices.Contains(validTokens, token)
}
