import React, { useState, useEffect, useContext } from 'react';
import { Link } from 'react-router-dom';
import reactLogo from '/src/assets/react.svg';
import { styled, alpha, useTheme } from '@mui/material/styles';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import InputBase from '@mui/material/InputBase';
import Badge from '@mui/material/Badge';
import MenuItem from '@mui/material/MenuItem';
import AccountCircle from '@mui/icons-material/AccountCircle';
import SearchIcon from '@mui/icons-material/Search';
import NotificationsIcon from '@mui/icons-material/Notifications';
import Tooltip from '@mui/material/Tooltip';
import PersonAdd from '@mui/icons-material/PersonAdd';
import Settings from '@mui/icons-material/Settings';
import Logout from '@mui/icons-material/Logout';
import ListItemIcon from '@mui/material/ListItemIcon';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemText from '@mui/material/ListItemText';
import CloseIcon from '@mui/icons-material/Close';
import { useTranslation } from 'react-i18next';
import { useUser } from '/src/context/UserContext';
import { FormControl, InputLabel, Select } from '@mui/material';
import { ThemeContext } from '/src/context/ThemeContext';
import getToken from '/src/utils/getToken'; // Import getToken helper

const Search = styled('div')(({ theme }) => ({
    position: 'relative',
    borderRadius: theme.shape.borderRadius,
    backgroundColor: alpha(theme.palette.text.primary, 0.05),
    '&:hover': {
        backgroundColor: alpha(theme.palette.text.primary, 0.1),
    },
    marginRight: theme.spacing(2),
    marginLeft: 0,
    width: '100%',
    [theme.breakpoints.up('sm')]: {
        marginLeft: theme.spacing(3),
        width: 'auto',
    },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
    padding: theme.spacing(0, 2),
    height: '100%',
    position: 'absolute',
    pointerEvents: 'none',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
    color: 'inherit',
    '& .MuiInputBase-input': {
        padding: theme.spacing(1, 1, 1, 0),
        paddingLeft: `calc(1em + ${theme.spacing(4)})`,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('md')]: {
            width: '20ch',
        },
    },
}));

function AccountMenu({ isLoggedIn, i18n, currentLang, handleLangChange, toggleTheme, mode, t }) {
    const [drawerOpen, setDrawerOpen] = useState(false);
    const { user } = useUser();
    const theme = useTheme();

    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }
        setDrawerOpen(open);
    };

    const menuItems = isLoggedIn ? [
        { text: t('mypage'), link: '/mypage', icon: <AccountCircle /> },
        { text: t('settings'), link: '/settings', icon: <Settings /> },
        { divider: true },
        { text: t('logout'), link: '/logout', icon: <Logout /> },
    ] : [
        { text: t('register'), link: '/auth/register', icon: <PersonAdd /> },
        { text: t('login'), link: '/auth/login', icon: <AccountCircle /> },
    ];

    return (
        <>
            <Tooltip title={t('account_settings')}>
                <IconButton
                    onClick={toggleDrawer(true)}
                    size="large"
                    color="inherit"
                >
                    <AccountCircle />
                </IconButton>
            </Tooltip>
            <Drawer anchor="right" open={drawerOpen} onClose={toggleDrawer(false)}>
                <Box
                    sx={{ width: 250, display: 'flex', flexDirection: 'column', height: '100%', bgcolor: 'background.default' }}
                    role="presentation"
                >
                    <Box sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                        {isLoggedIn && user ? (
                            <Box>
                                <Typography variant="subtitle1">{user.username}</Typography>
                            </Box>
                        ) : (
                            <Box></Box> // Empty box to maintain space-between for CloseIcon
                        )}
                        <IconButton onClick={toggleDrawer(false)} color="inherit">
                            <CloseIcon sx={{ color: 'text.secondary' }} />
                        </IconButton>
                    </Box>
                    <Divider />
                    <List sx={{ flexGrow: 1 }}>
                        {menuItems.map((item, index) => (
                            item.divider ? ( <Divider key={index} sx={{ my: 1 }}/> ) : (
                            <ListItem button component={Link} to={item.link} key={index} onClick={toggleDrawer(false)}>
                                <ListItemIcon sx={{ color: 'text.secondary' }}>{item.icon}</ListItemIcon>
                                <ListItemText primary={item.text} />
                            </ListItem>
                            )
                        ))}
                        <Divider sx={{ my: 1 }}/>
                        <ListItem>
                            <FormControl variant="outlined" fullWidth size="small" sx={{ minWidth: 90 }}>
                                <InputLabel id="lang-select-label-drawer" sx={{ color: 'text.secondary' }}>{t('language')}</InputLabel>
                                <Select
                                    labelId="lang-select-label-drawer"
                                    value={currentLang}
                                    onChange={handleLangChange}
                                    label={t('language')}
                                    sx={{
                                        color: 'text.primary',
                                        '.MuiOutlinedInput-notchedOutline': {
                                            borderColor: 'divider',
                                        },
                                        '&:hover .MuiOutlinedInput-notchedOutline': {
                                            borderColor: 'text.secondary',
                                        },
                                        '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                            borderColor: 'primary.main',
                                        },
                                        '.MuiSelect-icon': {
                                            color: 'text.secondary',
                                        },
                                    }}
                                >
                                    <MenuItem value="en">{t('english')}</MenuItem>
                                    <MenuItem value="ja">{t('japanese')}</MenuItem>
                                </Select>
                            </FormControl>
                        </ListItem>
                        <ListItem>
                            <FormControl variant="outlined" fullWidth size="small" sx={{ minWidth: 90 }}>
                                <InputLabel id="theme-select-label-drawer" sx={{ color: 'text.secondary' }}>{t('theme')}</InputLabel>
                                <Select
                                    labelId="theme-select-label-drawer"
                                    value={mode}
                                    onChange={toggleTheme}
                                    label={t('theme')}
                                    sx={{
                                        color: 'text.primary',
                                        '.MuiOutlinedInput-notchedOutline': {
                                            borderColor: 'divider',
                                        },
                                        '&:hover .MuiOutlinedInput-notchedOutline': {
                                            borderColor: 'text.secondary',
                                        },
                                        '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                            borderColor: 'primary.main',
                                        },
                                        '.MuiSelect-icon': {
                                            color: 'text.secondary',
                                        },
                                    }}
                                >
                                    <MenuItem value="light">{t('light')}</MenuItem>
                                    <MenuItem value="dark">{t('dark')}</MenuItem>
                                </Select>
                            </FormControl>
                        </ListItem>
                    </List>
                </Box>
            </Drawer>
        </>
    );
}

