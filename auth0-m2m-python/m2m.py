import requests

# Auth0 credentials
client_id = 'DNf8TlUyWLzliFvbAjxyXpIQZOE6Yh69'
client_secret = 'OCwNrz5KNSnsgMCsIcTK-InZbwonm9g_k5GPdG2mCTf6SbMDvv3FXcXG7qanR7og'
domain = 'real-time-chat-auth.eu.auth0.com'
audience = 'https://real-time-chat-app.com/auth-api'

# Get access token
token_url = f'https://{domain}/oauth/token'
data = {
    'client_id': client_id,
    'client_secret': client_secret,
    'audience': audience,
    'grant_type': 'client_credentials'
}

response = requests.post(token_url, data=data)
token = response.json().get('access_token')

# Make an API call (example: call Auth0 Management API)
api_url = f'https://{domain}/api/v2/users'
headers = {'Authorization': f'Bearer {token}'}
response = requests.get(api_url, headers=headers)

print(response.json())