import { useState } from 'react';
import { Box, Button, Typography, Card, CardContent, TextField, Checkbox, FormControlLabel, Alert, CircularProgress, List, ListItem, ListItemText, ListItemSecondaryAction, IconButton } from '@mui/material';
import { Download, Delete, Info } from '@mui/icons-material';
import { CMCLAPI, InstallOptions } from '../services/shellApi';

export default function VersionManager() {
  const [versions, setVersions] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [installVersion, setInstallVersion] = useState('');
  const [installOptions, setInstallOptions] = useState<InstallOptions>({
    fabric: false,
    forge: false,
    optifine: false,
    quilt: false,
    select: false
  });
  const [infoVersion, setInfoVersion] = useState<string | null>(null);
  const [versionInfo, setVersionInfo] = useState('');

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

  const handleInstall = async () => {
    if (!installVersion.trim()) {
      setError('Please enter a version to install');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.installVersion(installVersion, installOptions);
      await loadVersions();
      setInstallVersion('');
    } catch {
      setError('Failed to install version');
    }
    setIsLoading(false);
  };

  const handleUninstall = async (version: string) => {
    if (!window.confirm(`Are you sure you want to uninstall ${version}?`)) return;

    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.uninstallVersion(version);
      await loadVersions();
    } catch {
      setError('Failed to uninstall version');
    }
    setIsLoading(false);
  };

  const handleGetInfo = async (version: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await CMCLAPI.getVersionInfo(version);
      setInfoVersion(version);
      setVersionInfo(result.stdout || result.stderr || 'No info available');
    } catch {
      setError('Failed to get version info');
    }
    setIsLoading(false);
  };

  return (
    <Box sx={{ p: 4 }}>
      <Card>
        <CardContent>
          <Typography variant="h4" sx={{ mb: 4 }}>Version Management</Typography>

          {error && <Alert severity="error" sx={{ mb: 4 }}>{error}</Alert>}

          <Box sx={{ mb: 6 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>Install Version</Typography>
            <TextField
              label="Version (e.g., 1.19, 1.18.2)"
              value={installVersion}
              onChange={(e) => setInstallVersion(e.target.value)}
              sx={{ mr: 2, width: 200 }}
            />
            <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
              <FormControlLabel
                control={<Checkbox checked={installOptions.fabric} onChange={(e) => setInstallOptions({ ...installOptions, fabric: e.target.checked })} />}
                label="Fabric"
              />
              <FormControlLabel
                control={<Checkbox checked={installOptions.forge} onChange={(e) => setInstallOptions({ ...installOptions, forge: e.target.checked })} />}
                label="Forge"
              />
              <FormControlLabel
                control={<Checkbox checked={installOptions.optifine} onChange={(e) => setInstallOptions({ ...installOptions, optifine: e.target.checked })} />}
                label="OptiFine"
              />
              <FormControlLabel
                control={<Checkbox checked={installOptions.quilt} onChange={(e) => setInstallOptions({ ...installOptions, quilt: e.target.checked })} />}
                label="Quilt"
              />
              <FormControlLabel
                control={<Checkbox checked={installOptions.select} onChange={(e) => setInstallOptions({ ...installOptions, select: e.target.checked })} />}
                label="Select after install"
              />
            </Box>
            <Button
              variant="contained"
              startIcon={<Download />}
              onClick={handleInstall}
              disabled={isLoading}
              sx={{ mt: 2 }}
            >
              {isLoading ? <CircularProgress size={20} /> : 'Install'}
            </Button>
          </Box>

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>Installed Versions</Typography>
            <Button onClick={loadVersions} disabled={isLoading} sx={{ mb: 2 }}>
              {isLoading ? 'Loading...' : 'Refresh List'}
            </Button>
            <List>
              {versions.map((version, index) => (
                <ListItem key={index}>
                  <ListItemText primary={version} />
                  <ListItemSecondaryAction>
                    <IconButton onClick={() => handleGetInfo(version)}>
                      <Info />
                    </IconButton>
                    <IconButton onClick={() => handleUninstall(version)} color="error">
                      <Delete />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          </Box>

          {infoVersion && (
            <Box sx={{ mt: 4, p: 3, backgroundColor: '#f5f5f5', borderRadius: 2 }}>
              <Typography variant="h6">Info for {infoVersion}</Typography>
              <Box sx={{ mt: 2, whiteSpace: 'pre-wrap' }} component="pre">{versionInfo}</Box>
              <Button onClick={() => setInfoVersion(null)} sx={{ mt: 2 }}>Close</Button>
            </Box>
          )}
        </CardContent>
      </Card>
    </Box>
  );
}
