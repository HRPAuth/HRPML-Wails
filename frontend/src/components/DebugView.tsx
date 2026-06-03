import { useState } from 'react';
import { Box, Button, Typography, Card, CardContent, CircularProgress, Alert } from '@mui/material';
import { Refresh } from '@mui/icons-material';
import { CMCL_CONFIG } from '../variables';
import { ShellAPI } from '../services/shellApi';

export default function DebugView() {
  const [isLoading, setIsLoading] = useState(false);
  const [pingResult, setPingResult] = useState<{ success: boolean; responseTime: number; error?: string } | null>(null);

  const handlePing = async () => {
    setIsLoading(true);
    const result = await ShellAPI.ping();
    setPingResult(result);
    setIsLoading(false);
  };

  const fullBackendUrl = new URL(CMCL_CONFIG.shellApiUrl, window.location.origin).href;

  return (
    <Box sx={{ p: 4 }}>
      <Card>
        <CardContent>
          <Typography variant="h4" sx={{ mb: 4, fontWeight: 'bold', color: '#4CAF50' }}>
            Debug
          </Typography>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Backend URL
            </Typography>
            <Typography variant="body1" sx={{ fontFamily: 'monospace', backgroundColor: '#212529', p: 2, borderRadius: 1 }}>
              {fullBackendUrl}
            </Typography>
          </Box>

          <Box sx={{ mb: 4 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Ping-Pong Test
            </Typography>
            <Button
              variant="contained"
              startIcon={isLoading ? <CircularProgress size={20} /> : <Refresh />}
              onClick={handlePing}
              disabled={isLoading}
              sx={{ mb: 2, backgroundColor: '#4CAF50', '&:hover': { backgroundColor: '#45a049' } }}
            >
              {isLoading ? 'Pinging...' : 'Ping'}
            </Button>

            {pingResult && (
              <Alert severity={pingResult.success ? 'success' : 'error'}>
                {pingResult.success
                  ? `Pong! Response time: ${pingResult.responseTime}ms`
                  : `Ping failed: ${pingResult.error}`}
              </Alert>
            )}
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
