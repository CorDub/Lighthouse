import { createSignal, createEffect } from "solid-js";
import { Portal } from "solid-js/web";
import Dropdown from "./Dropdown.tsx";
import { BASE_URL } from "./helpers/config.ts";
import { useUser } from "./UserContext.tsx";

function LanguagePicker() {
  const [lang, setLang] = createSignal("EN")
  const [isDropdownOpen, setDropDownOpen] = createSignal(false)
  const [anchor, setAnchor] = createSignal({})
  const [isClassOpenAdded, setClassOpenAdded] = createSignal(false)
  const { user } = useUser()

  function openLanguagePickerDropdown(e: MouseEvent) {
    const rect = (e.currentTarget as HTMLButtonElement).getBoundingClientRect();
    setAnchor({ top: rect.top, right: rect.right })
    setDropDownOpen(true)
  }

  function closeDropDown() {
    setClassOpenAdded(false)
    setTimeout(() => setDropDownOpen(false), 250)
  }

  async function changeLanguage(langCode: string) {
    const newPreferredLanguage = lang() !== langCode
    setLang(langCode)
    closeDropDown()
    if (newPreferredLanguage) {
      try {
        if (!user()) {
          return
        }

        const response = await fetch(`${BASE_URL}/api/users/${user()?.id}`, {
          method: "POST",
          credentials: "include",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify({
            language: langCode,
          })
        })

        if (response.ok) {
          const data = await response.json()
          console.log(data)
        }

      } catch (error) {
        console.error(error)
      }
    }
  }

  createEffect(() => {
    if (isDropdownOpen()) {
      requestAnimationFrame(() => {
        setClassOpenAdded(true)
      })
    }
  })

  return(
    <div class="language-picker">
      <button
        class="black-button clickable"
        onClick={(e) => openLanguagePickerDropdown(e)}>
        {lang()}  
      </button>
      <Portal>
        <Dropdown
          openStatus={isDropdownOpen()}
          setOpenStatus={closeDropDown}
          anchor={anchor()}>
            <div class="black-dropdown"
              classList={{"open" : isClassOpenAdded()}}>
              <button class="grey-button clickable"
                onClick={() => changeLanguage("EN")}>EN</button>
              <button class="grey-button clickable"
                onClick={() => changeLanguage("ES")}>ES</button>
            </div>
        </Dropdown>
      </Portal>
    </div>
  )
}

export default LanguagePicker