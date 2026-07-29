import './styles/App.css'
import './styles/Form.css'
import './styles/Buttons.css'
import './styles/Links.css'
import './styles/Icons.css'
import './styles/Dropdowns.css'
import './styles/Modules.css'
import { Route, Router, Navigate } from "@solidjs/router";
import { UserProvider } from './UserContext.tsx';
import { DefaultsProvider } from "./DefaultsContext.tsx";
import ProtectedRoute from './ProtectedRoute.tsx';
import AgencyProtectedRoute from "./AgencyProtectedRoute.tsx";
import Login from "./Login.tsx";
import Home from "./Home.tsx";
import ForgottenPassword from "./ForgottenPassword.tsx";
import ResetPassword from "./ResetPassword.tsx";
import CreatorWelcome from "./CreatorWelcome.tsx";
import AgencyHome from "./AgencyHome.tsx";

function App() {
  return (
    <div class="app">
      <UserProvider>
        <DefaultsProvider>
          <Router>
            <Route path="/" component={() => <Navigate href="/login" />} />
            <Route path="/login" component={Login} />
            <Route path="/forgottenPassword" component={ForgottenPassword} />
            <Route path="/resetPassword" component={ResetPassword} />
            <Route path="/invite" component={CreatorWelcome} />
            
            <Route component={ProtectedRoute}>
              <Route path="/home" component={Home} />
            </Route>

            <Route component={AgencyProtectedRoute}>
              <Route path="/agencyHome" component={AgencyHome} />
            </Route>
          </Router>
        </DefaultsProvider>
      </UserProvider>
    </div>
  )
}

export default App
