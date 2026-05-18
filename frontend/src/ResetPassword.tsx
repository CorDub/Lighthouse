import { createSignal, Show } from "solid-js";
import { useSearchParams } from "@solidjs/router";
import Navbar from "./Navbar.tsx"

function ResetPassword() {
  const [searchParams] = useSearchParams()
  const paramsToken = searchParams.token
  const [newPassword, setNewPassword] = createSignal("")
  const [confirmed, setConfirmed] = createSignal(true)

  async function changePassword(e: Event) {
    e.preventDefault();

    //check if new password is at least 8 characters long
    if (newPassword().length < 8) {
      console.error("Passwords need to be at least 8 characters long")
      return
    } 

    // send the new password
    const response = await fetch('/api/changePassword', {
      method: "POST",
      credentials: "include",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        password: newPassword(),
        token: paramsToken,
      })
    })

    if (response.ok) {
      setConfirmed(true)
    }
  }

  return (
    <div class="reset-password">
      <Navbar />
      <Show when={!confirmed()}>
        <form 
          class="confirm-email-form"
          onSubmit={(e) => changePassword(e)}>
          <input
            class="form-input"
            type="text"
            placeholder="Enter your new password"
            onChange={(e) => setNewPassword(e.target.value)}/>
          <button
            class="green-button clickable"
            onClick={(e) => changePassword(e)}>
              Submit
          </button>
        </form>
      </Show>
      <Show when={confirmed()}>
        <div class="respas-confirmation">
          <p style={{ "margin-bottom": "1rem" }}>Your new password has been saved</p>
          <button class="green-button clickable">Return to login</button>
        </div>
      </Show>
    </div>
  )
}

export default ResetPassword;