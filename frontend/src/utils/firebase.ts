import { initializeApp } from "firebase/app";
import { getAuth, type Auth } from "firebase/auth";

import firebaseConfig from "../utils/firebaseConfig";

const firebaseApp = initializeApp(firebaseConfig);
const firebaseAuth: Auth = getAuth(firebaseApp);

export { firebaseApp, firebaseAuth };
