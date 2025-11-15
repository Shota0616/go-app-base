import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Box, Typography, Link, CssBaseline, Container, CircularProgress } from '@mui/material';
import { useNavigate } from 'react-router-dom'; // Import useNavigate

const Logout = () => {
    const { t } = useTranslation();
    const [logoutMessage, setLogoutMessage] = useState('');
    const [isLoggingOut, setIsLoggingOut] = useState(true); // New state for logging out status
    const navigate = useNavigate(); // Initialize useNavigate

    const handleLogout = async () => {
        setIsLoggingOut(true);
        try {
            localStorage.removeItem('token');
            localStorage.removeItem('refreshtoken');
            sessionStorage.removeItem('token');
            sessionStorage.removeItem('refreshtoken');
            window.dispatchEvent(new Event("storage"));
            setLogoutMessage(t('logout_successful'));
            setTimeout(() => {
                navigate('/auth/login'); // Redirect to login page after successful logout
            }, 1500); // Redirect after 1.5 seconds
        } catch (error) {
            setLogoutMessage(t('logout_failed'));
            setIsLoggingOut(false); // Stop loading on error
        }
    };

    useEffect(() => {
        handleLogout(); // Call handleLogout when component mounts
    }, []); // Empty dependency array means this runs once on mount

    return (
        <Container component="main" maxWidth="xs" sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh' }}>
            <CssBaseline />
            <Box
                sx={{
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                    justifyContent: 'center',
                    width: '100%',
                    mt: 8,
                }}
            >
                <Typography component="h1" variant="h5" sx={{ mb: 2 }}>
                    {t('logging_out')}
                </Typography>
                {isLoggingOut ? (
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
                        <CircularProgress />
                    </Box>
                ) : (
                    <Typography>{logoutMessage}</Typography>
                )}
            </Box>
        </Container>
    );
};

export default Logout;
