// HomePage.js
import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';
import LogoutButton from './LogoutButton';

const HomePage = () => {
  const { user, isAuthenticated } = useAuth0();

  if (!isAuthenticated) {
    return null;
  } else {
    return (
        <div
        style={{
            maxWidth: '400px',
            padding: '20px',
            textAlign: 'center',
            borderRadius: '12px',
            boxShadow: '0 4px 8px rgba(0,0,0,0.1)',
            backgroundColor: '#f9f9f9',
        }}
        >
        <h1 style={{ color: '#333' }}>Welcome, {user.name || 'User'}!</h1>
        <p style={{ color: '#555' }}>You have successfully logged in.</p>
        <LogoutButton />
        </div>
    );
  }
};

export default HomePage;
