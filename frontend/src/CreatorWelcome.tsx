import "./styles/CreatorWelcome.css";
import { createSignal, Show } from "solid-js";
import { useSearchParams, useNavigate } from "@solidjs/router";
import Navbar from "./Navbar.tsx";
import Errors from "./Errors.tsx";
import TextInput from "./TextInput.tsx";
import PasswordInput from "./PasswordInput.tsx";
import { checkForErrors } from "./helpers/checkForErrors.ts";
import { type ValueCheck } from "./types/helpersTypes.ts";
import { type ErrorKey } from "./Errors.tsx";
import creaWelText from "./translations/CreatorWelcome.json";
import Text from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";
import { getText } from "./helpers/getText.ts";
import { type LanguageCode } from "./types/langTypes.ts";
import { isLanguageCode } from "./types/langTypes.ts";
import LockableTextInput from "./LockableTextInput.tsx";
import { BASE_URL } from "./helpers/config.ts";
import { useUser } from "./UserContext.tsx";
import { UserSchema } from "./schemas/user.ts";

function CreatorWelcome() {
  const [searchParams] = useSearchParams<{ token: string, name: string }>()
  const [name, setName] = createSignal(searchParams.name ?? "")
  const [email, setEmail] = createSignal("")
  const [password, setPassword] = createSignal("")
  const [language, setLanguage] = createSignal<LanguageCode>("en")
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const { defaults } = useDefaults();
  const navigate = useNavigate();
  const { setUser } = useUser();

  async function createProfile(e: Event) {
    e.preventDefault();

    // check we got all values needed
    const checks: ValueCheck[] = [
      ["name", name()],
      ["email", email()],
      ["password", password()]
    ]
    const checkResults = checkForErrors(...checks)

    if (checkResults.length > 0) {
      setErrors(checkResults)
      return
    }

    // check browser language, set en by default
    const browserLanguage = navigator.language.split('-')[0];
    if (!isLanguageCode(browserLanguage)) {
      setLanguage("en")
    } else {
      setLanguage(browserLanguage)
    }

    const response = await fetch(`${BASE_URL}/api/users/creatorInvite`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        name: name(),
        email: email(),
        password: password(),
        language: language(),
        role: "creator"
      })
    })

    if (response.ok) {
      const data = await response.json();
      const parsedUser = UserSchema.parse(data);
      setUser(parsedUser);
      console.log("parseUser", parsedUser)
      navigate("/socialNetworks", { replace: true });
    }
  }

  return (
    <div class="creator-welcome">
      <Navbar />
      <form
        class="cw-form"
        onSubmit={(e) => createProfile(e)}>
        <h1 class="cw-heading">
          <Text
            value={creaWelText.welcomeHeading}
            lang={defaults().lang} />
        </h1>
        <p class="cw-message">
          <Text
            value={creaWelText.welcomeMessage}
            lang={defaults().lang} />
        </p>
        <div class="cw-input">
          <LockableTextInput 
            errors={errors()}
            errorsSetFn={setErrors}
            value={name()}
            valueSetFn={setName}/>
        </div>
        <div class="cw-input">
          <TextInput
            errors={errors()}
            errorsSetFn={setErrors}
            value={email()}
            valueSetFn={setEmail}
            placeholder={getText(creaWelText.emailInput, defaults().lang)}
            // name={getText(creaWelText.emailInput, defaults().lang)}
            bgColor={"var(--white)"}/>
        </div>
        <div class="cw-input">
          <PasswordInput
            errors={errors()}
            errorsSetFn={setErrors}
            value={password()}
            valueSetFn={setPassword}
            placeholder={getText(creaWelText.passwordInput, defaults().lang)}/>
        </div>
      
        <Show when={errors().length > 0}>
          <Errors 
            errors={errors()}
            margin={{
              marginTop: 0,
              marginBottom: 2
            }}/>
        </Show>
        
        <button class="green-button clickable" onClick={(e) => createProfile(e)}>
          <Text
            value={creaWelText.confirmButton}
            lang={defaults().lang} />
        </button>
      </form>

    </div>
  )
}

export default CreatorWelcome;
