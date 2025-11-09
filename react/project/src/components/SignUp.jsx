import React, { useState } from 'react';
import { useNavigate, Link as RouterLink } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next';

import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import { Alert, Card, CardContent } from '@mui/material';

export default function SignUp() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('error');

  const handleRegister = async (event) => {
    event.preventDefault();
    if (!email || !password) {
        setMessage(t('registration_failed'));
        setMessageType('error');
        return;
    }

    try {
        const response = await axios.post(`${import.meta.env.VITE_APP_API_URL}/api/register`, {
            email,
            password,
        });
        // Pass email to verify page
        navigate('/auth/verify', { state: { message: response.data.message, messageType: 'success', email: email } });
    } catch (error) {
        setMessage(error.response?.data?.error || t('registration_failed'));
        setMessageType('error');
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline />
      <Card sx={{ marginTop: 8, borderRadius: 2 }}>
        <CardContent sx={{ p: 4 }}>
          <Box
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
            }}
          >
            <Typography component="h1" variant="h5">
              {t('register')}
            </Typography>
            <Box component="form" noValidate onSubmit={handleRegister} sx={{ mt: 3 }}>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <TextField
                    required
                    fullWidth
                    id="email"
                    label={t('email')}
                    name="email"
                    autoComplete="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    required
                    fullWidth
                    name="password"
                    label={t('password')}
                    type="password"
                    id="password"
                    autoComplete="new-password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                  />
                </Grid>
              </Grid>
              {message && (
                <Alert severity={messageType} sx={{ mt: 2, width: '100%' }}>
                    {message}
                </Alert>
              )}
              <Button
                type="submit"
                fullWidth
                variant="contained"
                sx={{ mt: 3, mb: 2 }}
              >
                {t('register')}
              </Button>
              <Grid container justifyContent="flex-end">
                <Grid item>
                                <Link component={RouterLink} to="/auth/login" variant="body2">
                                  {t('signin_link_text')}
                                </Link>                </Grid>
              </Grid>
            </Box>
          </Box>
        </CardContent>
      </Card>
    </Container>
  );
}
