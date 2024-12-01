export const refreshAccessToken = async () => {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) return null;
  
    try {
      const response = await fetch('http://localhost:8080/api/refresh', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
  
      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        return data.access_token;
      } else {
        console.error('Failed to refresh token');
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        return null;
      }
    } catch (err) {
      console.error('Error refreshing token:', err);
      return null;
    }
  };
  
  export const fetchWithAuth = async (url, options = {}) => {
    let accessToken = localStorage.getItem('access_token');
    if (!accessToken) {
      accessToken = await refreshAccessToken();
      if (!accessToken) {
        window.location.href = '/';
        return null;
      }
    }
  
    const headers = options.headers || {};
    headers['Authorization'] = `${accessToken}`;
    options.headers = headers;
    console.log('Authorization Header:', headers['Authorization']);
    const response = await fetch(url, options);
  
    if (response.status === 401) {
      accessToken = await refreshAccessToken();
      if (!accessToken) {
        window.location.href = '/';
        return null;
      }
  
      headers['Authorization'] = `${accessToken}`;
      return fetch(url, options);
    }
  
    return response;
  };
  