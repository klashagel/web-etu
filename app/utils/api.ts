import axios from 'axios';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
});

export const login = async (username: string, password: string) => {
  const response = await api.post('/login', { username, password });
  return response.data;
}

export const getUser = async () => {
  const response = await api.get('/user');
  return response.data;
}

export const getControllers = async () => {
  const response = await api.get('/controllers');
  return response.data;
}

export const getController = async (id: string) => {
  const response = await api.get(`/controllers/${id}`);
  return response.data;
}