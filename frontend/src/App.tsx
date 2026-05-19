import './styles/App.css'
import './styles/Form.css'
import './styles/Buttons.css'
import './styles/Links.css'
import './styles/Icons.css'
import { Route, Router, Navigate } from "@solidjs/router";
import { UserProvider } from './UserContext.tsx';
import ProtectedRoute from './ProtectedRoute.tsx';
import Login from "./Login.tsx";
import Home from "./Home.tsx";
import ForgottenPassword from "./ForgottenPassword.tsx";
import ResetPassword from "./ResetPassword.tsx";

function App() {
  return (
    <div class="app">
      <UserProvider>
        <Router>
          <Route path="/" component={() => <Navigate href="/login" />} />
          <Route path="/login" component={Login} />
          <Route path="/forgottenPassword" component={ForgottenPassword} />
          <Route path="/resetPassword" component={ResetPassword} />
          
          <Route component={ProtectedRoute}>
            <Route path="/home" component={Home} />
          </Route>
        </Router>
      </UserProvider>
    </div>
  )
}

export default App
