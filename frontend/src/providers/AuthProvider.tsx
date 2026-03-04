import { AuthContext } from "../contexts/authContext";
import useFirebaseAuth from "../hooks/useFirebaseAuth";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const firebaseAuth = useFirebaseAuth();
  return (
    <AuthContext.Provider value={{ ...firebaseAuth }}>
      {children}
    </AuthContext.Provider>
  );
}
