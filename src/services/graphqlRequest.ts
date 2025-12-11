const GRAPHQL_ENDPOINT = import.meta.env.VITE_GRAPHQL_ENDPOINT!;

export async function graphqlRequest<T, V = Record<string, any>>(query: string, variables?: V): Promise<T> {
  const token = localStorage.getItem('authToken');

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(GRAPHQL_ENDPOINT, {
    method: "POST",
    headers,
    body: JSON.stringify({ query, variables }),
  });

  const json = await res.json();
  if (json.errors) throw new Error(json.errors.map((err: any) => err.message).join(", "));
  return json.data as T;
}
