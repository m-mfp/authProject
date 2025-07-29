import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import LoginButton from './LoginButton';
import HomePage from './HomePage';

const App = () => {
  const { isAuthenticated } = useAuth0();

  return (
    <div
      style={{
        display: 'flex',
        height: '100vh',
        justifyContent: 'center',
        alignItems: 'center',
        fontFamily: 'Arial, sans-serif',
      }}
    >
      {!isAuthenticated ? <LoginButton /> : <HomePage />}
    </div>
  );
};

export default App;
