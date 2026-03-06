import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import { User, Lock, Mail } from "lucide-react";

export default function AuthForms() {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    const url = isLogin ? "/login" : "/register";

    try {
      const res = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      const data = await res.json();

      if (!res.ok) {
        throw new Error(data.error || "Authentication failed");
      }

      login(data.token);
    } catch (err: any) {
      setError(err.message || "Network error occurred");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        minHeight: "80vh",
      }}
    >
      <div
        className="card"
        style={{
          width: "100%",
          maxWidth: "400px",
          display: "flex",
          flexDirection: "column",
          gap: "1.5rem",
        }}
      >
        <div style={{ textAlign: "center" }}>
          <div
            style={{
              display: "flex",
              justifyContent: "center",
              marginBottom: "1rem",
            }}
          >
            <div
              style={{
                width: "48px",
                height: "48px",
                borderRadius: "50%",
                backgroundColor: "var(--brand-primary)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                color: "white",
              }}
            >
              <User size={24} />
            </div>
          </div>
          <h2>{isLogin ? "Welcome Back" : "Create Account"}</h2>
          <p className="text-secondary" style={{ marginTop: "0.5rem" }}>
            {isLogin
              ? "Sign in to access your dashboard"
              : "Register to start testing"}
          </p>
        </div>

        {error && (
          <div
            style={{
              backgroundColor: "var(--error-bg)",
              color: "var(--error)",
              padding: "0.75rem",
              borderRadius: "var(--radius-md)",
              fontSize: "0.875rem",
            }}
          >
            {error}
          </div>
        )}

        <form
          onSubmit={handleSubmit}
          style={{ display: "flex", flexDirection: "column", gap: "1rem" }}
        >
          <div className="form-group">
            <label className="label">
              <Mail size={16} /> Email Address
            </label>
            <input
              type="email"
              required
              className="input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
            />
          </div>

          <div className="form-group">
            <label className="label">
              <Lock size={16} /> Password
            </label>
            <input
              type="password"
              required
              minLength={6}
              className="input"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
            />
          </div>

          <button
            type="submit"
            className="btn btn-primary"
            style={{
              width: "100%",
              marginTop: "0.5rem",
              justifyContent: "center",
            }}
            disabled={isLoading}
          >
            {isLoading ? "Processing..." : isLogin ? "Sign In" : "Register"}
          </button>
        </form>

        <div style={{ textAlign: "center", fontSize: "0.875rem" }}>
          <span className="text-secondary">
            {isLogin ? "Don't have an account? " : "Already have an account? "}
          </span>
          <button
            type="button"
            onClick={() => setIsLogin(!isLogin)}
            style={{
              background: "none",
              border: "none",
              color: "var(--brand-primary)",
              fontWeight: 500,
              cursor: "pointer",
              padding: 0,
            }}
          >
            {isLogin ? "Sign Up" : "Sign In"}
          </button>
        </div>
      </div>
    </div>
  );
}
