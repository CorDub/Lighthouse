import './App.css'
import { Route, Router, Navigate } from "@solidjs/router";
import { UserProvider } from './UserContext.tsx';
import ProtectedRoute from './ProtectedRoute.tsx';
import Login from "./Login.tsx";
import Home from "./Home.tsx";
import ForgottenPassword from "./ForgottenPassword.tsx";

function App() {
  return (
    <div class="app">
      <UserProvider>
        <Router>
          <Route path="/" component={() => <Navigate href="/login" />} />
          <Route path="/login" component={Login} />
          <Route path="/forgottenPassword" component={ForgottenPassword} />
          
          <Route component={ProtectedRoute}>
            <Route path="/home" component={Home} />
          </Route>
        </Router>
      </UserProvider>
    </div>
  )
}

export default App
