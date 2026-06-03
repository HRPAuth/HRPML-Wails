import { useState } from 'react';
import { Box, Button, Typography, Card, CardContent, TextField, Alert, CircularProgress, List, ListItem, ListItemButton, ListItemText, ListItemSecondaryAction } from '@mui/material';
import { Person, Login, Refresh } from '@mui/icons-material';
import { CMCLAPI } from '../services/shellApi';

export default function AccountManager() {
  const [accounts, setAccounts] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [offlineName, setOfflineName] = useState('');
  const [selectedAccount, setSelectedAccount] = useState<number | null>(null);

  const loadAccounts = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const result = await CMCLAPI.listAccounts();
      setAccounts(result);
    } catch {
      setError('Failed to load accounts');
    }
    setIsLoading(false);
  };

  const handleLoginOffline = async () => {
    if (!offlineName.trim()) {
      setError('Please enter a username');
      return;
    }

    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.loginOffline(offlineName);
      await loadAccounts();
      setOfflineName('');
    } catch {
      setError('Failed to create offline account');
    }
    setIsLoading(false);
  };

  const handleLoginMicrosoft = async () => {
    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.loginMicrosoft();
      await loadAccounts();
    } catch {
      setError('Failed to login with Microsoft');
    }
    setIsLoading(false);
  };

  const handleSelectAccount = async (index: number) => {
    setIsLoading(true);
    setError(null);
    try {
      await CMCLAPI.selectAccount(index + 1);
      setSelectedAccount(index);
    } catch {
      setError('Failed to select account');
    }
    setIsLoading(false);
  };

  return (
    <Box sx={{ p: 4 }}>
      <Card>
        <CardContent>
          <Typography variant="h4" sx={{ mb: 4 }}>Account Management</Typography>

          {error && <Alert severity="error" sx={{ mb: 4 }}>{error}</Alert>}

          <Box sx={{ mb: 6 }}>
            <Typography variant="h6" sx={{ mb: 2 }}>Add Account</Typography>
            
            <Box sx={{ mb: 4 }}>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Offline Account</Typography>
              <TextField
                label="Username"
                value={offlineName}
                onChange={(e) => setOfflineName(e.target.value)}
                sx={{ mr: 2, width: 200 }}
              />
              <Button
                variant="contained"
                startIcon={<Person />}
                onClick={handleLoginOffline}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Create Offline Account'}
              </Button>
            </Box>

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1 }}>Microsoft Account</Typography>
              <Button
                variant="contained"
                color="primary"
                startIcon={<Login />}
                onClick={handleLoginMicrosoft}
                disabled={isLoading}
              >
                {isLoading ? <CircularProgress size={20} /> : 'Login with Microsoft'}
              </Button>
            </Box>
          </Box>

          <Box>
            <Typography variant="h6" sx={{ mb: 2 }}>Accounts</Typography>
            <Button onClick={loadAccounts} disabled={isLoading} sx={{ mb: 2 }}>
              <Refresh sx={{ mr: 1 }} />
              {isLoading ? 'Loading...' : 'Refresh List'}
            </Button>
            <List>
              {accounts.map((account, index) => (
                <ListItem key={index}>
                  <ListItemButton onClick={() => handleSelectAccount(index)} selected={selectedAccount === index}>
                    <ListItemText primary={account} />
                    <ListItemSecondaryAction>
                      {selectedAccount === index && (
                        <Typography variant="body2" color="primary">Selected</Typography>
                      )}
                    </ListItemSecondaryAction>
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
            {accounts.length === 0 && (
              <Typography variant="body2" color="textSecondary">
                No accounts found. Add one above.
              </Typography>
            )}
          </Box>
        </CardContent>
      </Card>
    </Box>
  );
}
