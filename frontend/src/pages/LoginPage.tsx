import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../services/auth";

export function LoginPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setLoading(true);
    setError("");
    try {
      const data = await login({ email, password });
      localStorage.setItem("zilo_access_token", data.accessToken);
      localStorage.setItem("zilo_refresh_token", data.refreshToken);
      navigate("/console");
    } catch (err) {
      setError(err instanceof Error ? err.message : "登录失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={handleSubmit}>
        <h2>登录 Zilo</h2>
        <p className="auth-sub">登录后进入工作流控制台</p>
        <input
          className="auth-input"
          type="email"
          placeholder="邮箱"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          className="auth-input"
          type="password"
          placeholder="密码"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        {error && <div className="auth-error">{error}</div>}
        <button className="btn auth-submit" disabled={loading} type="submit">
          {loading ? "登录中..." : "登录"}
        </button>
        <p className="auth-foot">
          还没有账号？<Link to="/register">立即注册</Link>
        </p>
      </form>
    </div>
  );
}

