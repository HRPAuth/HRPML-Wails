import { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Button,
  Paper,
  Switch,
  Chip,
} from '@mui/material';
import {
  Add,
  Delete,
  FolderOpen,
} from '@mui/icons-material';

interface Mod {
  id: string;
  name: string;
  version: string;
  mcVersion: string;
  loaderType: 'fabric' | 'forge' | 'universal';
  enabled: boolean;
  filePath?: string;
}

interface ModManagerProps {
  onOpenModsFolder: () => void;
}

const STORAGE_KEY = 'installed_mods';

export default function ModManager({ onOpenModsFolder }: ModManagerProps) {
  const [mods, setMods] = useState<Mod[]>([]);

  useEffect(() => {
    loadMods();
  }, []);

  const loadMods = () => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        setMods(JSON.parse(stored));
      }
    } catch (err) {
      console.error('Failed to load mods:', err);
    }
  };

  const saveMods = (newMods: Mod[]) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(newMods));
      setMods(newMods);
    } catch (err) {
      console.error('Failed to save mods:', err);
    }
  };

  const handleToggleMod = (modId: string) => {
    const updated = mods.map(mod =>
      mod.id === modId ? { ...mod, enabled: !mod.enabled } : mod
    );
    saveMods(updated);
  };

  const handleDeleteMod = (modId: string) => {
    const updated = mods.filter(mod => mod.id !== modId);
    saveMods(updated);
  };

  const handleAddMod = () => {
    // Create a sample mod entry (in real app, this would open a file picker)
    const newMod: Mod = {
      id: Date.now().toString(),
      name: `Sample Mod ${mods.length + 1}`,
      version: '1.0.0',
      mcVersion: '1.20.4',
      loaderType: 'fabric',
      enabled: true,
    };
    saveMods([...mods, newMod]);
  };

  const getLoaderColor = (loader: string) => {
    switch (loader) {
      case 'fabric':
        return '#FF6B35';
      case 'forge':
        return '#8B5CF6';
      default:
        return '#9CA3AF';
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button
            variant="outlined"
            startIcon={<FolderOpen />}
            onClick={onOpenModsFolder}
            sx={{ borderColor: '#4a5568', color: '#9ca3af' }}
          >
            Open Mods Folder
          </Button>
          <Button
            variant="contained"
            startIcon={<Add />}
            onClick={handleAddMod}
            sx={{ backgroundColor: '#4CAF50' }}
          >
            Add Mod
          </Button>
        </Box>
      </Box>

      <Typography variant="body2" sx={{ color: '#9ca3af', mb: 2 }}>
        Enable or disable mods. Mods must be placed in the mods folder to work.
      </Typography>

      {mods.length === 0 ? (
        <Paper
          sx={{
            p: 4,
            textAlign: 'center',
            backgroundColor: '#2d3339',
            color: '#9ca3af',
          }}
        >
          <Typography variant="h6" sx={{ mb: 1 }}>
            No Mods
          </Typography>
          <Typography variant="body2" sx={{ mb: 2 }}>
            Click "Add Mod" to add mods or place .jar files in the mods folder
          </Typography>
        </Paper>
      ) : (
        <List sx={{ backgroundColor: '#282c34', borderRadius: 1 }}>
          {mods.map((mod) => (
            <ListItem
              key={mod.id}
              sx={{
                '&:hover': { backgroundColor: '#323842' },
              }}
            >
              <ListItemText
                primary={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <Typography sx={{ color: 'white' }}>{mod.name}</Typography>
                    <Chip
                      size="small"
                      label={mod.version}
                      sx={{
                        height: 20,
                        fontSize: '0.7rem',
                        backgroundColor: '#374151',
                        color: '#9ca3af',
                      }}
                    />
                  </Box>
                }
                secondary={
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5 }}>
                    <Chip
                      size="small"
                      label={`MC ${mod.mcVersion}`}
                      sx={{
                        height: 20,
                        fontSize: '0.7rem',
                        backgroundColor: '#374151',
                        color: '#9ca3af',
                      }}
                    />
                    <Chip
                      size="small"
                      label={mod.loaderType}
                      sx={{
                        height: 20,
                        fontSize: '0.7rem',
                        backgroundColor: getLoaderColor(mod.loaderType),
                        color: 'white',
                      }}
                    />
                  </Box>
                }
              />
              <ListItemSecondaryAction>
                <Switch
                  edge="end"
                  checked={mod.enabled}
                  onChange={() => handleToggleMod(mod.id)}
                  sx={{
                    '& .MuiSwitch-switchBase.Mui-checked': {
                      color: '#4CAF50',
                    },
                    '& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track': {
                      backgroundColor: '#4CAF50',
                    },
                  }}
                />
                <IconButton
                  edge="end"
                  onClick={() => handleDeleteMod(mod.id)}
                  sx={{ color: '#ef4444', ml: 1 }}
                >
                  <Delete />
                </IconButton>
              </ListItemSecondaryAction>
            </ListItem>
          ))}
        </List>
      )}
    </Box>
  );
}
