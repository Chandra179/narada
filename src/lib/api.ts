import type { Profile } from './calculations';

export async function fetchProfile(birthdate: string, birthtime?: string): Promise<Profile> {
  const body = JSON.stringify({ birthdate, birthtime });
  console.log('[api] POST /api/profile', body, 'birthtime param:', JSON.stringify(birthtime));
  const res = await fetch('/api/profile', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body,
  });
  if (!res.ok) {
    throw new Error(`server error: ${res.status}`);
  }
  return res.json();
}
