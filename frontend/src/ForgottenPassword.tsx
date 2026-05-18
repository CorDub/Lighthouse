import "./ForgottenPassword.css";
import { createSignal, Show } from "solid-js";
import Navbar from "./Navbar";
import { UserSchema } from "./schemas/user.ts"

function ForgottenPassword() {
  const [email, setEmail] = createSignal("")
  const [resetAccepted, setResetAccepted] = createSignal(false)

  async function resetPassword(e: Event) {
    try {
      e.preventDefault()

      //check if the email sent is a valid password
      const validityCheck = UserSchema.unwrap().shape.email.safeParse(email())
      if (!validityCheck.success) {
        console.error(validityCheck.error.message)
        return
      }

      const response = await fetch("/api/checkPassword", {
        method: "POST",
        credentials: "include",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({
          email: email()
        })
      })

      if (response.ok) {
        setResetAccepted(true)
      }
    } catch(error) {
      console.error(error)
    }
  }

  return (
    <div class="forgotten-password">
      <Navbar/>
      <Show when={resetAccepted()}>
        <div>
          <p>{`A reset email has been sent to ${email()}`}</p>
        </div>
      </Show>
      <Show when={!resetAccepted()}>
        <form 
          class="confirm-email-form"
          onSubmit={(e) => resetPassword(e)}>
          <input
            class="form-input"
            type="text"
            placeholder="Enter email address"
            onChange={(e) => setEmail(e.target.value)}/>
          <button 
            class="green-button clickable"
            onClick={(e) => resetPassword(e)}>
              Submit
          </button>
        </form>
      </Show>
    </div>
  )
}

export default ForgottenPassword;