package services

type OAuthInfoExtractor func(token string) (string, string, error)
