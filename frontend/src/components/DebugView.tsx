
import { useState } from 'react';
import { Box, Typography, Button, Card, CardContent, CircularProgress, Paper, Divider } from '@mui/material';
import { TestHealthCheck, TestPing } from '../../wailsjs/go/main/App';

interface ApiResponse {
  data: any;
  success: boolean;
  error?: string;
  time?: number;
}

export default function DebugView() {
  const [healthCheck, setHealthCheck] = useState<ApiResponse | null>(null);
  const [pingResult, setPingResult] = useState<ApiResponse | null>(null);
  const [loading, setLoading] = useState<{ health: boolean; ping: boolean }>({ health: false, ping: false });

  const apiBase = 'http://localhost:34501';

  const testHealthCheck = async () => {
    setLoading(prev => ({ ...prev, health: true }));
    try {
      const response = await TestHealthCheck();
      setHealthCheck({
        data: response.data,
        success: response.success,
        error: response.error,
        time: Number(response.time),
      });
    } catch (error) {
      setHealthCheck({
        data: null,
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      });
    } finally {
      setLoading(prev => ({ ...prev, health: false }));
    }
  };

  const testPing = async () => {
    setLoading(prev => ({ ...prev, ping: true }));
    try {
      const response = await TestPing();
      setPingResult({
        data: response.data,
        success: response.success,
        error: response.error,
        time: Number(response.time),
      });
    } catch (error) {
      setPingResult({
        data: null,
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error',
      });
    } finally {
      setLoading(prev => ({ ...prev, ping: false }));
    }
  };

  return (
    <Box sx={{ maxWidth: 900 }}>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Debug & API Testing
      </Typography>

      <Card sx={{ mb: 3, backgroundColor: '#212529' }}>
        <CardContent>
          <Typography variant="h5" gutterBottom sx={{ color: '#fff' }}>
            Backend API Index (Health Check)
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
            <Button
              variant="contained"
              onClick={testHealthCheck}
              disabled={loading.health}
              sx={{ backgroundColor: '#4CAF50', '&:hover': { backgroundColor: '#45a049' } }}
            >
              {loading.health ? <CircularProgress size={20} /> : 'Test Health Check'}
            </Button>
          </Box>
          {healthCheck && (
            <Paper sx={{ p: 2, backgroundColor: '#2d3137' }}>
              <Typography variant="subtitle1" sx={{ color: healthCheck.success ? '#4CAF50' : '#f44336' }}>
                {healthCheck.success ? 'Success' : 'Failed'}
                {healthCheck.time && ` (${healthCheck.time}ms)`}
              </Typography>
              {healthCheck.error && (
                <Typography variant="body2" sx={{ color: '#f44336', mt: 1 }}>
                  {healthCheck.error}
                </Typography>
              )}
              {healthCheck.data && (
                <Typography component="pre" sx={{ color: '#fff', mt: 2, whiteSpace: 'pre-wrap' }}>
                  {JSON.stringify(healthCheck.data, null, 2)}
                </Typography>
              )}
            </Paper>
          )}
        </CardContent>
      </Card>

      <Divider sx={{ my: 2, borderColor: '#3d4044' }} />

      <Card sx={{ mb: 3, backgroundColor: '#212529' }}>
        <CardContent>
          <Typography variant="h5" gutterBottom sx={{ color: '#fff' }}>
            Ping-Pong Test
          </Typography>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 2 }}>
            <Button
              variant="contained"
              onClick={testPing}
              disabled={loading.ping}
              sx={{ backgroundColor: '#2196F3', '&:hover': { backgroundColor: '#1976d2' } }}
            >
              {loading.ping ? <CircularProgress size={20} /> : 'Test Ping'}
            </Button>
          </Box>
          {pingResult && (
            <Paper sx={{ p: 2, backgroundColor: '#2d3137' }}>
              <Typography variant="subtitle1" sx={{ color: pingResult.success ? '#4CAF50' : '#f44336' }}>
                {pingResult.success ? 'Success' : 'Failed'}
                {pingResult.time && ` (${pingResult.time}ms)`}
              </Typography>
              {pingResult.error && (
                <Typography variant="body2" sx={{ color: '#f44336', mt: 1 }}>
                  {pingResult.error}
                </Typography>
              )}
              {pingResult.data && (
                <Typography component="pre" sx={{ color: '#fff', mt: 2, whiteSpace: 'pre-wrap' }}>
                  {JSON.stringify(pingResult.data, null, 2)}
                </Typography>
              )}
            </Paper>
          )}
        </CardContent>
      </Card>

      <Card sx={{ backgroundColor: '#212529' }}>
        <CardContent>
          <Typography variant="h5" gutterBottom sx={{ color: '#fff' }}>
            API Base URL
          </Typography>
          <Typography component="code" sx={{ color: '#81c784', fontSize: '1rem' }}>
            {apiBase}
          </Typography>
        </CardContent>
      </Card>
    </Box>
  );
}
