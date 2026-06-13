import { useState, useEffect } from 'react';
import { Box, Typography, Card, CardContent, Chip, Button, TextField, Dialog, DialogTitle, DialogContent, DialogActions, IconButton, Tooltip } from '@mui/material';
import { Add, Delete, Star, Edit, Check, X } from '@mui/icons-material';

interface JavaInstallation {
  id: number;
  path: string;
  friendly_name: string;
  version: string;
  is_default: boolean;
  created_at: string;
  updated_at: string;
}

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
  const [javaInstallations, setJavaInstallations] = useState<JavaInstallation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [openDialog, setOpenDialog] = useState(false);
  const [editingJava, setEditingJava] = useState<JavaInstallation | null>(null);
  const [newJavaPath, setNewJavaPath] = useState('');
  const [newJavaName, setNewJavaName] = useState('');
  const [dialogError, setDialogError] = useState('');

  useEffect(() => {
    fetchConfig();
    fetchJavaInstallations();
  }, []);

  const fetchConfig = async () => {
    try {
      const response = await fetch('http://localhost:34501/api/v1/launcher/config');
      const data = await response.json();
      if (data.success && data.data) {
        setConfig(data.data);
      }
    } catch (err) {
      console.error('Failed to fetch config:', err);
    }
  };

  const fetchJavaInstallations = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await fetch('http://localhost:34501/api/v1/db/java');
      const data = await response.json();
      if (Array.isArray(data)) {
        setJavaInstallations(data);
      } else if (data.error) {
        setError(data.error);
      }
    } catch (err) {
      setError('Failed to fetch Java installations');
    } finally {
      setLoading(false);
    }
  };

  const handleAddJava = async () => {
    setDialogError('');
    if (!newJavaPath.trim()) {
      setDialogError('Java path is required');
      return;
    }

    try {
      const response = await fetch('http://localhost:34501/api/v1/db/java', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          path: newJavaPath.trim(),
          friendly_name: newJavaName.trim()
        })
      });

      const data = await response.json();
      if (response.ok) {
        fetchJavaInstallations();
        setOpenDialog(false);
        setNewJavaPath('');
        setNewJavaName('');
      } else {
        setDialogError(data.error || 'Failed to add Java installation');
      }
    } catch (err) {
      setDialogError('Failed to add Java installation');
    }
  };

  const handleUpdateJava = async () => {
    setDialogError('');
    if (!editingJava) return;

    try {
      const response = await fetch(`http://localhost:34501/api/v1/db/java/${editingJava.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          path: newJavaPath.trim() || editingJava.path,
          friendly_name: newJavaName.trim()
        })
      });

      const data = await response.json();
      if (response.ok) {
        fetchJavaInstallations();
        setOpenDialog(false);
        setEditingJava(null);
        setNewJavaPath('');
        setNewJavaName('');
      } else {
        setDialogError(data.error || 'Failed to update Java installation');
      }
    } catch (err) {
      setDialogError('Failed to update Java installation');
    }
  };

  const handleDeleteJava = async (id: number) => {
    if (!confirm('Are you sure you want to delete this Java installation?')) return;

    try {
      const response = await fetch(`http://localhost:34501/api/v1/db/java/${id}`, {
        method: 'DELETE'
      });

      if (response.ok) {
        fetchJavaInstallations();
      } else {
        const data = await response.json();
        setError(data.error || 'Failed to delete Java installation');
      }
    } catch (err) {
      setError('Failed to delete Java installation');
    }
  };

  const handleSetDefault = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:34501/api/v1/db/java/${id}/default`, {
        method: 'POST'
      });

      if (response.ok) {
        fetchJavaInstallations();
      } else {
        const data = await response.json();
        setError(data.error || 'Failed to set default Java');
      }
    } catch (err) {
      setError('Failed to set default Java');
    }
  };

  const openAddDialog = () => {
    setEditingJava(null);
    setNewJavaPath('');
    setNewJavaName('');
    setDialogError('');
    setOpenDialog(true);
  };

  const openEditDialog = (java: JavaInstallation) => {
    setEditingJava(java);
    setNewJavaPath(java.path);
    setNewJavaName(java.friendly_name);
    setDialogError('');
    setOpenDialog(true);
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

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom sx={{ color: 'white' }}>
        Settings
      </Typography>

      {error && (
        <Typography sx={{ color: '#ef4444', marginBottom: 2 }}>{error}</Typography>
      )}

      <Card sx={{ backgroundColor: '#1f2937', marginBottom: 2 }}>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 2 }}>
            <Typography variant="h6" sx={{ color: 'white' }}>
              Java Installations
            </Typography>
            <Button
              variant="contained"
              color="primary"
              startIcon={<Add />}
              onClick={openAddDialog}
            >
              Add Java
            </Button>
          </Box>

          {javaInstallations.length === 0 ? (
            <Typography sx={{ color: '#9ca3af' }}>No Java installations found. Click "Add Java" to add one.</Typography>
          ) : (
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
              {javaInstallations.map((java) => (
                <Box
                  key={java.id}
                  sx={{
                    backgroundColor: '#374151',
                    padding: 2,
                    borderRadius: 1,
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 1
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Tooltip title={java.is_default ? 'Default Java' : 'Set as default'}>
                      <IconButton
                        onClick={() => handleSetDefault(java.id)}
                        sx={{ color: java.is_default ? '#fbbf24' : '#9ca3af' }}
                      >
                        <Star />
                      </IconButton>
                    </Tooltip>
                    <div style={{ flex: 1 }}>
                      <Typography sx={{ color: 'white', fontWeight: 'bold' }}>
                        {java.friendly_name || `Java ${java.version}`}
                      </Typography>
                      <Typography sx={{ color: '#9ca3af', fontSize: '0.875rem' }}>
                        {java.version}
                      </Typography>
                    </div>
                    <Box sx={{ display: 'flex', gap: 1 }}>
                      <Tooltip title="Edit">
                        <IconButton onClick={() => openEditDialog(java)} sx={{ color: '#9ca3af' }}>
                          <Edit />
                        </IconButton>
                      </Tooltip>
                      <Tooltip title="Delete">
                        <IconButton onClick={() => handleDeleteJava(java.id)} sx={{ color: '#ef4444' }}>
                          <Delete />
                        </IconButton>
                      </Tooltip>
                    </Box>
                  </Box>
                  <Chip
                    label={java.path}
                    sx={{
                      backgroundColor: '#1f2937',
                      color: '#d1d5db',
                      fontSize: '0.75rem',
                      maxWidth: '100%',
                      wordBreak: 'break-all'
                    }}
                  />
                </Box>
              ))}
            </Box>
          )}
        </CardContent>
      </Card>

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

      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="md">
        <DialogTitle sx={{ backgroundColor: '#1f2937', color: 'white' }}>
          {editingJava ? 'Edit Java Installation' : 'Add Java Installation'}
        </DialogTitle>
        <DialogContent sx={{ backgroundColor: '#1f2937', color: 'white' }}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, marginTop: 2 }}>
            <TextField
              label="Java Path"
              variant="outlined"
              fullWidth
              value={newJavaPath}
              onChange={(e) => setNewJavaPath(e.target.value)}
              placeholder="e.g., /usr/bin/java or C:\Program Files\Java\jdk17\bin\java.exe"
              sx={{
                input: { color: 'white' },
                label: { color: '#9ca3af' },
                '& .MuiOutlinedInput-root': {
                  '& fieldset': { borderColor: '#4b5563' },
                  '&:hover fieldset': { borderColor: '#60a5fa' },
                  '&.Mui-focused fieldset': { borderColor: '#3b82f6' }
                }
              }}
            />
            <TextField
              label="Friendly Name (Optional)"
              variant="outlined"
              fullWidth
              value={newJavaName}
              onChange={(e) => setNewJavaName(e.target.value)}
              placeholder="e.g., Java 17 (Forge)"
              sx={{
                input: { color: 'white' },
                label: { color: '#9ca3af' },
                '& .MuiOutlinedInput-root': {
                  '& fieldset': { borderColor: '#4b5563' },
                  '&:hover fieldset': { borderColor: '#60a5fa' },
                  '&.Mui-focused fieldset': { borderColor: '#3b82f6' }
                }
              }}
            />
            {dialogError && (
              <Typography sx={{ color: '#ef4444', fontSize: '0.875rem' }}>{dialogError}</Typography>
            )}
          </Box>
        </DialogContent>
        <DialogActions sx={{ backgroundColor: '#1f2937' }}>
          <Button
            onClick={() => {
              setOpenDialog(false);
              setEditingJava(null);
              setNewJavaPath('');
              setNewJavaName('');
              setDialogError('');
            }}
            sx={{ color: '#9ca3af' }}
            startIcon={<X />}
          >
            Cancel
          </Button>
          <Button
            onClick={editingJava ? handleUpdateJava : handleAddJava}
            variant="contained"
            color="primary"
            startIcon={<Check />}
          >
            {editingJava ? 'Update' : 'Add'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}