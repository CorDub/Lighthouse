import { createSignal, createEffect } from "solid-js";
import { Portal } from "solid-js/web";
import Dropdown from "./Dropdown.tsx";
import { BASE_URL } from "./helpers/config.ts";
import { useUser } from "./UserContext.tsx";
import type { LanguageCode } from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";

function LanguagePicker() {
  const [isDropdownOpen, setDropDownOpen] = createSignal(false)
  const [anchor, setAnchor] = createSignal({})
  const [isClassOpenAdded, setClassOpenAdded] = createSignal(false)
  const { user } = useUser()
  const { defaults, setDefaults } = useDefaults()

  function openLanguagePickerDropdown(e: MouseEvent) {
    const rect = (e.currentTarget as HTMLButtonElement).getBoundingClientRect();
    setAnchor({ top: rect.top, right: rect.right })
    setDropDownOpen(true)
  }

  function closeDropDown() {
    setClassOpenAdded(false)
    setTimeout(() => setDropDownOpen(false), 250)
  }

  async function changeLanguage(langCode: LanguageCode) {
    const newPreferredLanguage = defaults().lang !== langCode
    setDefaults(prev => ({
      ...prev,
      lang: langCode
    }))
    console.log("langCode", langCode)
    console.log("defaults.language", defaults().lang)
    closeDropDown()
    if (newPreferredLanguage) {
      try {
        if (!user()) {
          return
        }

        const response = await fetch(`${BASE_URL}/api/users/${user()?.id}/language`, {
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
        {defaults().lang.toUpperCase()}
      </button>
      <Portal>
        <Dropdown
          openStatus={isDropdownOpen()}
          setOpenStatus={closeDropDown}
          anchor={anchor()}>
            <div class="black-dropdown"
              classList={{"open" : isClassOpenAdded()}}>
              <button class="grey-button clickable"
                onClick={() => changeLanguage("en")}>EN</button>
              <button class="grey-button clickable"
                onClick={() => changeLanguage("es")}>ES</button>
            </div>
        </Dropdown>
      </Portal>
    </div>
  )
}

export default LanguagePicker
