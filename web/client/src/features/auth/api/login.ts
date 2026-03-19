export type LoginInput = {
  email: string;
  password: string;
};

export async function login(input: LoginInput): Promise<void> {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    const message = await res.text();
    throw new Error(message || "Login failed");
  }
}

