import { Link } from "react-router-dom";
import { VendorOrbitalShowcase } from "../components/VendorOrbitalShowcase";

export function HomePage() {
  return (
    <div className="home-page">
      <header className="top-nav">
        <div className="logo">Zilo</div>
        <nav>
          <a href="#features">功能</a>
          <a href="#providers">供应商</a>
          <Link to="/login">登录</Link>
          <Link to="/register" className="btn btn-sm">
            注册
          </Link>
        </nav>
      </header>

      <main className="hero">
        <section className="hero-left">
          <p className="badge">Eino + go-zero</p>
          <h1>搭建你的 AI 工作流中枢</h1>
          <p className="desc">
            以可视化方式编排节点、发布版本、追踪执行日志，快速构建可运营的智能流程平台。
          </p>
          <div className="cta-row">
            <Link to="/register" className="btn">
              立即注册
            </Link>
            <Link to="/login" className="btn btn-ghost">
              立即登录
            </Link>
            <Link to="/console" className="btn btn-ghost">
              进入控制台
            </Link>
          </div>
        </section>
        <section className="hero-right" id="providers">
          <VendorOrbitalShowcase />
        </section>
      </main>
    </div>
  );
}

