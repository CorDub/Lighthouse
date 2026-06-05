import './styles/Login.css';
import loginText from './translations/Login.json';
import { createEffect, createRenderEffect, createSignal, Show } from 'solid-js';
import { useNavigate } from "@solidjs/router"
import { type ValueCheck } from './types/helpersTypes.ts';
import Navbar from "./Navbar.tsx"
import Errors from "./Errors.tsx"
import TextInput from "./TextInput.tsx"
import PasswordInput from "./PasswordInput.tsx"
import { useUser } from "./UserContext.tsx"
import { UserSchema } from "./schemas/user.ts";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { BASE_URL } from './helpers/config.ts';
import Text from "./Text.tsx";
import { getText } from './helpers/getText.ts';
import { useDefaults } from "./DefaultsContext.tsx";
import { type ErrorKey } from "./Errors.tsx";

function Login() {
  const [email, setEmail] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const navigate = useNavigate();
  const { user, setUser } = useUser();
  const { defaults, setDefaults } = useDefaults();

  // navigate home directly without logging in if a user is found
  createRenderEffect(() => {
    if (user()) {
      console.log("user", user())
      const resUser = user()
      if (resUser) {
        switch(resUser.role) {
          case "visitor":
            navigate("/home", {replace: true})
            break;
          case "agency":
            navigate("/agencyHome", {replace: true})
            break;
          case "creator":
            navigate("/creatorHome", {replace: true})
            break;
          case "brand":
            navigate("/brandHome", {replace: true})
            break;
          default:
            console.error("could not find this switch case")
            break;
        }
      }
    }
  })

  //change language on user logging in 
  createEffect(() => {
    const loggedInLanguage = user()?.language;
    if (loggedInLanguage) {
      setDefaults(prev => ({
        ...prev,
        lang: loggedInLanguage
      })
    )}
  })

  async function submitForAuth(e: Event) {
    try {
      e.preventDefault()

      // check for errors before sending to server
      const checks: ValueCheck[] = [
        ["email", email()],
        ["password", password()]
      ]
      const checksResults = checkForErrors(...checks)

      if (checksResults.length > 0) {
        setErrors(checksResults)
        return
      }

      // if no errors send it through
      const response = await fetch(`${BASE_URL}/api/login`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          email: email(),
          password: password()
        })
      })

      if (response.ok) {
        const data = await response.json()
        const parsedUser = UserSchema.parse(data);
        setUser(parsedUser);
        if (parsedUser) {
          switch(parsedUser.role) {
            case "visitor":
              navigate("/home", {replace: true})
              break;
            case "agency":
              console.log("agency path reached")
              navigate("/agencyHome", {replace: true})
              break;
            case "creator":
              navigate("/creatorHome", {replace: true})
              break;
            case "brand":
              navigate("/brandHome", {replace: true})
              break;
            default:
              console.error("could not find this switch case")
              break;
          }
        }
      }
    } catch(error) {
      console.error(error)
    }
  }

  async function redirectToGoogle() {
    window.location.href="/api/google"
  }

  return (
    <div class="login">
      <Navbar />
      <form
        class="entry-form"
        onSubmit={(e) => submitForAuth(e)}>

        <div class="connect-google">
          <button
            type="button"
            class="google-button clickable"
            onClick={redirectToGoogle}>
            <svg
              width="25"
              height="25"
              viewBox="0 0 20 20"
              fill="none"
              xmlns="http://www.w3.org/2000/svg">
                <g clip-path="url(#clip0_9_516)">
                <path d="M19.8052 10.2304C19.8052 9.55059 19.7501 8.86714 19.6325 8.19839H10.2002V12.0492H15.6016C15.3775 13.2912 14.6573 14.3898 13.6027 15.088V17.5866H16.8252C18.7176 15.8449 19.8052 13.2728 19.8052 10.2304Z" fill="#4285F4" data-sentry-element="path" data-sentry-source-file="GoogleIcon.tsx"></path>
                <path d="M10.1999 20.0007C12.897 20.0007 15.1714 19.1151 16.8286 17.5866L13.6061 15.0879C12.7096 15.6979 11.5521 16.0433 10.2036 16.0433C7.59474 16.0433 5.38272 14.2832 4.58904 11.9169H1.26367V14.4927C2.96127 17.8695 6.41892 20.0007 10.1999 20.0007V20.0007Z" fill="#34A853" data-sentry-element="path" data-sentry-source-file="GoogleIcon.tsx"></path>
                <path d="M4.58565 11.9169C4.16676 10.675 4.16676 9.33011 4.58565 8.08814V5.51236H1.26395C-0.154389 8.33801 -0.154389 11.6671 1.26395 14.4927L4.58565 11.9169V11.9169Z" fill="#FBBC04" data-sentry-element="path" data-sentry-source-file="GoogleIcon.tsx"></path><path d="M10.1999 3.95805C11.6256 3.936 13.0035 4.47247 14.036 5.45722L16.8911 2.60218C15.0833 0.904587 12.6838 -0.0287217 10.1999 0.000673888C6.41892 0.000673888 2.96127 2.13185 1.26367 5.51234L4.58537 8.08813C5.37538 5.71811 7.59107 3.95805 10.1999 3.95805V3.95805Z" fill="#EA4335" data-sentry-element="path" data-sentry-source-file="GoogleIcon.tsx"></path>
                </g>
                <defs data-sentry-element="defs" data-sentry-source-file="GoogleIcon.tsx"><clipPath id="clip0_9_516" data-sentry-element="clipPath" data-sentry-source-file="GoogleIcon.tsx"><rect width="20" height="20" fill="white" data-sentry-element="rect" data-sentry-source-file="GoogleIcon.tsx"></rect></clipPath></defs>
                </svg>
            <p class="google-button-text"><Text value={loginText.google} lang={defaults().lang}/></p>
          </button>
        </div>

        <div class="or">
          <p style="font-size:14px"><Text value={loginText.or} lang={defaults().lang}/></p>
        </div>

        <div class="login-password">
          <div class="form-line">
            <TextInput
              errors={errors()}
              errorsSetFn={setErrors}
              value={email()}
              valueSetFn={setEmail}
              placeholder={getText(loginText.inputEmail, defaults().lang)}
              autofocus={true}/>
          </div>
          <div class="form-line"
            style={{"margin-bottom": "0.75rem"}}>
            <PasswordInput
              errors={errors()}
              errorsSetFn={setErrors}
              value={password()}
              valueSetFn={setPassword}
              placeholder={getText(loginText.inputPassword, defaults().lang)}/>
          </div>

          <Show when={errors().length > 0} >
            <Errors
              errors={errors()}
              margin={{marginBottom: 0.75}}/>
          </Show>

          <button
            type="submit"
            class="green-button clickable"
            onClick={(e) => submitForAuth(e)}>
              <Text value={loginText.submit} lang={defaults().lang} />
          </button>
        </div>
      </form>

      <div class="lost-password">
        <a
          class="subdued-link"
          onClick={() => navigate("/forgottenPassword")}>
            <Text value={loginText.forgottenPassword} lang={defaults().lang}/>
        </a>
      </div>

    </div>
  )
}

export default Login
