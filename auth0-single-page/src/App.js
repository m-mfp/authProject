import React from 'react';
import { Auth0Provider, useAuth0 } from '@auth0/auth0-react';

const App = () => {
  return (
    <Auth0Provider
      domain="real-time-chat-auth.eu.auth0.com"
      clientId="hGzms0BWRTLvA1Ekn1FBSJ71QJeCATtr"
      authorizationParams={{
        redirect_uri: window.location.origin + '/callback'
      }}
    >
      <AuthComponent />
    </Auth0Provider>
  );
};

const AuthComponent = () => {
  const { loginWithRedirect, logout, user, isAuthenticated } = useAuth0();

  return (
    <div>
      {!isAuthenticated ? (
        <button onClick={() => loginWithRedirect()}>Login with Auth0</button>
      ) : (
        <div>
          <pre>{JSON.stringify(user, null, 2)}</pre>
          <button onClick={() => logout({ returnTo: window.location.origin })}>
            Logout
          </button>
        </div>
      )}
    </div>
  );
};

export default App;