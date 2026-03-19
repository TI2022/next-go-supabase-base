"use client";

import { useEffect, useState } from "react";

export type CurrentUser = {
  id: string;
  email: string;
  name?: string;
};

type UseCurrentUserResult = {
  user: CurrentUser | null;
  isLoading: boolean;
  isUnauthorized: boolean;
  errorMessage: string | null;
};

export function useCurrentUser(): UseCurrentUserResult {
  const [user, setUser] = useState<CurrentUser | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [isUnauthorized, setIsUnauthorized] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function fetchMe() {
      setIsLoading(true);
      try {
        const res = await fetch("/api/me", {
          method: "GET",
          credentials: "include",
        });

        if (res.status === 401) {
          if (!cancelled) {
            setUser(null);
            setIsUnauthorized(true);
            setErrorMessage(null);
          }
          return;
        }

        if (!res.ok) {
          const message = await res.text();
          throw new Error(message || "Failed to fetch current user");
        }

        const data = (await res.json()) as CurrentUser;
        if (!cancelled) {
          setUser(data);
          setIsUnauthorized(false);
          setErrorMessage(null);
        }
      } catch (err) {
        if (!cancelled) {
          setUser(null);
          setIsUnauthorized(false);
          setErrorMessage(err instanceof Error ? err.message : "Unknown error");
        }
      } finally {
        if (!cancelled) {
          setIsLoading(false);
        }
      }
    }

    void fetchMe();

    return () => {
      cancelled = true;
    };
  }, []);

  return { user, isLoading, isUnauthorized, errorMessage };
}