export default function PrimarySearchAppBar() {
    const { t, i18n } = useTranslation();
    const { toggleTheme, mode } = useContext(ThemeContext);
    const [isLoggedIn, setIsLoggedIn] = useState(Boolean(getToken())); // Use getToken
    const [currentLang, setCurrentLang] = useState(i18n.language);

    useEffect(() => {
        const handleStorageChange = () => {
            setIsLoggedIn(Boolean(getToken())); // Use getToken
        };

        window.addEventListener('storage', handleStorageChange);
        window.dispatchEvent(new Event("storage"));

        return () => {
            window.removeEventListener('storage', handleStorageChange);
        };
    }, [isLoggedIn]);

    useEffect(() => {
        setCurrentLang(i18n.language);
    }, [i18n.language]);

    const handleLangChange = (event) => {
        const newLang = event.target.value;
        i18n.changeLanguage(newLang);
        setCurrentLang(newLang);
    };

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static" sx={{ bgcolor: 'background.paper', color: 'text.primary' }} elevation={1}>
                <Toolbar>
                    <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
                        <img src={reactLogo} alt="React Logo" style={{ height: 40, width: 40 }} />
                    </Box>
                    <Box sx={{ display: { xs: 'none', sm: 'block' } }}> {/* Hide Search on xs screens */}
                        <Search>
                            <SearchIconWrapper>
                                <SearchIcon />
                            </SearchIconWrapper>
                            <StyledInputBase
                                placeholder={t('search')}
                                inputProps={{ 'aria-label': 'search' }}
                            />
                        </Search>
                    </Box>
                    <Box sx={{ flexGrow: 1 }} />
                    <Box sx={{ display: { xs: 'none', sm: 'flex' }, alignItems: 'center' }}> {/* Show language/theme selectors on sm and up screens */}
                        <FormControl variant="outlined" size="small" sx={{ m: 1, minWidth: 90 }}>
                            <InputLabel id="lang-select-label" sx={{ color: 'text.secondary' }}>{t('language')}</InputLabel>
                            <Select
                                labelId="lang-select-label"
                                value={currentLang}
                                onChange={handleLangChange}
                                label={t('language')}
                                sx={{
                                    color: 'text.primary',
                                    '.MuiOutlinedInput-notchedOutline': {
                                        borderColor: 'divider',
                                    },
                                    '&:hover .MuiOutlinedInput-notchedOutline': {
                                        borderColor: 'text.secondary',
                                    },
                                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                        borderColor: 'primary.main',
                                    },
                                    '.MuiSelect-icon': {
                                        color: 'text.secondary',
                                    },
                                }}
                            >
                                <MenuItem value="en">{t('english')}</MenuItem>
                                <MenuItem value="ja">{t('japanese')}</MenuItem>
                            </Select>
                        </FormControl>
                        <FormControl variant="outlined" size="small" sx={{ m: 1, minWidth: 90 }}>
                            <InputLabel id="theme-select-label" sx={{ color: 'text.secondary' }}>{t('theme')}</InputLabel>
                            <Select
                                labelId="theme-select-label"
                                value={mode}
                                onChange={toggleTheme}
                                label={t('theme')}
                                sx={{
                                    color: 'text.primary',
                                    '.MuiOutlinedInput-notchedOutline': {
                                        borderColor: 'divider',
                                    },
                                    '&:hover .MuiOutlinedInput-notchedOutline': {
                                        borderColor: 'text.secondary',
                                    },
                                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                        borderColor: 'primary.main',
                                    },
                                    '.MuiSelect-icon': {
                                        color: 'text.secondary',
                                    },
                                }}
                            >
                                <MenuItem value="light">{t('light')}</MenuItem>
                                <MenuItem value="dark">{t('dark')}</MenuItem>
                            </Select>
                        </FormControl>
                    </Box>
                    <AccountMenu 
                        isLoggedIn={isLoggedIn} 
                        i18n={i18n} 
                        currentLang={currentLang} 
                        handleLangChange={handleLangChange} 
                        toggleTheme={toggleTheme} 
                        mode={mode} 
                        t={t} 
                    />
                </Toolbar>
            </AppBar>
        </Box>
    );
}
