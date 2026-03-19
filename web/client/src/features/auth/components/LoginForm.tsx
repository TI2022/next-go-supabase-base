"use client";

type LoginFormProps = {
  email: string;
  password: string;
  isSubmitting: boolean;
  errorMessage: string | null;
  onEmailChange: (value: string) => void;
  onPasswordChange: (value: string) => void;
  onSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
};

export function LoginForm({
  email,
  password,
  isSubmitting,
  errorMessage,
  onEmailChange,
  onPasswordChange,
  onSubmit,
}: LoginFormProps) {
  return (
    <form
      onSubmit={onSubmit}
      className="w-full max-w-sm rounded-lg border border-gray-200 bg-white p-6 shadow-sm"
    >
      <h1 className="mb-6 text-2xl font-bold text-gray-900">Login</h1>

      <label className="mb-2 block text-sm font-medium text-gray-700" htmlFor="email">
        Email
      </label>
      <input
        id="email"
        type="email"
        value={email}
        onChange={(e) => onEmailChange(e.target.value)}
        required
        className="mb-4 w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900"
        placeholder="demo@example.com"
      />

      <label className="mb-2 block text-sm font-medium text-gray-700" htmlFor="password">
        Password
      </label>
      <input
        id="password"
        type="password"
        value={password}
        onChange={(e) => onPasswordChange(e.target.value)}
        required
        className="mb-4 w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900"
        placeholder="plain-text-demo-password"
      />

      {errorMessage ? <p className="mb-4 text-sm text-red-600">{errorMessage}</p> : null}

      <button
        type="submit"
        disabled={isSubmitting}
        className="w-full rounded-md bg-black px-4 py-2 text-white disabled:opacity-60"
      >
        {isSubmitting ? "Signing in..." : "Sign in"}
      </button>
    </form>
  );
}

