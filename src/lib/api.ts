import type { Profile } from './calculations';

export async function fetchProfile(birthdate: string): Promise<Profile> {
  const res = await fetch('/api/profile', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ birthdate }),
  });
  if (!res.ok) {
    throw new Error(`server error: ${res.status}`);
  }
  return res.json();
}
