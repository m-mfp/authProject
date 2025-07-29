import React from 'react';
import { useAuth0 } from '@auth0/auth0-react';

const LogoutButton = () => {
  const { logout } = useAuth0();

  return (
    <button
      onClick={() => logout({ returnTo: window.location.origin })}
      style={{
        backgroundColor: '#E94E77',
        color: 'white',
        padding: '12px 24px',
        fontSize: '16px',
        borderRadius: '8px',
        border: 'none',
        cursor: 'pointer',
        boxShadow: '0 4px 6px rgba(0,0,0,0.1)',
        transition: 'background-color 0.3s ease',
      }}
      onMouseOver={e => (e.currentTarget.style.backgroundColor = '#D43962')}
      onMouseOut={e => (e.currentTarget.style.backgroundColor = '#E94E77')}
    >
      Log Out
    </button>
  );
};

export default LogoutButton;
