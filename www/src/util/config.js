let API_BACKEND_HOST = 'http://localhost:8081';

if (process.env.NODE_ENV === 'production') {
    API_BACKEND_HOST = "http://codemaps.liamz.co"
}

export {
    API_BACKEND_HOST
};