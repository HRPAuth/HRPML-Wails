import { useState } from 'react';
import { Box, Button, Typography, Card, CardContent, TextField, Alert, CircularProgress } from '@mui/material';
import { Search, Download } from '@mui/icons-material';
import { CMCLAPI, ShellResult } from '../services/shellApi';

export default function ModManager() {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [modUrl, setModUrl] = useState('');
  const [modName, setModName] = useState('');
  const [modInfo, setModInfo] = useState('');
  const [modpackUrl, setModpackUrl] = useState('');
  const [modpackName, setModpackName] = useState('');
  const [modpackInfo, setModpackInfo] = useState('');

  const handleInstallMod = async () => {
    if (!modUrl.trim()) {
      setError('Please enter a mod URL');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.installMod(modUrl);
      setModUrl('');
    } catch {
      setError('Failed to install mod');
    }
    setIsLoading(false);
  };

  const handleSearchMod = async () => {
    if (!modName.trim()) {
      setError('Please enter a mod name');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      const result: ShellResult = await CMCLAPI.searchMod(modName);
      setModInfo(result.stdout || result.stderr || 'No info found');
    } catch {
      setError('Failed to search mod');
    }
    setIsLoading(false);
  };

  const handleInstallModpack = async () => {
    if (!modpackUrl.trim()) {
      setError('Please enter a modpack URL');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.installModpack(modpackUrl);
      setModpackUrl('');
    } catch {
      setError('Failed to install modpack');
    }
    setIsLoading(false);
  };

  const handleSearchModpack = async () => {
    if (!modpackName.trim()) {
      setError('Please enter a modpack name');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      const result: ShellResult = await CMCLAPI.searchModpack(modpackName);
      setModpackInfo(result.stdout || result.stderr || 'No info found');
    } catch {
      setError('Failed to search modpack');
    }
    setIsLoading(false);
  };

  return (
    <Box sx={{ p: 4 }}>
      <Card>
        <CardContent>
          <Typography variant="h4" sx={{ mb: 4 }}>Mod Management</Typography>

          {error && <Alert severity="error" sx={{ mb: 4 }}>{error}</Alert>}

          <Box sx={{ mb: 6 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>Mods</Typography>
            
            <Box sx={{ mb: 4 }}>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Install Mod from URL</Typography>
              <TextField
                label="Mod URL"
                value={modUrl}
                onChange={(e) => setModUrl(e.target.value)}
                sx={{ mr: 2, width: 400 }}
              />
              <Button
                variant="contained"
                startIcon={<Download />}
                onClick={handleInstallMod}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Install'}
              </Button>
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Search Mod</Typography>
              <TextField
                label="Mod Name"
                value={modName}
                onChange={(e) => setModName(e.target.value)}
                sx={{ mr: 2, width: 250 }}
              />
              <Button
                variant="contained"
                startIcon={<Search />}
                onClick={handleSearchMod}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Search'}
              </Button>
              {modInfo && (
                <Box sx={{ mt: 2, p: 3, backgroundColor: '#f5f5f5', borderRadius: 2 }}>
                  <Typography variant="subtitle1">Mod Info:</Typography>
                  <Box sx={{ mt: 1, whiteSpace: 'pre-wrap' }} component="pre">{modInfo}</Box>
                  <Button onClick={() => setModInfo('')} sx={{ mt: 2 }}>Clear</Button>
                </Box>
              )}
            </Box>
          </Box>

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>Modpacks</Typography>
            
            <Box sx={{ mb: 4 }}>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Install Modpack from URL</Typography>
              <TextField
                label="Modpack URL"
                value={modpackUrl}
                onChange={(e) => setModpackUrl(e.target.value)}
                sx={{ mr: 2, width: 400 }}
              />
              <Button
                variant="contained"
                startIcon={<Download />}
                onClick={handleInstallModpack}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Install'}
              </Button>
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Search Modpack</Typography>
              <TextField
                label="Modpack Name"
                value={modpackName}
                onChange={(e) => setModpackName(e.target.value)}
                sx={{ mr: 2, width: 250 }}
              />
              <Button
                variant="contained"
                startIcon={<Search />}
                onClick={handleSearchModpack}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Search'}
              </Button>
              {modpackInfo && (
                <Box sx={{ mt: 2, p: 3, backgroundColor: '#f5f5f5', borderRadius: 2 }}>
                  <Typography variant="subtitle1">Modpack Info:</Typography>
                  <Box sx={{ mt: 1, whiteSpace: 'pre-wrap' }} component="pre">{modpackInfo}</Box>
                  <Button onClick={() => setModpackInfo('')} sx={{ mt: 2 }}>Clear</Button>
                </Box>
              )}
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
