
import { useState, useEffect } from 'react';
import { Box, Typography, Card, CardContent, Chip } from '@mui/material';

interface LauncherConfig {
  java_path: string;
  game_dir: string;
  max_memory: number;
  min_memory: number;
  game_width: number;
  game_height: number;
}

export default function SettingsView() {
  const [config, setConfig] = useState<LauncherConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch('http://localhost:34501/api/v1/launcher/config');
      const data = await response.json();
      if (data.success && data.data) {
        setConfig(data.data);
      } else {
        setError(data.error || 'Failed to fetch configuration');
      }
    } catch (err) {
      setError('Failed to fetch configuration');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Box>
        <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
          Settings
        </Typography>
        <Typography sx={{ color: '#9ca3af' }}>Loading...</Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Box>
        <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
          Settings
        </Typography>
        <Typography sx={{ color: '#ef4444' }}>{error}</Typography>
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Settings
      </Typography>
      
      <Card sx={{ backgroundColor: '#1f2937', marginBottom: 2 }}>
        <CardContent>
          <Typography variant="h6" sx={{ color: 'white', marginBottom: 2 }}>
            Path Configuration
          </Typography>
          
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
            <div>
              <Typography sx={{ color: '#9ca3af', fontSize: '0.875rem' }}>Java Path</Typography>
              <Chip
                label={config?.java_path || 'Not configured'}
                sx={{ 
                  backgroundColor: '#374151', 
                  color: 'white',
                  marginTop: 1,
                  maxWidth: '100%',
                  wordBreak: 'break-all'
                }}
              />
            </div>
            
            <div>
              <Typography sx={{ color: '#9ca3af', fontSize: '0.875rem' }}>Game Directory</Typography>
              <Chip
                label={config?.game_dir || 'Not configured'}
                sx={{ 
                  backgroundColor: '#374151', 
                  color: 'white',
                  marginTop: 1,
                  maxWidth: '100%',
                  wordBreak: 'break-all'
                }}
              />
            </div>
          </Box>
        </CardContent>
      </Card>

      <Card sx={{ backgroundColor: '#1f2937' }}>
        <CardContent>
          <Typography variant="h6" sx={{ color: 'white', marginBottom: 2 }}>
            Memory Settings
          </Typography>
          
          <Box sx={{ display: 'flex', gap: 4 }}>
            <div>
              <Typography sx={{ color: '#9ca3af', fontSize: '0.875rem' }}>Minimum Memory</Typography>
              <Typography sx={{ color: 'white' }}>{config?.min_memory || 0} MB</Typography>
            </div>
            <div>
              <Typography sx={{ color: '#9ca3af', fontSize: '0.875rem' }}>Maximum Memory</Typography>
              <Typography sx={{ color: 'white' }}>{config?.max_memory || 0} MB</Typography>
            </div>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
