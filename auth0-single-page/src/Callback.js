import React, { useEffect } from 'react';
import { useAuth0 } from '@auth0/auth0-react';

const Callback = () => {
  const { handleRedirectCallback } = useAuth0();

  useEffect(() => {
    handleRedirectCallback().then(() => {
      window.location.href = '/';
    });
  }, [handleRedirectCallback]);

  return <div>Loading...</div>;
};

export default Callback;