import { useCallback, useEffect, useState } from "react";
import type { IAuth } from "../contexts/authContext";
import type { User } from "../types/user";
import { firebaseAuth } from "../utils/firebase";
import {
  signOut as firebaseSignOut,
  GoogleAuthProvider,
  onAuthStateChanged,
  signInWithPopup,
  type User as FirebaseUser,
} from "firebase/auth";
import http from "../lib/http";
import type { LoginResponse } from "../types/auth";

const formatAuthState = (user: FirebaseUser): IAuth => ({
  uid: user.uid,
  email: user.email,
  name: user.displayName,
  avatar: user.photoURL,
  token: null,
});

function useFirebaseAuth() {
  const [auth, setAuth] = useState<IAuth | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [currentUser, setCurrentUser] = useState<User | null>(null);

  useEffect(() => {
    const handle = setInterval(async () => {
      const user = firebaseAuth.currentUser;
      if (user) await user.getIdToken(true);
    }, 10 * 60 * 1000);

    return () => clearInterval(handle);
  }, []);

  const handleAuthChange = useCallback(
    async (authState: FirebaseUser | null) => {
      if (!authState) {
        setLoading(false);
        return;
      }

      const formatedAuth: IAuth = formatAuthState(authState);

      try {
        formatedAuth.token = await authState.getIdToken();

        const loginResponse = await http.post<LoginResponse>(
          "/login",
          undefined,
          true
        );

        localStorage.setItem("jwtToken", loginResponse.token);
        setAuth(formatedAuth);
        setCurrentUser((prevUser) => ({ ...prevUser, ...loginResponse.user }));
      } catch (error: unknown) {
        console.error(
          "Error fetching ID token:",
          error instanceof Error ? error.message : "Login Failed"
        );
      } finally {
        setLoading(false);
      }
    },
    []
  );

  useEffect(() => {
    setLoading(true);
    const unsubscribe = onAuthStateChanged(firebaseAuth, handleAuthChange);
    return () => unsubscribe();
  }, [handleAuthChange]);

  const signOut = useCallback(async () => {
    await firebaseSignOut(firebaseAuth);
  }, []);

  const signInWithGoogle = async () => {
    setLoading(true);
    const provider = new GoogleAuthProvider();
    provider.setCustomParameters({ access_type: "offline" });

    try {
      const userCredential = await signInWithPopup(firebaseAuth, provider);

      return userCredential;
    } catch (error) {
      console.error("Error during Google sign-in:", error);
      setLoading(false);
    }
  };

  return {
    auth,
    loading,
    currentUser,
    signOut,
    signInWithGoogle,
  };
}

export default useFirebaseAuth;
