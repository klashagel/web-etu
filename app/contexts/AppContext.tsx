import React, { createContext, useContext, useReducer, useEffect } from 'react';
import { getUser } from '../utils/api'; // You'll need to implement this function

interface User {
  id: string;
  username: string;
  // Add other user properties as needed
}

interface AppState {
  user: User | null;
  isLoading: boolean;
}

type Action =
  | { type: 'SET_USER'; payload: User }
  | { type: 'CLEAR_USER' }
  | { type: 'SET_LOADING'; payload: boolean };

const initialState: AppState = {
  user: null,
  isLoading: true,
};

const AppContext = createContext<{
  state: AppState;
  dispatch: React.Dispatch<Action>;
}>({ state: initialState, dispatch: () => null });

const appReducer = (state: AppState, action: Action): AppState => {
  switch (action.type) {
    case 'SET_USER':
      return { ...state, user: action.payload, isLoading: false };
    case 'CLEAR_USER':
      return { ...state, user: null, isLoading: false };
    case 'SET_LOADING':
      return { ...state, isLoading: action.payload };
    default:
      return state;
  }
};

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(appReducer, initialState);

  useEffect(() => {
    const checkAuth = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        try {
          const userData = await getUser();
          dispatch({ type: 'SET_USER', payload: userData });
        } catch (error) {
          console.error('Error verifying token:', error);
          localStorage.removeItem('token');
        }
      }
      dispatch({ type: 'SET_LOADING', payload: false });
    };

    checkAuth();
  }, []);

  return (
    <AppContext.Provider value={{ state, dispatch }}>
      {children}
    </AppContext.Provider>
  );
};

export const useAppContext = () => useContext(AppContext);