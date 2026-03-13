import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";

interface AuthContextType {
  token: string | null;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType>({
  token: null,
  login: () => {},
  logout: () => {},
});

export const useAuth = () => useContext(AuthContext);

function isTokenExpired(token: string): boolean {
  try {
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.exp * 1000 < Date.now();
  } catch {
    return true;
  }
}

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    // Check local storage for an existing token on mount
    const savedToken = localStorage.getItem("auth_token");
    if (savedToken && !isTokenExpired(savedToken)) {
      setToken(savedToken);
    } else {
      localStorage.removeItem("auth_token");
    }
  }, []);

  const logout = () => {
    setToken(null);
    localStorage.removeItem("auth_token");
  };

  const login = (newToken: string) => {
    setToken(newToken);
    localStorage.setItem("auth_token", newToken);
  };

  useEffect(() => {
    const originalFetch = window.fetch;
    window.fetch = async (...args) => {
      const response = await originalFetch(...args);
      if (response.status === 401) {
        logout();
      }
      return response;
    };
    return () => {
      window.fetch = originalFetch;
    };
  }, []);

  return (
    <AuthContext.Provider value={{ token, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};
