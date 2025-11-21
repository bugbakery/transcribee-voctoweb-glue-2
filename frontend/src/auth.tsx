import { createContext, useContext, useEffect, useState } from 'react';
import { pb } from './pb';

export const AuthContext = createContext<{
  loggedIn: boolean | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
} | null>(null);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [loggedIn, setLoggedIn] = useState<boolean>(pb.authStore.isValid);

  useEffect(() => {
    if (!pb.authStore.isValid) return;

    pb.collection('users')
      .authRefresh()
      .then(() => {
        setLoggedIn(true);
        console.log('Refreshed auth');
      })
      .catch((e) => {
        if (e.isAbort) return;
        setLoggedIn(false);
        console.log('Failed to refresh auth', e);
      });
  }, []);

  const login = async (username: string, password: string) => {
    await pb.collection('users').authWithPassword(username, password);
    setLoggedIn(true);
  };

  const logout = async () => {
    try {
      pb.authStore.clear();
      setLoggedIn(false);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <AuthContext.Provider value={{ loggedIn, login, logout }}>{children}</AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  return context;
};
