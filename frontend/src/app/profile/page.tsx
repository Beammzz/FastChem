"use client";

import { useEffect, useState } from "react";
import ProfileContent from "@/components/ProfileContent";

export default function ProfilePage() {
  const [username, setUsername] = useState<string | null>(null);

  useEffect(() => {
    // Extract username from the URL path: /profile/<username>/
    const parts = window.location.pathname.split("/").filter(Boolean);
    // parts = ["profile", "<username>"]
    if (parts.length >= 2 && parts[0] === "profile") {
      setUsername(decodeURIComponent(parts[1]));
    }
  }, []);

  if (!username) {
    return (
      <div className="min-h-screen bg-[#0a0a1a] text-white flex items-center justify-center">
        <div className="w-10 h-10 border-3 border-violet-400 border-t-transparent rounded-full animate-spin" />
      </div>
    );
  }

  return <ProfileContent username={username} />;
}
