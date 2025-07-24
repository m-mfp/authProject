const express = require('express');
const { auth } = require('express-openid-connect');
const app = express();

const config = {
  authRequired: false,
  auth0Logout: true,
  baseURL: 'http://localhost:3000',
  clientID: 'tBNDKW2pyL1A4OsmWPXqOqYIE0H07bHb',
  issuerBaseURL: 'real-time-chat-auth.eu.auth0.com',
  secret: 'random-secret-string', // Required for session handling
};

app.use(auth(config));

app.get('/', (req, res) => {
  if (!req.oidc.isAuthenticated()) {
    res.send('<a href="/login">Login with Auth0</a>');
  } else {
    res.send(`
      <h1>Welcome</h1>
      <pre>${JSON.stringify(req.oidc.user, null, 2)}</pre>
      <a href="/logout">Logout</a>
    `);
  }
});

app.listen(3000, () => console.log('Native app running on http://localhost:3000'));