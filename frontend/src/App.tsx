import { useState, useEffect } from "react";
import { Moon, Sun, Beaker } from "lucide-react";

import TestForm from "./components/TestForm";
import TestList from "./components/TestList";
import AuthForms from "./components/AuthForms";
import { AuthProvider, useAuth } from "./context/AuthContext";

function MainApp() {
  const { token, logout } = useAuth();
  const [isDark, setIsDark] = useState(true);
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  useEffect(() => {
    // Check initial dark mode preference
    const root = window.document.documentElement;
    if (isDark) {
      root.classList.add("dark");
    } else {
      root.classList.remove("dark");
    }
  }, [isDark]);

  const toggleTheme = () => {
    setIsDark(!isDark);
  };

  const handleTestCreated = () => {
    setRefreshTrigger((prev) => prev + 1);
  };

  return (
    <div className="container">
      <header className="header">
        <div className="header-title">
          <Beaker className="text-secondary" size={32} />
          <div>
            <h1>Testrunner</h1>
            <p className="text-secondary">
              Run tests concurrently and flawlessly.
            </p>
          </div>
        </div>

        <div style={{ display: "flex", gap: "1rem", alignItems: "center" }}>
          {token && (
            <button
              className="btn btn-outline"
              onClick={logout}
              style={{ padding: "0.4rem 0.8rem", fontSize: "0.875rem" }}
            >
              Sign Out
            </button>
          )}

          <button
            className="icon-btn"
            onClick={toggleTheme}
            aria-label="Toggle Theme"
          >
            {isDark ? <Sun size={24} /> : <Moon size={24} />}
          </button>
        </div>
      </header>

      {token ? (
        <main
          style={{
            display: "grid",
            gridTemplateColumns: "minmax(350px, 1fr) minmax(400px, 2fr)",
            gap: "2rem",
            alignItems: "start",
          }}
        >
          <section>
            <TestForm onRefresh={handleTestCreated} />
          </section>

          <section>
            <TestList refreshTrigger={refreshTrigger} />
          </section>
        </main>
      ) : (
        <AuthForms />
      )}
    </div>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <MainApp />
    </AuthProvider>
  );
}
