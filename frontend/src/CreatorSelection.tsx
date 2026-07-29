import { createMemo, createSignal, For, onMount, Show } from "solid-js";
import TextInput from "./TextInput.tsx";
import type { ErrorKey } from "./Errors.tsx";
import AddButton from "./AddButton";
import "./styles/CreatorSelection.css";
import { BASE_URL } from "./helpers/config.ts";
import SearchResults from "./SearchResults.tsx";
import DeleteButton from "./DeleteButton.tsx";
import { useUser } from "./UserContext.tsx";
import Clipboard from "./Clipboard.tsx";

type SearchResultsCreator = {
  name: string,
  inDB: boolean,
  invite?: string
}

function CreatorSelection() {
  const { user } = useUser();
  const [errors, setErrors] = createSignal<ErrorKey[]>([]);
  const [creators, setCreators] = createSignal<SearchResultsCreator[]>([]);
  const [creatorsAvailable, setCreatorsAvailable] = createSignal<string[]>([]);
  const [creator, setCreator] = createSignal("");
  const [autocomplete, setAutocomplete] = createSignal("");
  const [searchResults, setSearchResults] = createSignal<string[]>([]);
  const searchResultsMemo = createMemo(() => 
    [...creators(), {name: searchResults()[0], inDB: true}]
  )

  onMount(
    async() => {
      const response = await fetch(`${BASE_URL}/api/users/creatorsAvailable`, {
        method:"GET",
        credentials: "include",
        headers: {"Content-Type":"application/json"}
      })

      if (response.ok) {
        const data = await response.json()
        setCreatorsAvailable(data)
      }
    }
  )

  async function getInviteLink(userId: string, name: string) {
    const response = await fetch(`${BASE_URL}/api/invite`, {
      method: "POST",
      credentials: "include",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({
        userId: userId,
        name: name
      })
    })

    if (response.ok) {
      const data = await response.json()
      return data
    }
  }

  function findCreator() {
    if (creator().length === 0) {
      setSearchResults([])
      return
    }

    let res = []
    for (const talent of creatorsAvailable()) {
      // add it if fits
      if (talent.toLowerCase().startsWith(creator().toLowerCase())) {
        // check if already selected
        let alreadySelected = false;
        for (const creatorSelected of creators()) {
          if (talent === creatorSelected.name) {
            alreadySelected = true
          }
        }
        if (alreadySelected) {
          continue
        }
        res.push(talent)
      }

      //max 5
      if (res.length >= 5) {
        break
      }
    }
    setSearchResults(res)

    // autocomplete behaviour
    if (res.length > 0) {
      const remainder = res[0].slice(creator().length); 
      setAutocomplete(creator() + remainder)
    } else {
      setAutocomplete("")
    }
  }

  async function selectCreator(newCreators: SearchResultsCreator[]) {
    // can't add an empty line
    if (creator().trim() === "") {
      return
    }

    // new creator if not in search results
    if (searchResults().length === 0 && creator().trim() !== "") {
      let userId = user()?.id
      if (!userId) {
        return
      }

      const inviteLink = await getInviteLink(userId, creator())
      setCreators([...creators(), {name: creator(), inDB: false, invite: inviteLink}])
      setCreator("");
      return
    }

    setCreators(newCreators);
    setSearchResults([]);
    setCreator("");
    setAutocomplete("");
  }

  function deselectCreator(creator: string) {
    let res = []
    for (const selectedCreator of creators()) {
      if (selectedCreator.name !== creator) {
        res.push(selectedCreator)
      }
    }
    setCreators(res);
  }

  return(
    <div class="creator-selection">
      <Show when={creators().length > 0}>
        <div class="cs-list">
          <For each={creators()}>{(creator, i) => 
            <div class="cs-list-item">
              <div class="csli-name">{creator.name}</div>
              <Show when={!creator.inDB}>
                <div class="csli-invite">
                  <p>{`This creator is not in our database. Please invite them with the following link: `}
                    {creator?.invite && <Clipboard link={creator.invite} />}
                  </p>
                </div>
              </Show>
              <div class="csli-delete">
                <DeleteButton 
                  clickFn={deselectCreator}
                  arg={creator.name}/>
              </div>
            </div>
          }</For>
        </div>
      </Show>
      <div class="cs-input">
        <div class="cs-topline">
          <div class="cs-text-input">
            <TextInput 
              errors={errors()}
              errorsSetFn={setErrors}
              value={creator()}
              valueSetFn={setCreator}
              autocompleteValue={autocomplete()}
              setAutocompleteFn={setAutocomplete}
              searchFn={findCreator}
              enterFn={selectCreator}
              enterFnArg={searchResultsMemo()}
              autofocus={false}
              name={"Add a creator"}
              bgColor={"var(--pale-green)"}/>
          </div>
          <div class="cs-add-button">
            <AddButton 
              clickFn={selectCreator} 
              arg={searchResultsMemo()}/>
          </div>
        </div>
        <div class="cs-searchresults">
          <Show when={searchResults().length != 0}>
            <SearchResults 
              searchResults={searchResults()}
              resultSelectFn={selectCreator}
              arg={creators()}/>
          </Show>
        </div>
      </div>
    </div>

  )
}

export default CreatorSelection;