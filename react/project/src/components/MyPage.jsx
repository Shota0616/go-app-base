import React, { useEffect, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTranslation } from 'react-i18next';
import { 
    Container, 
    Typography, 
    Box, 
    CircularProgress, 
    TextField, 
    Button, 
    Card,
    CardContent,
    Divider,
    IconButton,
    Link,
    Stack,
    Snackbar 
} from '@mui/material';
import MuiAlert from '@mui/material/Alert';
import { useUser } from '/src/context/UserContext';
import EditIcon from '@mui/icons-material/Edit';
import SaveIcon from '@mui/icons-material/Save';
import CancelIcon from '@mui/icons-material/Cancel';
import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined';
import PersonOutlineOutlinedIcon from '@mui/icons-material/PersonOutlineOutlined';
import getToken from '/src/utils/getToken';

// Custom Alert component for Snackbar
const Alert = React.forwardRef(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});


const MyPage = () => {
    const { t } = useTranslation();
    const { user, setUser } = useUser();
    const navigate = useNavigate();
    
    const [newUsername, setNewUsername] = useState('');
    const [newEmail, setNewEmail] = useState('');
    const [message, setMessage] = useState('');
    const [messageType, setMessageType] = useState('');
    const [snackbarOpen, setSnackbarOpen] = useState(false);
    const [isEditingMode, setIsEditingMode] = useState(false); // Global editing mode
    const [isLoading, setIsLoading] = useState(true);

    const fetchUser = useCallback(async (token) => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_APP_API_URL}/api/getuser`, {
                headers: { 'Authorization': token },
            });
            setUser(response.data);
            setNewUsername(response.data.username);
            setNewEmail(response.data.email);
        } catch (error) {
            localStorage.removeItem('token');
            localStorage.removeItem('refreshtoken');
            sessionStorage.removeItem('token');
            sessionStorage.removeItem('refreshtoken');
            window.dispatchEvent(new Event("storage"));
            navigate('/auth/login');
        } finally {
            setIsLoading(false);
        }
    }, [navigate, setUser]);

    useEffect(() => {
        const token = getToken();
        if (!token) {
            navigate('/auth/login');
            return;
        }
        if (!user) {
            fetchUser(token);
        } else {
            setNewUsername(user.username);
            setNewEmail(user.email);
            setIsLoading(false);
        }
    }, [navigate, setUser, user, fetchUser]);

    const handleSave = async () => {
        const token = getToken();
        if (!token) {
            navigate('/auth/login');
            return;
        }

        let usernameUpdated = false;
        let emailUpdated = false;

        // Update Username
        if (newUsername !== user.username) {
            try {
                const response = await axios.put(
                    `${import.meta.env.VITE_APP_API_URL}/api/user/username`,
                    { username: newUsername },
                    { headers: { 'Authorization': token } }
                );
                setMessage(response.data.message);
                setMessageType('success');
                setSnackbarOpen(true);
                setUser({ ...user, username: newUsername });
                usernameUpdated = true;
            } catch (error) {
                setMessage(error.response?.data?.error || t('failed_to_update_username'));
                setMessageType('error');
                setSnackbarOpen(true);
                return; // Stop if username update fails
            }
        }

        // Update Email
        if (newEmail !== user.email) {
            try {
                const response = await axios.put(
                    `${import.meta.env.VITE_APP_API_URL}/api/user/email`,
                    { newEmail: newEmail },
                    { headers: { 'Authorization': token } }
                );
                setMessage(response.data.message);
                setMessageType('success');
                setSnackbarOpen(true);
                // Email is updated, but user is now inactive, so redirect to verify
                navigate('/auth/verify', { state: { message: response.data.message, messageType: 'success', email: newEmail } });
                emailUpdated = true;
            } catch (error) {
                setMessage(error.response?.data?.error || t('failed_to_update_email'));
                setMessageType('error');
                setSnackbarOpen(true);
                setNewEmail(user.email); // Revert newEmail to original on failure
                return; // Stop if email update fails
            }
        }

        if (usernameUpdated || emailUpdated) {
            setIsEditingMode(false); // Exit editing mode if any update was successful
        } else {
            setMessage(t('no_changes_to_save'));
            setMessageType('info');
            setSnackbarOpen(true);
            setIsEditingMode(false); // Exit editing mode even if no changes
        }
    };

    const handleCancel = () => {
        setNewUsername(user.username);
        setNewEmail(user.email);
        setIsEditingMode(false);
        setMessage('');
        setSnackbarOpen(false);
    };

    const handleToggleEditMode = () => {
        setIsEditingMode((prev) => !prev);
        if (isEditingMode) { // If exiting edit mode, reset values
            setNewUsername(user.username);
            setNewEmail(user.email);
        }
        setMessage('');
        setSnackbarOpen(false);
    };

    const handleCloseSnackbar = (event, reason) => {
        if (reason === 'clickaway') {
            return;
        }
        setSnackbarOpen(false);
        setMessage('');
    };

    if (isLoading) {
        return (
            <Container maxWidth="sm" sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                <CircularProgress />
            </Container>
        );
    }

    return (
        <Container component="main" maxWidth="sm" sx={{ mt: 4 }}>
            <Card elevation={3} sx={{ borderRadius: 2 }}>
                <CardContent sx={{ p: 4, alignItems: 'flex-start' }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
                        <Typography variant="h5" component="h1">
                            {t('mypage_title')}
                        </Typography>
                        <IconButton onClick={handleToggleEditMode} size="small" color="primary">
                            {isEditingMode ? <CancelIcon /> : <EditIcon />}
                        </IconButton>
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    <Stack spacing={2} sx={{ mb: 2, alignItems: 'flex-start', width: '100%' }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <EmailOutlinedIcon color="action" fontSize="small" />
                            <Typography variant="subtitle2" color="text.secondary">Email</Typography>
                        </Box>
                        {isEditingMode ? (
                            <TextField
                                fullWidth
                                variant="outlined"
                                value={newEmail}
                                onChange={(e) => setNewEmail(e.target.value)}
                                sx={{ mb: 1 }}
                            />
                        ) : (
                            <Typography variant="body1" sx={{ textAlign: 'left' }}>{user?.email}</Typography>
                        )}
                    </Stack>

                    <Stack spacing={2} sx={{ mb: 2, alignItems: 'flex-start', width: '100%' }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <PersonOutlineOutlinedIcon color="action" fontSize="small" />
                            <Typography variant="subtitle2" color="text.secondary">Username</Typography>
                        </Box>
                        {isEditingMode ? (
                            <TextField
                                fullWidth
                            
                                variant="outlined"
                                value={newUsername}
                                onChange={(e) => setNewUsername(e.target.value)}
                                sx={{ mb: 1 }}
                            />
                        ) : (
                            <Typography variant="body1" sx={{ textAlign: 'left' }}>{user?.username}</Typography>
                        )}
                    </Stack>

                    {isEditingMode && (
                        <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 1, mt: 2 }}>
                            <Button onClick={handleCancel} startIcon={<CancelIcon />}>
                                {t('cancel')}
                            </Button>
                            <Button onClick={handleSave} variant="contained" color="primary" startIcon={<SaveIcon />}>
                                {t('save')}
                            </Button>
                        </Box>
                    )}
                </CardContent>
            </Card>
            <Snackbar open={snackbarOpen} autoHideDuration={6000} onClose={handleCloseSnackbar} anchorOrigin={{ vertical: 'top', horizontal: 'center' }}>
                <Alert onClose={handleCloseSnackbar} severity={messageType} sx={{ width: '100%' }}>
                    {message}
                </Alert>
            </Snackbar>
        </Container>
    );
};

export default MyPage;