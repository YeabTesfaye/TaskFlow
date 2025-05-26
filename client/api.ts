import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

const api = axios.create({
    baseURL: API_URL,
    headers: {
      'Content-Type': 'application/json',
    },
});


// Add request interceptor to attach token 
api.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if(token){
        config.headers.Authorization = `Bearer ${token}`
    }
    return config
})

// Auth endpoints 
export const auth = {
    login : async(email: string, password : string) => {
        const response = await api.post("/login", {email,password})
        return response.data
    },
    signup: async (email: string, password: string, name : string) => {
        const response = await api.post('/signup', { email, password, name });
        return response.data;
      }, 
}

export const user = {
    getProfile: async () => {
      const response = await api.get('/users/me');
      return response.data;
    },
    updateProfile: async (data: any) => {
      const response = await api.patch('/users/me', data);
      return response.data;
    },
    changePassword: async (currentPassword: string, newPassword: string) => {
      const response = await api.post('/users/password', {
        current_password: currentPassword,
        new_password: newPassword
      });
      return response.data;
    },
    deleteAccount: async () => {
      const response = await api.delete('/users/me');
      return response.data;
    }
  };
export default api