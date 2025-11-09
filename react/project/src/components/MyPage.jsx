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
    Alert, 
    Card,
    CardContent,
    Divider,
    IconButton,
    Link,
    Stack // Import Stack
} from '@mui/material';
import { useUser } from '/src/context/UserContext';
import EditIcon from '@mui/icons-material/Edit';
import SaveIcon from '@mui/icons-material/Save';
import CancelIcon from '@mui/icons-material/Cancel';
import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined'; // Import Email icon
import PersonOutlineOutlinedIcon from '@mui/icons-material/PersonOutlineOutlined'; // Import Person icon
import getToken from '/src/utils/getToken'; // Import getToken helper


const MyPage = () => {
    const { t } = useTranslation();
    const { user, setUser } = useUser();
    const navigate = useNavigate();
    
    const [newUsername, setNewUsername] = useState('');
    const [message, setMessage] = useState('');
    const [messageType, setMessageType] = useState(''); // 'success' or 'error'
    const [isEditing, setIsEditing] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    const fetchUser = useCallback(async (token) => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_APP_API_URL}/api/getuser`, {
                headers: { 'Authorization': token },
            });
            setUser(response.data);
            setNewUsername(response.data.username);
        } catch (error) {
            // If token is invalid, clear all tokens and redirect to login
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
        const token = getToken(); // Use getToken helper
        if (!token) {
            navigate('/auth/login');
            return;
        }
        if (!user) {
            fetchUser(token);
        } else {
            setNewUsername(user.username);
            setIsLoading(false);
        }
    }, [navigate, setUser, user, fetchUser]);

    const handleUpdateUsername = async (e) => {
        e.preventDefault();
        const token = getToken(); // Use getToken helper
        if (!token || newUsername === user.username) {
            setIsEditing(false);
            return;
        }

        try {
            const response = await axios.put(
                `${import.meta.env.VITE_APP_API_URL}/api/user/username`,
                { username: newUsername },
                { headers: { 'Authorization': token } }
            );
            setMessage(response.data.message);
            setMessageType('success');
            setUser({ ...user, username: newUsername });
            setIsEditing(false);
        } catch (error) {
            setMessage(error.response?.data?.error || 'Failed to update username.');
            setMessageType('error');
        }
    };

    const handleCancelEdit = () => {
        setIsEditing(false);
        setNewUsername(user.username);
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
                <CardContent sx={{ p: 4 }}>
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
                        <Typography variant="h5" component="h1">
                            {t('mypage_title')}
                        </Typography>
                        {!isEditing && (
                            <IconButton onClick={() => setIsEditing(true)} size="small" color="primary">
                                <EditIcon fontSize="small" />
                            </IconButton>
                        )}
                    </Box>

                    <Divider sx={{ my: 2 }} />

                    <Stack spacing={2} sx={{ mb: 2 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <EmailOutlinedIcon color="action" fontSize="small" />
                            <Typography variant="subtitle2" color="text.secondary">Email</Typography>
                        </Box>
                        <Typography variant="body1">{user?.email}</Typography>
                    </Stack>

                    <Stack spacing={2} sx={{ mb: 2 }}>
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                            <PersonOutlineOutlinedIcon color="action" fontSize="small" />
                            <Typography variant="subtitle2" color="text.secondary">Username</Typography>
                        </Box>
                        {isEditing ? (
                            <Box>
                                <TextField
                                    fullWidth
                                    variant="outlined"
                                    value={newUsername}
                                    onChange={(e) => setNewUsername(e.target.value)}
                                    sx={{ mb: 1 }}
                                />
                                <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
                                    <Button onClick={handleCancelEdit} startIcon={<CancelIcon />}>
                                        {t('cancel')}
                                    </Button>
                                    <Button type="submit" variant="contained" color="primary" startIcon={<SaveIcon />}>
                                        {t('save')}
                                    </Button>
                                </Box>
                            </Box>
                        ) : (
                            <Typography variant="body1">{user?.username}</Typography>
                        )}
                    </Stack>
                </CardContent>
                {message && (
                    <Box sx={{ px: 4, pb: 2 }}>
                        <Alert severity={messageType} onClose={() => setMessage('')}>
                            {message}
                        </Alert>
                    </Box>
                )}
            </Card>
        </Container>
    );
};

export default MyPage;