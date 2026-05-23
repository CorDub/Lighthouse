import { createSignal, createEffect } from "solid-js";
import { Portal } from "solid-js/web";
import Dropdown from "./Dropdown.tsx";

function LanguagePicker() {
  const [lang, setLang] = createSignal("EN")
  const [isDropdownOpen, setDropDownOpen] = createSignal(false)
  const [anchor, setAnchor] = createSignal({})
  const [isClassOpenAdded, setClassOpenAdded] = createSignal(false)

  function openLanguagePickerDropdown(e: MouseEvent) {
    const rect = (e.currentTarget as HTMLButtonElement).getBoundingClientRect();
    setAnchor({ top: rect.top, right: rect.right })
    setDropDownOpen(true)
  }

  function closeDropDown() {
    setClassOpenAdded(false)
    setTimeout(() => setDropDownOpen(false), 250)
  }

  function changeLanguage(langCode: string) {
    setLang(langCode)
    closeDropDown()
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