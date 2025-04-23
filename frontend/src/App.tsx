import React, { useState } from 'react';

interface AccessCheckResponse {
  vpn: boolean;
  production: boolean;
  configTool: boolean;
  currentProfile: string;
  missingGroups: string[];
  profileArn?: string;
}

export default function App() {
  const [email, setEmail] = useState('');
  const [access, setAccess] = useState<AccessCheckResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleCheck = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    if (!email) {
      setError('Please enter your email');
      return;
    }

    // Validate email format
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      setError('Please enter a valid email address');
      return;
    }

    setLoading(true);
    
    try {
      const requestBody = {
        email: email,
        environment: 'Production'
      };

      console.log('Sending request:', JSON.stringify(requestBody, null, 2));

      const response = await fetch('http://localhost:8080/api/check-access', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody)
      });

      const responseText = await response.text();
      console.log('Received response:', responseText);

      if (!response.ok) {
        try {
          const errorData = JSON.parse(responseText);
          throw new Error(errorData.error || 'Unknown error occurred');
        } catch {
          throw new Error(responseText || 'Request failed');
        }
      }

      const data = JSON.parse(responseText) as AccessCheckResponse;
      setAccess(data);
    } catch (err) {
      console.error('API Error:', err);
      setError(err instanceof Error ? err.message : 'Failed to check access');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '20px', maxWidth: '600px', margin: '0 auto' }}>
      <h1>On Call Access Checker</h1>
      
      <form onSubmit={handleCheck} style={{ marginBottom: '20px' }}>
        <div style={{ marginBottom: '10px' }}>
          <label>Environment:</label>
          <select disabled style={{ display: 'block', width: '100%', padding: '8px' }}>
            <option>Production</option>
          </select>
        </div>
        
        <div style={{ marginBottom: '10px' }}>
          <label>Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={{ display: 'block', width: '100%', padding: '8px' }}
            required
          />
        </div>
        
        <button 
          type="submit" 
          disabled={loading}
          style={{ padding: '8px 16px', background: loading ? '#ccc' : '#0066cc', color: 'white', border: 'none' }}
        >
          {loading ? 'Checking...' : 'Check'}
        </button>
        
        {error && <div style={{ color: 'red', marginTop: '10px' }}>{error}</div>}
      </form>

      {access && (
        <div>
          <h2>Access Checklist</h2>
          <div style={{ marginBottom: '10px' }}>
            <span style={{
              display: 'inline-block',
              width: '12px',
              height: '12px',
              borderRadius: '50%',
              backgroundColor: access.vpn ? 'green' : 'red',
              marginRight: '8px'
            }}></span>
            VPN
          </div>
          <div style={{ marginBottom: '10px' }}>
            <span style={{
              display: 'inline-block',
              width: '12px',
              height: '12px',
              borderRadius: '50%',
              backgroundColor: access.production ? 'green' : 'red',
              marginRight: '8px'
            }}></span>
            Production Group
          </div>
          <div style={{ marginBottom: '10px' }}>
            <span style={{
              display: 'inline-block',
              width: '12px',
              height: '12px',
              borderRadius: '50%',
              backgroundColor: access.configTool ? 'green' : 'red',
              marginRight: '8px'
            }}></span>
            Config Tool
          </div>
          
          <p style={{ color: access.missingGroups.length ? 'red' : 'green' }}>
            {access.missingGroups.length 
              ? `Please go to the access management tool and request your ${access.missingGroups.join(', ')}`
              : 'You have all the required group accesses.'}
          </p>
          
          <h2>Tool Profile Checker</h2>
          <p style={{ color: access.currentProfile === 'prod' ? 'green' : 'orange' }}>
            {access.currentProfile === 'prod'
              ? 'You are currently using the correct AWS profile - prod'
              : 'You are currently using the dev profile. Once access is granted, please switch to prod (link)'}
          </p>
        </div>
      )}
    </div>
  );
}
