import { useState } from "react";

const API = "";

// ─── Утилиты ────────────────────────────────────────────────────────────────

async function apiLogin(login, password) {
  const res = await fetch(`${API}/user/auth`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ login, password }),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "Ошибка входа");
  return data;
}

async function apiRegister(fields) {
  const res = await fetch(`${API}/user/register`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(fields),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.message || "Ошибка регистрации");
  return data;
}

// ─── Общие компоненты ────────────────────────────────────────────────────────

function Input({ label, type = "text", value, onChange, placeholder }) {
  return (
      <div style={styles.fieldWrap}>
        <label style={styles.label}>{label}</label>
        <input
            type={type}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            placeholder={placeholder}
            style={styles.input}
            onFocus={(e) => Object.assign(e.target.style, styles.inputFocus)}
            onBlur={(e) => Object.assign(e.target.style, styles.input)}
        />
      </div>
  );
}

function Button({ children, onClick, loading, variant = "primary" }) {
  return (
      <button
          onClick={onClick}
          disabled={loading}
          style={variant === "primary" ? styles.btn : styles.btnGhost}
          onMouseEnter={(e) => {
            if (!loading)
              e.target.style.transform = "translateY(-2px)";
            e.target.style.boxShadow = variant === "primary"
                ? "0 8px 30px rgba(99,209,148,0.45)"
                : "none";
          }}
          onMouseLeave={(e) => {
            e.target.style.transform = "translateY(0)";
            e.target.style.boxShadow = variant === "primary"
                ? "0 4px 20px rgba(99,209,148,0.25)"
                : "none";
          }}
      >
        {loading ? "..." : children}
      </button>
  );
}

function ErrorMsg({ text }) {
  if (!text) return null;
  return <div style={styles.error}>{text}</div>;
}

// ─── Страница входа ──────────────────────────────────────────────────────────

function LoginPage({ onSuccess, onGoRegister }) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleSubmit() {
    if (!login || !password) { setError("Заполните все поля"); return; }
    setError(""); setLoading(true);
    try {
      const data = await apiLogin(login, password);
      localStorage.setItem("access_token", data.access_token);
      onSuccess();
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
      <Card>
        <Logo />
        <h2 style={styles.title}>Добро пожаловать</h2>
        <p style={styles.subtitle}>Войдите в свой аккаунт</p>

        <Input label="Логин или email" value={login} onChange={setLogin} placeholder="username" />
        <Input label="Пароль" type="password" value={password} onChange={setPassword} placeholder="••••••••" />

        <ErrorMsg text={error} />

        <Button onClick={handleSubmit} loading={loading}>Войти</Button>

        <p style={styles.hint}>
          Нет аккаунта?{" "}
          <span style={styles.link} onClick={onGoRegister}>
          Зарегистрируйтесь
        </span>
        </p>
      </Card>
  );
}

// ─── Страница регистрации ─────────────────────────────────────────────────────

