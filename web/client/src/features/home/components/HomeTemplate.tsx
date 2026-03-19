"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { useCurrentUser } from "@/src/features/auth/hooks/useCurrentUser";
import { WelcomeCard } from "@/src/features/home/components/WelcomeCard";

const DEMO_ITEMS = [
  { id: 1, title: "MVP Task 1", description: "ログイン後ホーム表示を確認する" },
  { id: 2, title: "MVP Task 2", description: "認証状態に応じて画面遷移する" },
  { id: 3, title: "MVP Task 3", description: "次フェーズで API を拡張する" },
];

export function HomeTemplate() {
  const router = useRouter();
  const { user, isLoading, isUnauthorized, errorMessage } = useCurrentUser();

  useEffect(() => {
    if (isUnauthorized) {
      router.replace("/login");
    }
  }, [isUnauthorized, router]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 text-gray-600">
        Loading...
      </div>
    );
  }

  if (errorMessage) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50 p-6">
        <p className="rounded-md bg-red-50 p-4 text-sm text-red-700">{errorMessage}</p>
      </div>
    );
  }

  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <main className="mx-auto flex w-full max-w-3xl flex-col gap-6">
        <WelcomeCard userName={user.name ?? "User"} userEmail={user.email} />

        <section className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
          <h2 className="mb-4 text-lg font-semibold text-gray-900">ダミーデータ一覧</h2>
          <ul className="space-y-3">
            {DEMO_ITEMS.map((item) => (
              <li key={item.id} className="rounded-md border border-gray-100 p-3">
                <p className="font-medium text-gray-900">{item.title}</p>
                <p className="text-sm text-gray-600">{item.description}</p>
              </li>
            ))}
          </ul>
        </section>
      </main>
    </div>
  );
}

