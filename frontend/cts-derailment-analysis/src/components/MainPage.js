import React, { useEffect, useState } from 'react';
import { fetchWithAuth } from '../utils/auth';

const MainPage = () => {
  const [userData, setUserData] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetchWithAuth('http://localhost:8080/api/protected');
        if (response.ok) {
          const data = await response.json();
          setUserData(data);
        } else {
          console.error('Failed to fetch protected data');
        }
      } catch (err) {
        console.error('Error fetching data:', err);
      }
    };

    fetchData();
  }, []);

  return (
    <div className="container">
      <h1>Welcome to the Main Page</h1>
      {userData ? (
        <div>
          <p>User ID: {userData.user_id}</p>
          <p>Email: {userData.email}</p>
          <p>Role: {userData.role}</p>
        </div>
      ) : (
        <p>Loading...</p>
      )}
    </div>
  );
};

export default MainPage;
