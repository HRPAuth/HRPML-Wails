import { useState } from 'react';
import { Box, Button, Typography, Card, CardContent, TextField, Alert, CircularProgress } from '@mui/material';
import { Settings as SettingsIcon, Info, Refresh } from '@mui/icons-material';
import { CMCLAPI, ShellAPI, ShellResult } from '../services/shellApi';
import { CMCL_CONFIG, LAUNCHER_CONFIG } from '../variables';

interface MemorySettings {
  maxMemory: string;
  minMemory: string;
}

export default function SettingsView() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [config, setConfig] = useState<string | null>(null);
  const [aboutInfo, setAboutInfo] = useState('');
  const [memorySettings, setMemorySettings] = useState<MemorySettings>({
    maxMemory: LAUNCHER_CONFIG.defaultMaxMemory,
    minMemory: LAUNCHER_CONFIG.defaultMinMemory
  });
  const [commandInput, setCommandInput] = useState('');
  const [commandOutput, setCommandOutput] = useState<ShellResult | null>(null);

  const loadConfig = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result: ShellResult = await CMCLAPI.getConfig();
      setConfig(result.stdout || '');
    } catch {
      setError('Failed to load config');
    }
    setIsLoading(false);
  };

  const handleSetMemory = async () => {
    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.setConfig('maxMemory', memorySettings.maxMemory);
      await CMCLAPI.setConfig('minMemory', memorySettings.minMemory);
    } catch {
      setError('Failed to set memory settings');
    }
    setIsLoading(false);
  };

  const handleCheckUpdates = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result: ShellResult = await CMCLAPI.checkForUpdates();
      setAboutInfo(result.stdout || result.stderr || 'No update info available');
    } catch {
      setError('Failed to check for updates');
    }
    setIsLoading(false);
  };

  const handleGetAbout = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result: ShellResult = await CMCLAPI.getAboutInfo();
      setAboutInfo(result.stdout || 'No info available');
    } catch {
      setError('Failed to get about info');
    }
    setIsLoading(false);
  };

  const handleExecuteCommand = async () => {
    if (!commandInput.trim()) return;
    setIsLoading(true);
    setError(null);
    try {
      const result: ShellResult = await ShellAPI.execute(commandInput);
      setCommandOutput(result);
    } catch {
      setError('Failed to execute command');
    }
    setIsLoading(false);
  };

  return (
    <Box sx={{ p: 4 }}>
      <Card>
        <CardContent>
          <Typography variant="h4" sx={{ mb: 4 }}>Settings</Typography>

          {error && <Alert severity="error" sx={{ mb: 4 }}>{error}</Alert>}

          <Box sx={{ mb: 6 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>Launcher Configuration</Typography>
            
            <Box sx={{ display: 'flex', gap: 4, mb: 4 }}>
              <TextField
                label="CMCL Path"
                value={CMCL_CONFIG.cmclPath}
                disabled
                sx={{ width: 400 }}
              />
              <TextField
                label="Shell API URL"
                value={CMCL_CONFIG.shellApiUrl}
                disabled
                sx={{ width: 300 }}
              />
            </Box>

            <Typography variant="subtitle1" sx={{ mb: 2 }}>Memory Settings (MB)</Typography>
            <Box sx={{ display: 'flex', gap: 4 }}>
              <TextField
                label="Max Memory"
                type="number"
                value={memorySettings.maxMemory}
                onChange={(e) => setMemorySettings({ ...memorySettings, maxMemory: e.target.value })}
                sx={{ width: 150 }}
              />
              <TextField
                label="Min Memory"
                type="number"
                value={memorySettings.minMemory}
                onChange={(e) => setMemorySettings({ ...memorySettings, minMemory: e.target.value })}
                sx={{ width: 150 }}
              />
              <Button
                variant="contained"
                onClick={handleSetMemory}
                disabled={isLoading}
                sx={{ mt: 1 }}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Save'}
              </Button>
            </Box>
          </Box>

          <Box sx={{ mb: 6 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>CMCL Configuration</Typography>
            <Button
              variant="outlined"
              startIcon={<Refresh />}
              onClick={loadConfig}
              disabled={isLoading}
              sx={{ mb: 2 }}
            >
              {isLoading ? 'Loading...' : 'Load Config'}
            </Button>
            {config && (
              <Box sx={{ p: 3, backgroundColor: '#f5f5f5', borderRadius: 2, maxHeight: 300, overflowY: 'auto' }}>
                <pre>{config}</pre>
              </Box>
            )}
          </Box>

          <Box sx={{ mb: 6 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>About & Updates</Typography>
            <Box sx={{ display: 'flex', gap: 2 }}>
              <Button
                variant="contained"
                startIcon={<Info />}
                onClick={handleGetAbout}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'About CMCL'}
              </Button>
              <Button
                variant="contained"
                startIcon={<SettingsIcon />}
                onClick={handleCheckUpdates}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Check for Updates'}
              </Button>
            </Box>
            {aboutInfo && (
              <Box sx={{ mt: 2, p: 3, backgroundColor: '#f5f5f5', borderRadius: 2 }}>
                <pre>{aboutInfo}</pre>
              </Box>
            )}
          </Box>

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>Execute Command</Typography>
            <Box sx={{ display: 'flex', gap: 2, mb: 2 }}>
              <TextField
                label="Command"
                value={commandInput}
                onChange={(e) => setCommandInput(e.target.value)}
                placeholder="e.g., cmcl --help"
                fullWidth
                onKeyPress={(e) => e.key === 'Enter' && handleExecuteCommand()}
              />
              <Button
                variant="contained"
                onClick={handleExecuteCommand}
                disabled={isLoading || !commandInput.trim()}
                sx={{ minWidth: 120 }}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Execute'}
              </Button>
            </Box>
            {commandOutput && (
              <Box sx={{ p: 3, backgroundColor: '#f5f5f5', borderRadius: 2 }}>
                <Typography variant="subtitle2" sx={{ mb: 1 }}>
                  Status: {commandOutput.success ? 'Success' : 'Error'}
                </Typography>
                {commandOutput.stdout && (
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" sx={{ color: '#2e7d32' }}>
                      Standard Output:
                    </Typography>
                    <Box sx={{ mt: 0.5, whiteSpace: 'pre-wrap', fontSize: '0.875rem' }} component="pre">
                      {commandOutput.stdout}
                    </Box>
                  </Box>
                )}
                {commandOutput.stderr && (
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" sx={{ color: '#d32f2f' }}>
                      Standard Error:
                    </Typography>
                    <Box sx={{ mt: 0.5, whiteSpace: 'pre-wrap', fontSize: '0.875rem' }} component="pre">
                      {commandOutput.stderr}
                    </Box>
                  </Box>
                )}
                {commandOutput.error && (
                  <Typography variant="subtitle2" sx={{ color: '#d32f2f' }}>
                    Error: {commandOutput.error}
                  </Typography>
                )}
              </Box>
            )}
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
