import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  TextField,
  Paper,
  Divider,
  Alert,
  Snackbar,
} from '@mui/material';
import {
  Home,
  FolderOpen,
  Save,
  Refresh,
} from '@mui/icons-material';

const STORAGE_KEY = 'launcher_settings';

interface LauncherSettings {
  javaPath: string;
  gamePath: string;
  maxMemory: number;
  minMemory: number;
  javaArgs: string;
  gameWidth: number;
  gameHeight: number;
  selectedLoader: string;
  loaderVersion: string;
}

const DEFAULT_SETTINGS: LauncherSettings = {
  javaPath: '',
  gamePath: '',
  maxMemory: 4096,
  minMemory: 1024,
  javaArgs: '',
  gameWidth: 854,
  gameHeight: 480,
  selectedLoader: 'none',
  loaderVersion: '',
};

interface SettingsViewProps {
  onSettingsChange: (settings: LauncherSettings) => void;
}

export default function SettingsView({ onSettingsChange }: SettingsViewProps) {
  const [settings, setSettings] = useState<LauncherSettings>(DEFAULT_SETTINGS);
  const [saved, setSaved] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    loadSettings();
  }, []);

  const loadSettings = () => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored);
        setSettings({ ...DEFAULT_SETTINGS, ...parsed });
      } else {
        // Try to detect default paths
        detectDefaultPaths();
      }
    } catch (err) {
      console.error('Failed to load settings:', err);
      detectDefaultPaths();
    }
  };

  const detectDefaultPaths = () => {
    let gamePath = '';
    const homeDir = process.env.HOME || process.env.USERPROFILE || '/';
    
    if (process.platform === 'win32') {
      gamePath = `${process.env.APPDATA}\\.minecraft`;
    } else if (process.platform === 'darwin') {
      gamePath = `${homeDir}/Library/Application Support/minecraft`;
    } else {
      gamePath = `${homeDir}/.minecraft`;
    }

    // Try to find Java
    let javaPath = 'java';
    const possiblePaths = [
      '/usr/bin/java',
      '/usr/local/bin/java',
      '/opt/java/bin/java',
      `${homeDir}/.jdks/current/bin/java`,
    ];

    for (const path of possiblePaths) {
      // In browser environment, we can't check file existence
      // Just use the first one as default suggestion
      javaPath = path;
      break;
    }

    setSettings(prev => ({
      ...prev,
      gamePath,
      javaPath,
    }));
  };

  const saveSettings = () => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
      onSettingsChange(settings);
      setSaved(true);
      setError('');
      setTimeout(() => setSaved(false), 3000);
    } catch (err) {
      setError('Failed to save settings');
    }
  };

  const handleChange = (key: keyof LauncherSettings, value: string | number) => {
    setSettings(prev => ({
      ...prev,
      [key]: value,
    }));
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ color: 'white' }}>
          Settings
        </Typography>
        <Box>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={detectDefaultPaths}
            sx={{ borderColor: '#6B7280', color: '#9CA3AF', mr: 2 }}
          >
            Detect Paths
          </Button>
          <Button
            variant="contained"
            startIcon={<Save />}
            onClick={saveSettings}
            sx={{ backgroundColor: '#4CAF50', '&:hover': { backgroundColor: '#45A049' } }}
          >
            Save
          </Button>
        </Box>
      </Box>

      <Paper sx={{ p: 4, backgroundColor: '#1F2937', borderRadius: 2 }}>
        {/* Java Settings */}
        <Box sx={{ mb: 4 }}>
          <Typography variant="h6" sx={{ color: '#4CAF50', mb: 3, display: 'flex', alignItems: 'center' }}>
            <FolderOpen sx={{ mr: 2 }} />
            Java Settings
          </Typography>
          
          <TextField
            fullWidth
            label="Java Path"
            variant="filled"
            value={settings.javaPath}
            onChange={(e) => handleChange('javaPath', e.target.value)}
            sx={{
              mb: 3,
              '& .MuiFilledInput-root': {
                backgroundColor: '#374151',
                '&:hover': { backgroundColor: '#4B5563' },
              },
              '& .MuiInputLabel-root': { color: '#9CA3AF' },
              '& .MuiInputBase-input': { color: 'white' },
            }}
            placeholder="Path to Java executable (e.g., /usr/bin/java)"
          />

          <Box sx={{ display: 'flex', gap: 3 }}>
            <TextField
              label="Minimum Memory (MB)"
              type="number"
              variant="filled"
              value={settings.minMemory}
              onChange={(e) => handleChange('minMemory', parseInt(e.target.value) || 1024)}
              sx={{
                flex: 1,
                '& .MuiFilledInput-root': {
                  backgroundColor: '#374151',
                  '&:hover': { backgroundColor: '#4B5563' },
                },
                '& .MuiInputLabel-root': { color: '#9CA3AF' },
                '& .MuiInputBase-input': { color: 'white' },
              }}
            />
            <TextField
              label="Maximum Memory (MB)"
              type="number"
              variant="filled"
              value={settings.maxMemory}
              onChange={(e) => handleChange('maxMemory', parseInt(e.target.value) || 4096)}
              sx={{
                flex: 1,
                '& .MuiFilledInput-root': {
                  backgroundColor: '#374151',
                  '&:hover': { backgroundColor: '#4B5563' },
                },
                '& .MuiInputLabel-root': { color: '#9CA3AF' },
                '& .MuiInputBase-input': { color: 'white' },
              }}
            />
          </Box>

          <TextField
            fullWidth
            label="Additional Java Arguments"
            variant="filled"
            value={settings.javaArgs}
            onChange={(e) => handleChange('javaArgs', e.target.value)}
            sx={{
              mt: 3,
              '& .MuiFilledInput-root': {
                backgroundColor: '#374151',
                '&:hover': { backgroundColor: '#4B5563' },
              },
              '& .MuiInputLabel-root': { color: '#9CA3AF' },
              '& .MuiInputBase-input': { color: 'white' },
            }}
            placeholder="-XX:+UseG1GC -XX:-UseAdaptiveSizePolicy"
          />
        </Box>

        <Divider sx={{ backgroundColor: '#4B5563', my: 4 }} />

        {/* Game Settings */}
        <Box sx={{ mb: 4 }}>
          <Typography variant="h6" sx={{ color: '#4CAF50', mb: 3, display: 'flex', alignItems: 'center' }}>
            <Home sx={{ mr: 2 }} />
            Game Settings
          </Typography>

          <TextField
            fullWidth
            label="Game Directory"
            variant="filled"
            value={settings.gamePath}
            onChange={(e) => handleChange('gamePath', e.target.value)}
            sx={{
              mb: 3,
              '& .MuiFilledInput-root': {
                backgroundColor: '#374151',
                '&:hover': { backgroundColor: '#4B5563' },
              },
              '& .MuiInputLabel-root': { color: '#9CA3AF' },
              '& .MuiInputBase-input': { color: 'white' },
            }}
            placeholder="Path to Minecraft game directory"
          />

          <Box sx={{ display: 'flex', gap: 3 }}>
            <TextField
              label="Game Width"
              type="number"
              variant="filled"
              value={settings.gameWidth}
              onChange={(e) => handleChange('gameWidth', parseInt(e.target.value) || 854)}
              sx={{
                flex: 1,
                '& .MuiFilledInput-root': {
                  backgroundColor: '#374151',
                  '&:hover': { backgroundColor: '#4B5563' },
                },
                '& .MuiInputLabel-root': { color: '#9CA3AF' },
                '& .MuiInputBase-input': { color: 'white' },
              }}
            />
            <TextField
              label="Game Height"
              type="number"
              variant="filled"
              value={settings.gameHeight}
              onChange={(e) => handleChange('gameHeight', parseInt(e.target.value) || 480)}
              sx={{
                flex: 1,
                '& .MuiFilledInput-root': {
                  backgroundColor: '#374151',
                  '&:hover': { backgroundColor: '#4B5563' },
                },
                '& .MuiInputLabel-root': { color: '#9CA3AF' },
                '& .MuiInputBase-input': { color: 'white' },
              }}
            />
          </Box>
        </Box>
      </Paper>

      {/* Error Alert */}
      {error && (
        <Alert severity="error" sx={{ mt: 3 }} onClose={() => setError('')}>
          {error}
        </Alert>
      )}

      {/* Save Snackbar */}
      <Snackbar
        open={saved}
        autoHideDuration={3000}
        message="Settings saved successfully!"
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      />
    </Box>
  );
}
