import useAuth from "./hooks/useAuth";

function App() {
  const { signInWithGoogle } = useAuth();

  return (
    <>
      <button onClick={signInWithGoogle}>Login</button>
    </>
  );
}

export default App;
