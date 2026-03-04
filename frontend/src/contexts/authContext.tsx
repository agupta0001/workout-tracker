import { createContext, type Context } from "react";
import type { User } from "../types/user";
import type { UserCredential } from "firebase/auth";

export interface IAuth {
  uid: string;
  email: string | null;
  name: string | null;
  avatar: string | null;
  token: string | null;
}

interface IAuthContext {
  auth: IAuth | null;
  loading: boolean;
  signInWithGoogle: () => Promise<UserCredential | void>;
  signOut: () => Promise<void>;
  currentUser: User | null;
}

export const AuthContext: Context<IAuthContext> = createContext<IAuthContext>({
  auth: null,
  loading: true,
  signInWithGoogle: async () => {},
  signOut: async () => {},
  currentUser: null,
});
