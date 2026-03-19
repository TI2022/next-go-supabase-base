type WelcomeCardProps = {
  userName: string;
  userEmail: string;
};

export function WelcomeCard({ userName, userEmail }: WelcomeCardProps) {
  return (
    <section className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
      <h1 className="text-2xl font-bold text-gray-900">ようこそ, {userName}</h1>
      <p className="mt-2 text-sm text-gray-600">ログイン中: {userEmail}</p>
    </section>
  );
}

