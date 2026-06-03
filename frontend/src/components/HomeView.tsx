import { useState } from 'react';
import { Box, Button, Typography, Card, CardContent, CircularProgress, Alert } from '@mui/material';
import { PlayArrow, Refresh } from '@mui/icons-material';
import { CMCLAPI } from '../services/shellApi';

interface HomeViewProps {
  selectedVersion: string | null;
  onVersionSelect: (version: string) => void;
}

export default function HomeView({ selectedVersion, onVersionSelect }: HomeViewProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [versions, setVersions] = useState<string[]>([]);

  const loadVersions = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await CMCLAPI.listVersions();
      setVersions(result);
    } catch {
      setError('Failed to load versions');
    }
    setIsLoading(false);
  };

  const handlePlay = async () => {
    if (!selectedVersion) {
      setError('Please select a version first');
      return;
    }

    setIsPlaying(true);
    setError(null);
    try {
      await CMCLAPI.selectVersion(selectedVersion);
      await CMCLAPI.startGame();
    } catch {
      setError('Failed to start game');
    }
    setIsPlaying(false);
  };

  return (
    <Box sx={{ p: 4 }}>
      <Card>
        <CardContent sx={{ textAlign: 'center' }}>
          <Typography variant="h3" sx={{ mb: 2, fontWeight: 'bold', color: '#4CAF50' }}>
            Samuel Client
          </Typography>
          <Typography variant="body1" sx={{ mb: 4 }}>
            test
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 4 }}>
              {error}
            </Alert>
          )}

          <Box sx={{ mb: 4 }}>
            <Typography variant="subtitle1" sx={{ mb: 2 }}>
              Selected Version: {selectedVersion || 'None'}
            </Typography>
            
            <Button
              variant="outlined"
              startIcon={<Refresh />}
              onClick={loadVersions}
              disabled={isLoading}
              sx={{ mr: 2 }}
            >
              {isLoading ? 'Loading...' : 'Refresh Versions'}
            </Button>

            {versions.length > 0 && (
              <Box sx={{ mt: 2, maxHeight: 200, overflowY: 'auto' }}>
                {versions.map((version, index) => (
                  <Button
                    key={index}
                    variant={selectedVersion === version ? 'contained' : 'text'}
                    onClick={() => onVersionSelect(version)}
                    sx={{ display: 'block', width: '100%', textAlign: 'left' }}
                  >
                    {version}
                  </Button>
                ))}
              </Box>
            )}
          </Box>

          <Button
            variant="contained"
            size="large"
            startIcon={isPlaying ? <CircularProgress size={20} /> : <PlayArrow />}
            onClick={handlePlay}
            disabled={isPlaying || !selectedVersion}
            sx={{
              backgroundColor: '#4CAF50',
              '&:hover': { backgroundColor: '#45a049' },
              px: 8,
              py: 2,
              fontSize: '1.2rem',
              fontWeight: 'bold'
            }}
          >
            {isPlaying ? 'Launching...' : 'Play'}
          </Button>
        </CardContent>
      </Card>
    </Box>
  );
}
