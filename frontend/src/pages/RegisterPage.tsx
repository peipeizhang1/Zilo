import { FormEvent, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { register } from "../services/auth";

export function RegisterPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [nickname, setNickname] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setLoading(true);
    setError("");
    try {
      await register({ email, password, nickname });
      navigate("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "注册失败");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={handleSubmit}>
        <h2>创建账号</h2>
        <p className="auth-sub">注册后可直接配置你的工作流</p>
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
          type="text"
          placeholder="昵称（可选）"
          value={nickname}
          onChange={(e) => setNickname(e.target.value)}
        />
        <input
          className="auth-input"
          type="password"
          placeholder="密码（至少8位）"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={8}
        />
        {error && <div className="auth-error">{error}</div>}
        <button className="btn auth-submit" disabled={loading} type="submit">
          {loading ? "注册中..." : "注册"}
        </button>
        <p className="auth-foot">
          已有账号？<Link to="/login">去登录</Link>
        </p>
      </form>
    </div>
  );
}