function RegisterPage({ onSuccess, onGoLogin }) {
  const [form, setForm] = useState({
    username: "", password: "", email: "", name: "", surname: "", birth: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const set = (key) => (val) => setForm((f) => ({ ...f, [key]: val }));

  async function handleSubmit() {
    const { username, password, email, name, surname, birth } = form;
    if (!username || !password || !email) { setError("Заполните обязательные поля"); return; }
    setError(""); setLoading(true);
    try {
      await apiRegister({ username, password, email, name, surname, birth });
      onSuccess();
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
      <Card wide>
        <Logo />
        <h2 style={styles.title}>Создать аккаунт</h2>
        <p style={styles.subtitle}>Заполните данные для регистрации</p>

        <div style={styles.grid2}>
          <Input label="Имя" value={form.name} onChange={set("name")} placeholder="Иван" />
          <Input label="Фамилия" value={form.surname} onChange={set("surname")} placeholder="Иванов" />
        </div>

        <Input label="Логин *" value={form.username} onChange={set("username")} placeholder="username" />
        <Input label="Email *" type="email" value={form.email} onChange={set("email")} placeholder="ivan@mail.ru" />
        <Input label="Пароль *" type="password" value={form.password} onChange={set("password")} placeholder="••••••••" />
        <Input label="Дата рождения" type="date" value={form.birth} onChange={set("birth")} />

        <ErrorMsg text={error} />

        <Button onClick={handleSubmit} loading={loading}>Зарегистрироваться</Button>

        <p style={styles.hint}>
          Уже есть аккаунт?{" "}
          <span style={styles.link} onClick={onGoLogin}>
          Войти
        </span>
        </p>
      </Card>
  );
}

// ─── Страница успешной регистрации ────────────────────────────────────────────

function RegisteredPage({ onGoLogin }) {
  return (
      <Card>
        <div style={styles.successIcon}>✓</div>
        <h2 style={{ ...styles.title, color: "var(--green)" }}>Готово!</h2>
        <p style={styles.subtitle}>Вы успешно зарегистрировались</p>
        <Button onClick={onGoLogin}>Войти в аккаунт</Button>
      </Card>
  );
}

// ─── Страница успешного входа ─────────────────────────────────────────────────

function SuccessPage({ onLogout }) {
  return (
      <Card>
        <div style={{ ...styles.successIcon, background: "rgba(99,209,148,0.15)" }}>📅</div>
        <h2 style={{ ...styles.title, color: "var(--green)" }}>Успешный вход</h2>
        <p style={styles.subtitle}>Вы вошли в myCalendar</p>
        <Button variant="ghost" onClick={onLogout}>Выйти</Button>
      </Card>
  );
}

// ─── Обёртка Card ─────────────────────────────────────────────────────────────

function Card({ children, wide }) {
  return (
      <div style={{ ...styles.card, maxWidth: wide ? 480 : 380 }}>
        {children}
      </div>
  );
}

function Logo() {
  return (
      <div style={styles.logo}>
        <span style={styles.logoDot} />
        <span style={styles.logoText}>myCalendar</span>
      </div>
  );
}

// ─── Корневой роутер ──────────────────────────────────────────────────────────

export default function App() {
  const [page, setPage] = useState("login"); // login | register | registered | success

  return (
      <>
        <style>{globalCss}</style>
        <div style={styles.bg}>
          <div style={styles.glow1} />
          <div style={styles.glow2} />

          <div style={styles.center}>
            {page === "login" && (
                <LoginPage
                    onSuccess={() => setPage("success")}
                    onGoRegister={() => setPage("register")}
                />
            )}
            {page === "register" && (
                <RegisterPage
                    onSuccess={() => setPage("registered")}
                    onGoLogin={() => setPage("login")}
                />
            )}
            {page === "registered" && (
                <RegisteredPage onGoLogin={() => setPage("login")} />
            )}
            {page === "success" && (
                <SuccessPage onLogout={() => { localStorage.removeItem("access_token"); setPage("login"); }} />
            )}
          </div>
        </div>
      </>
  );
}

// ─── Стили ────────────────────────────────────────────────────────────────────

const globalCss = `
  @import url('https://fonts.googleapis.com/css2?family=DM+Serif+Display&family=DM+Sans:wght@300;400;500&display=swap');

  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --bg: #0d0f12;
    --surface: #14171c;
    --border: rgba(255,255,255,0.07);
    --text: #e8eaf0;
    --muted: #7a7f8e;
    --green: #63d194;
    --green-dim: rgba(99,209,148,0.12);
    --red: #f87171;
  }

  body { background: var(--bg); color: var(--text); font-family: 'DM Sans', sans-serif; }

  input[type='date']::-webkit-calendar-picker-indicator { filter: invert(0.6); cursor: pointer; }
`;

const styles = {
  bg: {
    minHeight: "100vh",
    background: "var(--bg)",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    position: "relative",
    overflow: "hidden",
  },
  glow1: {
    position: "absolute", top: "-10%", left: "20%",
    width: 500, height: 500, borderRadius: "50%",
    background: "radial-gradient(circle, rgba(99,209,148,0.08) 0%, transparent 70%)",
    pointerEvents: "none",
  },
  glow2: {
    position: "absolute", bottom: "-10%", right: "15%",
    width: 400, height: 400, borderRadius: "50%",
    background: "radial-gradient(circle, rgba(99,209,148,0.05) 0%, transparent 70%)",
    pointerEvents: "none",
  },
  center: {
    position: "relative", zIndex: 1,
    padding: "20px",
    width: "100%",
    display: "flex",
    justifyContent: "center",
  },
  card: {
    width: "100%",
    background: "var(--surface)",
    border: "1px solid var(--border)",
    borderRadius: 20,
    padding: "40px 36px",
    boxShadow: "0 24px 80px rgba(0,0,0,0.5)",
    animation: "fadeUp 0.4s ease both",
  },
  logo: {
    display: "flex", alignItems: "center", gap: 8,
    marginBottom: 28,
  },
  logoDot: {
    width: 8, height: 8, borderRadius: "50%",
    background: "var(--green)",
    boxShadow: "0 0 8px var(--green)",
    display: "inline-block",
  },
  logoText: {
    fontFamily: "'DM Serif Display', serif",
    fontSize: 18, color: "var(--text)", letterSpacing: "-0.3px",
  },
  title: {
    fontFamily: "'DM Serif Display', serif",
    fontSize: 26, fontWeight: 400,
    color: "var(--text)", marginBottom: 6,
    letterSpacing: "-0.5px",
  },
  subtitle: {
    fontSize: 14, color: "var(--muted)",
    marginBottom: 28, fontWeight: 300,
  },
  fieldWrap: { marginBottom: 16 },
  label: {
    display: "block", fontSize: 12,
    color: "var(--muted)", marginBottom: 6,
    fontWeight: 500, letterSpacing: "0.3px",
    textTransform: "uppercase",
  },
  input: {
    width: "100%", padding: "11px 14px",
    background: "rgba(255,255,255,0.04)",
    border: "1px solid var(--border)",
    borderRadius: 10, color: "var(--text)",
    fontSize: 14, outline: "none",
    transition: "border-color 0.2s, background 0.2s",
    fontFamily: "'DM Sans', sans-serif",
  },
  inputFocus: {
    width: "100%", padding: "11px 14px",
    background: "rgba(99,209,148,0.06)",
    border: "1px solid rgba(99,209,148,0.4)",
    borderRadius: 10, color: "var(--text)",
    fontSize: 14, outline: "none",
    transition: "border-color 0.2s, background 0.2s",
    fontFamily: "'DM Sans', sans-serif",
  },
  btn: {
    width: "100%", padding: "13px",
    background: "var(--green)",
    border: "none", borderRadius: 10,
    color: "#0d0f12", fontSize: 14,
    fontWeight: 600, cursor: "pointer",
    marginTop: 8,
    transition: "transform 0.15s, box-shadow 0.15s",
    boxShadow: "0 4px 20px rgba(99,209,148,0.25)",
    fontFamily: "'DM Sans', sans-serif",
    letterSpacing: "0.2px",
  },
  btnGhost: {
    width: "100%", padding: "13px",
    background: "transparent",
    border: "1px solid var(--border)",
    borderRadius: 10,
    color: "var(--muted)", fontSize: 14,
    fontWeight: 500, cursor: "pointer",
    marginTop: 8,
    transition: "transform 0.15s, color 0.15s",
    fontFamily: "'DM Sans', sans-serif",
  },
  error: {
    background: "rgba(248,113,113,0.1)",
    border: "1px solid rgba(248,113,113,0.25)",
    borderRadius: 8, padding: "10px 12px",
    fontSize: 13, color: "var(--red)",
    marginTop: 4, marginBottom: 4,
  },
  hint: {
    textAlign: "center", marginTop: 20,
    fontSize: 13, color: "var(--muted)",
  },
  link: {
    color: "var(--green)", cursor: "pointer",
    fontWeight: 500,
    textDecoration: "underline",
    textDecorationColor: "rgba(99,209,148,0.3)",
  },
  grid2: {
    display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12,
  },
  successIcon: {
    width: 64, height: 64, borderRadius: "50%",
    background: "rgba(99,209,148,0.12)",
    display: "flex", alignItems: "center", justifyContent: "center",
    fontSize: 28, color: "var(--green)",
    marginBottom: 20,
    border: "1px solid rgba(99,209,148,0.2)",
  },
};