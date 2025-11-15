// react/project/src/utils/getToken.js
const getToken = () => {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token');
    return token;
};

export default getToken;
