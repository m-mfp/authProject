import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';

const LoginButton = () => {
  const { loginWithRedirect } = useAuth0();

  return (
    <button
      onClick={() => loginWithRedirect()}
      style={{
        backgroundColor: '#4A90E2',
        color: 'white',
        padding: '12px 24px',
        fontSize: '16px',
        borderRadius: '8px',
        border: 'none',
        cursor: 'pointer',
        boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
        transition: 'background-color 0.3s ease',
      }}
      onMouseOver={e => (e.currentTarget.style.backgroundColor = '#357ABD')}
      onMouseOut={e => (e.currentTarget.style.backgroundColor = '#4A90E2')}
    >
      Log In
    </button>
  );
};

export default LoginButton;
