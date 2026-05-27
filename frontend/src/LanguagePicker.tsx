import { createSignal, createEffect, For } from "solid-js";
import { Portal } from "solid-js/web";
import Dropdown from "./Dropdown.tsx";
import { BASE_URL } from "./helpers/config.ts";
import { useUser } from "./UserContext.tsx";
import type { LanguageCode } from "./Text.tsx";
import { useDefaults } from "./DefaultsContext.tsx";
import Alert from "./Alert.tsx";
import { LANGUAGE_CODES } from "./Text.tsx";

function LanguagePicker() {
  const [isDropdownOpen, setDropDownOpen] = createSignal(false)
  const [anchor, setAnchor] = createSignal({})
  const [isClassOpenAdded, setClassOpenAdded] = createSignal(false)
  const { user } = useUser()
  const { defaults, setDefaults } = useDefaults()
  const [isAlertOpen, setAlertOpen] = createSignal(false)
  const [languageCodesList, setLanguageCodeList] = createSignal<LanguageCode[]>([...LANGUAGE_CODES])

  // change the order of the list to always have the currently chosen first
  createEffect(() => {
    if (defaults().lang == "en") {
      return
    }

    let newOrderedList: LanguageCode[] = []
    newOrderedList.push(defaults().lang)
    for (const lang of languageCodesList()) {
      if (newOrderedList.includes(lang)) {
        continue
      }
      newOrderedList.push(lang)
    }

    setLanguageCodeList(newOrderedList)
  })

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
    closeDropDown()
    if (newPreferredLanguage) {
      try {
        if (!user()) {
          return
        }

        const response = await fetch(`${BASE_URL}/api/users/${user()?.id}/language`, {
          method: "PATCH",
          credentials: "include",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify({
            language: langCode,
          })
        })

        if (response.ok) {
          setAlertOpen(true)
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
              <For each={languageCodesList()}>
                {(lang, _) => 
                <button class="grey-button clickable"
                  onClick={() => changeLanguage(lang)}>
                  {lang.toUpperCase()}
                </button>}
              </For>
            </div>
        </Dropdown>
      </Portal>
      <Alert 
        message={"languageChanged"}
        type={"confirmation"} 
        alertOpen={isAlertOpen()}
        setAlertOpen={setAlertOpen}/>
    </div>
  )
}

export default LanguagePicker
